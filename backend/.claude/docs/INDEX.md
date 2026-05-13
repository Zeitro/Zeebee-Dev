# Zeebee 后端知识库总索引

**创建时间**: 2025-11-24 | **最后更新**: 2026-05-12
**项目**: Zeebee 贷款管理系统后端
**技术栈**: Go + go-zero + gRPC + PostgreSQL
**文档类型**: 项目知识库导航中心

## 📋 概述

Zeebee Backend 是基于 go-zero 框架的微服务贷款管理系统后端，采用 DDD 分层架构，通过共享 domains 层实现跨服务代码复用。

**项目规模**: 4个微服务，27个领域模型，50+ API定义，3个 Proto 文件
**架构评分**: 8.5/10 ([2026最新分析](reports/architecture/architecture_deep_analysis_20260226.md))
**项目状态**: 开发阶段，核心功能已实现

---

## 🚀 快速导航

### 核心文档

- **[CLAUDE.md](../../CLAUDE.md)** - 开发规范和工作流程 ⭐
- **[架构总览](analysis/architecture.md)** - 系统架构设计必读文档 ⭐
- **[API文档索引](guides/api_documentation_index.md)** - 完整API接口导航
- **[项目结构指南](guides/project_structure_guide.md)** - 代码组织和架构说明
- **[域架构设计指南](guides/domain_architecture_guide.md)** - 共享数据层设计
- **[导航指南](guides/navigation-guide.md)** - 知识库使用指南

### 分析报告

- **[多域深度分析](reports/code-quality/multi_domain_analysis_20260226.md)** - 安全+性能+架构综合分析 ⭐
- **[架构深度分析](reports/architecture/architecture_deep_analysis_20260226.md)** - 架构专项评估
- **[安全深度分析](reports/security/security_deep_analysis_20260226.md)** - 安全专项评估
- **[性能深度分析](reports/performance/performance_deep_analysis_20260226.md)** - 性能专项评估

---

## 🏗️ 项目架构概览

### 四大微服务 (共享数据层)

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌───────────────┐
│ Borrower API │  │  Lender API  │  │  Admin API   │  │ PDF Generator │
│  :10000 HTTP │  │  :11000 HTTP │  │  :12000 HTTP │  │               │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘  └───────────────┘
       │                 │                 │
       ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│Borrower RPC  │  │ Lender RPC   │  │  Admin RPC   │
│  :10001      │  │  :11001      │  │  :12001      │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └────────────┬────┴─────────────────┘
                    ▼
          ┌─────────────────┐
          │  domains/zeebee │ ← 共享领域层
          │  • model (27)   │
          │  • repo  (27)   │
          │  • dto   (26)   │
          │  • logic        │
          └────────┬────────┘
                   ▼
          ┌─────────────────┐
          │  PostgreSQL     │
          │  + Redis        │
          │  + AWS S3       │
          └─────────────────┘
```

### 服务端口映射

| 服务 | HTTP 端口 | RPC 端口 | 入口文件 | 配置文件 |
|------|-----------|----------|----------|----------|
| Borrower | 10000 | 10001 | `app/borrower/cmd/main.go` | `app/borrower/etc/config-*.yaml` |
| Lender | 11000 | 11001 | `app/lender/cmd/main.go` | `app/lender/etc/config-*.yaml` |
| Admin | 12000 | 12001 | `app/admin/cmd/main.go` | `app/admin/etc/config-*.yaml` |
| PDF Generator | - | - | `app/pdf_generator/cmd/main.go` | `app/pdf_generator/etc/config-*.yaml` |

---

## 📂 目录结构详解

### 顶层结构

```
backend/
├── app/              # 4个微服务应用
├── domains/          # 共享领域层 (DDD)
├── common/           # 通用工具库
├── deploy/           # Docker 部署配置
├── scripts/          # 构建和管理脚本
├── third_party/      # 第三方 proto 定义
├── Makefile          # 代码生成和构建命令
├── go.mod            # Go 模块依赖
└── .golangci.yml     # 代码检查规则
```

### 服务内部结构 (以 borrower 为例)

```
app/borrower/
├── api/                          # HTTP API 层
│   ├── desc/                     # API 消息体定义
│   ├── etc/                      # API 服务配置
│   ├── swagger/                  # Swagger 文档
│   ├── statics/                  # 静态资源
│   └── internal/
│       ├── handler/              # HTTP 请求处理器
│       │   ├── borrower_applicant/
│       │   ├── borrower_entity/
│       │   ├── borrower_property/
│       │   ├── common/
│       │   ├── file/
│       │   ├── loan/
│       │   ├── loan_applicant/
│       │   ├── loan_document/
│       │   ├── organization/
│       │   ├── organization_email/
│       │   ├── organization_user/
│       │   ├── pdf/
│       │   ├── resource/
│       │   ├── role/
│       │   ├── storage/
│       │   ├── user/
│       │   └── zeitro/
│       ├── logic/                # API 业务逻辑
│       ├── middleware/           # 中间件
│       ├── svc/                  # 服务上下文
│       └── types/                # 请求/响应类型
├── rpc/                          # RPC 业务层
│   ├── pb/                       # Protocol Buffers 定义
│   ├── client/                   # RPC 客户端
│   └── internal/
│       ├── logic/                # RPC 业务逻辑
│       │   ├── auth/
│       │   ├── borrowerapplicant/
│       │   ├── borrowerentity/
│       │   ├── borrowerproperty/
│       │   ├── common/
│       │   ├── file/
│       │   ├── loan/
│       │   ├── loanapplicant/
│       │   ├── loandocument/
│       │   ├── organization/
│       │   ├── organizationemail/
│       │   ├── organizationuser/
│       │   ├── resource/
│       │   ├── role/
│       │   ├── storage/
│       │   ├── user/
│       │   └── zeitro/
│       ├── server/               # RPC 服务端
│       └── svc/                  # 服务上下文
├── cmd/main.go                   # 服务入口
├── etc/config-*.yaml             # 环境配置
└── config/                       # 配置结构体
```

---

## 🏛️ 领域模型 (domains/zeebee)

### 核心模型 (27个)

| 模型 | 文件 | 用途 |
|------|------|------|
| `BaseModel` | `model/base.go` | 基础模型（ID, CreatedAt, UpdatedAt, DeletedAt） |
| `User` | `model/user.go` | 用户 |
| `Role` | `model/role.go` | 角色 |
| `Resource` | `model/resource.go` | 权限资源 |
| `RoleAndResource` | `model/role_and_resource.go` | 角色-资源关联 |
| `Api` | `model/api.go` | API 权限定义 |
| `Organization` | `model/organization.go` | 组织/贷款机构 |
| `OrganizationUser` | `model/organization_user.go` | 组织-用户关联 |
| `OrganizationEmail` | `model/organization_email.go` | 组织邮箱 |
| `OrganizationSetting` | `model/organization_setting.go` | 组织设置 |
| `OrganizationTemplate` | `model/organization_template.go` | 组织模板 |
| `BaseTemplate` | `model/base_template.go` | 基础模板 |
| `Loan` | `model/loan.go` | 贷款申请 |
| `LoanBorrower` | `model/loan_borrower.go` | 贷款-借款人关联 |
| `LoanApplicant` | `model/loan_applicant.go` | 贷款申请人 |
| `LoanProperty` | `model/loan_property.go` | 贷款房产 |
| `LoanDocument` | `model/loan_document.go` | 贷款文档 |
| `LoanOperation` | `model/loan_operation.go` | 贷款操作记录 |
| `LoanExperience` | `model/loan_experience.go` | 贷款经验 |
| `LoanDemographic` | `model/loan_demographic.go` | 贷款人口统计 |
| `LoanThirdService` | `model/loan_third_service.go` | 贷款第三方服务 |
| `BorrowerApplicant` | `model/borrower_applicant.go` | 借款申请人 |
| `BorrowerEntity` | `model/borrower_entity.go` | 借款实体 |
| `BorrowerProperty` | `model/borrower_property.go` | 借款房产 |
| `File` | `model/file.go` | 文件存储 |
| `EmailAuthorization` | `model/email_authorization.go` | 邮箱验证 |
| `UsGeography` | `model/us_geography.go` | 美国地理数据 |

### 数据层架构

```
domains/zeebee/
├── model/           # 27个领域模型（GORM实体）
├── repo/            # 27个仓储（数据访问层）
│   └── provider.go  # 仓储提供者（依赖注入）
├── data/dto/        # 26个数据传输对象
├── logic/           # 共享业务逻辑
├── auth/            # 认证逻辑
├── constants/       # 业务常量
├── email/client/    # 邮件客户端
├── limiter/         # 限流器
├── templater/       # 模板引擎
│   ├── convertor/   # 模板转换器
│   └── loader/      # 模板加载器
└── zeitro/          # Zeitro API 集成
```

---

## 🔧 通用工具库 (common/)

| 模块 | 用途 |
|------|------|
| `aws_s3` | AWS S3 文件存储 |
| `conforming_loan_limit` | 贷款额度限制计算 |
| `copierx` | 泛型数据转换 (`copierx.Trans[T,U]`) |
| `email` | 邮件发送 |
| `enums` | 枚举定义 |
| `errorx` | 统一错误处理 |
| `gormx` | GORM 数据库扩展 |
| `limiter` | Redis 限流器 |
| `mapx` | Map 工具函数 |
| `preapproval_engine` | 预审批引擎 |
| `rbac` | 基于角色的访问控制 |
| `redisx` | Redis 客户端封装 |
| `response` | 统一响应格式 |
| `routex` | 路由工具 |
| `stringx` | 字符串工具 |
| `timex` | 时间工具 |
| `types` | 公共类型定义 |
| `utils` | 通用工具函数 |
| `uuid` | Sonyflake UUID 生成 |
| `validator` | 数据验证 |

---

## 📊 服务功能矩阵

### Borrower 服务 (18个 API 模块)

| 模块 | API Handlers | RPC Logic | 功能 |
|------|-------------|-----------|------|
| auth | - | ✅ | 认证登录 |
| borrower_applicant | ✅ | ✅ | 借款申请人 CRUD |
| borrower_entity | ✅ | ✅ | 借款实体 CRUD |
| borrower_property | ✅ | ✅ | 借款房产 CRUD |
| loan | ✅ | ✅ | 贷款申请和管理 |
| loan_applicant | ✅ | ✅ | 贷款申请人 |
| loan_document | ✅ | ✅ | 贷款文档管理 |
| organization | ✅ | ✅ | 组织查询 |
| organization_email | ✅ | ✅ | 组织邮箱 |
| organization_user | ✅ | ✅ | 组织用户 |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| file | ✅ | ✅ | 文件操作 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口（州/城市/选项） |
| pdf | ✅ | - | PDF 生成 |
| zeitro | ✅ | ✅ | Zeitro 房产搜索 |

### Lender 服务 (14个 API 模块)

| 模块 | API Handlers | RPC Logic | 功能 |
|------|-------------|-----------|------|
| auth | - | ✅ | 认证登录 |
| loan | ✅ | ✅ | 贷款审核和管理 |
| loan_document | ✅ | ✅ | 贷款文档管理 |
| organization | ✅ | ✅ | 组织管理 |
| organization_setting | ✅ | ✅ | 组织设置 |
| organization_template | ✅ | ✅ | 组织模板 |
| organization_user | ✅ | ✅ | 组织用户 |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| file | ✅ | ✅ | 文件操作 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口 |
| pdf | ✅ | - | PDF 生成 |

### Admin 服务 (12个 API 模块)

| 模块 | API Handlers | RPC Logic | 功能 |
|------|-------------|-----------|------|
| auth | - | ✅ | 认证登录 |
| api | ✅ | ✅ | API 权限管理 |
| base_template | ✅ | ✅ | 基础模板 CRUD |
| loan | ✅ | ✅ | 贷款管理（列表/详情/分配） |
| organization | ✅ | ✅ | 组织管理（含创建/删除/状态切换） |
| organization_user | ✅ | ✅ | 组织用户管理 |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口（州/城市/数据同步） |

---

## 📊 分析报告索引

### 2026年最新分析 (2026-02-26)

| 分析类型 | 评分 | 报告 | 状态 |
|---------|------|------|------|
| 🔒 安全性 | 6.5/10 | [安全深度分析](reports/security/security_deep_analysis_20260226.md) | 需改进 |
| ⚡ 性能 | 7.5/10 | [性能深度分析](reports/performance/performance_deep_analysis_20260226.md) | 良好 |
| 🏗️ 架构 | 8.5/10 | [架构深度分析](reports/architecture/architecture_deep_analysis_20260226.md) | 优秀 |

- **[多域深度分析](reports/code-quality/multi_domain_analysis_20260226.md)** - 综合分析 ⭐
- **[后端代码质量分析](reports/code-quality/backend_quality_analysis_20260226.md)** - 全面质量评估

### 🔴 关键发现

1. **CORS配置漏洞**: `Access-Control-Allow-Origin: *` + credentials 启用 → CSRF 风险
2. **权限中间件未实现**: `AuthorityMiddleware` 空实现 → 无权限控制
3. **Redis TLS证书跳过**: `InsecureSkipVerify: true` → 中间人攻击风险

---

## 🔍 数据流向

```
HTTP Request
    → API Handler (参数校验、路由)
    → API Logic (数据整合、转发)
    → RPC Client (gRPC 调用)
    → RPC Logic (核心业务逻辑)
    → Repository (数据访问)
    → Model (GORM 实体)
    → PostgreSQL
```

**数据转换链**: `DTO → Model → BO → Proto Entity` (使用 `copierx.Trans`)

---

## 🗺️ 知识地图

### 按功能域

```
贷款管理系统
├── 借款方 (Borrower) ← 最复杂的业务域
│   ├── 申请人管理 (BorrowerApplicant)
│   ├── 实体管理 (BorrowerEntity)
│   ├── 房产管理 (BorrowerProperty)
│   ├── 贷款申请 (Loan + LoanBorrower/Property/Document)
│   └── Zeitro 房产搜索集成
├── 贷款方 (Lender)
│   ├── 贷款审核和状态管理
│   ├── 文档管理
│   ├── 组织设置和模板
│   └── 审批流程
├── 管理后台 (Admin)
│   ├── 组织管理 (CRUD + 状态切换)
│   ├── 基础模板管理
│   ├── 贷款分配
│   └── API 权限管理
└── 通用 (Common)
    ├── 用户认证 (JWT + bcrypt)
    ├── RBAC 权限 (Role/Resource)
    ├── 文件存储 (AWS S3)
    ├── 邮件服务
    ├── 预审批引擎
    └── PDF 生成
```

---

## 📋 Proto 文件

| 服务 | Proto 文件 | 位置 |
|------|-----------|------|
| Admin | `admin.proto` | `app/admin/rpc/pb/` |
| Borrower | `borrower.proto` | `app/borrower/rpc/pb/` |
| Lender | `lender.proto` | `app/lender/rpc/pb/` |

第三方 Proto: `third_party/google/api/` (annotations, http, etc.)

---

## 🔗 部署配置

| 服务 | Docker Compose | 环境 |
|------|---------------|------|
| Borrower | `deploy/borrower/docker-compose.yml` | local/staging/production |
| Lender | `deploy/lender/docker-compose.yml` | local/staging/production |
| Admin | `deploy/admin/docker-compose.yml` | local/staging/production |
| PDF Generator | `deploy/pdfgenerator/docker-compose.yml` | local/staging/production |

---

## 📚 外部资源

- [go-zero 官方文档](https://go-zero.dev/)
- [GORM 文档](https://gorm.io/docs/)
- [Protocol Buffers 指南](https://developers.google.com/protocol-buffers)
- [golangci-lint](https://golangci-lint.run/)

---

*本文档由 /sc:index 自动生成。最后更新: 2026-05-12 | 版本: v3.0*
