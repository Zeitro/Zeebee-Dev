# Zeebee 后端项目结构指南

> 更新日期：2026-06-24 | 本文档反映项目当前实际结构

## 1. 架构总览

Zeebee 是基于 **Go + go-zero** 的贷款管理微服务系统，采用 **DDD 分层架构**。

```
┌──────────┬──────────┬──────────┬───────────────┐
│ Borrower │  Lender  │  Admin   │ PDF Generator │
│ HTTP:10K │ HTTP:11K │ HTTP:12K │  (standalone) │
│  RPC:10K1│  RPC:11K1│  RPC:12K1│  (HTTP:13K) │
├──────────┴──────────┴──────────┴───────────────┤
│              API Layer (HTTP)                    │
│       handler / logic / middleware / svc         │
├─────────────────────────────────────────────────┤
│              RPC Layer (gRPC)                    │
│        pb / client / internal/logic / server     │
├─────────────────────────────────────────────────┤
│           Domain Layer (shared)                  │
│    domains/zeebee/ - model / repo / logic / dto  │
├─────────────────────────────────────────────────┤
│           Common Layer (shared)                  │
│          common/ - 22 个通用工具包                 │
├─────────────────────────────────────────────────┤
│           Infrastructure                         │
│      PostgreSQL (GORM) / Redis / AWS S3          │
└─────────────────────────────────────────────────┘
```

**数据流**: `HTTP Request → Handler → API Logic → RPC Client → RPC Logic → Repo → Model → PostgreSQL`

**数据转换链** (`copierx.Trans`): `DTO → Model → BO → Proto Entity`

## 2. 顶层目录

```
backend/
├── app/                  # 4个微服务
├── domains/zeebee/       # 共享领域层
├── common/               # 20个通用工具包
├── deploy/               # Docker Compose 部署
├── scripts/              # 构建/管理脚本
└── third_party/          # 第三方 proto 定义
```

## 3. 服务内部结构

每个标准服务结构一致：

```
app/{service}/
├── cmd/main.go                    # 入口
├── etc/
│   ├── config-local.yaml
│   ├── config-staging.yaml
│   └── config-production.yaml
├── config/                        # 配置结构体
├── api/                           # HTTP API 层
│   ├── desc/                      # .api 定义文件
│   ├── etc/                       # API 配置
│   ├── swagger/                   # Swagger 文档
│   └── internal/
│       ├── handler/{module}/      # HTTP 处理器
│       ├── logic/{module}/        # API 逻辑
│       ├── middleware/            # 中间件
│       ├── svc/                   # ServiceContext
│       └── types/                 # 请求/响应类型
└── rpc/                           # RPC 层
    ├── pb/                        # Protobuf
    ├── client/                    # RPC 客户端
    └── internal/
        ├── logic/{module}/        # 核心业务逻辑
        ├── server/                # gRPC 服务端
        └── svc/                   # ServiceContext
```

### 服务端口

| 服务 | HTTP | RPC |
|------|------|-----|
| borrower | 10000 | 10001 |
| lender | 11000 | 11001 |
| admin | 12000 | 12001 |
| pdf_generator | 13000 | - |

### 各服务模块一览

**Borrower API (17模块)**: borrower_applicant, borrower_entity, borrower_property, common, file, loan, loan_applicant, loan_document, organization, organization_email, organization_user, pdf, resource, role, storage, user, zeitro

**Borrower RPC (18模块)**: auth + 上述17个

**Lender API (13模块)**: common, file, loan, loan_document, organization, organization_setting, organization_template, organization_user, pdf, resource, role, storage, user

**Lender RPC (14模块)**: auth + 上述13个

**Admin API (10模块)**: api, base_template, common, loan, organization, organization_user, resource, role, storage, user

**Admin RPC (12模块)**: auth, api, basetemplate, common, loan, organization, organizationuser, resource, role, storage, user

## 4. 领域层 (domains/zeebee/)

```
domains/zeebee/
├── model/          # 38个 GORM 数据模型
├── repo/           # 38个 Repository + provider.go (DI)
├── data/dto/       # 26个 DTO
├── logic/          # 共享业务逻辑 + provider.go
├── auth/           # 认证逻辑
├── constants/      # 业务常量
├── email/client/   # 邮件客户端
├── limiter/        # 限流器
├── templater/      # 模板引擎 (convertor/, loader/)
└── zeitro/         # Zeitro API 集成
```

### 核心模型

| 类别 | 模型 |
|------|------|
| 基础 | BaseModel |
| 用户 | User, Role, Resource, RoleAndResource, Api |
| 组织 | Organization, OrganizationUser, OrganizationEmail, OrganizationSetting, OrganizationTemplate |
| 借款 | BorrowerApplicant, BorrowerEntity, BorrowerProperty |
| 贷款 | Loan, LoanBorrower, LoanApplicant, LoanProperty, LoanDocument, LoanOperation, LoanExperience, LoanDemographic, LoanThirdService |
| 其他 | File, EmailAuthorization, UsGeography, BaseTemplate, PricingPlan, PlanPolicy, OrganizationSubscription |

## 5. 通用工具包 (common/)

| 包名 | 用途 |
|------|------|
| `aws_s3` | AWS S3 文件存储 |
| `conforming_loan_limit` | 贷款额度计算 |
| `copierx` | 泛型数据转换（核心） |
| `cron` | 定时任务引擎 |
| `email` | 邮件发送 |
| `enums` | 枚举定义 |
| `errorx` | 统一错误处理 |
| `gormx` | GORM 扩展（BaseRepo） |
| `limiter` | Redis 限流器 |
| `mapx` | Map 工具 |
| `preapproval_engine` | 预审批引擎 |
| `rbac` | 角色访问控制 |
| `redisx` | Redis 封装 |
| `response` | 统一响应 |
| `routex` | 路由工具 |
| `stringx` | 字符串工具 |
| `timex` | 时间解析与格式化（Duration ↔ String） |
| `types` | 类型定义 |
| `utils` | 通用工具 |
| `uuid` | Sonyflake ID 生成 |
| `validator` | 数据验证 |

### copierx 使用示例

```go
// DTO → Model
data, err := copierx.Trans[dto.XxxDto, model.Xxx](&req)
// Model → BO
result, err := copierx.Trans[model.Xxx, bo.XxxBo](modelData)
// BO → Proto（带时间格式化）
protoData, err := copierx.Trans[bo.XxxBo, admin.XxxEntity](boData, copierx.TimeToRFC3339Str())
```

## 6. 依赖方向规则

```
app/{service}/api/  →  app/{service}/rpc/  (通过 RPC Client)
app/{service}/rpc/  →  domains/zeebee/     (直接引用)
app/{service}/rpc/  →  common/             (直接引用)
domains/zeebee/     →  common/             (直接引用)

禁止反向依赖: rpc → api, domain → rpc, common → domain
```

## 7. 快速定位

| 需求 | 位置 |
|------|------|
| 新增 HTTP 接口 | `app/{service}/api/desc/` → 生成 handler/logic |
| 新增业务逻辑 | `app/{service}/rpc/internal/logic/{module}/` |
| 新增共享逻辑 | `domains/zeebee/logic/` |
| 新增数据表 | `domains/zeebee/model/` + `repo/` |
| 新增 DTO | `domains/zeebee/data/dto/` |
| 修改认证 | `domains/zeebee/auth/` |
| 修改错误处理 | `common/errorx/` |
| 查看配置 | `app/{service}/etc/` |
| 查看 Swagger | `app/{service}/api/swagger/` |
| 查看部署 | `deploy/{service}/docker-compose.yml` |
