# Feature Specification: Lender Organization Management

**Feature Branch**: `003-lender-organization`
**Created**: 2025-01-09
**Status**: Draft
**Input**: User description: "在后端实现lender的组织修改功能开发,包含api与rpc服务"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 查看当前组织信息 (Priority: P1)

作为 lender 端用户，我需要能够查看我当前所在组织的详细信息，以便了解组织的基本情况和状态。

**Why this priority**: 这是基础功能，用户必须能够查看组织信息才能进行后续操作。没有这个功能，用户无法确认自己所在组织的状态。

**Independent Test**: 可以通过调用 API 获取当前用户的组织信息来独立测试。验证返回的组织信息包含所有必要字段（名称、状态、联系方式等）。

**Acceptance Scenarios**:

1. **Given** 用户已登录并属于某个组织，**When** 用户请求查看组织信息，**Then** 系统返回该组织的完整信息（包括名称、状态、所有者、联系方式、地址等）
2. **Given** 用户已登录但未加入任何组织，**When** 用户请求查看组织信息，**Then** 系统返回适当的错误提示
3. **Given** 用户登录已过期，**When** 用户请求查看组织信息，**Then** 系统返回认证错误

---

### User Story 2 - 更新组织信息 (Priority: P2)

作为组织的 owner，我需要能够更新组织的基本信息，以便保持组织资料的准确性和时效性。

**Why this priority**: 组织信息可能会发生变化（如联系方式、地址、logo 等），owner 需要能够更新这些信息。优先级低于查看功能，因为只有在查看后才能进行更新。

**Independent Test**: 可以由组织 owner 登录后独立测试更新功能。验证只有 owner 能够更新，其他用户无法更新。验证更新后的信息能够正确保存和显示。

**Acceptance Scenarios**:

1. **Given** 用户是组织的 owner，**When** 用户提交组织信息更新请求，**Then** 系统成功更新组织信息并返回更新后的数据
2. **Given** 用户不是组织的 owner，**When** 用户尝试更新组织信息，**Then** 系统返回权限不足错误
3. **Given** 用户是组织 owner 但提供的数据格式不正确，**When** 用户提交更新请求，**Then** 系统返回数据验证错误
4. **Given** 用户同时属于多个组织，**When** 用户更新组织信息，**Then** 系统只更新当前激活组织的对应信息

---

### User Story 3 - 创建新组织 (Priority: P3)

作为 lender 端用户，我需要能够创建新的组织，以便管理不同的业务实体或项目。

**Why this priority**: 这是扩展性功能，允许用户在需要时创建新的组织。优先级最低，因为用户可以先使用现有组织，后续再创建新组织。

**Independent Test**: 可以由已登录用户独立测试创建功能。验证创建成功后，用户自动成为新组织的 owner，组织状态为 Inactive，owner_type 为 ThirdParty。

**Acceptance Scenarios**:

1. **Given** 用户已登录，**When** 用户提交创建新组织的请求，**Then** 系统成功创建组织，用户成为 owner，组织状态为 Inactive，owner_type 为 ThirdParty
2. **Given** 用户已登录但提供的必填字段不完整，**When** 用户提交创建请求，**Then** 系统返回数据验证错误
3. **Given** 用户已登录，**When** 用户成功创建新组织后，**Then** 用户可以在组织列表中看到新创建的组织
4. **Given** 用户创建新组织，**When** 组织创建成功，**Then** 系统自动将用户与新组织关联，角色为 owner

---

### Edge Cases

- **并发更新**：多个 owner 同时尝试更新组织信息时如何处理？
- **组织状态限制**：当组织状态为 Inactive 时，是否允许更新信息？
- **多组织用户**：用户同时属于多个组织时，如何确定操作哪个组织？
- **组织删除**：删除组织时如何处理关联的用户和角色？
- **数据一致性**：更新组织信息时，如何确保关联数据的完整性（如 organization_email）？
- **权限检查**：用户声称自己是 owner 但实际不是时的安全处理
- **组织名称重复**：是否允许同一用户创建多个同名组织？
- **联系信息验证**：是否需要验证 email 和 phone 的真实性？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须能够从用户上下文（ctx）中获取当前用户信息
- **FR-002**: 系统必须能够根据用户信息查询其所在的组织
- **FR-003**: 系统必须能够返回组织的完整信息，包括但不限于：名称、状态、所有者、联系方式、地址、logo、网站、标语等
- **FR-004**: 系统必须在用户尝试查看组织信息时验证用户是否已加入组织
- **FR-005**: 系统必须验证请求更新组织信息的用户是否为该组织的 owner
- **FR-006**: 系统必须在更新组织信息前验证输入数据的有效性（email 格式、phone 长度、必填字段等）
- **FR-007**: 系统必须允许更新以下字段：contact_email、contact_phone、contact_name、logo、website、tagline、metadata（lender_metadata）
- **FR-008**: 系统必须在创建新组织时自动将创建者设置为组织的 owner
- **FR-009**: 系统必须在创建新组织时将组织状态设置为 Inactive
- **FR-010**: 系统必须在创建新组织时将 owner_type 设置为 ThirdParty
- **FR-011**: 系统必须在创建组织时验证必填字段：name、contact_email、contact_phone、contact_name
- **FR-012**: 系统必须在创建组织成功后自动创建用户与组织的关联关系
- **FR-013**: 系统必须提供清晰的错误消息，包括但不限于：权限不足、数据验证失败、组织不存在等
- **FR-014**: 系统必须确保所有操作都有适当的日志记录

### Key Entities

- **Organization（组织）**：
  - 基本属性：名称、状态、所有者ID
  - 联系信息：联系邮箱、联系电话、联系人名称
  - 品牌信息：logo、网站、标语
  - 元数据：lender 相关的元数据（贷款金额范围、期限范围、利率范围、积分范围）
  - 地址信息
  - 所有者类型：Self、ThirdParty、Partner

- **OrganizationUser（组织用户关联）**：
  - 用户与组织的多对多关系
  - 用户在组织中的状态
  - 用户在组织中的角色

- **User（用户）**：
  - 用户基本信息
  - 用户所属的组织列表

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 用户能够在 1 秒内成功获取当前组织的完整信息
- **SC-002**: 组织 owner 能够在 2 秒内成功更新组织信息
- **SC-003**: 用户能够在 3 秒内成功创建新组织并获得响应
- **SC-004**: 100% 的权限检查能够正确识别非 owner 用户并拒绝访问
- **SC-005**: 100% 的数据验证能够捕获格式错误的输入并返回清晰的错误消息
- **SC-006**: 创建的新组织 100% 符合规范：状态为 Inactive，owner_type 为 ThirdParty，创建者为 owner
- **SC-007**: 系统能够支持至少 1000 个并发用户查询组织信息而不出现性能下降
- **SC-008**: 所有组织信息更新操作的事务成功率达到 99.9% 以上

## Assumptions

1. **用户认证**：假设所有 API 请求都经过 JWT 认证，用户信息已包含在请求上下文中
2. **当前组织上下文**：假设用户在请求时会自动使用当前激活的组织 ID
3. **数据模型**：假设使用现有的 Organization 和 OrganizationUser 数据模型
4. **API 框架**：假设使用 go-zero 框架构建 API 和 RPC 服务
5. **错误处理**：假设使用 common/errorx 包进行统一的错误处理
6. **数据转换**：假设使用 copierx.Trans 进行数据类型转换
7. **组织唯一性**：假设同一用户可以创建多个组织，组织名称不需要全局唯一
8. **激活流程**：假设组织创建后为 Inactive 状态，需要通过其他流程激活
9. **权限模型**：假设只有 owner 可以更新组织基本信息，其他用户无权更新
10. **RPC 服务**：假设 lender RPC 服务已存在并包含必要的组织相关方法

## Clarifications

### Session 2025-01-09

- **Q: 如何在 RPC 服务中获取当前登录用户信息？** → **A**: 使用 `l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)` 方法从 gRPC metadata 中获取 JWT token 并解析用户信息（包含 UserId、OrganizationId、Username 等）。该方法会从 `metadata.FromIncomingContext(ctx)` 中提取 token，使用 `utils.ParseToken` 解析，返回 `utils.JwtUser` 结构体。
- **Q: 如何获取组织详情信息？** → **A**: 使用 `l.baseLogic.OrganizationLogic.Detail(ctx, organizationId, userType)` 方法，该方法会调用 `OrganizationRepo.GetOrganization` 获取组织数据，并包含预加载的 OrganizationEmails 关联数据，返回业务对象 `bo.OrganizationBo`。
- **Q: 如何进行数据类型转换？** → **A**: 使用 `copierx.Trans[SourceType, TargetType](source, copierx.TimeToRFC3339Str())` 进行类型转换，其中 `TimeToRFC3339Str()` 会将时间字段转换为 RFC3339 字符串格式。

## Dependencies

1. **现有 RPC 服务**：lender RPC 服务需要提供组织查询、更新和创建的底层实现
2. **用户上下文中间件**：API 服务需要能够从 JWT 中提取用户信息和当前组织信息
3. **数据访问层**：需要 Organization 和 OrganizationUser 的 repository 实现
4. **验证器**：需要 common/validator 包提供数据验证功能
5. **错误处理**：需要 common/errorx 包提供统一的错误处理
6. **UUID 生成器**：需要 common/uuid 包提供 UUID 生成功能
7. **数据库**：需要数据库表 organizations 和 organization_users 已正确创建
