# API 文档索引

> Zeebee 后端 API 接口文档索引
> 最后更新：2026-07-30

## 概述

3 个 API/RPC 双层服务（borrower / lender / admin，各含 HTTP + gRPC）+ 1 个 PDF 生成服务（pdf_generator，仅 HTTP）。端口以 `app/<service>/etc/config-<env>.yaml` 为准。

| 服务 | HTTP 端口 | RPC 端口 | Swagger |
|------|-----------|----------|---------|
| Borrower | 10000 | 10001 | `app/borrower/api/swagger/` |
| Lender | 10100 | 10101 | `app/lender/api/swagger/` |
| Admin | 10200 | 10201 | `app/admin/api/swagger/` |
| PDF Generator | 10300 | - | `app/pdf_generator/api/swagger/` |

> 下述各服务端点目录由 `app/<svc>/api/internal/logic/` 实际文件枚举；模块级权威清单（含 API/RPC 对照）见 [INDEX](../INDEX.md) §服务功能矩阵，端点定义以 `app/<svc>/api/desc/*.api` 为准。

---

## Borrower API (Port 10000)

**API 定义**: `app/borrower/api/desc/`
**Handler 目录**: `app/borrower/api/internal/handler/`
**Logic 目录**: `app/borrower/api/internal/logic/`（20 个模块）

### 模块和端点

| 模块 | 端点 | Handler |
|------|------|---------|
| **borrower_applicant** | all, check_existed, create, delete, detail, page, update | `handler/borrower_applicant/` |
| **borrower_entity** | all, check_existed, create, delete, detail, page, update | `handler/borrower_entity/` |
| **borrower_property** | all, borrower_property_transactions, create, delete, detail, detail_summary, page, update | `handler/borrower_property/` |
| **borrower_property_transaction** | [transaction] all / create / delete / detail / page / switch_status / update / update_checklist；[document] document_create / document_delete_file / document_upload_file；[operation] operation_list / operation_stats；[parse] parse_file / parse_file_result | `handler/borrower_property_transaction/` |
| **common** | get_us_cities, get_us_states, get_uuid, select_options, sync_api, sync_data, sync_resource | `handler/common/` |
| **document** | document_files | `handler/document/` |
| **drive** | drive_auth_status, drive_auth_url, drive_authorize, drive_disconnect, drive_get_sync_transaction_status, drive_list_drives, drive_sync_transaction, drive_update_folder, drive_update_selection | `handler/drive/` |
| **file** | delete | `handler/file/` |
| **loan** | apply, create, create_draft, dashboard, delete_loan, detail, download_import_template, import, lender_companies, page, pre_approval, send_pre_approval_loan_notification, send_quote_email, update, update_apply, update_draft, update_pending_status, update_stage | `handler/loan/` |
| **loan_applicant** | borrower_names | `handler/loan_applicant/` |
| **loan_document** | add_files | `handler/loan_document/` |
| **organization** | organization_lender_page, organization_lenders, organization_profile | `handler/organization/` |
| **organization_email** | authorize, generate_auth_code_url, save_email | `handler/organization_email/` |
| **organization_user** | accept_invitation, check_invitation, all, assign_roles, create, delete, detail, page, send_email, update, update_status | `handler/organization_user/` |
| **pdf** | apply_form_pdf | `handler/pdf/` |
| **resource** | resource_tree | `handler/resource/` |
| **role** | all, assign_organization_users, create, delete, detail, page, update, update_status | `handler/role/` |
| **storage** | create_aws_s3_presign_url | `handler/storage/` |
| **user** | change_password, edit_profile, info, login, refresh_token, register, reset_password, send_reset_password, send_verify_email, switch_organization, user_menu, verify_email | `handler/user/` |
| **zeitro** | zeitro_listing_search | `handler/zeitro/` |

**额外路由**: `handler/extra_routes.go`

---

## Lender API (Port 10100)

**API 定义**: `app/lender/api/desc/`
**Handler 目录**: `app/lender/api/internal/handler/`
**Logic 目录**: `app/lender/api/internal/logic/`（13 个模块）

### 模块和端点

| 模块 | 端点 | Handler |
|------|------|---------|
| **common** | get_us_cities, get_us_states, select_options, sync_api, sync_data, sync_resource | `handler/common/` |
| **file** | delete | `handler/file/` |
| **loan** | detail, loan_documents_opts, page, review_pending_status, update, update_pending_status, update_stage | `handler/loan/` |
| **loan_document** | add_files, create_loan_document, delete_loan_document_request, send_loan_document_requests, update_comment | `handler/loan_document/` |
| **organization** | organization_create, organization_profile, organization_update_profile | `handler/organization/` |
| **organization_setting** | detail_by_purpose, switch_status | `handler/organization_setting/` |
| **organization_template** | detail, detail_by_purpose, page, preview, reset, update | `handler/organization_template/` |
| **organization_user** | accept_invitation, all, assign_roles, check_invitation, create, delete, detail, page, send_email, switch_status, update | `handler/organization_user/` |
| **pdf** | html_to_pdf | `handler/pdf/` |
| **resource** | resource_tree | `handler/resource/` |
| **role** | all, assign_organization_users, create, delete, detail, page, switch_status, update | `handler/role/` |
| **storage** | create_aws_s3_presign_url | `handler/storage/` |
| **user** | info, login, refresh_token, register, reset_password, send_reset_password, send_verify_email, switch_organization, user_menu, verify_email | `handler/user/` |

---

## Admin API (Port 10200)

**API 定义**: `app/admin/api/desc/`
**Handler 目录**: `app/admin/api/internal/handler/`
**Logic 目录**: `app/admin/api/internal/logic/`（13 个模块）

### 模块和端点

| 模块 | 端点 | Handler |
|------|------|---------|
| **api** | create, delete, detail, tree, update | `handler/api/` |
| **base_template** | create, delete, detail, page, preview, reset, switch_status, update | `handler/base_template/` |
| **common** | get_us_cities, get_us_states, select_options, sync_api, sync_data, sync_resource | `handler/common/` |
| **loan** | lender_companies, loan_detail, loan_page, loan_reassign | `handler/loan/` |
| **organization** | create_lender, delete, detail, discover_commit, discover_get_settings, discover_latest_run, discover_preview, discover_run, discover_update_settings, page, select_option, switch_status, update_lender, update_order | `handler/organization/` |
| **organization_subscription** | create, detail, page, renew, switch_status, update_expire_date | `handler/organization_subscription/` |
| **organization_user** | create, delete, detail, page, switch_status, update | `handler/organization_user/` |
| **pricing_plan** | all, bind_resources, create, delete, detail, page, resource_list, switch_status, update | `handler/pricing_plan/` |
| **pricing_plan_policy** | create, delete, update | `handler/pricing_plan_policy/` |
| **resource** | create, delete, detail, tree, update | `handler/resource/` |
| **role** | all, assign_organization_users, create, delete, detail, page, switch_status, update | `handler/role/` |
| **storage** | create_aws_s3_presign_url | `handler/storage/` |
| **user** | edit_profile, info, login, refresh_token, reset_password, user_menu | `handler/user/` |

---

## gRPC Proto 定义

| 服务 | Proto 文件 | RPC 目录 |
|------|-----------|----------|
| Admin | `app/admin/rpc/pb/admin.proto` | `app/admin/rpc/internal/logic/` |
| Borrower | `app/borrower/rpc/pb/borrower.proto` | `app/borrower/rpc/internal/logic/` |
| Lender | `app/lender/rpc/pb/lender.proto` | `app/lender/rpc/internal/logic/` |

> 各服务 RPC 模块清单（含 API/RPC 对照）见 [INDEX](../INDEX.md) §服务功能矩阵；以 `app/<svc>/rpc/internal/logic/` 实际目录为准（borrower 21 / lender 14 / admin 14）。

---

## API 定义文件 (.api)

各服务的 HTTP 接口通过 `.api` 文件定义：

- **Borrower**: `app/borrower/api/desc/`（common.api、loan.api、user.api、property.api 等）
- **Lender**: `app/lender/api/desc/`
- **Admin**: `app/admin/api/desc/`

---

## 开发流程

新增端点的完整步骤（含代码生成命令、端点实现范例）见 [architecture](architecture.md) §开发流程。
