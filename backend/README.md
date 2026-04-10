# Zeebee Backend

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![go-zero](https://img.shields.io/badge/go--zero-1.9.3-green.svg)](https://go-zero.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/your-org/zeebee-backend/actions)
[![Coverage](https://img.shields.io/badge/coverage-<2%25-red.svg)](https://github.com/your-org/zeebee-backend)

A modern microservices architecture loan management system backend built on the go-zero framework, providing
comprehensive digital solutions for the lending industry.

## 🌟 Features

- **Microservices Architecture**: Scalable microservices design based on go-zero
- **Shared Data Layer**: Unified domain models and data access layer for code reuse
- **Modern Technology Stack**: Go 1.24+, PostgreSQL, Redis, AWS S3
- **Complete Business Flow**: Support for the entire loan lifecycle from application to approval and disbursement
- **Secure Authentication**: JWT authentication, password encryption, API permission control
- **Comprehensive Documentation**: Auto-generated Swagger docs and complete API index

## 📋 项目概述

这是一个基于 Go 语言的微服务贷款管理系统，使用 go-zero 框架构建。系统包含多个不同业务领域的服务（borrower、lender、admin），并具有共享的通用组件。

### 项目规模

- **Go 文件数量**: 775 个
    - app 服务层: 94 个 Go 文件 (borrower + lender + admin 服务)
    - domains 共享层: 637 个 Go 文件 (共享数据模型、仓储和业务逻辑)
    - common 通用组件: 44 个 Go 文件
    - 测试文件: 12 个 (*_test.go)
- **API 定义文件**: 53 个 (.api)
- **Proto 文件**: 3 个 (.proto，仅 app 目录下项目自身文件)
- **代码复用**: domains 层实现 borrower、lender 与 admin 服务代码复用

### 业务领域

- **贷款生命周期管理**: 完整的贷款申请、审批、放款流程
- **借款方管理**: 借款实体、房产、申请人管理
- **贷款方管理**: 贷款审核、文档管理、审批流程
- **组织用户管理**: 组织架构、用户角色、权限管理
- **预审批引擎**: 基于经验评分的自动化预审批系统
- **系统管理**: 借款方管理、权限管理、组织架构、用户管理

## 🏗️ 系统架构

### 服务结构

- **微服务架构**：每个业务域都有自己的 API 服务和 RPC 服务
- **双层服务**：API 层处理 HTTP 端点，RPC 层处理服务间通信
- **共享组件**：通用库用于工具、常量、数据库操作和业务逻辑

### 服务领域

- **borrower**：处理借款方操作（实体、房产、贷款、申请）
- **lender**：管理贷款方操作（贷款审核、文档、批准流程）
- **admin**：管理后台操作（借款方管理、权限管理、用户管理）
- **domains/zeebee/**：共享数据层，实现 borrower、lender 与 admin 代码复用
    - `model/`：统一的业务数据模型（Base、Loan、Borrower、Property 等）
    - `repo/`：共享的数据访问层，提供通用的 CRUD 和查询操作
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
- **业务服务引用**: app/borrower 和 app/lender 服务通过 import 引用 domains 层实现代码复用

### Architecture Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                    Zeebee Backend System                    │
├─────────────────┬─────────────────┬─────────────────────────────────────┤
│   Borrower API  │    Lender API   │          Admin API             │
│   (Port: 10000)  │    (Port: 11000)   │        (Port: 12000)            │
└─────────┬───────┴─────────────────┘─────────────┬─────────────────┘
          │                            │                      │
          ▼                            ▼                      ▼
┌─────────────────────────────┬─────────────────────┐         │
│       Borrower RPC           │     Lender RPC     │         │
│       (Port: 10001)           │     (Port: 11001)    │         │
└─────────────────┬───────────────┴─────────────────┘         │
                  │                                         │
                  └────────────────┬─────────────────────┘
                                 ▼
                  ┌─────────────────────────────┐
                  │      Shared Domain Layer      │
                  │    (domains/zeebee/)          │
                  │  • Models                │
                  │  • Repositories          │
                  │  • Business Logic        │
                  │  • Constants             │
                  └─────────────────┬───────────┘
                                ▼
                  ┌─────────────────────────────┐
                  │     Common Components      │
                  │  • Utils                  │
                  │  • Error Handling        │
                  │  • Database (GORM)       │
                  │  • Cache (Redis)          │
                  │  • Storage (AWS S3)       │
                  └─────────────────┬───────────┘
                                │
                  ┌─────────────────────────────┐
                  │        Infrastructure       │
                  │  • PostgreSQL (14+)       │
                  │  • Redis (7+)              │
                  │  • AWS S3                  │
                  │  • Docker/K8s             │
                  │  • Monitoring             │
                  └─────────────────────────────┘
```

### Core Services

| Service          | Port  | Port Type | Primary Functions            |
|------------------|-------|-----------|------------------------------|
| **borrower-api** | 10000 | HTTP API  | Borrower business interfaces |
| **borrower-rpc** | 10001 | gRPC      | Borrower business logic      |
| **lender-api**   | 11000 | HTTP API  | Lender management interfaces |
| **lender-rpc**   | 11001 | gRPC      | Lender review logic          |
| **admin-api**    | 12000 | HTTP API  | Admin backend interfaces     |
| **admin-rpc**    | 12001 | gRPC      | Management function logic    |

## 🚀 快速开始

### 1. 环境准备

- Go 1.24.0+
- Docker 20.10+
- PostgreSQL 14+
- Redis 7+

### 2. 开发环境设置

```bash
# 安装 goctl CLI 工具
go install github.com/zeromicro/go-zero/tools/goctl@latest
goctl env check --install --verbose --force

# 安装验证插件
go install github.com/linabellbiu/goctl-validate@latest

# 安装代码质量工具
go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest

# 验证安装
goctl version
```

### 3. 服务启动

```bash
# 启动 Borrower 服务
cd app/borrower
go run cmd/main.go -f etc/config-local.yaml  # (Http Port: 10000, RPC Port: 10001)

# 启动 Lender 服务
cd app/lender
go run cmd/main.go -f etc/config-local.yaml  # (Http Port: 11000, RPC Port: 11001)

# 启动 Admin 服务
cd app/admin
go run cmd/main.go -f etc/config-local.yaml  # (Http Port: 12000, RPC Port: 12001)
```

### 4. 代码生成

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

# 生成 API 服务、RPC 服务、Swagger 文档
make create-service service-name
```

## 🔧 开发工作流程

### 架构分层原则

- **API 服务 (HTTP 层)**: 负责数据校验、参数整合和请求转发
- **RPC 服务 (业务层)**: 负责核心业务逻辑处理和数据操作
- **领域服务 (共享层)**: 提供跨服务的共享业务逻辑和数据模型
- **数据流向**: API → Handler → Logic → RPC → Logic → Domain Logic → Repo → Model → Database

### API 服务开发流程

1. **消息体定义**: 在 `app/{service}/api/desc/` 目录下对应的业务域 `.api` 文件中定义请求/响应消息体
2. **HTTP 服务定义**: 在 `app/{service}/api/` 目录下的主要 `.api` 文件中引用消息体并定义 HTTP 服务和路由
3. **代码生成**: 使用 `make create-api-service {service}` 生成基础代码
4. **逻辑实现**: 在 `app/{service}/api/internal/logic/` 目录找到对应的 logic 文件
5. **功能实现**: 实现 API 层的业务逻辑，主要是：

- 请求数据验证
- 参数格式化
- 调用 RPC 服务
- 响应数据整合

### API 文件定义规则

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

### API 文件示例

#### 消息体定义 (app/borrower/api/desc/loan.api)

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

#### HTTP 服务定义 (app/borrower/api/borrower.api)

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

### RPC 服务开发流程

1. **接口定义**: 在 `app/{service}/rpc/pb/` 目录下修改或创建 `.proto` 文件
2. **代码生成**: 使用 `make create-rpc-service {service}` 生成基础代码
3. **逻辑实现**: 在 `app/{service}/rpc/internal/logic/` 目录找到对应的 logic 文件
4. **功能实现**: 实现 RPC 层的核心业务逻辑，包括：

- 业务规则处理
- 数据库操作
- 复杂计算逻辑
- 事务管理

### 领域模型开发流程

1. **模型定义**: 在 `domains/zeebee/model/` 目录创建或修改共享领域模型文件
2. **仓储定义**: 在 `domains/zeebee/repo/` 目录创建或修改共享仓储操作文件
3. **逻辑定义**: 在 `domains/zeebee/logic/` 目录创建或修改共享业务逻辑文件
4. **数据验证**: 使用 `common/validator/` 进行输入验证
5. **业务引用**: 在 app/{service}/rpc/internal/logic 中引用 domains 层实现业务逻辑
6. **代码复用**: borrower、lender 和 admin 服务通过 import 共享相同的 domain 模型、repo 和 logic

## 📂 目录结构

```
zeebee-backend/
├── app/                                # Application Services Layer
│   ├── borrower/                       # Borrower Service
│   │   ├── api/                        # HTTP API Layer
│   │   │   ├── desc/                   # API definition files
│   │   │   ├── internal/               # Internal implementation
│   │   │   │   ├── handler/            # HTTP handlers
│   │   │   │   ├── logic/              # API business logic
│   │   │   │   └── svc/                # Service context
│   │   │   ├── etc/                     # Configuration files
│   │   │   └── swagger/                # API documentation
│   │   ├── rpc/                         # RPC Business Layer
│   │   │   ├── pb/                     # Protocol Buffers
│   │   │   ├── internal/               # RPC implementation
│   │   │   └── etc/                    # RPC configuration
│   │   ├── cmd/                        # Service entry
│   │   └── docs/                       # Service documentation
│   ├── lender/                         # Lender Service
│   └── admin/                          # Admin Backend
├── domains/zeebee/                     # Shared Domain Layer
│   ├── model/                          # Data models
│   │   ├── base.go                   # Base model
│   │   ├── user.go                   # User model
│   │   ├── loan.go                   # Loan model
│   │   ├── borrower.go              # Borrower model
│   │   └── ...                       # Other models
│   ├── repo/                           # Data repository layer
│   │   ├── user.go                   # User repository
│   │   ├── loan.go                   # Loan repository
│   │   └── ...                       # Other repositories
│   ├── logic/                          # Business logic layer
│   │   ├── user_logic.go             # User business logic
│   │   └── ...                       # Other business logic
│   └── constants/                      # Business constants
├── common/                             # Common Components
│   ├── utils/                          # Utility functions
│   │   ├── jwt.go                    # JWT handling
│   │   ├── password.go                # Password encryption
│   │   └── ...                       # Other utilities
│   ├── gormx/                          # Database operations
│   │   ├── base_repo.go              # Base repository
│   │   ├── page.go                   # Pagination handling
│   │   └── option.go                # Query options
│   ├── errorx/                         # Error handling
│   │   ├── error.go                   # Error definitions
│   │   └── error_group.go            # Error groups
│   ├── response/                       # Response format
│   ├── validator/                      # Data validation
│   ├── aws_s3/                         # AWS S3 integration
│   └── preapproval_engine/             # Pre-approval engine
├── specs/                             # Feature specifications
│   └── 001-lender-role-org/           # Lender role and org management
│       ├── spec.md                    # Feature specification
│       ├── plan.md                    # Implementation plan
│       ├── tasks.md                   # Task list
│       └── ...
├── Makefile                           # Build scripts
├── go.mod                             # Go module
├── CLAUDE.md                          # AI development standards
└── README.md                          # Project documentation
```

### 完整数据流向

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
// borrower 服务 logic 层引用 domains
import "zeebee-backend/domains/zeebee/repo"
import "zeebee-backend/domains/zeebee/model"

// lender 服务 logic 层也引用相同的 domains
import "zeebee-backend/domains/zeebee/repo"
import "zeebee-backend/domains/zeebee/model"

// admin 服务 logic 层也引用相同的 domains
import "zeebee-backend/domains/zeebee/repo"
import "zeebee-backend/domains/zeebee/model"
```

### 数据层详细说明 (共享架构)

- **Shared Model**: 统一的数据模型定义，包含数据结构和数据库映射关系
- **Shared Repository**: 共享的数据访问层，提供通用的 CRUD 和业务查询方法
- **Shared Logic**: 共享的业务逻辑层，提供通用的业务处理方法
- **代码复用**: Borrower、Lender 和 Admin 服务通过 import 引用相同的 domains 层实现业务逻辑复用
- **关系**: 服务层 Logic 通过 domains 层的 Repository、Model 和 Logic 进行数据操作和业务处理

### 需要了解的关键文件和目录

#### 项目配置文件

- `Makefile`: 代码生成命令和管理
- `go.mod`: 依赖项和 Go 版本管理 (当前: 1.23.0)
- `CLAUDE.md`: AI 开发规范和指导

#### 服务配置文件

- `app/borrower/etc/config.yaml`: Borrower 服务本地 API 与 RPC 配置
- `app/lender/etc/config.yaml`: Lender 服务本地 API 与 RPC 配置
- `app/admin/etc/config.yaml`: Admin 服务本地 API 与 RPC 配置

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

#### 开发工具文件

- `third_party/goctl_template/`: go-zero 代码生成模板
- `.golangci.yml`: Go 代码检查配置

## 🧪 测试和代码质量

### 运行测试

```bash
# 运行所有测试（添加 -v 获取详细输出）
go test ./...

# 运行特定包测试
go test ./common/utils/...

# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 运行性能测试
go test -bench=. ./...
```

### 代码检查

```bash
# 运行 golangci-lint 检查
golangci-lint run

# 格式化代码
go fmt ./...

# 检查依赖
go mod tidy
```

### 测试覆盖率

当前测试覆盖率: <2% (目标: 80%+)

| Module                  | Coverage | Status               |
|-------------------------|----------|----------------------|
| User Management         | 5%       | 🔴 Needs Improvement |
| Loan Management         | 3%       | 🔴 Needs Improvement |
| Organization Management | 8%       | 🟡 In Progress       |
| Utility Functions       | 45%      | 🟢 Good              |
| Error Handling          | 30%      | 🟢 Good              |

## 💡 开发实践指南

### 错误处理模式

```go
// 统一错误处理
defer func() {
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

### 数据库操作模式

#### Model 文件示例 (domains/zeebee/model/loan.go)

```go
// Model 定义数据结构和数据库映射
type Loan struct {
    Base        // 继承基础字段 (Id, CreatedAt, UpdatedAt, DeletedAt)
    Code        string       `gorm:"column:code;size:100;not null" json:"code"`
    BorrowerId  int64        `gorm:"column:borrower_id;not null" json:"borrower_id"`
    TotalAmount float64      `gorm:"column:total_amount;type:numeric(10,2)" json:"total_amount"`
    // ... 其他字段
}

func (Loan) TableName() string {
    return "loans"
}
```

#### Repo 文件示例 (domains/zeebee/repo/loan.go)

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

### RPC 调用模式

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

### 数据转换模式

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

## 📊 项目状态和改进建议

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

## 🔍 常见问题解决

### 代码生成问题

```bash
# 如果 goctl 命令找不到，重新安装
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 检查环境
goctl env check --verbose

# 如果生成失败，检查模板路径
ls third_party/goctl_template/
```

### 服务启动问题

```bash
# 检查配置文件是否存在
ls app/borrower/etc/config.yaml

# 检查端口是否被占用
netstat -tlnp | grep 10000

# 查看服务日志
go run cmd/main.go -f etc/config.yaml
```

### 数据库连接问题

```bash
# 检查数据库配置
cat app/borrower/etc/config.yaml | grep -A 5 DbConf

# 测试数据库连接
psql "host=localhost user=dev_user password=dev_password dbname=dev_db port=5432"
```

## 🐛 调试和故障排除

### 日志查看

```go
// 在代码中添加调试日志
logx.Infow("debug_info",
    logx.Field("user_id", userId),
    logx.Field("loan_id", loanId))
```

### 性能分析

```bash
# 启用 pprof
go tool pprof http://localhost:10000/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:10000/debug/pprof/heap
```

### API 测试

```bash
# 使用 curl 测试 API
curl -X POST http://localhost:10000/api/loan/create \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{"loan_purpose":"Purchase","total_amount":100000}'
```

## 📚 相关文档

### 项目文档

- **项目文档索引**: `.claude/docs/INDEX.md` - 完整项目知识库导航中心
- **架构分析报告**: `.claude/docs/analysis/analysis.md` - 系统架构深度分析
- **项目研究分析**: `.claude/docs/analysis/project-research/research_project_analysis_20251106.md` - 项目现状调研
- **文档导航指南**: `.claude/docs/guides/navigation-guide.md` - 知识库导航和搜索指南

### API 文档

- **Borrower API**: `app/borrower/api/swagger/` - Swagger 文档
- **Lender API**: `app/lender/api/swagger/` - Swagger 文档
- **Admin API**: `app/admin/api/swagger/` - Swagger 文档

### 其他文档

- **数据库设计**: `app/borrower/docs/sql/`
- **架构图**: `app/borrower/docs/diagram/`
- **开发规范**: `CLAUDE.md` - AI 开发指导文档

## 🚀 Deployment

### 开发环境

```bash
# 使用 Docker Compose
docker-compose -f deploy/service/docker-compose.yml up -d
```

### 环境配置

| Environment | Database      | Redis            | Monitoring |
|-------------|---------------|------------------|------------|
| Development | Local         | Local            | Basic      |
| Testing     | Test DB       | Test Redis       | Full       |
| Staging     | Staging DB    | Staging Redis    | Full       |
| Production  | Production DB | Production Redis | Full       |

## 🛡️ Security

- **Authentication & Authorization**: JWT Token authentication + RBAC permission control
- **Password Security**: bcrypt encryption storage
- **Data Protection**: GORM prevents SQL injection
- **API Security**: Request validation and rate limiting
- **File Security**: AWS S3 signed URLs

## 📊 Monitoring

### 性能指标

- API 响应时间: <200ms (95% 请求)
- 数据库查询时间: <100ms
- 内存使用: <512MB
- 并发处理: 1000+ 请求/秒

### 日志记录

- 结构化日志
- 分布式追踪 (计划中)
- 错误告警和通知
- 性能监控和分析

## 🌿 Branch Strategy

This project follows a Gitflow-inspired branching model defined in the project root README.md:

### Branch Types

| Branch Type | Naming Convention                         | Purpose                                            |
|-------------|-------------------------------------------|----------------------------------------------------|
| Master      | `master`                                  | Stores all production code (frontend + backend)    |
| Staging     | `staging`                                 | Stores all staging phase code (frontend + backend) |
| Feature     | `feat-{borrower/lender/admin}-{feature}`  | New feature development                            |
| Environment | `staging-{borrower/lender/admin}-backend` | Staging preparation for backend features           |
| Environment | `prod-{borrower/lender/admin}--backend`   | Direct deployment branch for backend features      |

### Backend Branch Naming

- **Feature branches**: `feat-{business-area}-{feature-name}`
  - Example: `feat-borrower-user-auth`
  - Example: `feat-lender-loan-approval`
  - Example: `feat-admin-role-management`

- **Staging branches**: `staging-{business-area}-backend`
  - Example: `staging-borrower-backend`
  - Example: `staging-lender-backend`
  - Example: `staging-admin-backend`

- **Production branches**: `prod-{business-area}--backend`
  - Example: `prod-borrower--backend`
  - Example: `prod-lender--backend`
  - Example: `prod-admin--backend`

**Note**: Production branches use double dash (`--`) separator

### Merge Sequence

```
feat-{business-area}-{feature}
    ↓
staging-{business-area}-backend
    ↓
staging (contains all staging code)
    ↓ (after staging tests pass)
prod-{business-area}--backend (direct deployment branch)
    ↓
master (contains all production code)
```

## 🤝 Contributing

### 开发流程

1. Fork 项目
2. 创建功能分支
3. 编写代码和测试
4. 提交 Pull Request
5. 代码审查
6. 合并到主分支

### 代码规范

- 遵循 Go 官方代码规范
- 使用 go-zero 框架约定
- 编写单元测试
- 添加必要的注释

## 📝 Changelog

### [v1.0.0] (2025-12-22)

#### Added

- 完整的微服务架构实现
- Borrower 和 lender 核心功能
- 用户认证和权限管理
- 文件存储和管理系统
- 预审批引擎基础实现

#### Improved

- 更新到 Go 1.24.0 和 go-zero 1.9.3
- 完善的文档和 API 索引
- 自动化测试和部署脚本
- 安全加固和性能优化

#### Technical Debt

- 提高测试覆盖率 (目标: 80%+)
- 完善监控和日志系统
- 实现自动化 CI/CD 流水线
- 优化数据库查询性能

---

**Documentation maintained by**: Development team
**Last updated**: 2026-01-09
**Project homepage**: [https://github.com/Zeitro/Zeebee-Dev](https://github.com/Zeitro/Zeebee-Dev)
