# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## constitution

导入 `.specify/memory/constitution.md` 文件作为宪法文件, 所有的行为都需要遵循该宪法，并确保所有行为符合该宪法。

## 角色定位

你是一个资深的 Go 语言后端开发工程师，专精于微服务架构和 go-zero 框架。你的核心职责是：

### 主要任务

- **代码开发**: 根据需求创建高质量、可维护的代码，并给出清晰的英文注释
- **问题诊断**: 识别和修复代码中的错误，提出改进建议
- **架构优化**: 在现有架构基础上进行合理的优化和扩展
- **最佳实践**: 确保代码符合 Go 语言规范和项目开发标准

### 工作原则

- **最小修改**: 每次修改只做必要的最小改动，不过度设计
- **专注高效**: 不过度扩散思维，专注于当前任务
- **严格规范**: 严格遵守项目开发规范和代码标准
- **及时反馈**: 发现错误或规范问题时及时提出并给出解决方案

### 工作目录

本项目采用的时monorepo 结构，作为后端开发，你只关心项目根目录下的`backend`文件夹，这里存放了所有的后端代码，你的工作范围应该仅限于
`backend` 文件夹下的代码。

## 语言和交流规范

### 语言使用规范

- **代码注释**：使用英文，保持简洁明了
- **Git 提交信息**：使用英文，遵循 conventional commits 规范
- **通信交流**：使用中文，确保表达清晰准确
- **文档文件**：使用中文编写，包括 .md 文件和 API 文档

### AI 工具文档管理

#### 本地 AI 文档策略

- **目录范围**：`.claude/`、`.specify/` 和 `specs/` 目录
- **Git 状态**：已在 `.gitignore` 中配置忽略，不提交到代码仓库
- **使用目的**：仅用于本地开发和 AI 辅助，保留个人开发上下文
- **Claude 访问权限**：
    - ✅ 完全可读取：Claude 可以读取这些目录下的所有文件
    - ✅ 完全可写入：Claude 可以创建和更新这些目录下的文件
    - ✅ 支持 @ 引用：在对话中使用 `@` 符号可以引入这些目录下的文件作为上下文
- **目录用途**：
    - `.claude/`: Claude Code 生成的分析文档、项目索引、进度跟踪等
    - `.specify/`: SpecKit 工具生成的功能规格、计划、任务清单等
    - `specs/`: 功能规格文档目录，包含 spec.md、plan.md、tasks.md 等实施文档

#### AI 生成的文档类型

所有 AI 生成的文档统一放在本地目录中，包括但不限于：

- Claude 原生生成的分析文档
- 通过 `/sc` 命令产生的项目索引
- 任务清单和任务进度文档
- 开发工具规范和指导文档
- 项目状态和决策记录文档
- SpecKit 功能规格、实施计划、研究文档等

#### ./claude/docs 目录组织结构

```
./claude/docs/
├── analysis/                    # 架构分析类文档
│   ├── analysis.md              # 系统架构分析报告
│   └── project-research/         # 项目深度研究分析
│       └── research_project_analysis_20251106.md
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

- **./claude/docs/analysis/**: 架构分析、系统评估、技术决策记录
- **./claude/docs/planning/**: 项目规划、任务管理、进度跟踪
- **./claude/docs/guides/**: 开发指导、最佳实践、问题解决方案
- **./claude/docs/reports/**: 各类专业分析报告
- **./claude/docs/tools/**: 开发工具、工作流程、自动化脚本
- **./claude/docs/INDEX.md**: 整个项目文档的索引和导航中心

## 项目概述

这是一个基于 Go 语言的微服务贷款管理系统，使用 go-zero 框架构建。系统包含多个不同业务领域的服务（borrower、lender、admin），并具有共享的通用组件。

### 项目规模

- **Go 文件数量**: 775 个
    - app服务层: 94个Go文件 (borrower + lender + admin服务)
    - domains共享层: 637个Go文件 (共享数据模型、仓储和业务逻辑)
    - common通用组件: 44个Go文件
    - 测试文件: 12个 (*_test.go)
- **API 定义文件**: 53 个 (.api)
- **Proto 文件**: 3 个 (.proto，仅app目录下项目自身文件)
- **代码复用**: domains层实现borrower、lender与admin服务代码复用

### 业务领域

- **贷款生命周期管理**: 完整的贷款申请、审批、放款流程
- **借款方管理**: 借款实体、房产、申请人管理
- **贷款方管理**: 贷款审核、文档管理、审批流程
- **组织用户管理**: 组织架构、用户角色、权限管理
- **预审批引擎**: 基于经验评分的自动化预审批系统
- **系统管理**: 借款方管理、权限管理、组织架构、用户管理

## 系统架构

### 服务结构

- **微服务架构**：每个业务域都有自己的 API 服务和 RPC 服务
- **双层服务**：API 层处理 HTTP 端点，RPC 层处理服务间通信
- **共享组件**：通用库用于工具、常量、数据库操作和业务逻辑

### 服务领域

- **borrower**：处理借款方操作（实体、房产、贷款、申请）
- **lender**：管理贷款方操作（贷款审核、文档、批准流程）
- **admin**：管理后台操作（借款方管理、权限管理、用户管理）
- **domains/zeebee/**：共享数据层，实现borrower、lender与admin代码复用
    - `model/`：统一的业务数据模型（Base、Loan、Borrower、Property等）
    - `repo/`：共享的数据访问层，提供通用的CRUD和查询操作
    - `logic/`：共享的业务逻辑层
- **common/**：所有服务的共享组件：
    - `constants/`：贷款、房产、用户等业务常量
    - `errorx/`：错误处理工具
    - `response/`：标准化 HTTP 响应格式
    - `gormx/`：数据库操作和分页
    - `utils/`：密码、JWT、验证工具
    - `aws_s3/`：S3 文件存储操作
    - `preapproval_engine/`：贷款预审批业务逻辑

### 代码生成模式

此项目大量使用 go-zero 代码生成：

- `.api` 文件中的 API 定义生成 HTTP 处理器和路由
- `.proto` 文件中的 Proto 定义生成 gRPC 服务
- 使用 `third_party/goctl_template/` 中的自定义模板保持代码风格一致

### 数据库层

- 使用 GORM ORM 配合 MySQL/PostgreSQL 驱动
- 在 `common/gormx/base_repo.go` 中使用基础仓库模式
- 通过 GORM 插件启用软删除功能
- **统一数据模型**: 所有数据模型定义在 `domains/zeebee/model/` 目录中
- **共享数据访问层**: 所有数据访问操作定义在 `domains/zeebee/repo/` 目录中
- **业务服务引用**: app/borrower和app/lender服务通过import引用domains层实现代码复用

## 开发命令

### 环境设置

```bash
# 安装 goctl CLI 工具
go install github.com/zeromicro/go-zero/tools/goctl@latest
goctl env check --install --verbose --force

# 安装验证插件
go install github.com/linabellbiu/goctl-validate@latest
```

### 服务管理

```bash
# 生成新的 API 定义文件
make create-api service-name

# 从定义生成 API 服务
make create-api-service service-name

# 生成 proto 文件
make create-proto service-name

# 从 proto 生成 RPC 服务
make create-rpc-service service-name

# 生成 Swagger 文档
make create-swagger service-name

# 生成API服务,RPC服务,Swagger文档
make create-service service-name
```

### 运行服务

```bash
# 运行API与RPC服务
# borrower (配置在 app/borrower/etc/config-local.yaml)
# lender (配置在 app/lender/etc/config-local.yaml)
# admin (配置在 app/admin/etc/config-local.yaml)

# 从服务目录运行：
go run cmd/main.go -f etc/config-local.yaml
```

### 开发工作流程

#### 架构分层原则

- **API 服务 (HTTP 层)**: 负责数据校验、参数整合和请求转发
- **RPC 服务 (业务层)**: 负责核心业务逻辑处理和数据操作
- **领域服务 (共享层)**: 提供跨服务的共享业务逻辑和数据模型
- **数据流向**: API →Handler →Logic → RPC →Logic →Domain Logic →Repo → Model → Database

#### API 服务开发流程

1. **消息体定义**: 在 `app/{service}/api/desc/` 目录下对应的业务域 `.api` 文件中定义请求/响应消息体
2. **HTTP 服务定义**: 在 `app/{service}/api/` 目录下的主要 `.api` 文件中引用消息体并定义 HTTP 服务和路由
3. **代码生成**: 使用 `make create-api-service {service}` 生成基础代码
4. **逻辑实现**: 在 `app/{service}/api/internal/logic/` 目录找到对应的 logic 文件
5. **功能实现**: 实现 API 层的业务逻辑，主要是：
    - 请求数据验证
    - 参数格式化
    - 调用 RPC 服务
    - 响应数据整合

#### API 文件定义规则

```
app/{service}/api/
├── borrower.api                    # 主要 API 文件：定义服务、路由、引用消息体
└── desc/                          # 消息体定义目录
    ├── common.api                  # 通用消息体
    ├── user.api                    # 用户相关消息体
    ├── loan.api                    # 贷款相关消息体
    ├── property.api                # 房产相关消息体
    └── ...                         # 其他业务域消息体
```

#### API 文件示例

##### 消息体定义 (app/borrower/api/desc/loan.api)

```go
syntax = "v1"

// 贷款状态查询请求
type LoanStatusRequest {
Id int64 `path:"id"`
}

// 贷款状态查询响应
type LoanStatusResponse {
Id     int64  `json:"id"`
Status string `json:"status"`
Stage  string `json:"stage"`
}
```

##### HTTP 服务定义 (app/borrower/api/borrower.api)

```go
syntax = "v1"

import "loan.api"
import "user.api"
import "property.api"

@server(
prefix: /api
group:  borrower
middleware: AuthMiddleware
)
service borrower-api {
@handler getLoanStatus
get /loan/:id/status (LoanStatusRequest) returns (LoanStatusResponse)
}
```

#### RPC 服务开发流程

1. **接口定义**: 在 `app/{service}/rpc/pb/` 目录下修改或创建 `.proto` 文件
2. **代码生成**: 使用 `make create-rpc-service {service}` 生成基础代码
3. **逻辑实现**: 在 `app/{service}/rpc/internal/logic/` 目录找到对应的 logic 文件
4. **功能实现**: 实现 RPC 层的核心业务逻辑，包括：
    - 业务规则处理
    - 数据库操作
    - 复杂计算逻辑
    - 事务管理

#### 领域模型开发流程

1. **模型定义**: 在 `domains/zeebee/model/` 目录创建或修改共享领域模型文件
2. **仓储定义**: 在 `domains/zeebee/repo/` 目录创建或修改共享仓储操作文件
3. **逻辑定义**: 在 `domains/zeebee/logic/` 目录创建或修改共享业务逻辑文件
4. **数据验证**: 使用 `common/validator/` 进行输入验证
5. **业务引用**: 在app/{service}/rpc/internal/logic中引用domains层实现业务逻辑
6. **代码复用**: borrower、lender和admin服务通过import共享相同的domain模型、repo和logic

### 需要了解的关键文件和目录

#### 项目配置文件

- `Makefile`: 代码生成命令和管理
- `go.mod`: 依赖项和 Go 版本管理 (当前: 1.23.0)
- `CLAUDE.md`: 本文档，开发规范和指导

#### 服务配置文件

- `app/borrower/etc/config.yaml`: Borrower 服务本地API与RPC配置
- `app/lender/etc/config.yaml`: Lender 服务本地API与RPC配置
- `app/admin/etc/config.yaml`: Admin 服务本地API与RPC配置

#### 领域模型库 (`domains/zeebee/`)

- `domains/zeebee/model/`: 统一的领域数据模型
    - `base.go`: 基础模型和通用结构
    - `loan.go`: 贷款相关模型
    - `borrower.go`: 借款方相关模型
    - `property.go`: 房产相关模型
    - `user.go`: 用户模型
    - `organization.go`: 组织相关模型
    - `role.go`: 角色相关模型
- `domains/zeebee/repo/`: 共享的数据仓储层
    - `loan.go`: 贷款数据仓储操作
    - `borrower.go`: 借款方数据仓储操作
    - `property.go`: 房产数据仓储操作
    - `user.go`: 用户数据仓储操作
    - `organization.go`: 组织数据仓储操作
- `domains/zeebee/logic/`: 共享的业务逻辑层
    - `user_logic.go`: 用户业务逻辑
    - `organization_user_logic.go`: 组织用户业务逻辑
    - `role_logic.go`: 角色业务逻辑

#### 通用组件库 (`common/`)

- `common/constants/`: 业务常量定义 (贷款状态、用户类型等)
- `common/errorx/`: 自定义错误处理工具
- `common/response/`: 标准化 HTTP 响应格式
- `common/gormx/`: 数据库操作和分页 (Repository 模式)
- `common/utils/`: 通用工具函数 (JWT、密码、验证等)
- `common/validator/`: 数据验证工具
- `common/aws_s3/`: AWS S3 文件存储操作
- `common/preapproval_engine/`: 贷款预审批引擎

#### 服务目录结构 (新架构)

```
app/{service}/
├── api/                    # HTTP API 服务
│   ├── desc/              # API 定义文件 (.api)
│   ├── internal/          # 内部实现
│   │   ├── handler/       # HTTP 处理器 (接收HTTP请求)
│   │   ├── logic/         # API 层业务逻辑 (数据校验、参数整合)
│   │   └── svc/           # 服务上下文
│   └── etc/               # API 配置文件
├── rpc/                   # RPC 服务
│   ├── pb/                # Protocol Buffers (.proto)
│   ├── internal/          # 内部实现
│   │   ├── logic/         # RPC 层业务逻辑 (核心业务处理)
│   │   └── svc/           # 服务上下文
│   └── etc/               # RPC 配置文件
├── docs/                  # 服务文档
│   ├── sql/              # 数据库设计文档
│   ├── diagram/          # 架构图和流程图
│   └── README.md         # 服务说明文档
└── cmd/                   # 服务入口
    └── main.go           # 主程序

domains/zeebee/              # 共享领域层
├── model/                # 统一领域数据模型
│   ├── base.go           # 基础模型
│   ├── loan.go           # 贷款模型
│   ├── borrower.go      # 借款方模型
│   ├── property.go      # 房产模型
│   ├── organization.go  # 组织模型
│   ├── role.go          # 角色模型
│   └── ...              # 其他业务模型
├── repo/                 # 共享数据仓储层
│   ├── loan.go          # 贷款仓储
│   ├── borrower.go     # 借款方仓储
│   ├── property.go     # 房产仓储
│   ├── organization.go # 组织仓储
│   └── ...              # 其他业务仓储
└── logic/                # 共享业务逻辑层
    ├── user_logic.go    # 用户逻辑
    ├── organization_user_logic.go # 组织用户逻辑
    ├── role_logic.go    # 角色逻辑
    └── ...              # 其他业务逻辑
```

#### 完整数据流向 (新架构)

```
HTTP Request
    ↓
API Handler (app/{service}/api/internal/handler/)
    ↓
API Logic (app/{service}/api/internal/logic/)
    ↓
RPC Call (gRPC Protocol)
    ↓
RPC Logic (app/{service}/rpc/internal/logic/)
    ↓
Domain Logic (domains/zeebee/logic/*.go) [可选]
    ↓
Shared Repository (domains/zeebee/repo/*.go)
    ↓
Shared Model (domains/zeebee/model/*.go)
    ↓
Database (PostgreSQL/MySQL)

// 代码复用示例
// borrower服务logic层引用domains
import "zeebee-backend/domains/zeebee/repo"
import "zeebee-backend/domains/zeebee/model"

// lender服务logic层也引用相同的domains
import "zeebee-backend/domains/zeebee/repo"
import "zeebee-backend/domains/zeebee/model"

// admin服务logic层也引用相同的domains
import "zeebee-backend/domains/zeebee/repo"
import "zeebee-backend/domains/zeebee/model"
```

#### 数据层详细说明 (共享架构)

- **Shared Model**: 统一的数据模型定义，包含数据结构和数据库映射关系
- **Shared Repository**: 共享的数据访问层，提供通用的CRUD和业务查询方法
- **Shared Logic**: 共享的业务逻辑层，提供通用的业务处理方法
- **代码复用**: Borrower、Lender和Admin服务通过import引用相同的domains层实现业务逻辑复用
- **关系**: 服务层Logic通过domains层的Repository、Model和Logic进行数据操作和业务处理

#### 开发工具文件

- `third_party/goctl_template/`: go-zero 代码生成模板
- `.golangci.yml`: Go 代码检查配置

### 测试和代码质量

#### 运行测试

```bash
# 运行所有测试（添加 -v 获取详细输出）
go test ./...

# 运行特定包测试
go test ./common/utils/...

# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### 代码检查

```bash
# 运行 golangci-lint 检查
golangci-lint run

# 格式化代码
go fmt ./...

# 检查依赖
go mod tidy
```

### 开发实践指南

#### 错误处理模式

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

#### 数据库操作模式

##### Model 文件示例 (app/{service}/rpc/internal/repo/model/loan.go)

```go
// Model 定义数据结构和数据库映射
type Loan struct {
Base // 继承基础字段 (Id, CreatedAt, UpdatedAt, DeletedAt)
Code        string       `gorm:"column:code;size:100;not null" json:"code"`
BorrowerId  int64        `gorm:"column:borrower_id;not null" json:"borrower_id"`
TotalAmount float64      `gorm:"column:total_amount;type:numeric(10,2)" json:"total_amount"`
// ... 其他字段
}

func (Loan) TableName() string {
return "loans"
}
```

##### Repo 文件示例 (app/{service}/rpc/internal/repo/loan.go)

```go
// Repo 定义数据库操作方法
type LoanRepo struct {
*gormx.BaseRepo[*model.Loan]  // 使用泛型基础仓库
m *model.Loan
}

func NewLoanRepo(db *gorm.DB) *LoanRepo {
return &LoanRepo{
BaseRepo: gormx.NewBaseRepo[*model.Loan](db),
m:        new(model.Loan),
}
}

func (repo *LoanRepo) GetLoan(ctx context.Context, id int64) (*model.Loan, error) {
var loan model.Loan
err := repo.Client(ctx).
Where("id = ?", id).
Preload("Applicants"). // 预加载关联数据
First(&loan).Error
return &loan, err
}
```

#### RPC 调用模式

```go
// API 层调用 RPC 服务
func (l *CreateLogic) Create(req *types.LoanCreateRequest) (*types.LoanCreateResponse, error) {
// 调用 RPC 服务
rpcResp, err := l.svcCtx.BorrowerClient.CreateLoan(l.ctx, &borrower.LoanCreateRequest{
// 参数转换
})
if err != nil {
return nil, err
}

// 响应转换
return &types.LoanCreateResponse{
Id: rpcResp.Id,
}, nil
}
```

#### 数据转换模式

**项目统一使用 `copierx.Trans` 进行数据类型转换**

```go
// copierx.Trans 泛型函数签名
// func Trans[T, U any](src *T, opts ...copier.TypeConverter) (dist *U, err error)
//
// 注意：泛型参数 T 和 U 不包括指针符号 *
//       src 参数本身是指针类型 *T

// 单个对象转换
data, err := copierx.Trans[bo.OrganizationUserBo, lender.OrganizationUserEntity](result, copierx.TimeToRFC3339Str())
if err != nil {
return nil, errorx.GetFailedTransRespErr(err)
}

// 切表转换
data, err := copierx.Trans[[]*model.OrganizationUser, []*lender.OrganizationUserEntity](&results, copierx.TimeToRFC3339Str())
if err != nil {
return nil, errorx.GetFailedTransRespErr(err)
}

// 返回值是指针类型，使用时需要解引用
return &lender.OrganizationUserAllResponse{
Data: *data, // 解引用
}, nil
```

**类型转换规则**：

- BO (Business Object) → Proto Entity: `bo.XxxBo → lender.XxxEntity`
- Model → Proto Entity: `model.Xxx → lender.XxxEntity`
- 时间字段转换: 使用 `copierx.TimeToRFC3339Str()` 作为选项
- 切表转换: 泛型参数包含指针，如 `[]*model.Xxx → []*lender.XxxEntity`

## 代码风格和约定

### 代码规范

- 遵循 go-zero 框架约定和生成的代码
- 文件命名使用 snake_case，生成代码使用 go_zero 风格
- 错误处理使用 `common/errorx/` 中的自定义错误类型
- 数据库操作使用带有基础模型的仓库模式
- 每个服务配置 JWT 认证
- 通过 `common/uuid/` 使用 Sonyflake 进行 UUID 生成

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

# 其他
backend：表示后端服务backend表示后端服务，包括API和RPC
```

## 项目状态和改进建议

### 当前架构评估

- **整体评分**: 8.2/10 (基于 2025-11-24 的架构分析)
- **优势**: 清晰的微服务边界、现代化技术栈、良好的代码组织
- **挑战**: 监控体系不完善、服务治理机制缺失、测试覆盖率偏低(当前<20%)

### 技术债务优先级

#### 🔴 高优先级 (立即处理)

1. **测试覆盖率**: 提升至 80%+ (当前覆盖率较低)
2. **代码质量**: 完善错误处理和输入验证
3. **文档完善**: API 文档和开发文档补充

#### 🟡 中优先级 (业务稳定后考虑)

1. **服务治理**: 熔断降级机制 (当前项目阶段暂不需要)
2. **服务注册与发现**: Etcd 或 Consul 集成 (当前阶段不需要)
3. **监控系统**: Prometheus + Grafana 监控 (项目稳定后考虑)

#### 🟢 低优先级 (长期规划)

1. **链路追踪**: Jaeger/Zipkin 分布式追踪 (项目刚启动暂不考虑)
2. **API 网关**: 统一入口和路由管理 (业务规模扩大后考虑)
3. **CI/CD 完善**: GitHub Actions 工作流 (开发流程稳定后实施)
4. **微服务细化**: 按业务域进一步拆分服务 (业务扩展后考虑)
5. **数据库优化**: 读写分离和分库分表 (数据量增长后考虑)

### 开发注意事项 (当前阶段重点)

1. **基础功能**: 专注于核心业务逻辑的正确实现
2. **错误处理**: 统一使用 `common/errorx/` 包进行错误处理
3. **数据验证**: 使用 `common/validator/` 进行输入验证，确保数据安全
4. **日志规范**: 使用 go-zero 的 logx 记录关键操作日志
5. **配置管理**: 当前使用 YAML 配置文件，避免硬编码敏感信息
6. **代码质量**: 遵循项目规范，保持代码清晰易读

### 后续阶段考虑 (项目稳定后)

- **监控和追踪**: 业务稳定后考虑集成监控和链路追踪系统
- **服务治理**: 业务规模扩大后考虑服务注册与发现
- **性能优化**: 根据实际性能表现进行针对性优化

### 常见问题解决

#### 代码生成问题

```bash
# 如果 goctl 命令找不到，重新安装
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 检查环境
goctl env check --verbose

# 如果生成失败，检查模板路径
ls third_party/goctl_template/
```

#### 服务启动问题

```bash
# 检查配置文件是否存在
ls app/borrower/etc/config.yaml

# 检查端口是否被占用
netstat -tlnp | grep 10000

# 查看服务日志
go run cmd/main.go -f etc/config.yaml
```

#### 数据库连接问题

```bash
# 检查数据库配置
cat app/borrower/etc/config.yaml | grep -A 5 DbConf

# 测试数据库连接
psql "host=localhost user=dev_user password=dev_password dbname=dev_db port=5432"
```

### 调试和故障排除

#### 日志查看

```go
// 在代码中添加调试日志
logx.Infow("debug_info",
logx.Field("user_id", userId),
logx.Field("loan_id", loanId))
```

#### 性能分析

```bash
# 启用 pprof
go tool pprof http://localhost:10000/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:10000/debug/pprof/heap
```

#### API 测试

```bash
# 使用 curl 测试 API
curl -X POST http://localhost:10000/api/loan/create \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"loan_purpose":"Purchase","total_amount":100000}'
```

### 相关文档

- **项目文档索引**: `docs/INDEX.md` - 完整项目知识库导航中心
- **架构分析报告**: `docs/analysis/analysis.md` - 系统架构深度分析
- **项目研究分析**: `docs/analysis/project-research/research_project_analysis_20251106.md` - 项目现状调研
- **文档导航指南**: `docs/guides/navigation-guide.md` - 知识库导航和搜索指南
- **API 文档**: 各服务的 swagger 目录
    - `app/borrower/api/swagger/` - Borrower API Swagger 文档
    - `app/lender/api/swagger/` - Lender API Swagger 文档
    - `app/admin/api/swagger/` - Admin API Swagger 文档
- **数据库设计**: `app/borrower/docs/sql/`
- **架构图**: `app/borrower/docs/diagram/`

### ai 工具范围限制

当前项目采用的是monorepo，目前专注于后端开发，后端项目位于项目根目录的`backend`下，使用specify-kit工具生成文件或者读取文件都需要在
`backend`目录下，不要跃升到项目根目录下。