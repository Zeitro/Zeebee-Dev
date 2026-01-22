# Feature Specification: 管理后台组织管理 CRUD

**Feature Branch**: `001-admin-org-crud`
**Created**: 2025-01-12
**Status**: Draft
**Input**: User description: "@.claude/docs/tasks/20260112_admin_organization_crud.md 实现文档中的后端功能开发"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 创建 Lender 组织 (Priority: P1)

系统管理员需要通过管理后台创建新的 Lender（贷款方）组织，以便让贷款方能够使用系统进行贷款业务管理。

**Why this priority**: 这是核心功能，没有组织创建能力，整个系统无法接入新的贷款方客户，直接影响业务扩展。

**Independent Test**: 可以通过管理员登录后台，成功创建一个贷款方组织并验证组织状态为 Active，同时验证关联的用户和关系记录自动创建完成。

**Acceptance Scenarios**:

1. **Given** 管理员已登录后台，**When** 填写完整的组织信息并提交创建，**Then** 系统创建组织记录，组织状态为
   Active，同时创建关联的默认用户和关系记录
2. **Given** 管理员提交的组织信息缺少必填字段，**When** 尝试创建组织，**Then** 系统返回明确的错误提示，要求补充缺失信息
3. **Given** 组织名称已存在，**When** 尝试创建同名组织，**Then** 系统返回组织已存在的错误提示

---

### User Story 2 - 查看组织详情 (Priority: P1)

系统管理员需要查看组织的完整信息，包括组织基本信息、关联用户、状态等，以便进行审核和管理。

**Why this priority**: 管理员必须能够查看组织详情才能进行后续的编辑、状态管理或删除操作，这是管理流程的基础。

**Independent Test**: 可以通过选择已存在的组织并查看详情页面，验证显示的组织信息完整且准确。

**Acceptance Scenarios**:

1. **Given** 系统中存在组织记录，**When** 管理员选择查看组织详情，**Then** 系统展示组织的完整信息
2. **Given** 组织不存在或已被删除，**When** 尝试查看组织详情，**Then** 系统返回组织不存在的错误提示

---

### User Story 3 - 编辑组织信息 (Priority: P2)

系统管理员需要修改组织的可编辑信息（如名称、联系方式等），以便保持组织信息的准确性。

**Why this priority**: 组织信息会随时间变化，管理员需要能够更新信息，但优先级低于创建和查看功能，因为组织可以先用初始信息创建。

**Independent Test**: 可以通过修改组织信息并保存，然后重新查看详情验证修改已生效。

**Acceptance Scenarios**:

1. **Given** 组织已存在，**When** 管理员修改允许编辑的字段并保存，**Then** 系统更新组织信息并保存成功
2. **Given** 管理员尝试修改不允许编辑的字段，**When** 提交修改，**Then** 系统忽略不允许修改的字段或返回错误提示

---

### User Story 4 - 切换组织状态 (Priority: P2)

系统管理员需要启用或禁用组织，以便控制组织的访问权限，同时系统需要自动同步更新组织成员的状态。

**Why this priority**: 状态管理是运营控制的重要功能，但可以在组织创建和查看之后实现，优先级低于基础 CRUD。

**Independent Test**: 可以通过切换组织状态，验证组织状态变更的同时，组织成员（organization_user）的状态也同步更新。

**Acceptance Scenarios**:

1. **Given** 组织状态为 Disabled，**When** 管理员将状态切换为 Active，**Then** 组织状态更新为 Active，同时该组织的所有
   organization_user 记录状态从 Disabled 更新为 Active
2. **Given** 组织状态为 Active，**When** 管理员将状态切换为 Disabled，**Then** 组织状态更新为 Disabled，同时该组织的所有
   organization_user 记录状态从 Active 更新为 Disabled
3. **Given** 组织下没有成员，**When** 管理员切换组织状态，**Then** 仅更新组织状态，不影响不存在的成员记录

---

### User Story 5 - 删除组织 (Priority: P3)

系统管理员需要删除不再需要的组织，以便清理系统数据。删除前需要先删除组织下的用户、角色及关联关系。

**Why this priority**: 删除是数据清理功能，重要性低于创建、查看、编辑和状态管理。优先级最低是因为可以通过禁用来替代删除。

**Independent Test**: 可以通过删除一个测试组织，验证组织及其关联的用户、角色和关系记录被正确删除。

**Acceptance Scenarios**:

1. **Given** 组织存在且没有关联数据，**When** 管理员执行删除操作，**Then** 系统删除组织记录并返回成功
2. **Given** 组织存在但有关联的用户、角色等数据，**When** 管理员执行删除操作，**Then**
   系统先删除关联数据（用户、角色、关系），然后删除组织，整个操作在事务中完成
3. **Given** 组织正在被使用或有重要数据，**When** 管理员执行删除操作，**Then** 系统显示删除确认提示，警告删除后数据无法恢复

---

### Edge Cases

- 当组织删除失败时（如数据库约束错误），系统如何处理？
- 当组织状态切换时，如果有用户的个人状态已经是 Disabled，应该如何处理？
- 当并发编辑同一组织时，如何保证数据一致性？
- 当组织下有大量用户（如 1000+）时，状态切换操作的性能如何保证？
- 当组织删除操作被中断时（如系统故障），如何保证数据完整性？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须允许管理员创建新的 Lender 组织，组织创建后状态自动为 Active
- **FR-002**: 系统必须在创建组织时自动创建默认管理员用户，用户状态为 Active
- **FR-003**: 系统必须在创建组织时自动创建 organization_user 关联记录，状态为 Active
- **FR-004**: 系统必须允许管理员查看组织的完整信息详情
- **FR-005**: 系统必须允许管理员编辑组织的可编辑字段（仅 organizations 表字段）
- **FR-006**: 系统必须在管理员切换组织状态为 Active 时，自动将该组织所有 organization_user 记录状态从 Disabled 更新为
  Active
- **FR-007**: 系统必须在管理员切换组织状态为 Disabled 时，自动将该组织所有 organization_user 记录状态从 Active 更新为
  Disabled
- **FR-008**: 系统必须在删除组织前先删除该组织下的所有用户
- **FR-009**: 系统必须在删除组织前先删除该组织下的所有角色
- **FR-010**: 系统必须在删除组织前先删除该组织的所有关联关系（organization_user、组织角色关联等）
- **FR-011**: 系统必须确保删除操作在单个事务中完成，保证数据一致性
- **FR-012**: 系统必须对组织名称进行唯一性验证
- **FR-013**: 系统必须对必填字段进行验证，缺失时返回明确错误提示

### Key Entities

- **Organization（组织）**: 代表一个 Lender（贷款方）实体，包含组织基本信息、状态等属性
- **User（用户）**: 组织下的用户账号，与组织通过 organization_user 关联
- **OrganizationUser（组织用户关系）**: 用户与组织的多对多关联关系，包含用户在组织中的状态
- **Role（角色）**: 组织下的角色定义，用于权限管理

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 管理员可以在 30 秒内完成一个新组织的创建流程
- **SC-002**: 组织详情查询响应时间在 1 秒以内
- **SC-003**: 组织状态切换操作在 3 秒内完成，包括所有关联用户状态的更新
- **SC-004**: 组织删除操作（含关联数据）在 5 秒内完成
- **SC-005**: 所有组织操作的成功率达到 99% 以上（排除数据验证错误）
- **SC-006**: 组织名称重复验证准确率达到 100%
- **SC-007**: 状态切换操作确保组织成员状态同步准确率达到 100%

## Assumptions

- 管理员已通过身份认证并有权限访问组织管理功能
- 系统已有基本的用户认证和授权基础设施
- 数据库表结构（organizations、users、organization_user、roles 等）已存在
- 参考代码 `organization_page_logic.go` 和 `page_logic.go` 的模式和风格
- 参考代码 `RegisterWithOrganization` 的创建逻辑模式
- 模板代码已在 `app/admin/api/internal/logic/organization/` 和 `app/admin/rpc/internal/logic/organization/` 生成

## Dependencies

- 前端组织管理界面（用于调用后端 API）
- 数据库表和索引已正确配置
- 常量和枚举定义（如组织状态、用户状态等）已存在
- 共享领域逻辑 `domains/zeebee/` 中有相关的数据访问方法
