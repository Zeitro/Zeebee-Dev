# Borrower Property Transaction 文件内容提取 方案设计（v2 — 异步 + 多场景）

> 来源：`tasks/20260715_borrower_property_transaction_v2.md` §声明（2026-07-17）。
> 本版对 v1 的真实增量：**异步加进度/去重**（v1 upload 已是 `threading.GoSafe` 异步，但无进度/去重/手动触发）、**单 parser 配置 → 多规则配置**、**单文件 → 多场景（上传/整目录/子目录）**、**rpa 改按签字上下文择优**。
> v1 同步实现（`app/borrower/rpc/internal/biz_helper/file_parse.go`）**未上线**，仅作为「提取原语」参考；编排层与配置结构重写，数据结构可自由调整。

---

## 0. 变更声明（v1 → v2）

| 维度 | v1（已实现，未上线） | v2（本设计） |
|------|------|------|
| 触发方式 | upload 已异步（`threading.GoSafe`）但无进度/去重 | 上传 + 手动均异步，加 Redis 进度 + 去重 + requestId |
| 输入范围 | 单文档的单/多文件 | 单文件 \| 整目录递归 \| 指定子目录 |
| 多文件处理 | 取首个非空 | escrow 单文件；rpa 逐文件提 `price`，**冲突时 LLM 按上下文仲裁**（不预设信号） |
| 配置结构 | `{parser}` 单值 | `{rules[]}` 三模式（on_upload/manual_tree/manual_select）+ `selector`/`tags` |
| 进度可见性 | 无 | Redis `requestId` → 百分比，每文件滑动续期；终态独立长 TTL key |
| 重复触发 | 无保护 | 同 transaction + 同规则进行中 → **报错** |
| 文件类型 | PDF | PDF（本期）+ Word/图片（接口预留，延后实现） |
| 详情接口 | 返回 escrows + 6 列 | checklist 节点（`ChecklistNode`）携带 `parse_rules`，前端按 `manual_*` 规则渲染按钮（`document_id` 取自节点） |

> v1 的提取内核（pdfium 文本/渲染、LLM prompt、skip-empty merge、escrow upsert）**保留复用**；变更集中在「何时解析、解析哪些文件、如何综合、如何反馈进度」。

---

## 1. 背景与目标

### 1.1 业务背景

borrower 为 transaction 上传文件后，文件内含 transaction 信息（purchase_amount、escrow 公司/officer、各类日期等），目前需手动填写。v1 仅在「单文件上传」时同步提取两个字段，覆盖面窄；且大文件同步解析阻塞上传响应、前端无法获知解析进度。

### 1.2 目标

- **场景扩展**：除「上传即解析」外，支持「解析整目录并综合」「解析指定子目录并综合」，因为同一信息可能散落在多份文件中、需跨文件去重/择优。
- **异步化**：所有解析异步执行；手动触发返回 `requestId`，前端轮询进度；重复触发同一规则被拒绝。
- **配置驱动**：`document_templates` 中的 `file_parse_config` 升级为多规则结构，表达「同一文件夹可有多种解析需求」与「某文件夹被排除在父级递归之外」。
- **降低手动录入**：在保持 6 个手动列（`title_holder` 等）仍由人工填写的前提下，自动回写 `purchase_amount` 与 escrow 行。

### 1.3 范围边界

- **两个 parser 不变**：`rpa`（文本路径 → 5 字段：`coe_date`/`purchase_amount`/`emd_amount`/`emd_due_date`/`contingency_removal_date`）、`escrow_info`（视觉路径 → escrow 行字段）。**仅 `title_holder`/`contract_date`/`tc`/`cost_allocation` 4 列纯手动**（rpa 不触碰）。
- **三种触发/范围**：`on_upload` + `self`（场景 A，单文件）、`manual` + `self`（场景 C，文件夹直系文件综合）、`manual` + `tree`（场景 B，整目录递归综合）。
- **异步 + Redis 进度**：进度仅存 Redis，TTL 600s **每文件滑动续期**，终态独立长 TTL（1h）；**不建持久化 job 表**，进程重启进行中任务丢失（fail-silent，可重触发）。
- **文件类型**：本期仅 PDF；Word/图片通过 `FileExtractor` 接口预留，不实现。
- **fail-silent 保留**：解析失败不影响上传响应或手动触发响应；失败仅记日志 + 进度标记 `failed`。

---

## 2. 锁定决策

1. **Parser 不变 + 三模式触发映射。**
   - `escrow_info` → escrow 行；**mode = `on_upload`**（Wire Instructions 叶子，上传即解析**自身**文件，vision 提取）。上传传本次全部文件（同 v1），dedup 只防同文档并发。
   - `rpa` → **5 字段**（`coe_date`、`purchase_amount`、`emd_amount`、`emd_due_date`、`contingency_removal_date`）；**mode = `manual_tree`**（规则挂 Contract 文件夹，手动触发，解析该文件夹 + 所有子文件夹的文件）。综合：逐 PDF 读**前 10 页**文本拼接（`maxRPAPages=10`，非整文件）+ `extractRPAFields`（5 字段 prompt，扩展 v1）。**逐字段**收集非空/非 0 值——一致即取；冲突则 LLM 仲裁（各 `{fields, 前 10 页文本}` 喂 LLM，按上下文逐字段择优，**不预设信号**——签字时间仅例；无法定夺：金额取最大、日期取首个非空）。
   - **回写 = 逐字段权威替换**：提取非空/非 0 即 UPDATE 覆盖原值（含手填）；空/0 跳过。
   - Original RPA 叶子**不带 rule**（在 Contract 子树内，随 `manual_tree` 一并解析）。
2. **配置 schema：三模式 + selector/tags（替代 `excluded_parsers`）。** `FileParseConfig = { Rules []FileParseRule }`，rule = `{parser, mode, selector?}`，`mode ∈ {on_upload, manual_tree, manual_select}`：
   - `on_upload`：上传事件解析本次上传 doc 自身文件。
   - `manual_tree`：手动触发，解析触发节点 + **所有后代**文件（整子树，沿 `parent_id` 递归；无 selector、无排除）。
   - `manual_select`：手动触发 + `selector.tags`，按 AND 圈定**部分**文件（可跨树）。
   **取消 `manual_node` 与 `excluded_parsers`**：「整子树」用 `manual_tree`；「部分/含排除/跨树」用 `manual_select` + tags（纳入 = 打 tag，排除 = 不打）。目标 doc 打 `tags`，生成期拷进 `metadata.tags`。详见 §5。legacy `{parser}` normalize 为 `on_upload`。
3. **异步触发模型。** goroutine per job 经 go-zero `threading.GoSafe`（已有 panic 恢复；upload v1 已是此模式）+ Redis 去重锁 + Redis 进度。不引入 worker pool 或持久化队列。
4. **requestId = uuid**（`google/uuid`，确认 go.mod）。
5. **去重键** `bpt:parse:lock:{transactionId}:{ruleSig}`，`ruleSig = parser|mode|document_id`。`SetNX` 失败 → HTTP 409「解析进行中」。完成或失败后 `Del`。
6. **进度键** `bpt:parse:progress:{requestId}` = JSON `{percent, status, error?}`，TTL 600s **每处理完一个文件滑动续期**（`Expire` 续 600s），防慢任务中途过期。终态另用 `bpt:parse:done:{requestId}`（`{status,error}`，TTL 1h）供晚来查询、与 404 区分。`status ∈ {running, succeeded, failed}`。
7. **LLM client 单例。** `llm.NewClient` 在 ServiceContext 初始化一次，upload/parse 复用（v1 每次上传 new client 是浪费）。
8. **文件类型扩展。** `FileExtractor` 接口；本期仅 `pdfExtractor`（pdfium）。Word/Image 后续注册。
9. **进度模型 Redis-only。** 无 job 表、无重启恢复；lock/progress TTL 后允许重触发（fail-silent 接受重启丢任务）。
10. **Escrow 1:N。** upsert 主 escrow 行（min id），create 已 seed 故实际更新 seed 行。保留 `no` btree。
11. **PUT escrow 全量替换 / DELETE 级联 / GET detail 预加载。** 沿用 v1。
12. **Template 配置迁移走 lender `common/sync_data`。** 在 `app/lender/rpc/internal/logic/common/sync_data_logic.go`（现有 backfill 钩子）加一步：加载 `BorrowerPropertyTransaction` template → 改写 `structure` 中 `file_parse_config` → 存回。非裸 SQL。
13. **既有 transaction checklist 重建。** 部署后借 `sync_data`（或一次性脚本）遍历 transaction 重建 `checklist`，使 `ChecklistNode.parse_rules` 生效（rpa 按钮显现）。
14. **IDOR 校验。** `/parse` 的 `document_id` 须验：`document → entity_document_id → EntityDocument.EntityId == transactionId` 且 `EntityType == BorrowerPropertyTransaction` 且 org 属当前用户。

---

## 3. 数据模型

### 3.1 既有结构（保留，不改）

- `borrower_property_transactions`：6 列已存在；其中 `emd_due_date`/`contingency_removal_date` **改由 rpa 回写**（不再纯手动），`title_holder`/`contract_date`/`tc`/`cost_allocation` 仍纯手动。rpa 回写 `coe_date`/`purchase_amount`/`emd_amount`/`emd_due_date`/`contingency_removal_date` 5 列（均已存在）。`checklist` jsonb 保留。
- `property_transaction_escrows`：表与索引保留；escrow_info 综合结果写入主行（min id）。`checklist` jsonb（`dto.PropertyTransactionEscrowChecklist`）保留。**注意**：`create_logic.go:104` 在建交易时已 seed 一行空 escrow，故 escrow_info 解析实际总是更新该 seed 行；"无则建" 分支保留为兜底。

### 3.2 配置来源与流转（template → doc metadata）

`file_parse_config` 的**单一事实来源是 `document_templates.structure`**（全局行，`type=BorrowerPropertyTransaction`）。每个 template 节点（`DocumentTemplateStructure`）携带 `FileParseConfig`，声明该节点的解析规则。

- **生成期拷贝**：`create_logic.createDocumentNode` 把每个节点的 `node.FileParseConfig` → `metadata.file_parse_config`、`node.Tags` → `metadata.tags`（`tags` 为新增拷贝）。
- **template 改配置** → 只影响**之后新建**的 transaction（其 doc metadata 拿到新形态）。
- **既有 transaction** 的 doc metadata 是建交易时的快照，不随 template 变更自动刷新；靠 `effectiveRules` normalize 兼容旧 `{parser}` 形态，直到文档树重建。
- **运行期读取**：upload / parse logic 从 `doc.Metadata.file_parse_config` 读取（**不从 template 读**），保证每个 transaction 用建交易时锁定的配置。
- **Go 类型单一**：`dto.FileParseConfig`（`Rules`/legacy `Parser`）+ `dto.TagSelector` 同时用于 `DocumentTemplateStructure`（template）与 Document `metadata` 反序列化——一处定义，两处复用。
- **tags 同流转**：`DocumentTemplateStructure.Tags`（目标节点分类）同样由 `createDocumentNode` 拷进 `doc.metadata.tags`，供 `manual_select` 的 selector 圈定。

**配置形态变更**：`{parser}` → `{rules[]}`（`manual_select` 配 `tags`），迁移走 sync_data（§5.5）。`DocumentTemplateStructure` 升级（新 `FileParseConfig` + 新增 `Tags`），保留 `Parser` 供 normalize。

### 3.3 Redis 键（无新 DB 表）

| 用途 | key | value | TTL |
|------|-----|-------|-----|
| 去重锁 | `bpt:parse:lock:{transactionId}:{ruleSig}` | `requestId` | 600s（SetNX，每文件 `Expire` 续期） |
| 运行进度 | `bpt:parse:progress:{requestId}` | `{"percent":int,"status":"running","error":""}` | 600s（每文件处理完滑动续期） |
| 终态 | `bpt:parse:done:{requestId}` | `{"status":"succeeded\|failed","error":""}` | 3600s（晚来可查，与 404 区分） |

> 进度 `percent` 粗粒度：多文件 = `processed/total*100`；单文件 = 下载/提取 50% → 回写 100%。终态用独立长 TTL key，避免运行 key 过期后 getProgress 误返 404。

---

## 4. 接口设计

| 端点 | 改动 |
|------|------|
| create transaction | 无逻辑改动；新 schema `file_parse_config` 透传到 doc metadata |
| upload-file | 上传事务提交**后**：读规则，对 `trigger=on_upload` 的 rule **异步派发**（fire-and-forget，requestId 入日志，不阻塞响应、不改返回体） |
| GET detail | 返回 checklist（`ChecklistNode` 带 `parse_rules`）+ escrows + 6 列。前端遍历 checklist，节点 `parse_rules` 含 `manual_tree`/`manual_select` → 渲染按钮，`document_id` 取自 `ChecklistNode.DocumentId` |
| PUT update | 不变（escrow 全量替换、checklist 原样持久化） |
| DELETE | 不变（级联删 escrows） |
| **NEW** `POST /api/v1/borrower_property_transaction/:id/parse` | body `{document_id, parser, mode}`（document_id = 文件夹 Document 行 PK）→ 校验归属+规则 → 去重 → 异步派发 → `{request_id}`；进行中 → 409 |
| **NEW** `GET /api/v1/borrower_property_transaction/parse/:request_id/progress` | → `{percent, status, error?}`（Redis 命中；过期 → 404） |

### 4.1 Proto 变更

- `domains/zeebee/proto/entity.proto`
  - `ChecklistNode`：新增 `repeated FileParseRule parse_rules`（从该节点 Document 的 `metadata.file_parse_config` 经 `effectiveRules` 填入）。前端按 `manual_tree`/`manual_select` 规则渲染按钮，`document_id` 取自节点既有字段——**不另起 `parse_actions`**。
  - 新 message `FileParseRule { parser; mode; selector }` + `TagSelector { tags }`；`ChecklistNode` 可暴露 `tags` 供前端展示。
- `app/borrower/rpc/pb/borrower_borrower_property_transaction.proto`
  - 新 RPC `ParseFileExtract`：请求 `{transaction_id, document_id, parser, mode}`，响应 `{request_id}`。
  - 新 RPC `GetFileExtractProgress`：请求 `{request_id}`，响应 `{percent, status, error}`。

> Proto 改完执行 `make create-service borrower` 重生成。

---

## 5. 配置设计（`file_parse_config` + `tags`）

规则附**触发节点**（按钮位置 + 触发声明）；输入范围由 `mode` 决定，`manual_select` 再用 `selector`+`tags` 圈定部分集。

### 5.1 Model

```go
type FileParseConfig struct {
    Rules  []FileParseRule `json:"rules,omitempty"`
    Parser string          `json:"parser,omitempty"` // legacy 单 parser，normalize 为 on_upload
}
type FileParseRule struct {
    Parser   string       `json:"parser"`            // rpa | escrow_info
    Mode     string       `json:"mode"`              // on_upload | manual_tree | manual_select
    Selector *TagSelector `json:"selector,omitempty"` // 仅 manual_select
}
type TagSelector struct {
    Tags []string `json:"tags"` // 目标 doc 须含全部这些 tag（AND）
}
// DocumentTemplateStructure 新增：Tags []string `json:"tags,omitempty"`（目标节点分类 → doc.metadata.tags）
```

字段说明：

- `parser`：提取器名，**一对一隐含输出目标**（`rpa`→`purchase_amount`，`escrow_info`→escrow 行）。
- `mode`：三值（见 §5.2）。
- `selector`（仅 `manual_select`）：`{tags:[...]}` AND 语义。
- `tags`（目标节点）：`manual_select` 的 selector 据此圈定；生成期拷进 `metadata.tags`。
- **类型归属**：放 `dto`（避免 `model↔dto` 循环）。
- **位置式，不含 id**：rule 附触发节点；`document_id` 是运行期句柄（触发节点 PK），不进 config。

### 5.2 三模式 → 场景映射

| mode | 触发 | 输入集 |
|------|------|------|
| `on_upload` | 上传文件事件 | 本次上传 doc 自身文件 |
| `manual_tree` | 手动按钮（rule 所在节点） | 触发节点 + **所有后代**文件（整子树，沿 `parent_id` 递归） |
| `manual_select` | 手动按钮 + `selector` | 本 transaction 所有 `metadata.tags ⊇ selector.tags` 的 doc 的文件（部分/可跨树） |

> 「整子树」= `manual_tree`（简单、无排除）；「部分/含排除/跨树」= `manual_select` + tags（纳入 = 打 tag，排除 = 不打）。无 `manual_node`、无 `excluded_parsers`。同节点可多 rule。

### 5.3 示例

```jsonc
// 触发叶子（Wire Instructions）：escrow 上传即解析自身
{"file_parse_config":{"rules":[{"parser":"escrow_info","mode":"on_upload"}]}}

// 触发节点（Contract 文件夹）：rpa 手动解析该文件夹 + 所有子文件夹文件（冲突 LLM 仲裁）
{"file_parse_config":{"rules":[{"parser":"rpa","mode":"manual_tree"}]}}

// 部分解析（manual_select）：目标叶子打 tag，规则用 selector 圈定
//   目标：{"tags":["emd"]}
{"file_parse_config":{"rules":[{"parser":"<future>","mode":"manual_select","selector":{"tags":["emd"]}}]}}
```

> Original RPA 叶子在 Contract 子树内，随 `manual_tree` 一并解析，不需自带 rule。排除某目录 = 改用 `manual_select` + tags 精确圈定，而非整树。

### 5.4 normalize（兼容既有数据）

```go
func effectiveRules(cfg *FileParseConfig) []FileParseRule {
    if cfg == nil { return nil }
    if len(cfg.Rules) > 0 { return cfg.Rules }
    if cfg.Parser != "" {
        return []FileParseRule{{Parser: cfg.Parser, Mode: "on_upload"}} // legacy = 自身上传解析
    }
    return nil
}
```

> 既有 transaction doc metadata 是旧 `{parser}`；normalize 为 on_upload（自身）。新 transaction 写 `{rules}` + `tags`。

### 5.5 配置迁移（走 lender `common/sync_data`）

不写裸 SQL。在 `sync_data_logic.go` 加一步：加载 `BorrowerPropertyTransaction` template → Go 改写 `structure` → 存回。改写：

```jsonc
// Wire Instructions 叶子：{"parser":"escrow_info"} →
{"file_parse_config":{"rules":[{"parser":"escrow_info","mode":"on_upload"}]}}

// Original RPA 叶子：{"parser":"rpa"} → 移除 rule（随 Contract manual_tree 解析）

// Contract 文件夹：原无 → 新增
{"file_parse_config":{"rules":[{"parser":"rpa","mode":"manual_tree"}]}}
```

> `manual_select` 场景按需另加 `tags` + `selector` 规则。部署后调一次 `POST /common/sync_data`（幂等：已是新形态则跳过）。仅影响新建 transaction；既有靠 normalize + checklist 重建（§2 决策 13）。

---

## 6. 核心流程

### 6.1 on_upload 异步（场景 A）

1. 文件事务提交落库（`file.FileUrl` 已生成）。
2. 读 `doc.Metadata.file_parse_config` → `effectiveRules` → 过滤 `mode=on_upload`。**当前仅 `escrow_info` 是 `on_upload`；`rpa` 是 `manual_tree`，永不在上传时触发**（综合 Contract 全量 PDF 太重，交给手动按钮）。
3. 对每条匹配 rule：`SetNX` 去重锁；成功则 `threading.GoSafe` 跑 `run(requestId, rule, files=本次上传的全部文件)`（同 v1 传全部，dedup 只防同文档并发），requestId 入日志；失败（已有锁）则跳过（上传触发的重复不报错给客户端）。
4. 立即返回上传响应（不等待解析）。

### 6.2 手动触发（场景 B/C）

> **目标目录用 `document_id` 标识，不用字符串路径。** documents 按邻接表（`parent_id`）物化为行：文件夹 = `DocumentTypeBranch` 行，文件挂在 `DocumentTypeLeaf` 行上。每行有稳定主键 `id`，`file_parse_config` 在创建时从 template 拷进**每行** `metadata`。前端从 detail 拿到带 `id` 的文档树，按钮直接携带目标文件夹的 `document_id`——稳定、唯一、不依赖名称/路径。

1. `POST /:id/parse` 用 `{document_id, parser, mode}` 校验（IDOR）：`document → entity_document_id → EntityDocument.EntityId == transactionId` 且 `EntityType == BorrowerPropertyTransaction` 且 org 属当前用户；`document.type == DocumentTypeBranch`；`metadata.file_parse_config` 含匹配 `{parser, mode}` rule（防伪造）。
2. `ruleSig = parser|mode|document_id`；`SetNX(lock)`：
   - 失败 → 409 `parse in progress`。
   - 成功 → 生成 `requestId`，`threading.GoSafe(func(){ run(...) })`，返回 `{request_id}`。
3. `run` 内按 `mode` 收集文件（仅 PDF）：`manual_tree` 沿 `parent_id` 递归收集触发节点 + 所有后代 Leaf 文件；`manual_select` 用 `collectBySelector(本 transaction 全部 doc, rule.selector)` 筛 `metadata.tags ⊇ selector.tags` 的 doc 文件。

### 6.3 提取 + 综合

两个 parser 的提取策略不同（parser 自决，orchestrator 只按 mode 交付文件集）：

1. **escrow_info（on_upload，本次全部文件）**：各文件跑 `extractEscrowFields`（pdfium 渲染首页 PNG → LLM vision），逐字段 skip-empty merge 进主 escrow 行（create 已 seed，实际更新 seed 行）。复用 v1。
2. **rpa（manual_tree，综合）**：Contract 子树逐 PDF 读**前 10 页**文本拼接（`maxRPAPages=10`）+ `extractRPAFields`（5 字段：`coe_date`/`purchase_amount`/`emd_amount`/`emd_due_date`/`contingency_removal_date`，扩展 v1 prompt）→ 收集 `{name, fields, text}`（跳过空/0 字段）。**逐字段**：跨文件非空值一致 → 取；冲突 → **LLM 仲裁**（各 `{fields, 前 10 页文本}` 喂 LLM，按上下文逐字段择优，不预设信号；无法定夺：金额取最大、日期取首个非空）。
3. **rpa 回写（逐字段权威替换）**：5 字段提取非空/非 0 即 UPDATE 覆盖原值（含手填）；空/0 跳过。
4. **escrow 回写**：skip-empty（空字段保留原值）。

### 6.4 进度 + 去重生命周期

```
SetNX(lock, requestId, 600s)        → 失败: 409
SET(progress, {0,running}, 600s)
for file in scope:
    extract...
    SET(progress, {percent,running})
    EXPIRE(lock,600s) + EXPIRE(progress,600s)   # 每文件滑动续期
arbitrate(rpa) + write-back
SET(done, {succeeded|failed, err?}, 3600s)      # 终态长 TTL
DEL(progress); DEL(lock)
```

- 经 `threading.GoSafe` 跑（go-zero panic 恢复）；任意 error → 终态 `failed` + 错误信息 → `Del(lock)`。
- 滑动续期防慢任务（rpa 综合 + 仲裁可能久）中途 lock/progress 过期；终态独立 key 供晚来查询，与 404 区分。

### 6.5 回写总则

- **escrow**：skip-empty——提取 `""`/`null` 保留原字段（不覆盖真实值）。
- **rpa（5 字段）**：逐字段权威替换——提取非空/非 0 覆盖原值（含手填）；空/0 跳过。与 escrow 因 parser 不同不冲突。

---

## 7. 异步与进度模型

- **执行单元**：goroutine per job 经 `threading.GoSafe`（go-zero panic 恢复）。无 worker pool、无持久化队列。
- **LLM client**：ServiceContext 单例（`llm.NewClient` 一次），upload/parse 复用。
- **并发上限**：提取受 pdfium pool（`MaxTotal=5`）天然限制；下载/LLM 调用受各自 HTTP 超时（60s）约束。高并发风暴的兜底：上传触发的去重锁按 `document+rule` 抑制同文档并发；手动触发按 `transaction+rule` 抑制。
- **requestId**：`google/uuid`（已是常见依赖，确认 go.mod）。
- **可观测**：每次 run 记 `transactionId / requestId / ruleSig / percent / status / error`；失败带 `errorx.GetErrorStack`。
- **重启行为**：进程重启 → 进行中 goroutine 丢失；lock/progress 在 Redis 存活至 TTL。期间去重仍生效（防重复触发）；TTL 后可重触发。**无恢复、无重放**——符合 fail-silent。

---

## 8. 文件类型扩展（延后）

```go
// FileExtractor 按文件类型提取结构化字段；本期仅 pdfExtractor。
type FileExtractor interface {
    Supports(fileName string) bool
    // rpa：逐文件提取 5 字段，冲突字段内部 LLM 仲裁，返回最终字段集
    ExtractRPA(ctx context.Context, readers []io.ReadSeeker) (fields RPAFields, err error)
    // escrow 单文件 vision 提取
    ExtractEscrow(ctx context.Context, r io.ReadSeeker) (fields escrowFields, err error)
}
```

- `pdfExtractor`：pdfium pool + `extractRPAFields`（rpa 5 字段，前 10 页文本）/ `extractEscrowFields`（escrow，首页 vision）。
- Word/Image：后续实现 `docxExtractor` / `imageExtractor` 并在 registry 注册；orchestrator 按 `Supports` 派发，无需改流程。**本期不实现**（声明明确「先可以不实现，后续再实现」）。

---

## 9. 与现有代码的关系

**复用（提取原语，已验证）**
- `common/llm.Client`（`ChatJSON` / `VisionJSON`）。
- pdfium pool（`getPdfPool` / `withPdfInstance`）、`downloadPDF`（10MB cap）、`openPDF`、`readPageText`、`renderPagePNG`。
- `extractEscrowFields` / `escrowPrompt` / `parseEscrowFields`；rpa 侧 v1 `extractPurchasePrice` → 扩展为 `extractRPAFields`（5 字段）+ `parseRPAFields` + `rpaPrompt`（5 字段）+ `arbitratePrompt`。
- `applyEscrowFields`（skip-empty）、escrow upsert（min id 主行）。
- `common/redisx`（`SetNX` / `SetExpire` / `Get` / `Del`）。
- escrow model / dto / repo（`PropertyTransactionEscrowRepo`，含 `DeleteByIdsNotInWithTx`）。
- 6 列迁移、template structure 既有形态。

**拆分 + 重写**（`biz_helper/file_parse.go` 481 行职责过重 → 拆 `biz_helper/parse_file/` 子包，按职责分文件）：
- `parse_file/pdf.go`：pdfium pool、`downloadPDF`（10MB cap）、`openPDF`、`readPageText`、`readFirstNPages`（rpa 前 10 页，`maxRPAPages=10`）、`renderPagePNG`（escrow 首页）、`isPDF`。
- `parse_file/prompts.go`：`extractSystemPrompt`、`rpaPrompt`、`escrowPrompt`、`arbitratePrompt`。
- `parse_file/extract.go`：`escrowFields`、`rpaFields`（5 字段）、`parseEscrowFields`/`parseRPAFields`、`extractEscrowFields`/`extractRPAFields`。
- `parse_file/arbitrate.go`：rpa 冲突 LLM 仲裁。
- `parse_file/selector.go`：`TagSelector`、`effectiveRules`、`collectSubtree`（manual_tree）、`collectBySelector`（manual_select）、读 `metadata.file_parse_config`/`tags`。
- `parse_file/progress.go`：Redis 锁/进度/终态（滑动续期）。
- `parse_file/parser.go`：编排——`DispatchOnUpload` / `TriggerManual` / `run`（按 mode 收集 + 综合 + 回写 + 进度）。
- `domains/zeebee/model/document_template.go`：`FileParseConfig` 单 `Parser` → `Rules`（保留 `Parser` 供 normalize）+ `DocumentTemplateStructure.Tags`。
- `upload_document_file_logic.go`：调用点改指 `parse_file` 包（已异步）。
- `buildChecklistNode`（`biz_helper/borrower_property_transaction.go`）：`effectiveRules` → `ChecklistNode.parse_rules`。

**新增**
- `POST /parse` + `GET /parse/:request_id/progress`（API + RPC）。
- ServiceContext：`LLM llm.Client` 单例。
- lender `sync_data_logic.go`：template structure 改写（含 `tags`）+ checklist 重建。

**弃用**
- `readParserName`、`excluded_parsers`、`manual_node`、`merge`。

---

## 10. 风险与权衡

- **Redis-only 进度 / 无 job 表**：重启丢任务；滑动续期防中途过期、终态 key 供晚来查询。权衡：实现简；代价：无持久历史、无失败重试。若需审计/重试，再加 `property_transaction_parse_jobs` 表。
- **goroutine per job / 无背压**：突发大量上传或手动触发可能起大量 goroutine。权衡：简单；兜底靠 pdfium pool（提取端 5 并发）+ 去重锁（同文档/同规则串行）。若吞吐成问题，加 buffered channel + N worker。
- **rpa 仲裁依赖 LLM 判断**：多文件 price 冲突时由 LLM 按上下文定夺，不预设信号（签字时间仅例）。权衡：灵活、适应现实不确定；代价：仲裁结果不可解释/偶有偏差，回退最大值兜底。
- **配置迁移兼容**：既有 transaction 的 doc metadata 是旧 `{parser}`，靠 `effectiveRules` 兼容到文档刷新。权衡：零停机；代价：reader 双形态分支保留一周期。
- **去重粒度**：按 `ruleSig`（含 mode/node_path），同 parser 不同 mode 不互斥——符合「同一文件夹多种解析需求」的设计意图。
- **进度粗粒度**：percent 按「文件数」计，单文件内下载/LLM 耗时不细分。权衡：简单；代价：大文件时进度长时间停在某一档。
- **fail-silent 可见性**：上传触发的解析失败对用户不可见（仅日志 + 进度 failed，但上传场景前端不持有 requestId）。权衡：保证上传可用；代价：用户偶需手动补录。

---

## 11. 实现要点 / 测试

### 11.1 文件清单

**Model / Config**
- `domains/zeebee/model/document_template.go`：`FileParseConfig`（`Rules`，保留 `Parser`）+ `DocumentTemplateStructure.Tags` + `effectiveRules`（`parse_file/selector.go`）。

**Orchestrator（拆分到 `biz_helper/parse_file/` 子包，见 §9）**
- `DispatchOnUpload(ctx, transactionId, doc, files)`：过滤 `on_upload` rule，`SetNX` 去重，`threading.GoSafe(run)`，不阻塞。
- `TriggerManual(ctx, transactionId, documentId, rule) (requestId, error)`：校验 `document_id` 归属（Branch + 本 transaction）+ 规则命中，`SetNX` 失败返错，否则派发。
- `run(ctx, requestId, transactionId, rule, modeFiles)`：进度推进（每文件滑动续期）+ per-file extract + rpa 逐字段冲突 LLM 仲裁 + 独立事务回写（5 字段权威替换）+ 终态 key + `Del(lock)`，经 `threading.GoSafe`。
- `collectSubtree(rootDoc, anchorDocId)`（manual_tree，沿 `parent_id` 递归）+ `collectBySelector(txDocs, selector)`（manual_select，筛 `metadata.tags ⊇ selector.tags`）。
- 进度/锁封装（`biz_helper/parse_progress.go`）：`acquireLock` / `setProgress` / `getProgress` / `releaseLock`。

**Endpoints（API + RPC）**
- `POST /:id/parse` + `GET /parse/:request_id/progress`：handler + logic（API 转发 RPC）。
- RPC `ParseFileExtract` / `GetFileExtractProgress` logic。
- Proto：`ChecklistNode.parse_rules` + `FileParseRule` message + 两个新 RPC 请求/响应。
- `buildChecklistNode` 填 `parse_rules`。

**Migration（走 sync_data，非 SQL）**
- lender `sync_data_logic.go`：加 `BorrowerPropertyTransaction` template `structure` 改写步（`{parser}`→`{rules}`，rpa 上移 Contract）+ 遍历 transaction 重建 `checklist` 步。部署后调一次 `POST /common/sync_data`。

### 11.2 构建顺序

1. Config：`FileParseConfig`（`Rules` + `Tags`，去 `merge`/`excluded_parsers`）+ `effectiveRules`。
2. 进度/锁/终态封装（`parse_file/progress.go`，滑动续期）+ 单测（mock redisx 或 miniredis）。
3. Orchestrator 重写（`parse_file/` 子包）：`run` + `collectSubtree`/`collectBySelector` + rpa 逐文件提取+冲突 LLM 仲裁 + escrow 单/多文件（复用 v1）+ 单测（mock `llm.Client`）。
4. ServiceContext LLM 单例 + `DispatchOnUpload` / `TriggerManual`（`threading.GoSafe`）+ upload 调用点（已异步，改传全部 files）。
5. Proto（`ChecklistNode.parse_rules` + 新 RPC）+ `make create-service borrower` 重生成。
6. 新端点（API + RPC）+ `buildChecklistNode.parse_rules`。
7. 端到端接线、build、lint。
8. lender `sync_data`：template 改写 + checklist 重建步。

### 11.3 测试点

- `effectiveRules`：`{rules}` 直通；`{parser}` 归一化为 `on_upload`；空 config 返回 nil。
- `collectSubtree`（manual_tree）：触发节点 + 所有后代；`collectBySelector`（manual_select）：`metadata.tags ⊇ selector.tags` 命中、部分/跨树场景。
- rpa 仲裁：逐字段——单值直取；多冲突值触发 LLM 仲裁（mock 返回权威 5 字段）；非空/非 0 字段权威替换。
- 进度/锁：`SetNX` 冲突 → 409；每文件滑动续期；完成 → 终态 key + `Del(lock)`；panic → `failed`。
- parser：mock `llm.Client`，断言 rpa 5 字段解析 + 空/0 跳过、escrow 字段解析 + 空 跳过。pdfium 用 wazero，CI 无 CGO。

### 11.4 Open 项（非阻塞）

1. ~~checklist 位置~~ **已定**：`ChecklistNode.parse_rules`（不另起 `parse_actions`）。
2. ~~rpa 触发~~ **已定**：rpa = `manual_tree`（Contract 整子树手动）；escrow = `on_upload`；第三模式 `manual_select` + selector/tags 供部分/跨树解析。
3. ~~rpa 综合~~ **已定**：逐文件提取，冲突 LLM 仲裁（不预设信号），权威替换。
4. ~~进度 TTL/终态/merge/LLM 单例/IDOR/migration~~ **已定**（见 §2）。
5. **进度 `percent` 粒度**：当前按文件数。若前端要更平滑（含下载/LLM/仲裁阶段），细化状态机。
6. **uuid 依赖**：确认 `google/uuid` 在 go.mod；否则 go-zero 自带 id 生成。
