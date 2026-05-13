# API 文档索引

> Zeebee 后端 API 接口文档索引
> 最后更新：2026-05-12

## 概述

3个 HTTP API 服务 + 3个 gRPC 服务，1个 PDF 生成服务。

| 服务 | HTTP 端口 | RPC 端口 | Swagger |
|------|-----------|----------|---------|
| Borrower | 10000 | 10001 | `app/borrower/api/swagger/` |
| Lender | 11000 | 11001 | `app/lender/api/swagger/` |
| Admin | 12000 | 12001 | `app/admin/api/swagger/` |
| PDF Generator | - | - | `app/pdf_generator/api/swagger/` |

---

## Borrower API (Port 10000)

**API 定义**: `app/borrower/api/desc/`
**Handler 目录**: `app/borrower/api/internal/handler/`
**Logic 目录**: `app/borrower/api/internal/logic/`

### 模块和端点

| 模块 | 端点 | Handler 文件 |
|------|------|-------------|
| **user** | register, login, get_profile, update_profile, change_password, verify_email, user_menu | `handler/user/` |
| **borrower_applicant** | create, update, delete, all, check_existed | `handler/borrower_applicant/` |
| **borrower_entity** | create, update, delete, all, check_existed | `handler/borrower_entity/` |
| **borrower_property** | property_create, property_update, property_delete, property_detail, property_page | `handler/borrower_property/` |
| **loan** | create, apply, detail, list, update, update_pending_status, dashboard, send_quote_email | `handler/loan/` |
| **loan_applicant** | borrower_names | `handler/loan_applicant/` |
| **loan_document** | add_files | `handler/loan_document/` |
| **organization** | organization_lender_page, organization_lenders | `handler/organization/` |
| **organization_email** | create, update, delete, detail, page, authorize, generate_auth_code_url, save_email | `handler/organization_email/` |
| **organization_user** | create, update, delete, detail, page, switch_role, assign_roles, send_email, accept_invitation | `handler/organization_user/` |
| **common** | get_us_states, get_us_cities, get_uuid, select_options | `handler/common/` |
| **file** | delete | `handler/file/` |
| **pdf** | generate_loan_pdf | `handler/pdf/` |
| **storage** | * | `handler/storage/` |
| **role** | * | `handler/role/` |
| **resource** | * | `handler/resource/` |
| **zeitro** | zeitro_listing_search | `handler/zeitro/` |

**额外路由**: `handler/extra_routes.go`

---

## Lender API (Port 11000)

**API 定义**: `app/lender/api/desc/`
**Handler 目录**: `app/lender/api/internal/handler/`

### 模块和端点

| 模块 | 端点 | Handler 文件 |
|------|------|-------------|
| **user** | register, login, get_profile, update_profile, change_password, user_menu | `handler/user/` |
| **loan** | page, detail, update, update_pending_status, update_stage, review_pending_status, loan_documents_opts | `handler/loan/` |
| **loan_document** | add_files, create_loan_document, delete_loan_document_request | `handler/loan_document/` |
| **organization** | * | `handler/organization/` |
| **organization_user** | * | `handler/organization_user/` |
| **organization_setting** | * | `handler/organization_setting/` |
| **organization_template** | * | `handler/organization_template/` |
| **common** | get_us_states, get_us_cities, select_options, sync_data | `handler/common/` |
| **file** | delete | `handler/file/` |
| **pdf** | html_to_pdf | `handler/pdf/` |
| **storage** | * | `handler/storage/` |
| **role** | * | `handler/role/` |
| **resource** | * | `handler/resource/` |

---

## Admin API (Port 12000)

**API 定义**: `app/admin/api/desc/`
**Handler 目录**: `app/admin/api/internal/handler/`

### 模块和端点

| 模块 | 端点 | Handler 文件 |
|------|------|-------------|
| **user** | register, login, get_profile, update_profile, change_password, user_menu | `handler/user/` |
| **api** | create, update, delete, detail, tree | `handler/api/` |
| **base_template** | create, update, delete, detail, page, preview, reset, switch_status | `handler/base_template/` |
| **organization** | create_lender, update_lender, delete, detail, page, select_option, switch_status | `handler/organization/` |
| **organization_user** | * | `handler/organization_user/` |
| **loan** | loan_page, loan_detail, loan_reassign, lender_companies | `handler/loan/` |
| **common** | get_us_states, get_us_cities, select_options, sync_data | `handler/common/` |
| **storage** | * | `handler/storage/` |
| **role** | * | `handler/role/` |
| **resource** | * | `handler/resource/` |

---

## gRPC Proto 定义

| 服务 | Proto 文件 | RPC 目录 |
|------|-----------|----------|
| Admin | `app/admin/rpc/pb/admin.proto` | `app/admin/rpc/internal/logic/` |
| Borrower | `app/borrower/rpc/pb/borrower.proto` | `app/borrower/rpc/internal/logic/` |
| Lender | `app/lender/rpc/pb/lender.proto` | `app/lender/rpc/internal/logic/` |

### Borrower RPC 模块 (18个)

auth, borrowerapplicant, borrowerentity, borrowerproperty, common, file, loan, loanapplicant, loandocument, organization, organizationemail, organizationuser, resource, role, storage, user, zeitro

### Lender RPC 模块 (14个)

auth, common, file, loan, loandocument, organization, organizationsetting, organizationtemplate, organizationuser, resource, role, storage, user

### Admin RPC 模块 (12个)

api, auth, basetemplate, common, loan, organization, organizationuser, resource, role, storage, user

---

## API 定义文件 (.api)

各服务的 HTTP 接口通过 `.api` 文件定义：

- **Borrower**: `app/borrower/api/desc/` (common.api, loan.api, user.api, property.api 等)
- **Lender**: `app/lender/api/desc/`
- **Admin**: `app/admin/api/desc/`

---

## 开发流程

```
1. 定义 .api 文件 → app/{service}/api/desc/
2. 生成代码 → make create-api-service
3. 实现 Handler → app/{service}/api/internal/handler/
4. 实现 Logic → app/{service}/api/internal/logic/
5. 调用 RPC → RPC Logic 处理业务
```
