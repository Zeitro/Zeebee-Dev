# Zeebee 后端知识库索引

> 本文件为知识库导航中心。架构设计、分层职责、数据流详见 [architecture.md](analysis/architecture.md)；核心原则、技术栈、开发标准详见 `.specify/memory/constitution.md`。
>
> **最后更新**: 2026-06-24 | **项目状态**: 开发阶段

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
| [API 文档索引](guides/api_documentation_index.md) | API 接口导航 |
| [项目结构指南](guides/project_structure_guide.md) | 代码组织和架构说明 |
| [域架构设计指南](guides/domain_architecture_guide.md) | 共享数据层设计 |
| [技术架构总览](guides/technical_architecture_overview.md) | 技术架构详解 |
| [导航指南](guides/navigation-guide.md) | 知识库使用指南 |

---

## 服务功能矩阵

| 服务 | HTTP | RPC | LoC |
|------|------|-----|-----|
| Borrower | 10000 | 10001 | ~68k |
| Lender | 11000 | 11001 | ~36k |
| Admin | 12000 | 12001 | ~43k |
| PDF Generator | 13000 | — | ~1k |

### Borrower 服务（18 个模块）

| 模块 | API | RPC | 功能 |
|------|-----|-----|------|
| auth | - | ✅ | 认证登录 |
| borrower_applicant | ✅ | ✅ | 借款申请人 CRUD |
| borrower_entity | ✅ | ✅ | 借款实体 CRUD |
| borrower_property | ✅ | ✅ | 借款房产 CRUD |
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

### Admin 服务（12 个模块）

| 模块 | API | RPC | 功能 |
|------|-----|-----|------|
| auth | - | ✅ | 认证登录 |
| api | ✅ | ✅ | API 权限管理 |
| base_template | ✅ | ✅ | 基础模板 CRUD |
| loan | ✅ | ✅ | 贷款管理（列表/详情/分配） |
| organization | ✅ | ✅ | 组织管理（创建/删除/状态切换） |
| organization_user | ✅ | ✅ | 组织用户管理 |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口（州/城市/数据同步） |

---

## 领域模型（38 个）

| 类别 | 文件 | 模型 |
|------|------|------|
| 基础 | `model/base.go` | BaseModel |
| 用户/权限 | `model/user.go`, `model/role.go`, `model/resource.go`, `model/role_and_resource.go`, `model/api.go` | User, Role, Resource, RoleAndResource, Api |
| 组织 | `model/organization.go`, `model/organization_user.go`, `model/organization_email.go`, `model/organization_setting.go`, `model/organization_template.go`, `model/organization_subscription.go` | Organization 及相关 |
| 模板 | `model/base_template.go`, `model/organization_template.go` | BaseTemplate |
| 贷款 | `model/loan.go`, `model/loan_borrower.go`, `model/loan_applicant.go`, `model/loan_property.go`, `model/loan_document.go`, `model/loan_operation.go`, `model/loan_experience.go`, `model/loan_demographic.go`, `model/loan_third_service.go` | Loan 聚合 |
| 借款 | `model/borrower_applicant.go`, `model/borrower_entity.go`, `model/borrower_property.go` | Borrower 聚合 |
| 其他 | `model/file.go`, `model/email_authorization.go`, `model/us_geography.go`, `model/pricing_plan.go`, `model/plan_policy.go` | File, EmailAuthorization, UsGeography 等 |
| 仓储 | `repo/` | 38 个 Repository + `provider.go` (DI) |

## DTO

`domains/zeebee/data/dto/` — 26+ 个请求/响应 DTO

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

*本文档由 /sc:index 生成，2026-06-24 更新。*
