# Zeebee Backend Architecture

> **必读文档** — 开始开发前必须理解本文档内容。
> 技术栈版本、服务端口详见 `.specify/memory/constitution.md` §架构治理。

## 架构概览

```
Client (Web/Mobile)
    │ HTTP
    ▼
┌─────────────────────────────────────────────┐
│ API Layer (HTTP)                             │
│  borrower-api : lender-api : admin-api       │
└──────────────┬──────────────────────────────┘
               │ gRPC（直连 Endpoints，dev 无 Etcd）
               ▼
┌─────────────────────────────────────────────┐
│ RPC Layer (Business Logic)                   │
│  borrower-rpc : lender-rpc : admin-rpc       │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│ Shared Domain Layer (domains/zeebee/)        │
│  model → repo → logic                        │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│ Infrastructure (PostgreSQL / Redis / S3)     │
└─────────────────────────────────────────────┘

注：app/pdf_generator 为独立二进制服务（cmd/bin/client），不走 API/RPC 双层。
```

## 服务矩阵

| 服务 | HTTP | RPC | 职责 | 备注 |
|------|------|-----|------|------|
| borrower | 10000 | 10001 | 借款方操作（实体、房产、贷款、申请） | API/RPC 双层 |
| lender | 11000 | 11001 | 贷款方操作（审核、文档、批准） | API/RPC 双层 |
| admin | 12000 | 12001 | 管理后台（组织、模板、权限） | API/RPC 双层 |
| pdf_generator | — | — | PDF 生成 | 独立 cmd 程序 |

gRPC 服务发现：API 配置 `RpcConf.Endpoints` 直连 RPC 端口（`NonBlock: true`）。Etcd 在 dev 配置中默认注释关闭。

## 分层职责

### API 层（HTTP）

- 请求参数验证、格式化
- 调用 RPC 服务
- 响应数据整合
- **禁止**包含业务逻辑

目录：`app/{service}/api/internal/logic/`

### RPC 层（gRPC）

- 核心业务逻辑
- 数据库操作（通过 Repository）
- 事务管理

目录：`app/{service}/rpc/internal/logic/`

### 共享领域层

- 统一数据模型、DTO、业务逻辑
- 所有服务共享，避免重复代码
- 依赖方向：只能向下（API → RPC → Domain → Infrastructure）

目录：`domains/zeebee/`

### 通用组件层

- 错误处理（`errorx`）、响应格式（`response`）、数据库操作（`gormx`）
- 类型转换（`copierx`）、数据验证（`validator`）、文件存储（`aws_s3`）、API 文档（`swagger`）

目录：`common/`

## 数据流

```
Request
  → API Handler（路由分发）
  → API Logic（参数验证、整合）
  → RPC Client（gRPC 调用）
  → RPC Logic（业务逻辑）
  → Repository（数据访问）
  → Model（GORM 实体）
  → PostgreSQL
```

类型转换链：`DTO → Model → BO → Proto Entity`，统一使用 `copierx.Trans`（详见 CLAUDE.md §数据转换规范）。

## 代码组织

### app/ — 应用服务

```
backend/app/
├── borrower/                 # 借款方服务
│   ├── api/                  # HTTP API 层
│   │   ├── borrower.api      # API 定义入口
│   │   ├── borrower.go       # main（生成）
│   │   ├── desc/             # API 消息体定义
│   │   ├── etc/              # 配置文件
│   │   ├── statics/          # 静态资源
│   │   ├── swagger/          # Swagger 文档产物
│   │   └── internal/
│   │       ├── config/       # 配置结构
│   │       ├── handler/      # HTTP 处理器（生成）
│   │       ├── logic/        # API 逻辑
│   │       ├── middleware/   # 中间件
│   │       ├── pdf/          # PDF 相关
│   │       ├── svc/          # 服务上下文（依赖注入）
│   │       ├── types/        # 类型定义（生成）
│   │       └── validate_rules/ # 校验规则
│   └── rpc/                  # gRPC 业务层
│       ├── borrower.go       # main（生成）
│       ├── client/           # gRPC client（生成）
│       ├── etc/              # 配置文件
│       ├── pb/               # Proto 定义（生成）
│       └── internal/
│           ├── config/       # 配置结构
│           ├── email_tpl/    # 邮件模板
│           ├── global/       # 全局状态
│           ├── logic/        # RPC 逻辑
│           ├── server/       # gRPC 服务器
│           ├── statics/      # 静态资源
│           ├── svc/          # 服务上下文
│           └── types/        # 类型定义
├── lender/                   # 贷款方服务（结构同 borrower）
├── admin/                    # 管理后台服务（结构同 borrower）
└── pdf_generator/            # PDF 生成（独立程序）
    ├── api/                  # API schema
    ├── bin/                  # 构建产物
    ├── client/               # 调用客户端
    ├── cmd/                  # 程序入口（main）
    ├── config/               # 配置结构
    └── etc/                  # 配置文件
```

### domains/zeebee/ — 共享领域层

```
domains/zeebee/
├── model/                    # GORM 数据模型
├── repo/                     # Repository（数据访问）
├── logic/                    # 共享业务逻辑
├── data/                     # DTO 定义（data/dto/）
├── constants/                # 业务常量
├── auth/                     # 认证相关
├── email/                    # 邮件服务
├── limiter/                  # 限流逻辑
├── templater/                # 模板引擎
└── zeitro/                   # Zeitro 特定逻辑
```

### common/ — 通用组件（21 包）

```
common/
├── aws_s3/                   # 文件存储（S3）
├── conforming_loan_limit/    # 贷款限额
├── copierx/                  # 类型转换（Trans）
├── email/                    # 邮件工具
├── enums/                    # 枚举定义
├── errorx/                   # 错误处理
├── gormx/                    # 数据库操作（BaseRepo）
├── limiter/                  # 限流
├── mapx/                     # Map 工具
├── preapproval_engine/       # 预审批引擎
├── rbac/                     # 角色权限
├── redisx/                   # Redis 操作
├── response/                 # 统一响应格式
├── routex/                   # 路由工具
├── stringx/                  # 字符串工具
├── swagger/                  # Swagger 文档生成
├── timex/                    # 时间工具
├── types/                    # 通用类型
├── utils/                    # 其他工具
├── uuid/                     # UUID 生成
└── validator/                # 数据验证
```

### 其他目录

```
backend/
├── deploy/                   # 部署配置（admin/borrower/lender/pdfgenerator/zeebee）
├── scripts/                  # 脚本工具
├── specs/                    # 功能规格文档
└── third_party/              # 第三方模板
```

## Repository 模式

```go
// 泛型基础仓储，T 必须实现 schema.Tabler
type BaseRepo[T schema.Tabler] struct { db *gorm.DB }

// 具体仓储嵌入基础仓储，添加业务查询方法
type LoanRepo struct {
    *gormx.BaseRepo[*model.Loan]
}
```

BaseRepo 通用方法（成对提供 `WithTx` 事务版本）：

| 类别 | 方法 |
|------|------|
| 连接/事务 | `Client`、`GetDB`、`BuildQuery`、`Transaction` |
| 写 | `Create`、`UpdateOne` / `UpdateOneById`、`Delete` / `DeleteByID` / `DeleteByIDs` |
| 读-单条 | `FindOneByColumns`、`FindColumnByColumns`、`PluckByColumns` |
| 读-列表 | `FindAll`、`Page`（分页，返回 `PageResponse`）、`FindStrDistinctColumn` |
| 存在性 | `IsExistedByColumns`、`IsExistedWithFilter` |

查询选项由 `QueryOption`（`Filter []WhereRaw`、`Order`、`Preloads`、`Fields`）驱动，通过 `BuildQuery` 组装。所有错误经 `errorx.WithStack` 包装。

## 服务通信

- **API → RPC**：gRPC 直连（类型安全、高性能）
- **服务间依赖**：`borrower-api → borrower-rpc`，`lender-api → lender-rpc`，`admin-api → admin-rpc`
- **共享层依赖**：所有 RPC 服务 → `domains/zeebee/`；所有服务 → `common/`
- **当前无跨服务 RPC 调用**（lender-rpc 不直接调用 borrower-rpc）

## 关键约束

1. **依赖方向严格向下**：API → RPC → Domain → Infrastructure，禁止反向依赖
2. **共享数据库**：所有服务使用同一 PostgreSQL 实例，统一模型定义在 `domains/zeebee/model/`
3. **类型转换统一**：所有 `DTO ↔ Model ↔ BO ↔ Proto` 转换使用 `copierx.Trans`
4. **错误处理统一**：使用 `common/errorx/`，禁止裸 `errors.New`
5. **代码生成**：使用 `Makefile` 命令生成 API/RPC 代码，不手写 handler/server
6. **数据库访问**：RPC logic 不直接操作数据库，必须经 `domains/zeebee/repo` 层

## 开发流程

### 确定任务归属

1. **确定服务**：任务属于哪个业务域？（borrower = 借款方操作，lender = 贷款方操作，admin = 管理后台）
2. **确定层次**：
   - 新增 HTTP 端点 → API 层 + RPC 层
   - 新增业务逻辑 → RPC 层（可能涉及 Domain 层）
   - 新增数据模型 → Domain 层（model + repo）
   - 新增工具函数 → `common/`
3. **涉及文件位置**：参考 §代码组织 目录结构

### 代码生成命令

```bash
# 查看所有可用命令
make help

# 一键生成某服务 API + RPC 代码（推荐）
make create-service <service>

# 单独生成
make create-api-service <service>     # 从 .api 定义生成
make create-rpc-service <service>     # 从 .proto 定义生成
make create-api <service>             # 仅 API
make create-proto <service>           # 仅 Proto
make create-swagger <service>         # Swagger 文档
make create-sql <name>                # SQL 脚本
```

### 新增端点典型步骤

1. 在 `app/{service}/api/desc/` 添加 API 消息体定义
2. 在 `app/{service}/rpc/pb/` 添加 Proto 消息和方法定义
3. 运行 `make create-service <service>`（或分别 `create-api-service` + `create-rpc-service`）生成代码
4. 在 `app/{service}/api/internal/logic/` 实现 API 逻辑（参数验证 + RPC 调用）
5. 在 `app/{service}/rpc/internal/logic/` 实现 RPC 逻辑（业务逻辑 + repo 调用）
6. 如需新数据模型，在 `domains/zeebee/model/` 添加模型，`domains/zeebee/repo/` 添加仓储

---

**最后更新**: 2026-06-16
