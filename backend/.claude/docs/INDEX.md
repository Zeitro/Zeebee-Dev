# Zeebee 后端知识库索引

> 本文件为知识库导航中心。架构设计、分层职责、数据流详见 [architecture.md](analysis/architecture.md)；核心原则、技术栈、开发标准详见 `.specify/memory/constitution.md`。
>
> **最后更新**: 2026-06-25 | **项目状态**: 开发阶段

---

## 快速开始路径

| 角色 / 场景 | 阅读路径 |
|-------------|----------|
| 新手入门 | `INDEX.md` → [architecture.md](analysis/architecture.md) → `CLAUDE.md` |
| API 开发 | [API 文档索引](guides/api_documentation_index.md) → `CLAUDE.md`（数据转换规范） → [architecture.md](analysis/architecture.md) |
| 领域建模 | [architecture.md](analysis/architecture.md) §代码组织 → `domains/zeebee/model/` → `domains/zeebee/data/dto/` |
| 安全/性能审查 | `reports/security/security_deep_analysis_20260226.md`、`reports/performance/performance_deep_analysis_20260226.md` |
| 审计/全量体检 | [analysis.md](analysis/analysis.md) — 90 个已验证发现 |

---

## 知识库目录结构

`.claude/docs/` 下文档按职能分目录：

| 目录 | 内容 |
|------|------|
| `analysis/` | 系统分析：[architecture.md](analysis/architecture.md)（必读架构）、[analysis.md](analysis/analysis.md)（全量审计）、历史快照 |
| `guides/` | 开发指导：[API 文档索引](guides/api_documentation_index.md)、项目结构、域架构、技术架构总览 |
| `reports/` | 时间点分析报告：`architecture/`、`code-quality/`、`security/`、`performance/` |
| `planning/` | 项目规划：进度跟踪、路线图 |
| `tasks/` | 功能开发任务记录 |

> 分析报告按日期后缀区分版本，优先阅读最新。

---

## 核心文档

| 文档 | 用途 |
|------|------|
| [CLAUDE.md](../../CLAUDE.md) | 开发规范、代码示例、文档阅读流程 |
| [architecture.md](analysis/architecture.md) | 架构设计、分层职责、目录结构、开发流程 |
| [analysis.md](analysis/analysis.md) | 全量深度审计报告（2026-06-24） |
| `.specify/memory/constitution.md` | 项目宪章、核心原则、技术栈、Git 规范 |

## 分析报告

### 2026-06-24 全量深度审计（当前）

| 分析类型 | 评分 | 报告 |
|---------|------|------|
| 安全性 | 4.0/10 | [全量深度审计](analysis/analysis.md) |
| 架构 | 6.2/10 | [全量深度审计](analysis/analysis.md) |
| 可靠性 | 4.0/10 | [全量深度审计](analysis/analysis.md) |
| 性能 | 6.0/10 | [全量深度审计](analysis/analysis.md) |
| 可维护性 | 6.0/10 | [全量深度审计](analysis/analysis.md) |

- [全量深度审计](analysis/analysis.md) — 7 区域多 agent 并行分析，90 个已验证发现（5 Critical / 23 High / 44 Medium）⭐

### 2026-02-26 综合分析

| 分析类型 | 评分 | 报告 |
|---------|------|------|
| 安全性 | 6.5/10 | [安全深度分析](reports/security/security_deep_analysis_20260226.md) |
| 性能 | 7.5/10 | [性能深度分析](reports/performance/performance_deep_analysis_20260226.md) |
| 架构 | 8.5/10 | [架构深度分析](reports/architecture/architecture_deep_analysis_20260226.md) |

### 关键发现（最新）

4. **SQL 注入**: 9 个 repo `Page()` 方法未校验排序参数，认证后可利用
5. **pprof 端点公开**: 3 个服务在生产端口暴露 heap/goroutine profile
6. **数据竞争**: pdf_generator Chrome 池 `allocCtx` 无锁读写
7. **静默错误吞没**: `_ =` 模式在 6 处丢弃关键错误导致隐匿数据损坏

---

## 指南文档

| 文档 | 用途 |
|------|------|
| [API 文档索引](guides/api_documentation_index.md) | API 接口导航（含模块端点清单） |

> 架构、目录结构、共享层、ADR 已统一收编至 [architecture.md](analysis/architecture.md)。原 project_structure / domain_architecture / technical_architecture 三份 guide 已归档至 `reports/archive/guides/`。

---

## 服务功能矩阵

| 服务 | HTTP | RPC | LoC |
|------|------|-----|-----|
| Borrower | 10000 | 10001 | ~68k |
| Lender | 11000 | 11001 | ~36k |
| Admin | 12000 | 12001 | ~43k |
| PDF Generator | 13000 | — | ~1k |

### Borrower 服务（21 个模块）

| 模块 | API | RPC | 功能 |
|------|-----|-----|------|
| auth | - | ✅ | 认证登录 |
| borrower_applicant | ✅ | ✅ | 借款申请人 CRUD |
| borrower_entity | ✅ | ✅ | 借款实体 CRUD |
| borrower_property | ✅ | ✅ | 借款房产 CRUD |
| borrower_property_transaction | ✅ | ✅ | 借款房产交易 |
| document | ✅ | ✅ | 文档管理 |
| drive | ✅ | ✅ | Drive 集成 |
| loan | ✅ | ✅ | 贷款申请和管理 |
| loan_applicant | ✅ | ✅ | 贷款申请人 |
| loan_document | ✅ | ✅ | 贷款文档管理 |
| organization | ✅ | ✅ | 组织查询 |
| organization_email | ✅ | ✅ | 组织邮箱 |
| organization_user | ✅ | ✅ | 组织用户（含邀请接受/检查） |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| file | ✅ | ✅ | 文件操作 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口（州/城市/选项） |
| pdf | ✅ | - | PDF 生成 |
| zeitro | ✅ | ✅ | Zeitro 房产搜索 |

### Lender 服务（14 个模块）

| 模块 | API | RPC | 功能 |
|------|-----|-----|------|
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

### Admin 服务（14 个模块）

| 模块 | API | RPC | 功能 |
|------|-----|-----|------|
| auth | - | ✅ | 认证登录 |
| api | ✅ | ✅ | API 权限管理 |
| base_template | ✅ | ✅ | 基础模板 CRUD |
| loan | ✅ | ✅ | 贷款管理（列表/详情/分配） |
| organization | ✅ | ✅ | 组织管理（创建/删除/状态切换） |
| organization_subscription | ✅ | ✅ | 组织订阅管理 |
| organization_user | ✅ | ✅ | 组织用户管理 |
| pricing_plan | ✅ | ✅ | 定价计划管理 |
| pricing_plan_policy | ✅ | ✅ | 定价计划策略 |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口（州/城市/数据同步） |

---

## 领域模型（38 个）

| 类别 | 文件 |
|------|------|
| 基础 | `model/base.go` |
| 用户/权限 | `model/user.go`, `model/role.go`, `model/resource.go`, `model/role_and_resource.go`, `model/api.go` |
| 组织 | `model/organization.go`, `model/organization_user.go`, `model/organization_email.go`, `model/organization_setting.go`, `model/organization_template.go`, `model/organization_subscription.go` |
| 模板 | `model/base_template.go`, `model/document_template.go` |
| 贷款 | `model/loan.go`, `model/loan_borrower.go`, `model/loan_applicant.go`, `model/loan_property.go`, `model/loan_document.go`, `model/loan_operation.go`, `model/loan_experience.go`, `model/loan_demographic.go`, `model/loan_third_service.go` |
| 借款 | `model/borrower_applicant.go`, `model/borrower_entity.go`, `model/borrower_property.go`, `model/borrower_property_transaction.go` |
| 文档/实体 | `model/document.go`, `model/entity_document.go`, `model/file.go` |
| Drive | `model/drive_authorization.go`, `model/drive_sync_job.go` |
| 定价 | `model/pricing_plan.go`, `model/pricing_plan_policy.go`, `model/pricing_plan_and_resource.go` |
| 其他 | `model/email_authorization.go`, `model/operation.go`, `model/us_geography.go` |

仓储：`repo/` — 38 个 Repository + `provider.go`（DI 容器）

## DTO

`domains/zeebee/data/dto/` — 34 个请求/响应 DTO

## Proto 定义

| 服务 | 文件 | 位置 |
|------|------|------|
| Admin | `admin.proto` | `app/admin/rpc/pb/` |
| Borrower | `borrower.proto` | `app/borrower/rpc/pb/` |
| Lender | `lender.proto` | `app/lender/rpc/pb/` |

## 通用组件（22 个）

```
common/
├── aws_s3/            # AWS S3 文件存储
├── conforming_loan_limit/ # 贷款限额计算
├── copierx/           # 泛型数据转换（核心）
├── cron/              # 定时任务引擎
├── email/             # 邮件工具
├── enums/             # 枚举定义
├── errorx/            # 统一错误处理
├── gormx/             # GORM 扩展（BaseRepo）
├── limiter/           # Redis 限流器
├── mapx/              # Map 工具
├── preapproval_engine/ # 预审批规则引擎
├── rbac/              # 角色访问控制
├── redisx/            # Redis 封装
├── response/          # 统一响应格式
├── routex/            # 路由工具
├── stringx/           # 字符串工具
├── swagger/           # Swagger 文档
├── timex/             # 时间解析与格式化
├── types/             # 类型定义
├── utils/             # JWT / 工具
├── uuid/              # ID 生成
└── validator/         # 数据验证
```

## 部署配置

| 服务 | Docker Compose | 环境 |
|------|---------------|------|
| Borrower | `deploy/borrower/docker-compose.yml` | local/staging/production |
| Lender | `deploy/lender/docker-compose.yml` | local/staging/production |
| Admin | `deploy/admin/docker-compose.yml` | local/staging/production |
| PDF Generator | `deploy/pdfgenerator/docker-compose.yml` | local/staging/production |

---

## 按目录定位速查

| 需求 | 位置 |
|------|------|
| 接口定义 | `app/{service}/api/desc/` |
| API 业务逻辑 | `app/{service}/api/internal/logic/` |
| RPC 核心业务逻辑 | `app/{service}/rpc/internal/logic/` |
| 数据模型 | `domains/zeebee/model/` |
| 数据访问（Repository） | `domains/zeebee/repo/` |
| DTO | `domains/zeebee/data/dto/` |
| 共享业务逻辑 | `domains/zeebee/logic/` |
| 认证 / RBAC | `domains/zeebee/auth/` |
| 错误处理工具 | `common/errorx/` |
| 类型转换工具 | `common/copierx/` |
| 配置文件 | `app/{service}/etc/` |
| Swagger 文档 | `app/{service}/api/swagger/` |
| 部署配置 | `deploy/{service}/docker-compose.yml` |

---

*知识库导航中心。维护：保持与代码统计同步（见 [architecture.md](analysis/architecture.md) §服务矩阵）。*
