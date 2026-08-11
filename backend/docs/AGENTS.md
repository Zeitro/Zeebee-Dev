# AGENTS.md

> Zeebee Backend 工作入口。**本文件为 Claude 与 Codex 共用的唯一项目规范文档**：Claude 经根目录 `CLAUDE.md` 入口加载，Codex 经根目录 `AGENTS.md` 入口加载，切换工具无需重复编写规范文档。
>
> **代码风格**（项目无关）见 constitution，**结构**见 architecture，**事实**见 INDEX。本文件仅含项目特定规范与可执行示例。
> 改任何 `.api`/`.proto` 后须 `make create-service <service-name>` 重生成。
>
> `constitution.md` 为必须加载并遵循的代码风格宪章。

## 文档指针

- **代码风格（项目无关，可迁移）**：`constitution.md`
- **架构结构**：`guides/architecture.md`（分层职责、数据流、目录组织、服务矩阵、ADR，required reading）
- **事实清单**：`INDEX.md`（服务模块矩阵、模型/DTO 清单、通用组件、目录速查）
- **规则**（技术栈、命名、目录）= 本文件
- **项目配置定义文件**：`app/<service>/config/config.go`
- **项目配置文件**：`app/<service>/etc/config-<env>.yaml`

### 配置文件加载路径（必须遵守）

- 适用于所有服务：配置定义先看 `app/<service>/config/config.go`，运行时配置再看该服务实际加载的 `app/<service>/etc/config-<env>.yaml`。
- 当前服务入口为 `app/<service>/cmd/main.go`；各入口默认 `-f etc/config.yaml`，配置路径按进程当前工作目录解析，不能只凭目录名称推断。
- 从服务目录启动时使用 `-f etc/config-local.yaml`；从 backend 根目录启动时使用 `-f app/<service>/etc/config-local.yaml`。staging/prod 同理替换环境名。
- 检查任何服务前，必须先读取对应的 `cmd/main.go` 和实际启动命令，确认 `-f` 指向的文件；包括 borrower、lender、admin、tenant、pdf_generator。
- **不要把 `app/<service>/rpc/etc/`、`app/<service>/api/etc/` 或代码生成产生的其他 `etc/` 目录当作运行时配置来源**，除非对应入口明确通过 `-f` 加载它。
- 修改或检查配置项（包括 LogoUrl、FooterAddress）时，必须以入口实际加载的配置文件为准，并同时检查是否存在其他服务/目录中的重复配置。

### Document Reading Order

首次会话必读（按序）：

1. **本文件** — 项目规范、代码示例
2. **`constitution.md`** — 代码风格（项目无关）
3. **`guides/architecture.md`** — 架构结构

### 任务路由（按需加载，避免全量读取）

| 任务 | 只读 |
|------|------|
| 需求开发 / 功能新增 | 本文件 §需求开发流程（SpecKit）|
| 加/改 HTTP 端点 | 本文件 §核心原则 + [architecture](guides/architecture.md) §开发流程 / §端点实现范例 |
| 定位代码（model/repo/dto/模块） | [INDEX](INDEX.md) §按目录定位速查 / §服务功能矩阵 |
| 代码风格、命名、注释 | [constitution](constitution.md) |
| 数据转换 / 错误处理写法 | 本文件 §Data Conversion / §Error Handling |
| 服务端口 / 职责 | [architecture](guides/architecture.md) §服务矩阵 |
| 模型 / DTO / 组件清单 | [INDEX](INDEX.md) §领域模型 / §DTO / §通用组件 |
| 提交格式 | 本文件 §Git Commit + [constitution](constitution.md) §Git 提交 |

> 按任务取用。规则类事实单一来源（见各表项）；计数 / 版本等易变值以代码为准，文档不重复维护。

## 项目角色

资深 Go 后端工程师，专精 go-zero 微服务架构与 goroutine/channel、GORM、gRPC/Protobuf、PostgreSQL 调优、pprof/trace 性能分析。

## 核心原则

### I. 共享领域优先

功能从共享领域层开始设计。领域模型与仓储放 `domains/zeebee/`，共享代码独立可测试且有文档。不允许仅为了组织而抽象。

### II. API/RPC 分离

双层架构：API 层（HTTP）处理验证和路由，RPC 层（gRPC）处理业务逻辑。数据流：`API Handler → API Logic → RPC Call → RPC Logic → Repository → Model → Database`。

### III. 代码生成强制

使用 `Makefile` 封装的命令生成代码，如 `make create-service lender`。不手写 handler/server。

### IV. 数据转换标准

所有类型转换使用 `copierx.Trans[T, U]`（T、U 为非指针类型，src 为 `*T`）。时间字段需 `copierx.TimeToRFC3339Str()` 选项。

### V. 类型安全与测试

强类型贯穿全代码库。测试要求（单元/集成/覆盖率）见 [constitution](constitution.md) §测试原则。

## 技术栈

- **语言**：Go 1.26 | **框架**：go-zero 1.10.1
- **数据库**：PostgreSQL 14+ (GORM) | **缓存**：Redis 7+ | **存储**：AWS S3
- **协议**：HTTP (REST) + gRPC | **认证**：JWT with RBAC

## 服务与共享层

服务矩阵（端口/职责/规模）、模块清单、`domains/zeebee/` 模型仓储、`common/` 22 工具包 — 单一来源，不在此重复维护：

- [`guides/architecture.md`](guides/architecture.md) §服务矩阵 / §代码组织（结构与目录）
- [`INDEX.md`](INDEX.md) §服务功能矩阵（模块级）/ §通用组件（common/ 清单）/ §按目录定位速查

> 事实型数据以 `go.mod` 与实际代码统计为准。

## 代码约束

- **数据库访问**：RPC logic 不直接操作数据库，须经 `baseLogic` → `domains/zeebee/logic/` → `repo`（详见 [architecture](guides/architecture.md) §关键约束）
- **错误处理**：`common/errorx/`
- **校验**：`common/validator/`
- **类型转换**：`common/copierx`
- **日志**：go-zero logx，禁止硬编码敏感信息
- **lint**：PR 必过 golangci-lint（`.golangci.yml`）
- **配置**：YAML 文件，禁止硬编码敏感信息
- **函数签名**：≤4 参，≤3 返回值（风格细节见 constitution §函数约束）
- **注释**：简洁扼要、只写「为什么」（详见 [constitution](constitution.md) §注释）

### 配置中的密钥（HARDCODED_SECRETS 例外）

本仓库为**私有仓库**，`app/*/etc/*.yaml` 是服务启动时加载的运行时配置文件。将 LLM / 第三方服务的 `ApiKey` 直接写在 `etc/*.yaml` 的 `Providers[*].ApiKey` 字段是**有意为之**的部署模式（非泄露），安全扫描对该模式报 `HARDCODED_SECRETS` 可视为**预期告警、按已接受处理**，无需每次修复或上报。

- 仅适用于 `app/*/etc/*.yaml` 的 `ApiKey`（及等价凭证字段）；其他位置的硬编码密钥仍禁止。
- `config-local.yaml` 通常不入库或保留空值；`config-prod.yaml` / `config-staging.yaml` 由部署流程或人工填入真实 key。
- 若确有 key 被误提交且仓库可见性变更（转公开），须立刻轮换。
- 日志中仍禁止打印密钥（见上一条「日志」约束），`SENSITIVE-TO-OBSERVABILITY` 类告警与此例外无关，单独评估。

## Working Directory

Monorepo。工作范围限 `backend/`。Specs 入 `backend/specs/`（非项目根 `specs/`）。

## 文档策略

- `docs/`：项目知识库（constitution、architecture、INDEX、designs），纳入版本管理，工具无关
- `.agents/`（Codex 技能）、`.claude/`（Claude 技能）、`.specify/`：AI agent 与插件层（skills、speckit 工具链、本地分析报告），本地保留，不入版本管理
- `backend/specs/`：SpecKit 实施规格（非项目根 `specs/`；工具无关，Claude/Codex 共用），本地保留，不入版本管理

## Data Conversion Convention — `copierx.Trans`

所有类型转换统一用 `copierx.Trans`（**唯一示例源**，constitution/architecture 仅指针引用此处）。

```go
// func Trans[T, U any](src *T, opts ...copier.TypeConverter) (dist *U, err error)
// ⚠️ T 和 U 必须是非指针类型，src 是 *T，返回 *U

// DTO → Model
data, err := copierx.Trans[dto.BaseTemplateCreateDto, model.BaseTemplate](&req)

// Model → BO
result, err := copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](modelData)

// BO → Proto Entity（带时间格式化）
protoData, err := copierx.Trans[bo.BaseTemplateBo, admin.BaseTemplateEntity](boData, copierx.TimeToRFC3339Str())

// Slice 转换同样适用
items, err := copierx.Trans[[]dto.Item, []model.Item](&dtoSlice)
```

| Direction   | T (Source)   | U (Target)        | src         | Return             |
|-------------|--------------|-------------------|-------------|--------------------|
| DTO → Model | `dto.XxxDto` | `model.Xxx`       | `&dtoData`  | `*model.Xxx`       |
| Model → DTO | `model.Xxx`  | `dto.XxxDto`      | `modelData` | `*dto.XxxDto`      |
| DTO → Proto | `dto.XxxDto` | `admin.XxxEntity` | `dtoData`   | `*admin.XxxEntity` |

## Error Handling Convention

```go
defer func() {
    if err != nil {
        logx.WithContext(r.Context()).
            WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
            Error(err)
    }
}()

if err != nil {
    return nil, errorx.NewDefaultError("invalid input parameters")
}
```

禁止裸 `errors.New`，统一用 `common/errorx/`。

## Git Commit Convention

提交 commit 前先检查是否在项目根目录，如果不在的根目录，请先 cd 到项目根目录。

格式与通用约束（conventional commits、co-author 禁令、多 commit 拆分）见 `constitution.md` §Git 提交。

本项目 scope 格式：

```
<type>(<scope>-<service-name>-backend): <description>

feat | fix | refactor | docs | style | perf | test
# 示例
feat(applicant-borrower-backend): add loan application endpoint
fix(document-lender-backend): fix document upload validation error
```

## code extract
codebase-memory-mcp > codegraph > grep
