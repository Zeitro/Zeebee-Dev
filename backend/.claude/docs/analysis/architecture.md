# Zeebee Backend Architecture

> **文档版本**: v1.0.0
> **最后更新**: 2026-01-09
> **架构风格**: 微服务架构 + 领域驱动设计 (DDD)
> **技术栈**: Go 1.24+ / go-zero 1.9.3 / PostgreSQL / Redis

## 目录

- [架构概览](#架构概览)
- [核心设计原则](#核心设计原则)
- [分层架构](#分层架构)
- [微服务拆分](#微服务拆分)
- [数据层设计](#数据层设计)
- [服务通信](#服务通信)
- [代码组织](#代码组织)
- [关键设计决策](#关键设计决策)
- [技术栈选型](#技术栈选型)

---

## 架构概览

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Layer                             │
│                    (Web App / Mobile App)                        │
└───────────────────────────┬─────────────────────────────────────┘
                            │ HTTP/HTTPS
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Borrower API │  │  Lender API  │  │   Admin API  │          │
│  │   (10000)    │  │   (11000)    │  │   (12000)    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          │ gRPC             │ gRPC             │ gRPC
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Business Logic Layer                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ Borrower RPC │  │  Lender RPC  │  │   Admin RPC  │          │
│  │   (10001)    │  │   (11001)    │  │   (12001)    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          └──────────────────┼──────────────────┘
                             │ shared access
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Shared Domain Layer                           │
│  ┌────────────────────────────────────────────────────────┐     │
│  │              domains/zeebee/                            │     │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │     │
│  │  │   Model     │ │   Repo      │ │   Logic     │       │     │
│  │  │ (数据模型)   │ │ (数据访问)  │ │ (业务逻辑)  │       │     │
│  │  └─────────────┘ └─────────────┘ └─────────────┘       │     │
│  └────────────────────────────────────────────────────────┘     │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  PostgreSQL  │  │    Redis     │  │   AWS S3     │          │
│  │  (Primary DB)│  │   (Cache)    │  │  (Storage)   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

### 架构特点

1. **微服务架构**：按业务域拆分服务，每个服务独立部署
2. **双层服务**：API 层（HTTP）+ RPC 层（gRPC）分离
3. **共享领域层**：通过 `domains/zeebee/` 实现代码复用
4. **统一基础设施**：共享数据库、缓存、存储等基础设施

---

## 核心设计原则

### 1. 微服务边界清晰

- **按业务域拆分**：borrower（借款方）、lender（贷款方）、admin（管理后台）
- **独立部署**：每个服务有独立的 HTTP 和 RPC 端口
- **职责单一**：每个服务专注于自己的业务领域

### 2. 代码复用优先

- **共享数据模型**：所有服务使用相同的领域模型
- **共享数据访问**：统一的 Repository 层避免重复代码
- **共享业务逻辑**：通用业务逻辑放在 `domains/zeebee/logic/`

### 3. 分层架构严格

```
API Layer (HTTP)     → 数据校验、参数整合、请求转发
RPC Layer (Business) → 核心业务逻辑、数据处理、事务管理
Domain Layer         → 共享业务逻辑、数据访问、领域模型
Infrastructure       → 数据库、缓存、存储等基础设施
```

### 4. 依赖方向明确

```
┌─────────────┐
│ API Layer   │
└──────┬──────┘
       │ depends on
       ▼
┌─────────────┐
│ RPC Layer   │
└──────┬──────┘
       │ depends on
       ▼
┌─────────────┐
│Domain Layer │
└──────┬──────┘
       │ depends on
       ▼
┌─────────────┐
│Infrastructure│
└─────────────┘
```

**重要**：依赖方向只能向下，不能向上依赖。

---

## 分层架构

### API 服务层 (HTTP Layer)

**职责**：

- 处理 HTTP 请求和响应
- 请求参数验证和格式化
- 调用 RPC 服务
- 响应数据整合

**目录结构**：

```
app/{service}/api/
├── desc/                    # API 定义文件（消息体）
│   ├── common.api          # 通用消息体
│   ├── user.api            # 用户相关消息体
│   └── loan.api            # 贷款相关消息体
├── internal/
│   ├── handler/            # HTTP 处理器（自动生成）
│   ├── logic/              # API 业务逻辑
│   ├── middleware/         # 中间件（认证、日志等）
│   ├── svc/                # 服务上下文
│   └── types/              # 类型定义（自动生成）
├── etc/                    # 配置文件
├── swagger/                # Swagger 文档
└── {service}.api           # 主 API 定义文件
```

**代码示例**：

```go
// app/borrower/api/internal/logic/loan_create_logic.go
func (l *CreateLoanLogic) CreateLoan(req *types.LoanCreateRequest) (*types.LoanCreateResponse, error) {
    // 1. 参数验证
    if req.TotalAmount <= 0 {
        return nil, errorx.NewInvalidParameterError("total_amount")
    }

    // 2. 调用 RPC 服务
    rpcResp, err := l.svcCtx.BorrowerRpcClient.CreateLoan(l.ctx, &borrower.CreateLoanRequest{
        LoanPurpose:  req.LoopPurpose,
        TotalAmount:  req.TotalAmount,
        BorrowerId:   req.BorrowerId,
    })
    if err != nil {
        return nil, errorx.NewInternalError("failed to create loan")
    }

    // 3. 响应转换
    return &types.LoanCreateResponse{
        Id:     rpcResp.Id,
        Status: rpcResp.Status,
    }, nil
}
```

### RPC 服务层 (Business Layer)

**职责**：

- 实现核心业务逻辑
- 数据库操作
- 事务管理
- 复杂计算逻辑

**目录结构**：

```
app/{service}/rpc/
├── pb/                      # Proto 定义文件
│   └── {service}.proto     # gRPC 服务定义
├── internal/
│   ├── logic/              # RPC 业务逻辑
│   └── server/             # gRPC 服务器
└── etc/                    # RPC 配置文件
```

**代码示例**：

```go
// app/borrower/rpc/internal/logic/create_loan_logic.go
func (l *CreateLoanLogic) CreateLoan(in *borrower.CreateLoanRequest) (*borrower.CreateLoanResponse, error) {
    // 1. 业务规则验证
    if in.TotalAmount > 1000000 {
        return nil, errorx.NewBusinessError("loan amount exceeds limit")
    }

    // 2. 数据库事务操作
    loan := &model.Loan{
        Code:        generateLoanCode(),
        BorrowerId:  in.BorrowerId,
        TotalAmount: in.TotalAmount,
        Status:      constant.LoanStatusDraft,
    }

    err := l.svcCtx.LoanRepo.Create(l.ctx, loan)
    if err != nil {
        return nil, errorx.NewInternalError("failed to save loan")
    }

    // 3. 返回结果
    return &borrower.CreateLoanResponse{
        Id:     loan.Id,
        Code:   loan.Code,
        Status: loan.Status,
    }, nil
}
```

### 共享领域层 (Domain Layer)

**职责**：

- 定义统一的数据模型
- 提供共享的数据访问方法
- 实现跨服务的业务逻辑

**目录结构**：

```
domains/zeebee/
├── model/                   # 数据模型（GORM 实体）
│   ├── base.go             # 基础模型
│   ├── user.go             # 用户模型
│   ├── loan.go             # 贷款模型
│   ├── borrower.go         # 借款方模型
│   └── property.go         # 房产模型
├── repo/                    # 数据访问层（Repository）
│   ├── user.go             # 用户仓储
│   ├── loan.go             # 贷款仓储
│   └── ...
├── logic/                   # 业务逻辑层
│   ├── user_logic.go       # 用户业务逻辑
│   └── ...
└── constants/               # 业务常量
    ├── loan.go             # 贷款常量
    └── user.go             # 用户常量
```

**代码示例**：

```go
// domains/zeebee/model/loan.go
type Loan struct {
    Base
    Code        string  `gorm:"column:code;size:100;not null" json:"code"`
    BorrowerId  int64   `gorm:"column:borrower_id;not null" json:"borrower_id"`
    TotalAmount float64 `gorm:"column:total_amount;type:numeric(10,2)" json:"total_amount"`
    Status      string  `gorm:"column:status;size:50;not null" json:"status"`
    // ... 其他字段
}

func (Loan) TableName() string {
    return "loans"
}

// domains/zeebee/repo/loan.go
type LoanRepo struct {
    *gormx.BaseRepo[*model.Loan]
    m *model.Loan
}

func NewLoanRepo(db *gorm.DB) *LoanRepo {
    return &LoanRepo{
        BaseRepo: gormx.NewBaseRepo[*model.Loan](db),
        m:        new(model.Loan),
    }
}

func (r *LoanRepo) GetByCode(ctx context.Context, code string) (*model.Loan, error) {
    var loan model.Loan
    err := r.Client(ctx).
        Where("code = ?", code).
        Preload("Borrower").
        First(&loan).Error
    return &loan, err
}
```

### 通用组件层 (Common Layer)

**职责**：

- 提供通用工具函数
- 统一错误处理
- 标准化响应格式
- 数据库操作封装

**目录结构**：

```
common/
├── errorx/                  # 错误处理
│   ├── error.go            # 错误定义
│   └── error_group.go      # 错误分组
├── response/                # 响应格式
│   └── response.go         # 统一响应结构
├── gormx/                   # 数据库操作
│   ├── base_repo.go        # 基础仓储
│   ├── page.go             # 分页处理
│   └── option.go           # 查询选项
├── utils/                   # 工具函数
│   ├── jwt.go              # JWT 处理
│   ├── password.go         # 密码加密
│   └── validator.go        # 数据验证
├── copierx/                 # 数据转换
│   └── copier.go           # 类型转换工具
├── aws_s3/                  # 文件存储
│   └── s3.go               # S3 操作
└── preapproval_engine/      # 预审批引擎
    └── engine.go           # 预审批逻辑
```

---

## 微服务拆分

### 服务列表

| 服务名称             | HTTP 端口 | RPC 端口 | 职责描述         |
|------------------|---------|--------|--------------|
| **borrower-api** | 10000   | -      | 借款方 HTTP 接口  |
| **borrower-rpc** | -       | 10001  | 借款方业务逻辑      |
| **lender-api**   | 11000   | -      | 贷款方 HTTP 接口  |
| **lender-rpc**   | -       | 11001  | 贷款方业务逻辑      |
| **admin-api**    | 12000   | -      | 管理后台 HTTP 接口 |
| **admin-rpc**    | -       | 12001  | 管理后台业务逻辑     |

### Borrower 服务

**业务域**：借款方管理

**核心功能**：

- 借款实体（Borrower Entity）管理
- 房产（Property）管理
- 贷款（Loan）申请
- 申请人（Applicant）管理

**API 端点示例**：

```
POST   /api/borrower              # 创建借款实体
GET    /api/borrower/:id          # 获取借款实体详情
PUT    /api/borrower/:id          # 更新借款实体
POST   /api/property              # 添加房产
POST   /api/loan/application      # 提交贷款申请
```

### Lender 服务

**业务域**：贷款方管理

**核心功能**：

- 贷款审核
- 文档管理
- 审批流程
- 条件审批（Condition）

**API 端点示例**：

```
POST   /api/lender/review         # 提交审核
GET    /api/lender/review/:id     # 获取审核详情
POST   /api/lender/document       # 上传文档
GET    /api/lender/conditions     # 获取审批条件
```

### Admin 服务

**业务域**：后台管理

**核心功能**：

- 用户管理
- 角色权限管理
- 组织架构管理
- 系统配置

**API 端点示例**：

```
POST   /api/admin/user            # 创建用户
GET    /api/admin/users           # 用户列表
POST   /api/admin/role            # 创建角色
GET    /api/admin/organizations   # 组织列表
```

---

## 数据层设计

### 数据库架构

**设计理念**：

- **共享数据库**：所有服务共享同一个 PostgreSQL 数据库
- **统一数据模型**：通过 `domains/zeebee/model/` 定义统一的数据结构
- **Repository 模式**：封装数据访问逻辑，提供统一的查询接口

### 核心数据表

```sql
-- 用户表
users (id, email, password_hash, role, created_at, updated_at, deleted_at)

-- 借款实体表
borrower_entities (id, user_id, entity_name, entity_type, ...)

-- 房产表
properties (id, borrower_entity_id, address, property_type, ...)

-- 贷款表
loans (id, borrower_entity_id, total_amount, status, stage, ...)

-- 申请人表
applicants (id, loan_id, first_name, last_name, email, ...)

-- 组织表
organizations (id, name, type, status, ...)

-- 组织用户表
organization_users (id, organization_id, user_id, role, ...)
```

### 数据访问模式

**Repository 模式**：

```go
// 基础仓储（泛型）
type BaseRepo[T any] struct {
    db *gorm.DB
}

// 通用 CRUD 方法
func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error
func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error
func (r *BaseRepo[T]) Delete(ctx context.Context, id int64) error
func (r *BaseRepo[T]) GetByID(ctx context.Context, id int64) (*T, error)
func (r *BaseRepo[T]) List(ctx context.Context, opts ...Option) ([]*T, error)

// 具体仓储（继承基础仓储）
type LoanRepo struct {
    *gormx.BaseRepo[*model.Loan]
}

// 业务特定的查询方法
func (r *LoanRepo) GetByBorrowerID(ctx context.Context, borrowerID int64) ([]*model.Loan, error)
func (r *LoanRepo) GetByStatus(ctx context.Context, status string) ([]*model.Loan, error)
```

### 数据流向

```
API Logic (参数验证)
    ↓
RPC Logic (业务逻辑)
    ↓
Domain Logic (可选，共享业务逻辑)
    ↓
Repository (数据访问)
    ↓
Model (数据模型)
    ↓
Database (PostgreSQL)
```

---

## 服务通信

### API → RPC 通信

**协议**：gRPC

**特点**：

- 类型安全：基于 Proto 定义
- 高性能：二进制序列化
- 双向流：支持流式传输

**示例**：

```go
// API 层调用 RPC
func (l *CreateLogic) Create(req *types.CreateRequest) (*types.CreateResponse, error) {
    // 调用 RPC 服务
    rpcResp, err := l.svcCtx.BorrowerRpcClient.CreateLoan(l.ctx, &borrower.CreateLoanRequest{
        BorrowerId:  req.BorrowerId,
        TotalAmount: req.TotalAmount,
    })

    return &types.CreateResponse{
        Id: rpcResp.Id,
    }, nil
}
```

### 服务间依赖

```
borrower-api → borrower-rpc
lender-api   → lender-rpc
admin-api    → admin-rpc

所有 RPC 服务 → domains/zeebee (共享领域层)
所有服务      → common (通用组件层)
```

### 未来扩展

**服务间调用**（当前未实现，未来可能需要）：

- lender-rpc → borrower-rpc（查询借款方信息）
- admin-rpc → borrower-rpc（管理借款方数据）
- admin-rpc → lender-rpc（管理贷款方数据）

**实现方式**：

- 使用 gRPC 客户端
- 服务发现（Etcd/Consul）
- 负载均衡

---

## 代码组织

### 目录结构总览

```
backend/
├── app/                         # 应用服务层
│   ├── borrower/                # 借款方服务
│   │   ├── api/                 # HTTP API 层
│   │   ├── rpc/                 # RPC 业务层
│   │   ├── cmd/                 # 服务入口
│   │   └── etc/                 # 配置文件
│   ├── lender/                  # 贷款方服务
│   └── admin/                   # 管理后台服务
├── domains/                     # 领域层
│   └── zeebee/                  # Zeebee 领域
│       ├── model/               # 数据模型
│       ├── repo/                # 数据访问
│       ├── logic/               # 业务逻辑
│       └── constants/           # 业务常量
├── common/                      # 通用组件
│   ├── errorx/                  # 错误处理
│   ├── response/                # 响应格式
│   ├── gormx/                   # 数据库操作
│   ├── utils/                   # 工具函数
│   ├── copierx/                 # 数据转换
│   └── aws_s3/                  # 文件存储
├── deploy/                      # 部署配置
├── scripts/                     # 脚本工具
├── specs/                       # 功能规格文档
├── third_party/                 # 第三方模板
├── Makefile                     # 构建脚本
├── go.mod                       # Go 模块定义
├── CLAUDE.md                    # AI 开发规范
├── ARCHITECTURE.md              # 架构文档（本文件）
└── README.md                    # 项目说明
```

### 代码生成

**API 定义**（`.api` 文件）：

```go
// app/borrower/api/desc/loan.api
syntax = "v1"

type LoanCreateRequest {
    LoanPurpose  string  `json:"loan_purpose"`
    TotalAmount  float64 `json:"total_amount"`
    BorrowerId   int64   `path:"borrower_id"`
}

type LoanCreateResponse {
    Id     int64  `json:"id"`
    Status string `json:"status"`
}
```

**生成命令**：

```bash
# 生成 API 服务代码
make create-api-service borrower

# 生成 RPC 服务代码
make create-rpc-service borrower

# 生成 Swagger 文档
make create-swagger borrower
```

---

## 关键设计决策

### 1. 为什么选择共享数据库？

**决策**：所有服务共享同一个 PostgreSQL 数据库

**原因**：

- **简化架构**：避免分布式事务的复杂性
- **数据一致性**：更容易保证数据完整性
- **开发效率**：减少服务间协调成本
- **当前阶段适合**：项目早期，业务规模可控

**权衡**：

- ✅ 优势：简单、高效、易维护
- ⚠️ 劣势：服务耦合度高，未来可能需要拆分

**未来扩展**：

- 数据量增长后考虑读写分离
- 业务规模扩大后考虑分库分表
- 服务独立化后考虑数据库拆分

### 2. 为什么有共享领域层？

**决策**：通过 `domains/zeebee/` 实现代码复用

**原因**：

- **避免重复代码**：borrower、lender、admin 都需要操作相同的数据
- **统一业务逻辑**：保证业务规则的一致性
- **提高开发效率**：共享的数据访问层减少重复工作

**权衡**：

- ✅ 优势：代码复用、逻辑统一、维护简单
- ⚠️ 劣势：服务间耦合，但通过接口可以缓解

### 3. 为什么采用双层服务（API + RPC）？

**决策**：每个服务都有 API 层和 RPC 层

**原因**：

- **职责分离**：API 层负责 HTTP，RPC 层负责业务逻辑
- **协议适配**：API 层对外提供 HTTP，RPC 层内部使用 gRPC
- **性能优化**：服务间通信使用 gRPC 更高效
- **扩展灵活**：未来可以独立扩展 API 或 RPC 层

**权衡**：

- ✅ 优势：职责清晰、性能好、易扩展
- ⚠️ 劣势：架构复杂度增加

### 4. 为什么使用 go-zero 框架？

**决策**：基于 go-zero 框架构建

**原因**：

- **工具链完善**：goctl 代码生成工具提高开发效率
- **微服务友好**：原生支持 API 和 RPC 服务
- **生产就绪**：经过大规模生产环境验证
- **社区活跃**：框架更新和维护活跃

**权衡**：

- ✅ 优势：开发效率高、工具完善、文档齐全
- ⚠️ 劣势：学习成本、框架依赖

---

## 技术栈选型

### 后端框架

| 组件         | 技术选型      | 版本    | 用途    |
|------------|-----------|-------|-------|
| **Web 框架** | go-zero   | 1.9.3 | 微服务框架 |
| **ORM**    | GORM      | 1.25+ | 数据库操作 |
| **RPC 协议** | gRPC      | 1.59+ | 服务间通信 |
| **API 协议** | HTTP/REST | 1.1   | 对外接口  |

### 数据存储

| 组件       | 技术选型       | 版本  | 用途      |
|----------|------------|-----|---------|
| **主数据库** | PostgreSQL | 14+ | 关系型数据存储 |
| **缓存**   | Redis      | 7+  | 缓存和会话存储 |
| **文件存储** | AWS S3     | -   | 文档和附件存储 |

### 开发工具

| 工具                | 版本     | 用途     |
|-------------------|--------|--------|
| **Go**            | 1.24+  | 编程语言   |
| **goctl**         | latest | 代码生成工具 |
| **golangci-lint** | latest | 代码检查工具 |
| **Docker**        | 20.10+ | 容器化部署  |

---

## 架构演进

### 当前阶段（v1.0）

- ✅ 微服务架构基础
- ✅ 共享数据库
- ✅ 共享领域层
- ✅ 双层服务（API + RPC）
- ✅ 基础监控和日志

### 下一阶段（v2.0）

**计划改进**：

- 🔲 服务注册与发现（Etcd）
- 🔲 负载均衡
- 🔲 熔断降级机制
- 🔲 分布式追踪（Jaeger）
- 🔲 监控告警（Prometheus + Grafana）

### 长期规划（v3.0）

**考虑方向**：

- 🔲 数据库拆分（每个服务独立数据库）
- 🔲 API 网关（统一入口）
- 🔲 消息队列（异步处理）
- 🔲 CQRS（读写分离）
- 🔲 领域事件（事件驱动）

---

## 总结

### 架构优势

1. **清晰的服务边界**：按业务域拆分，职责明确
2. **高效的代码复用**：共享领域层避免重复开发
3. **灵活的扩展能力**：微服务架构支持独立扩展
4. **现代化技术栈**：Go + go-zero 高性能开发

### 架构挑战

1. **服务耦合**：共享数据库和领域层导致服务间耦合
2. **分布式复杂性**：服务间协调和事务管理
3. **运维成本**：多服务部署和监控

### 最佳实践

1. **严格分层**：遵循依赖方向，避免反向依赖
2. **接口优先**：通过接口解耦服务
3. **数据验证**：API 层和 RPC 层都需要验证
4. **错误处理**：使用统一的错误处理机制
5. **日志追踪**：使用 correlation ID 追踪请求链路

---

**文档维护者**: 开发团队
**最后审查**: 2026-01-09
**下次更新**: 架构变更时及时更新
