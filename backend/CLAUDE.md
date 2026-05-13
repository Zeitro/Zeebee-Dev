# CLAUDE.md

本文件为 Claude Code 在此代码仓库中工作时提供指导。

## Constitution

本项目的所有开发行为受 `.specify/memory/constitution.md` 约束，该文件定义了核心原则、开发标准、Git 提交规范、语言规范和行为准则。开始工作前必须阅读。

## 角色与架构

- **角色定位**：详见 `.specify/memory/constitution.md` §角色定位
- **架构设计**：详见 `.claude/docs/analysis/architecture.md`（必读）
- **服务与端口**：详见 `.specify/memory/constitution.md` §架构治理
- **文档规范**：详见 `.specify/memory/constitution.md` §语言规范

## 工作前检查

- [ ] 已阅读 `.specify/memory/constitution.md` 了解核心原则
- [ ] 已阅读 `.claude/docs/analysis/architecture.md` 理解架构
- [ ] 确认当前任务涉及的层次（API/RPC/Domain）和代码位置

## 工作目录

Monorepo 结构，工作范围限于 `backend/` 文件夹。specs 文件放 `backend/specs/` 目录。

## 数据转换规范

项目统一使用 `copierx.Trans` 进行数据类型转换。

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

| 转换方向 | T (源) | U (目标) | src | 返回值 |
|---------|--------|----------|-----|--------|
| DTO → Model | `dto.XxxDto` | `model.Xxx` | `&dtoData` | `*model.Xxx` |
| Model → BO | `model.Xxx` | `bo.XxxBo` | `modelData` | `*bo.XxxBo` |
| BO → Proto | `bo.XxxBo` | `admin.XxxEntity` | `boData` | `*admin.XxxEntity` |

## 错误处理规范

```go
defer func() {
    if err != nil {
        logx.WithContext(r.Context()).
            WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
            Error(err)
    }
}()

if err != nil {
    return nil, errorx.NewDefaultError("Invalid input parameters")
}
```

## 开发注意事项

1. **基础功能**：专注核心业务逻辑正确实现
2. **错误处理**：统一使用 `common/errorx/`
3. **数据验证**：使用 `common/validator/`
4. **日志**：使用 go-zero logx
5. **配置**：YAML 配置文件，禁止硬编码敏感信息
6. **代码质量**：遵循项目规范，清晰易读

## 相关文档

- **项目宪章**：`.specify/memory/constitution.md` — 核心原则、开发标准、行为准则
- **架构文档**：`.claude/docs/analysis/architecture.md` — 系统架构设计（必读）
- **文档索引**：`.claude/docs/INDEX.md` — 完整知识库导航
- **项目 README**：`README.md` — 项目概述和快速开始
