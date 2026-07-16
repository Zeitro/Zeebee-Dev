# Borrower Property Transaction 文件内容提取 方案设计
> 来源：tasks/20260715_borrower_property_transaction_v2_plan.md（重组为规范设计）

## 1. 背景与目标

### 1.1 业务背景

当前上传交易文件后，文件本身已保存，但其中蕴含的信息（如 `purchase_amount`、escrow 公司与 officer 信息）需要用户在界面上手动填写。这一步既增加录入成本，也容易出错或遗漏。

### 1.2 目标

在文件上传事务提交之后，依据 `document_templates` 中配置的 `file_parse_config`，自动识别文件类型并调用对应解析器提取关键字段，回写到 `borrower_property_transaction` 或其关联的 escrow 行，减少用户手动录入。

### 1.3 范围边界

- 仅支持两种解析器：`rpa`（文本路径，提取 `purchase_amount`）与 `escrow_info`（视觉路径，提取 escrow 行字段）。
- 6 个新增交易列（`title_holder`、`contract_date`、`emd_due_date`、`contingency_removal_date`、`tc`、`cost_allocation`）**仅手动填写**，任何解析器都不触碰。
- escrow `checklist` 由 PUT 请求体原样持久化，无重算逻辑。
- 解析失败采用 **fail-silent** 策略：文件已落库即视为上传成功，解析错误仅记日志、回写缺失由用户手动补全。
- 本次仅接入 openai provider；`claude` provider 配置保留但不接线；**不**改动 `lenderdiscover`。

## 2. 总体方案

锁定决策清单（全部来自源文档的 Locked decisions，逐条保留）：

1. **解析范围。** 仅两个解析器：
   - `rpa`（文本路径）→ 仅 `purchase_amount`。移植 `just-run/parse_pdf/main.go` 中的 `readFirstPage` + `extractPurchasePrice`。
   - `escrow_info`（视觉路径）→ escrow 行（`no`、`escrow_company`、`officer_name`、`officer_phone`、`officer_email`）。移植 `renderFirstPage` + `extractEscrowFields`。
   - 6 个新增交易列（`title_holder`、`contract_date`、`emd_due_date`、`contingency_removal_date`、`tc`、`cost_allocation`）**仅手动**，无解析器触及。

2. **Merge rule。** 当提取值为 `0` / `null` / `""` 时保留原值（跳过）。绝不以空提取覆盖真实值。

3. **Escrow（1:N，当前 1，未来 N）。** 解析路径：按 `borrower_property_transaction_id` 查找（取第一行，最小 id）；不存在则创建，存在则更新该首行。保留 `no` 的 btree 索引以备未来 N 场景。

4. **同步解析，位于事务外。** 上传先提交文件事务，**之后**再运行解析。提取字段的回写发生在**独立事务**中。解析失败为 **fail-silent**：仅记日志、返回成功（文件已保存；用户手动编辑）。

5. **LLM client。** 新建共享 package，使用 `openai-go/v3` SDK（已是依赖），支持文本 + 多模态/vision。定义 `Provider` 接口以便未来 `claude` 无需改动调用方即可接入。配置中保留两个 provider；当前仅接线 `TransactionFileParseConf`（openai）。**不要**改动 `lenderdiscover`。

6. **PDF。** pdfium WASM pool 单例（webassembly、wazero、无 CGO），`MaxTotal=5` 硬编码。通过既有 `common/aws_s3` 从 S3 下载，缓冲至内存，上限 **10 MB**（信任 `Content-Length`，超过 10 MB 拒绝），向 pdfium 传入带已知大小的 `bytes.Reader`。

7. **Escrow `checklist`。** 手动，从 PUT 请求体原样持久化。无重算 helper。由 parser 创建的 escrow → `checklist = {}`。

8. **PUT escrow 同步 = 全量替换。** 删除请求列表中缺失 id 的行；按 id 更新；无 id 的行新建。

9. **GET detail** 预加载并返回 escrows 与 6 个新列。

10. **DELETE transaction** 级联删除 escrows。

11. **Template SQL** 为针对单一全局 `document_templates` 行的一次性部署迁移（已核实：该表为全局，无 org scope）。

## 3. 数据模型

### 3.1 escrow 表与 6 个新列

迁移脚本 `deploy/zeebee/sql/20260714173922_property_transaction_doc.sql` 已起草，包含：
- 6 个新增列：`title_holder`、`contract_date`、`emd_due_date`、`contingency_removal_date`、`tc`、`cost_allocation`。
- escrow 表：`property_transaction_escrow`。
- 针对 `document_templates` 全局行的 UPDATE（写入 `file_parse_config`）。

校验点：列类型须与 model 对齐——`contract_date` 等为可空 date，model 使用 `*sql.NullTime`。

> 源文档未提供完整 SQL DDL 原文；如需具体列定义与约束，以迁移脚本文件为单一事实来源。

### 3.2 document_templates 的 file_parse_config

`document_templates` 表的全局行上将承载 `file_parse_config`，由 create logic 写入到 doc metadata，upload logic 在解析时从 `doc.Metadata` 读取以判断应派发哪个 parser（`rpa` / `escrow_info`，或两者皆无则跳过）。

## 4. 接口设计

受影响的 5 个端点及其改动范围：

| 端点 | 改动 |
|------|------|
| create | 无改动（已将 `file_parse_config` 写入 doc metadata）；create 阶段不创建 escrow |
| upload-file | 文件事务提交**之后**：读 `file_parse_config` → 从 S3 下载（10 MB 上限）→ 按类型派发 `rpa` / `escrow_info` → 独立事务回写；整体 `recover`，失败仅 `logx.Error` 并返回成功 |
| GET detail | 预加载并返回 escrows + 6 个新列 |
| PUT update | 事务内全量替换 escrows：删除缺失 id 行 / 按 id 更新 / 无 id 新建；逐行持久化请求中的 `checklist` |
| DELETE | 事务删除前/同时级联删除该 transaction 的 escrows |

### 4.1 Proto 新增字段

- `domains/zeebee/proto/entity.proto` — `BorrowerPropertyTransactionEntity`：新增 `title_holder`、`contract_date`、`emd_due_date`、`contingency_removal_date`、`tc`、`cost_allocation`、`repeated PropertyTransactionEscrowEntity escrows`。新增 message `PropertyTransactionEscrowEntity`（id、borrower_property_transaction_id、no、escrow_company、officer_name、officer_phone、officer_email、checklist）。
- `app/borrower/rpc/pb/borrower_borrower_property_transaction.proto` — create 与 update 请求 data：新增 `repeated PropertyTransactionEscrowData escrows`（id 可选）。

> Proto 改完后须执行 `make create-service borrower` 重生成。

## 5. 核心流程

### 5.1 upload 后 commit-then-parse

1. 文件事务提交落库（`file.FileUrl` 已生成）。
2. 读取 `doc.Metadata` 中的 `file_parse_config`，判断是否配置 parser。
3. 若配置：经 `common/aws_s3` 下载文件字节，内存缓冲上限 10 MB（信任 `Content-Length` 预检，`io.LimitReader` 兜底）。
4. 按 parser 类型派发，结果在**独立事务**中回写。
5. 全流程 `recover` 包裹；任何错误 → `logx.Error`，返回成功（fail-silent）。

### 5.2 rpa 文本路径

- 调用 `pdfparse.ExtractRPA(ctx, pdfBytes) (float64, error)`（内部走 pdfium 文本 + LLM 文本提取）。
- 若 price > 0：在独立事务更新 `purchase_amount`。
- 提取值为 0 → 按 merge rule 跳过（保留原值）。

### 5.3 escrow_info 视觉路径

- 调用 `pdfparse.ExtractEscrow(ctx, pdfBytes) (EscrowFields, error)`（pdfium 渲染首页 PNG → LLM vision）。
- Upsert escrow：按 `borrower_property_transaction_id` 查找首行（min id），无则创建、有则更新该行。
- 空字段（`""` / `null`）按 merge rule 跳过；parser 新建的 escrow → `checklist = {}`。

### 5.4 merge / 回写流程

- 回写与文件上传事务完全解耦：上传事务先提交，回写为独立事务。
- 任意一步解析或回写失败均不影响上传响应；失败仅入日志，用户手动补全。

## 6. 风险与权衡

- **fail-silent**：解析失败对用户不可见，可能造成字段静默未回写。权衡为优先保证上传可用性，代价是用户偶尔需要手动补录；以日志支撑排查。
- **10 MB cap**：超过 10 MB 的 PDF 会被拒绝解析（信任 `Content-Length` 预检 + `io.LimitReader` 兜底）。权衡为控制内存占用，代价是大文件不进入自动提取。
- **pdfium pool 上限 `MaxTotal=5`**：硬编码上限，并发解析超过 5 路时将排队/阻塞。权衡为避免 WASM 实例失控，代价是高并发下吞吐受限。
- **配置缺失**：`LlmConf` / `LlmBizConf` / `LlmProvider` 当前无 `yaml:` tag，无法 unmarshal，必须补齐后才能从 yaml 加载；缺失会导致 ServiceContext 初始化失败。
- **并发**：escrow 1:N 当前按“首行（min id）”upsert，并发上传同一 transaction 的多个 escrow_info 文件可能竞争同一行；未来 N 场景通过 `no` btree 索引扩展。

## 7. 实现要点 / 测试

### 7.1 文件清单

**New — 共享 infra**

- `domains/zeebee/llm/client.go` — `Provider` 接口：
  ```go
  type Provider interface {
      ChatJSON(ctx, system, user string) (string, error)   // text
      VisionJSON(ctx, system, text string, imgPNG []byte) (string, error)
  }
  func NewClient(cfg config.LlmConf, biz config.LlmBizConf) (Provider, error)
  ```
  Factory resolves `biz.Provider` against `cfg.Providers` by name.
- `domains/zeebee/llm/openai.go` — openai-go impl of `Provider` (chat +
  image content part), request timeout 60s, base-url `/v1` normalization
  (copy from reference).
- `app/borrower/rpc/internal/pdfparse/parser.go` — pdfium pool singleton
  (`Init` once, `MaxTotal=5`) + `ExtractRPA(ctx, pdfBytes) (float64, error)`
  + `ExtractEscrow(ctx, pdfBytes) (EscrowFields, error)`. Uses injected
  `llm.Provider` (testable). S3 download + 10 MB cap lives in the upload
  logic, not here (parser takes bytes).

**New — repo**

- `domains/zeebee/repo/property_transaction_escrow.go` —
  `PropertyTransactionEscrowRepo { *gormx.BaseRepo[*model.PropertyTransactionEscrow] }`
  + `FindByTransactionId`, `FindFirstByTransactionId` (min id), bulk
  delete-by-ids. Register in the `AllRepos` aggregate + its constructor.

**Config**

- `domains/zeebee/config/llm.go` — **add `yaml:` tags** to `LlmConf`,
  `LlmProvider`, `LlmBizConf` (currently none → won't unmarshal).
- `app/borrower/config/config.go` — `LlmConf` field already present (line 52).
  Confirm import + that `RpcConf` (or top-level) carries it through to
  ServiceContext.
- `app/borrower/rpc/.../etc/*.yaml` — add the `LlmConf` block (template from
  task doc).

**ServiceContext**

- `svc/service_context.go` — add `LLM llm.Provider` and `PdfPool` (or keep
  pool inside `pdfparse` package as a singleton init'd from ServiceContext).
  Init `llm.NewClient(cfg.LlmConf, cfg.TransactionFileParseConf)` in
  `NewServiceContext`.

**Proto (then `make create-service borrower` to regen)**

- `domains/zeebee/proto/entity.proto` — `BorrowerPropertyTransactionEntity`:
  add `title_holder`, `contract_date`, `emd_due_date`,
  `contingency_removal_date`, `tc`, `cost_allocation`, `repeated
  PropertyTransactionEscrowEntity escrows`. New message
  `PropertyTransactionEscrowEntity` (id, borrower_property_transaction_id, no,
  escrow_company, officer_name, officer_phone, officer_email, checklist).
- `app/borrower/rpc/pb/borrower_borrower_property_transaction.proto` —
  create + update request data: add `repeated PropertyTransactionEscrowData
  escrows` (id optional).

**Logic**

- `create_logic.go` — no change (already writes `file_parse_config` into doc
  metadata). No escrow at create.
- `upload_document_file_logic.go` — **after** the file tx commits: read
  `file_parse_config` from `doc.Metadata`; if parser set, download file bytes
  from S3 (`file.FileUrl`, 10 MB cap), dispatch:
  - `rpa` → `pdfparse.ExtractRPA` → if price > 0, update
    `purchase_amount` in a new tx (skip 0).
  - `escrow_info` → `pdfparse.ExtractEscrow` → upsert escrow (create if none /
    update first by transaction id); skip empty fields per merge rule.
  - Wrap whole parse in recover; on any error → `logx.Error`, return success.
- `detail_logic.go` — load escrows via repo, map to proto entity.
- `update_logic.go` — full-replace escrows inside the tx: delete missing ids,
  update by id, create no-id; persist each escrow's `checklist` from request.
- `delete_logic.go` — delete escrows for the transaction before/with tx delete.

**Migration**

- `deploy/zeebee/sql/20260714173922_property_transaction_doc.sql` — already
  drafted (6 columns + escrow table + template UPDATE). Verify column types
  match the model (`contract_date` etc. as nullable date — model uses
  `*sql.NullTime`).

### 7.2 构建顺序

1. Config: yaml tags + config.go wiring + ServiceContext LLM init.
2. `domains/zeebee/llm` (Provider + openai impl) + unit test with stub.
3. `pdfparse` package (pool + RPA + escrow extractors) + test with mock Provider.
4. Escrow repo + AllRepos registration.
5. Proto edits + `make create-service borrower` regen.
6. Logic: upload parse, detail, update full-replace, delete cascade.
7. Wire end-to-end, build, lint.

### 7.3 测试点

- `llm` openai impl: stub HTTP / SDK fake — not in scope to hit real API.
- `pdfparse`: inject mock `Provider`, assert (a) RPA price parse + 0-skip,
  (b) escrow field parse + empty-skip. Pure-Go pdfium (wazero) so no CGO in CI.
- Merge rule: extracted 0/null preserves prior `purchase_amount` / escrow field.

### 7.4 Open 项（非阻塞）

- Escrow `checklist` exact JSON shape already modeled
  (`dto.PropertyTransactionEscrowChecklist`); confirm frontend sends matching
  keys.
- 10 MB cap: enforce via `Content-Length` pre-check + `io.LimitReader` fallback.
