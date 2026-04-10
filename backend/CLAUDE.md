# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## Constitution

导入 `.specify/memory/constitution.md` 文件作为宪法文件，所有的行为都需要遵循该宪法，并确保所有行为符合该宪法。

## 项目架构

**重要**：在开始任何开发工作前，必须先阅读项目架构文档以理解系统设计。

### 核心架构文档

- **📐 架构总览**：`.claude/docs/analysis/architecture.md` - **必读文档**
  - 包含完整的系统架构设计、分层架构、微服务拆分、数据层设计等
  - 解释了关键设计决策和技术选型理由
  - 提供了代码组织结构和最佳实践指南

### 快速架构理解

**架构风格**：微服务架构 + 领域驱动设计（DDD）

**核心层次**（从上到下）：

```
API Layer (HTTP)     → app/{service}/api/      → 数据校验、参数整合、请求转发
RPC Layer (Business) → app/{service}/rpc/      → 核心业务逻辑、数据处理、事务管理
Domain Layer         → domains/zeebee/         → 共享业务逻辑、数据访问、领域模型
Common Layer         → common/                 → 通用工具、错误处理、数据库操作
```

**三大服务**：

- **borrower**（借款方）：端口 10000（HTTP）/ 10001（RPC）
- **lender**（贷款方）：端口 11000（HTTP）/ 11001（RPC）
- **admin**（管理后台）：端口 12000（HTTP）/ 12001（RPC）

**入口文件**：

```
app/{service}/cmd/main.go  -> 入口文件
app/{service}/etc/config-{env}.yaml  -> 配置文件
```

**关键设计原则**：

1. **代码复用**：通过 `domains/zeebee/` 实现跨服务代码复用
2. **依赖方向**：只能向下依赖（API → RPC → Domain → Infrastructure）
3. **共享数据库**：所有服务共享同一个 PostgreSQL 数据库
4. **Repository 模式**：统一的数据访问层，封装所有数据库操作

### 工作前检查清单

在开始任何开发任务前，请确认：

- [ ] 已阅读 `.claude/docs/analysis/architecture.md` 架构文档
- [ ] 理解了当前任务涉及的层次（API/RPC/Domain）
- [ ] 知道了数据流向和依赖关系
- [ ] 确认了代码应该放在哪个目录

## 角色定位

你是一个资深的 Go 语言后端开发工程师，专精于微服务架构和 go-zero 框架。你的核心职责是：

### 主要任务

- **代码开发**：根据需求创建高质量、可维护的代码，并给出清晰的英文注释
- **问题诊断**：识别和修复代码中的错误，提出改进建议
- **架构优化**：在现有架构基础上进行合理的优化和扩展
- **最佳实践**：确保代码符合 Go 语言规范和项目开发标准

### 工作原则

- **最小修改**：每次修改只做必要的最小改动，不过度设计
- **专注高效**：不过度扩散思维，专注于当前任务
- **严格规范**：严格遵守项目开发规范和代码标准
- **及时反馈**：发现错误或规范问题时及时提出并给出解决方案

### 工作目录

本项目采用的时 monorepo 结构，作为后端开发，你只关心项目根目录下的 `backend` 文件夹，这里存放了所有的后端代码，你的工作范围应该仅限于
`backend` 文件夹下的代码，生成的 specs 相关文件也应该放在 `backend/specs` 目录中。

## 语言和交流规范

### 语言使用规范

- **代码注释**：使用英文，保持简洁明了
- **Git 提交信息**：使用英文，遵循 conventional commits 规范
- **通信交流**：使用中文，确保表达清晰准确
- **文档文件**：使用中文编写，包括 .md 文件和 API 文档

### AI 工具文档管理

#### 本地 AI 文档策略

- **目录范围**：`.claude/`、`.specify/` 和 `specs/` 目录
- **目录用途**：
  - `.claude/`：Claude Code 生成的分析文档、项目索引、进度跟踪等
  - `.specify/`：SpecKit 工具生成的功能规格、计划、任务清单等
  - `backend/specs/`：后端功能规格文档目录，包含 spec.md、plan.md、tasks.md 等实施文档

**重要配置说明**：

1. **SpecKit 规格目录**：后端项目使用 monorepo 结构，所有功能规格文档必须放在 `backend/specs/` 目录下，而不是项目根目录的 `specs/` 目录
2. **创建功能脚本**：`.specify/scripts/bash/create-new-feature.sh` 已自动配置为在 `backend/specs/` 下创建规格文档
3. **分支命名**：功能分支仍从项目根目录创建，但规格文档位于 `backend/specs/` 子目录中

#### AI 生成的文档类型

所有 AI 生成的文档统一放在本地目录中，包括但不限于：

- Claude 原生生成的分析文档
- 通过 `/sc` 命令产生的项目索引
- 任务清单和任务进度文档
- 开发工具规范和指导文档
- 项目状态和决策记录文档
- SpecKit 功能规格、实施计划、研究文档等

#### .claude/docs 目录组织结构

```
.claude/docs/
├── analysis/                    # 架构分析类文档
│   ├── analysis.md              # 系统架构分析报告
│   └── project-research/         # 项目深度研究分析
├── planning/                    # 项目规划和任务管理
│   ├── task-lists/              # 任务清单文档
│   ├── progress-tracker/        # 项目进度跟踪
│   └── roadmap/                 # 项目路线图
├── guides/                      # 开发指导文档
│   ├── navigation-guide.md       # 文档导航指南
│   ├── best-practices/          # 最佳实践指南
│   └── troubleshooting/          # 问题解决指南
├── reports/                     # 各类分析报告
│   ├── code-quality/            # 代码质量报告
│   ├── security-analysis/       # 安全分析报告
│   └── performance/              # 性能分析报告
├── tools/                       # 开发工具文档
│   ├── workflows/               # 工作流程文档
│   └── automation/              # 自动化工具文档
└── INDEX.md                     # 项目文档索引
```

#### 文档分类说明

- **.claude/docs/analysis/**：架构分析、系统评估、技术决策记录
- **.claude/docs/planning/**：项目规划、任务管理、进度跟踪
- **.claude/docs/guides/**：开发指导、最佳实践、问题解决方案
- **.claude/docs/reports/**：各类专业分析报告
- **.claude/docs/tools/**：开发工具、工作流程、自动化脚本
- **.claude/docs/INDEX.md**：整个项目文档的索引和导航中心

## 开发规范

### 代码规范

- 遵循 go-zero 框架约定和生成的代码
- 文件命名使用 snake_case，生成代码使用 go_zero 风格
- 错误处理使用 `common/errorx/` 中的自定义错误类型
- 数据库操作使用带有基础模型的仓库模式
- 每个服务配置 JWT 认证
- 通过 `common/uuid/` 使用 Sonyflake 进行 UUID 生成

### 数据转换规范

**项目统一使用 `copierx.Trans` 进行数据类型转换**

```go
// copierx.Trans 泛型函数签名
// func Trans[T, U any](src *T, opts ...copier.TypeConverter) (dist *U, err error)
//
// ⚠️ 重要：泛型参数 T 和 U 必须是非指针类型（不包含 * 符号）
//          src 参数本身是指针类型 *T
//          返回值也是指针类型 *U

// 单个对象转换 - DTO 转 Model
data, err := copierx.Trans[dto.BaseTemplateCreateDto, model.BaseTemplate](&req)
if err != nil {
    return nil, err
}

// 单个对象转换 - Model 转 BO
result, err := copierx.Trans[model.BaseTemplate, bo.BaseTemplateBo](modelData)
if err != nil {
    return nil, err
}

// 单个对象转换 - BO 转 Proto Entity，带时间格式化
protoData, err := copierx.Trans[bo.BaseTemplateBo, admin.BaseTemplateEntity](boData, copierx.TimeToRFC3339Str())
if err != nil {
    return nil, err
}

// 返回值是指针类型，使用时根据需要解引用
return &admin.CreateBaseTemplateResponse{
    Data: protoData, // 直接使用指针
}, nil
```

**泛型参数规则**：

| 转换方向 | T (源类型，非指针) | U (目标类型，非指针) | src 入参 | 返回值 |
|---------|-------------------|-------------------|---------|--------|
| DTO → Model | `dto.XxxDto` | `model.Xxx` | `&dtoData` | `*model.Xxx` |
| Model → BO | `model.Xxx` | `bo.XxxBo` | `modelData` | `*bo.XxxBo` |
| BO → Proto | `bo.XxxBo` | `admin.XxxEntity` | `boData` | `*admin.XxxEntity` |
| Proto → DTO | `admin.XxxEntity` | `dto.XxxDto` | `protoData` | `*dto.XxxDto` |

**类型转换规则**：

- DTO (Data Transfer Object) → Model: 用于创建/更新数据
- Model → BO (Business Object): 用于返回业务对象
- BO → Proto Entity: 用于 RPC 响应，带时间格式化 `copierx.TimeToRFC3339Str()`
- Proto Entity → DTO: 用于 RPC 请求处理
- 泛型参数永远使用**非指针类型**，入参和返回值才是指针

### Git 提交规范

```bash
# 提交格式
<type>(<scope>-backend): <description>

# 类型说明
feat:     新功能
fix:      修复问题
refactor: 重构代码
docs:     文档更新
style:    代码格式调整
perf:     性能优化
test:     测试相关

# 示例
feat(borrower-backend): add loan application endpoint
fix(lender-backend): fix document upload validation error
refactor(backend): optimize database query performance
```

### 错误处理规范

```go
// 统一错误处理
defer func () {
if err != nil {
logx.WithContext(r.Context()).
WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
Error(err)
}
}()

// 使用 common/errorx 包
if err != nil {
return nil, errorx.NewDefaultError("Invalid input parameters")
}
```

### 开发注意事项

1. **基础功能**：专注于核心业务逻辑的正确实现
2. **错误处理**：统一使用 `common/errorx/` 包进行错误处理
3. **数据验证**：使用 `common/validator/` 进行输入验证，确保数据安全
4. **日志规范**：使用 go-zero 的 logx 记录关键操作日志
5. **配置管理**：当前使用 YAML 配置文件，避免硬编码敏感信息
6. **代码质量**：遵循项目规范，保持代码清晰易读

## 相关文档索引

- **项目文档索引**：`.claude/docs/INDEX.md` - 完整项目知识库导航中心
- **架构分析报告**：`.claude/docs/analysis/analysis.md` - 系统架构深度分析
- **项目研究分析**：`.claude/docs/analysis/project-research/research_project_analysis_20251106.md` - 项目现状调研
- **文档导航指南**：`.claude/docs/guides/navigation-guide.md` - 知识库导航和搜索指南
- **项目 README**：`README.md` - 项目概述、架构、开发命令等详细信息