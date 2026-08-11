# Tenant 模块迁移评审（v3 — codex 二次修改后）

分支 `feat/backend/tenant` vs `master`，working tree 含 codex 未提交修改。v3 基于当前文件实态。

**原则**：tenant 独占 13 模块（user/auth/organization/organizationuser/organizationemail/organizationtemplate/organizationsetting/role/pricingplan/pricingplanpolicy/organizationsubscription/api/resource），edge 不得直连其 repo/logic/auth，为连库拆分服务。命名保留 tenant。

---

## 一、本 round codex 修复确认

| v2 项 | 状态 | 证据 |
|------|------|------|
| #6 lender base_logic 直写 `users.default_organization_id` | ✅ | grep 无 `default_organization_id/UserRepo/UpdateColumnsById`，已移除 |
| #7 borrower register 本地铸 JWT | ✅ | `register_logic.go:95` 改 `GetUserClaimFromCtx(ctx)`；无 `GenerateJwtToken/UserRepo/OrganizationUserRepo` |
| #9 lender sync_data 直连 org/orgsetting | ✅ | 改调 `TenantClient.OrganizationSetting.SyncLenderDefaults`（纯 proxy） |
| #11 borrower loan 7 处 OrganizationRepo 读 | ✅ | `send_quote_email` 改 `TenantClient.Organization.DetailWithType`（:57,:130）；其他 loan 文件 grep 空 |
| #13 borrower verify_email 冗余读 | ✅ | 无 `UserRepo/GetUserByEmail`，已删 |
| #14 admin/lender base_logic orguser 读（operator 名） | ✅ | 改 `TenantClient.User.GetClaim`（admin:83 / lender:79） |
| #3 部分：orgemail 授权数据 | ✅ 半 | tenant 加 `OrganizationEmail.GetAuthorization` RPC；borrower `send_quote_email:181` 走 RPC 取授权，不再直读 `EmailAuthorizationRepo` |

**新增 tenant RPC**：`OrganizationEmail.GetAuthorization`、`User.GetClaim`、`OrganizationSetting.SyncLenderDefaults`。

---

## 二、仍未修（边界违规，按优先）

### 🔴 结构性根因（最高，未动）

1. **edge `svc/service_context.go` 仍嵌 `*repo.AllRepos`** — admin(:17,:35)、lender(:16,:27)、borrower(:24,:81)。给 edge 全部 13 tenant 表。
2. **edge `base_logic.go` 仍嵌 `*logic.AllLogics` + JwtConfig** — admin(:33,:45,:61,:72)、lender(:31,:43,:56,:68)、borrower(:27,:39,:52,:65)。给 edge 全部 tenant 逻辑 + 本地铸 token 能力。

不拆这俩，per-file 代理是表面功夫，DB 拆即崩。

### 🔴 多表事务写（admin discover，未动）

3. **`app/admin/rpc/internal/schedule/lender_discover_scheduler.go:31-39` NewRunner** 仍 `Repos: svcCtx.AllRepos, TemplateLogic: bl.OrganizationTemplateLogic, SettingLogic: bl.OrganizationSettingLogic`。`domains/zeebee/lenderdiscover/runner.go` 在单事务内写 Organization + OrganizationTemplate + OrganizationSetting。**单点拆分阻断，最难迁**（LLM discover 管线 admin 拥有，持久化却落 tenant 表）。

### 🔴 BaseTemplate 级联写（admin，未动）

4. **`app/admin/rpc/internal/logic/basetemplate/{create,update,delete,switch_status}_logic.go`** 仍调 `baseLogic.BaseTemplateLogic.*`，domain 逻辑级联写 `organization_template` + `organization`（CascadeStatus/CascadeDelete/BatchCreate）。仅 `reset_logic.go` 走 `TenantClient.OrganizationTemplate.ResetByBaseTemplateId`。tenant 缺 4 个级联 RPC。

### 🟠 orgemail 派发端（borrower，半修）

5. **授权数据已中心化（GetAuthorization RPC），但 OAuth client 缓存 + 发送仍在 borrower**：`loan/send_quote_email_logic.go:146,193,213` 用 `global.EmailAuthManager.BuildKey/GetClient/Register`；`:168 input.emailClient.SendEmail(ctx, msg)` 本地发。`global/email.go` 本地 EmailAuthManager + Gmail OAuth 配置仍在。**send 路径仍阻断 orgemail 拆分**。tenant 缺 `SendEmail`（或 `GetAuthorizedClient`）RPC。

### 🟠 admin discover 设置（未动）

6. **`organization/discover_update_settings_logic.go:35,46`** 仍 `lenderdiscover.LoadConfig/SaveConfig(l.ctx, l.svcCtx.OrganizationSettingRepo, ...)` 直读+直写 organization_setting。`discover_get_settings` + `admin.go` cron 同。tenant 缺 `OrganizationSetting.GetLenderDiscoverConfig`/`SaveLenderDiscoverConfig`。

### 🟡 lender orgtemplate preview（未动）

7. **`lender/rpc/internal/logic/organizationtemplate/preview_logic.go:43`** 仍 `baseLogic.BaseTemplateLogic.Preview(l.ctx, *input)` 直连 domain。tenant 缺 `OrganizationTemplate.Preview` RPC。**lender 唯一残留直连 domain 的 feature 文件**。

---

## 三、tenant 缺失 RPC（阻断完全收口，按优先）

1. `OrganizationEmail.SendEmail`（或 `GetAuthorizedClient`）— 收口 borrower send 路径，移 `global/email.go`。**orgemail 真拆分唯一剩阻断**。
2. `Organization.CreateLenderFromDiscover` bundle（单事务 org+template+setting）— admin scheduler。
3. BaseTemplate 级联 RPC ×4（`CascadeCreate/Update/Delete/SwitchStatusFromBaseTemplate`）— admin basetemplate。
4. `OrganizationSetting.GetLenderDiscoverConfig` + `SaveLenderDiscoverConfig` — admin discover 设置 + cron。
5. `OrganizationTemplate.Preview` — lender preview。

---

## 四、干净确认（无需改）

- 13 tenant 模块 × 3 edge 的 **handler 逻辑层全代理** `TenantClient.*`。
- 3 个 **API 层零 domain 导入**。
- edge **本地 casbin/enforcer 全删**，走 `tenant.AuthCheck`。
- **无本地铸 token**（borrower register 改 GetUserClaimFromCtx；accept_invitation 用 tenant token）。
- v1 P0 全修（3 Org RPC、authcode、品牌 config、policy 审计、孤儿 config）。
- lender `sync_data` 纯 proxy；borrower loan 的 org/orgemail **读**已切 tenant RPC；borrower register/verify 直连已删。

---

## 五、优先级清单（v3）

| # | 级别 | 项 | 文件 |
|---|------|-----|------|
| 1 | 🔴 结构 | 拆 edge `svc` 的 `*repo.AllRepos`（3 服务） | `app/{admin,lender,borrower}/rpc/internal/svc/service_context.go` |
| 2 | 🔴 结构 | 拆 edge `base_logic` 的 `*logic.AllLogics`+JwtConfig（3 服务） | `app/{admin,lender,borrower}/rpc/internal/logic/base_logic.go` |
| 3 | 🔴 阻断 | tenant 加 `OrganizationEmail.SendEmail`/`GetAuthorizedClient`；移 borrower `global/email.go` + `send_quote_email` 本地 EmailAuthManager/SendEmail | `app/borrower/rpc/internal/global/email.go`、`loan/send_quote_email_logic.go:146-213` |
| 4 | 🔴 阻断 | tenant 加 discover bundle RPC；admin scheduler `NewRunner` 改走 tenant（不再传 AllRepos/TemplateLogic/SettingLogic） | `app/admin/rpc/internal/schedule/lender_discover_scheduler.go:31-39`、`domains/zeebee/lenderdiscover/runner.go` |
| 5 | 🔴 高 | tenant 加 BaseTemplate 级联 RPC ×4；admin basetemplate create/update/delete/switch 改走 tenant | `app/admin/rpc/internal/logic/basetemplate/{create,update,delete,switch_status}_logic.go` |
| 6 | 🟠 高 | tenant 加 discover config RPC；admin discover_update/get_settings + admin.go cron 改走 tenant | `app/admin/rpc/internal/logic/organization/discover_*_settings_logic.go`、`app/admin/rpc/admin.go` |
| 7 | 🟡 中 | tenant 加 `OrganizationTemplate.Preview`；lender preview 改走 tenant | `app/lender/rpc/internal/logic/organizationtemplate/preview_logic.go:43` |
| — | 运维 | 播种 3 平台 DB 鉴权白名单；确认 prod `Tenant.Endpoints`；subscription cron 策略 | — |

---

## 六、判定

本 round codex **清掉了所有「读」类违规 + borrower register 铸 token + lender user 直写**（v2 的 #6/#7/#9/#11/#13/#14 + #3 半）。handler/读路径边界基本干净。

**剩 3 类硬骨头未动**：
1. 结构性 `AllRepos`/`AllLogics` 嵌入（#1/#2）——所有直连的根，最后拆。
2. admin 侧 3 处「写 tenant 表」业务路径：discover scheduler 多表事务（#4）、basetemplate 级联（#5）、discover settings（#6）——需 tenant 补 bundle/级联 RPC。
3. borrower email **发送端**本地化（#3）——授权数据已中心化，但 EmailAuthManager + SendEmail 仍在 borrower，需 tenant `SendEmail` RPC。

**结论**：读路径可拆，**写路径未可拆**。按 #3→#4→#5→#6 补 tenant RPC 并切 proxy，最后拆 #1/#2 结构嵌入。
