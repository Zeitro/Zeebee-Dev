# CLAUDE.md

> LLM 工作入口。**规则**见 constitution，**结构**见 architecture，**事实**见 INDEX。本文件仅含可执行示例与工作约束，不重复上述文档内容。
> 修改任何 `.api`/`.proto` 后须 `make create-service <service-name>` 重生成。

## Constitution & Architecture（指针引用，不复制内容）

- **规则**：`.specify/memory/constitution.md` — 核心原则、技术栈、Git 提交格式、语言规范、代码行为准则
- **架构**：`.claude/docs/analysis/architecture.md` — 分层职责、数据流、目录组织、服务矩阵、ADR（required reading）
- **事实清单**：`.claude/docs/INDEX.md` — 服务模块矩阵、模型/DTO 清单、通用组件、按目录定位速查
- **最新审计**：`.claude/docs/analysis/analysis.md` — 全量深度审计（2026-06-24，90 个发现）

## Document Reading Order

首次会话必读（按序）：

1. **本文件** — 开发约束、代码示例
2. **`.specify/memory/constitution.md`** — 规则（what must be done）
3. **`.claude/docs/analysis/architecture.md`** — 结构（how code is organized）

文档分工：
- `constitution.md` = 规则 | `architecture.md` = 结构 | **本文件 = 约定 + 示例（code style）**
- `PROJECT_INDEX.md` = repo 索引快照 | `README.md` = 对外概览（Quick Start）

## Working Directory

Monorepo。工作范围限 `backend/`。Specs 入 `backend/specs/`（非项目根 `specs/`）。

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

## Code Constraints

- **数据库操作**：logic 不直接操作数据库，必须经 `domains/zeebee/repo` 层
- **错误处理**：`common/errorx/` | **校验**：`common/validator/` | **日志**：go-zero logx | **类型转换**：`common/copierx`
- **函数签名**：≤4 参，≤3 返回值
- **配置**：YAML 文件，禁止硬编码敏感信息
- **TDD**：先写测试再实现；`common/` 工具须有单测，领域变更须集成测试
- **代码生成**：改 `.api`/`.proto` 后必须 `make create-service <service>` 重生成，不手写 handler/server

## Git Commit Convention
提交commit前先检查是否在项目根目录，如果不在的根目录，请先cd到项目根目录


格式与 scope 规则：见 `.specify/memory/constitution.md` §Git 提交格式。

```
<type>(<scope>-<service-name>-backend): <description>

feat | fix | refactor | docs | style | perf | test
# 示例
feat(applicant-borrower-backend): add loan application endpoint
fix(document-lender-backend): fix document upload validation error
```

**严格约束**：
- **禁止将 LLM（Claude、Copilot 等）列为 committer 或 co-author**
- **多功能 → 多 commit**：按功能分批次提交，不要拆得过细。同一类功能（如纯文档整理、纯重构）的所有相关文件合并为一个 commit，在 body 中逐项列出文件变化
  - 反例：把一次文档整理拆成 3 个 `docs` commit（归档、入口重构、加提醒）→ 应合并为 1 个，body 列各文件
  - 正例：一个 `feat` 代码改动 + 一个 `docs` 文档整理 = 2 个 commit

## Language Rules

- 代码注释：英文，简洁
- Git 提交：英文，conventional commits
- 交流沟通：中文
- Markdown 文档：中文（`README.md` 除外，用英文）
- LLM 生成的中间文档（plan、analysis、design）：中文
