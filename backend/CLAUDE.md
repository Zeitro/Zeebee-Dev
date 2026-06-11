No file at that path — the compressed content was provided inline. Here's the fixed version:

---

# CLAUDE.md

Guide for Claude Code in this repo.

## Constitution

All dev work governed by `.specify/memory/constitution.md` — defines core principles, dev standards, Git commit format, language rules, code of conduct. Read before starting.

## Role & Architecture

- **Role**: see `.specify/memory/constitution.md` §角色定位
- **Architecture**: see `.claude/docs/analysis/architecture.md` (required reading)
- **Services & Ports**: see `.specify/memory/constitution.md` §架构治理
- **Doc conventions**: see `.specify/memory/constitution.md` §语言规范

## Document Reading Order

Read in this order (first session required, subsequent on demand):

1. **This file (CLAUDE.md)** — dev standards, code examples (you are reading)
2. **`.specify/memory/constitution.md`** — project constitution: core principles, tech stack, Git conventions, code of conduct
3. **`.claude/docs/analysis/architecture.md`** — architecture design: layer responsibilities, data flow, directory structure (required)

Document relationships:
- `constitution.md` defines **rules** (what must be done), this file provides **code examples** (how)
- `architecture.md` defines **structure** (code organization), this file defines **conventions** (code style)
- `.claude/docs/INDEX.md` is knowledge base **navigation index** (optional, consult on demand)

## Working Directory

Monorepo. Work scope limited to `backend/` folder. Specs go in `backend/specs/`.

## Data Conversion Convention

Project uses `copierx.Trans` for all data type conversion.

```go
// func Trans[T, U any](src *T, opts ...copier.TypeConverter) (dist *U, err error)
// ⚠️ T 和 U 必须是非指针类型，src 是 *T，返回 *U

// DTO → Model
data, err := copierx.Trans[dto.BaseTemplateCreateDto, model.BaseTemplate](&req)

// Model → BO
result, err := copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](modelData)

// BO → Proto Entity（带时间格式化）
protoData, err := copierx.Trans[bo.BaseTemplateBo, admin.BaseTemplateEntity](boData, copierx.TimeToRFC3339Str())
```

| Direction | T (Source) | U (Target) | src | Return |
|-----------|--------|----------|-----|--------|
| DTO → Model | `dto.XxxDto` | `model.Xxx` | `&dtoData` | `*model.Xxx` |
| Model → BO | `model.Xxx` | `bo.XxxBo` | `modelData` | `*bo.XxxBo` |
| BO → Proto | `bo.XxxBo` | `admin.XxxEntity` | `boData` | `*admin.XxxEntity` |

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

## Git Commit Convention

- Format: `<type>(<scope>-<service-name>-backend): <description>`, see constitution.md §Git 提交格式
- Multiple features → multiple commits
- **Never add LLM (Claude, Copilot, etc.) as committer or co-author**

## Dev Notes

1. **Core functionality**: focus on correct business logic
2. **Error handling**: use `common/errorx/`
3. **Validation**: use `common/validator/`
4. **Logging**: use go-zero logx
5. **Config**: YAML files, no hardcoded secrets
6. **Code quality**: follow project conventions, clear and readable
7. **TDD**: write tests first, then implementation
8. **API/Proto changes**: after modifying `.api` or `.proto`, run `make create-service <service-name>` to regenerate
9. **Slice conversion**: also use `copierx.Trans`, e.g. `copierx.Trans[[]dto.Item, []model.Item](&dtoSlice)`
10. **Dev doc language**: LLM-generated intermediate docs (plans, analysis, design) use Chinese
11. **Functions/methods**: max 4 params, max 3 return values
12. **数据库操作**: logic中不直接操作数据库，而是通过repo层进行操作

## Related Docs

| Doc | Purpose | When to Read |
|------|------|----------|
| `.specify/memory/constitution.md` | Project constitution: core principles, dev standards | First session required |
| `.claude/docs/analysis/architecture.md` | Architecture: layers, data flow, directory structure | First session required |
| `.claude/docs/INDEX.md` | Knowledge base nav: service matrix, model list, analysis reports | On demand |
| `README.md` | Project overview and quick start | On demand |

Fixed: restored original Chinese comments in data conversion code block (`T 和 U 必须是非指针类型` and `带时间格式化`) — the compress step had translated them inside the code fence, which violates "preserve code blocks exactly."