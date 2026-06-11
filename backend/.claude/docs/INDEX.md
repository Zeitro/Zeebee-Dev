# Zeebee 后端知识库索引

> 本文件为知识库导航中心。架构设计、分层职责、数据流详见 [architecture.md](analysis/architecture.md)；核心原则、技术栈、开发标准详见 `.specify/memory/constitution.md`。
>
> **最后更新**: 2026-05-13 | **项目状态**: 开发阶段

---

## 核心文档

| 文档 | 用途 |
|------|------|
| [CLAUDE.md](../../CLAUDE.md) | 开发规范、代码示例、文档阅读流程 |
| [architecture.md](analysis/architecture.md) | 架构设计、分层职责、目录结构、开发流程 |
| `.specify/memory/constitution.md` | 项目宪章、核心原则、技术栈、Git 规范 |

## 分析报告

### 2026-02-26 综合分析

| 分析类型 | 评分 | 报告 |
|---------|------|------|
| 🔒 安全性 | 6.5/10 | [安全深度分析](reports/security/security_deep_analysis_20260226.md) |
| ⚡ 性能 | 7.5/10 | [性能深度分析](reports/performance/performance_deep_analysis_20260226.md) |
| 🏗️ 架构 | 8.5/10 | [架构深度分析](reports/architecture/architecture_deep_analysis_20260226.md) |

- [多域深度分析](reports/code-quality/multi_domain_analysis_20260226.md) — 安全+性能+架构综合 ⭐
- [后端代码质量分析](reports/code-quality/backend_quality_analysis_20260226.md) — 全面质量评估

### 关键发现

1. **CORS 配置漏洞**: `Access-Control-Allow-Origin: *` + credentials 启用 → CSRF 风险
2. **权限中间件未实现**: `AuthorityMiddleware` 空实现 → 无权限控制
3. **Redis TLS 证书跳过**: `InsecureSkipVerify: true` → 中间人攻击风险

## 指南文档

| 文档 | 用途 |
|------|------|
| [API 文档索引](guides/api_documentation_index.md) | 完整 API 接口导航 |
| [项目结构指南](guides/project_structure_guide.md) | 代码组织和架构说明 |
| [域架构设计指南](guides/domain_architecture_guide.md) | 共享数据层设计 |
| [导航指南](guides/navigation-guide.md) | 知识库使用指南 |

---

## 服务功能矩阵

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
| organization_user | ✅ | ✅ | 组织用户 |
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
| organization | ✅ | ✅ | 组织管理（含创建/删除/状态切换） |
| organization_user | ✅ | ✅ | 组织用户管理 |
| user | ✅ | ✅ | 用户管理 |
| role | ✅ | ✅ | 角色管理 |
| resource | ✅ | ✅ | 资源管理 |
| storage | ✅ | ✅ | 存储管理 |
| common | ✅ | ✅ | 公共接口（州/城市/数据同步） |

---

## 领域模型（27 个）

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

---

## Proto 文件

| 服务 | 文件 | 位置 |
|------|------|------|
| Admin | `admin.proto` | `app/admin/rpc/pb/` |
| Borrower | `borrower.proto` | `app/borrower/rpc/pb/` |
| Lender | `lender.proto` | `app/lender/rpc/pb/` |

第三方 Proto: `third_party/google/api/`（annotations, http 等）

## 部署配置

| 服务 | Docker Compose | 环境 |
|------|---------------|------|
| Borrower | `deploy/borrower/docker-compose.yml` | local/staging/production |
| Lender | `deploy/lender/docker-compose.yml` | local/staging/production |
| Admin | `deploy/admin/docker-compose.yml` | local/staging/production |
| PDF Generator | `deploy/pdfgenerator/docker-compose.yml` | local/staging/production |

---

*本文档由 /sc:index 生成，2026-05-13 重构。*
