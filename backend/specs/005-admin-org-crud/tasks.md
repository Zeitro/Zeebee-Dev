# Tasks: 管理后台组织管理 CRUD

**Input**: Design documents from `backend/specs/005-admin-org-crud/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/organization.api

**Tests**: 本功能规格未明确要求 TDD，测试任务可选

**Organization**: 任务按用户故事分组，支持独立实施和测试

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 任务所属用户故事（US1-US5）
- 包含精确文件路径

## Path Conventions

- **Zeebee Backend (Microservices)**: `app/admin/api/`, `app/admin/rpc/`, `domains/zeebee/`, `common/`
- **Services**: borrower, lender, admin
- **Shared layer**: `domains/zeebee/model/`, `domains/zeebee/repo/`
- **Utilities**: `common/errorx/`, `common/copierx/`, `common/validator/`, `common/gormx/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化和基础结构

- [ ] T001 确认当前分支为 `005-admin-org-crud`，从 `feat-backend-lender-register` 分支创建
- [ ] T002 确认 API 定义文件 `app/admin/api/desc/organization.api` 已存在（包含 CRUD 接口定义）
- [ ] T003 确认 RPC 定义文件已生成（goctl 自动生成）
- [ ] T004 确认 API Logic 模板文件已生成在 `app/admin/api/internal/logic/organization/`
- [ ] T005 确认 RPC Logic 模板文件已生成在 `app/admin/rpc/internal/logic/organization/`
- [ ] T006 [P] 配置 golangci-lint 代码检查工具（如未配置）

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有用户故事依赖的核心基础设施

**⚠️ CRITICAL**: 完成本阶段前，不能开始任何用户故事开发

- [ ] T007 在 `domains/zeebee/repo/organization_user_repo.go` 中添加 `UpdateStatusByOrganizationId` 方法用于批量更新组织用户状态
- [ ] T008 在 `domains/zeebee/repo/organization_user_repo.go` 中添加 `DeleteByOrganizationIdWithTx` 方法用于删除组织的所有用户关系
- [ ] T009 在 `domains/zeebee/repo/role_repo.go` 中添加 `DeleteByOrganizationIdWithTx` 方法用于删除组织的所有角色
- [ ] T010 在 `domains/zeebee/repo/user_repo.go` 中添加 `DeleteByOrganizationIdWithTx` 方法用于删除组织的所有用户

**Checkpoint**: 基础设施就绪 - 用户故事开发现在可以并行开始

---

## Phase 3: User Story 1 - 创建 Lender 组织 (Priority: P1) 🎯 MVP

**Goal**: 管理员可以创建新的 Lender 组织，组织自动激活并创建关联的用户

**Independent Test**: 管理员登录后台 → 填写组织信息 → 提交 → 验证组织状态为 Active 且关联用户已创建

### Implementation for User Story 1

- [ ] T011 [US1] 实现 `app/admin/rpc/internal/logic/organization/create_lender_logic.go` 的 `CreateLender` 方法，参考
  `user_logic.go:481-576` 的 `RegisterWithOrganization` 逻辑
- [ ] T012 [US1] 在创建逻辑中添加组织名称唯一性验证（Lender 类型）：调用 `OrganizationRepo.IsExistedByColumns` 检查
  `name` + `type`
- [ ] T013 [US1] 在创建逻辑中添加 owner 邮箱唯一性验证（Lender 类型）：调用 `UserRepo.IsExistedByColumns` 检查 `email` +
  `type`
- [ ] T014 [US1] 使用事务创建组织、用户和关联关系，状态全部设为 `Active`（不同于 `RegisterWithOrganization` 的激活流程）
- [ ] T015 [US1] 实现 `app/admin/api/internal/logic/organization/organization_create_lender_logic.go` 的
  `OrganizationCreateLender` 方法，使用 `copierx.Trans[T, U]` 进行类型转换
- [ ] T016 [US1] 添加错误处理使用 `common/errorx/` 包，使用 `errorx.WithStack` 包装错误
- [ ] T017 [US1] 添加日志记录使用 go-zero 的 `logx.WithContext`

**Checkpoint**: 此时，用户故事 1 应该完全功能且可独立测试

---

## Phase 4: User Story 2 - 查看组织详情 (Priority: P1)

**Goal**: 管理员可以查看组织的完整信息

**Independent Test**: 选择已存在的组织 → 查看详情 → 验证显示完整且准确

### Implementation for User Story 2

- [ ] T018 [US2] 实现 `app/admin/rpc/internal/logic/organization/detail_logic.go` 的 `Detail` 方法
- [ ] T019 [US2] 使用 `OrganizationLogic.GetOrganizationById` 获取组织详情，预加载关联数据（Owner, OrganizationEmails）
- [ ] T020 [US2] 实现 `app/admin/api/internal/logic/organization/organization_detail_logic.go` 的 `OrganizationDetail`
  方法，使用 `copierx.Trans[T, U]` 转换响应数据
- [ ] T021 [US2] 添加组织不存在时的错误处理，返回明确的错误提示

**Checkpoint**: 此时，用户故事 1 和 2 都应该独立工作

---

## Phase 5: User Story 3 - 编辑组织信息 (Priority: P2)

**Goal**: 管理员可以编辑组织的可编辑字段（组织名称不可修改）

**Independent Test**: 修改组织可编辑字段 → 保存 → 重新查看详情验证修改已生效

### Implementation for User Story 3

- [ ] T022 [US3] 实现 `app/admin/rpc/internal/logic/organization/update_lender_logic.go` 的 `UpdateLender` 方法
- [ ] T023 [US3] 确认只更新可编辑字段：ContactEmail, ContactPhone, ContactName, Logo, Website, Tagline, LenderMetadata
- [ ] T024 [US3] 组织名称、类型、owner_id、status 字段不允许编辑（API 定义中不包含这些字段）
- [ ] T025 [US3] 实现 `app/admin/api/internal/logic/organization/organization_update_lender_logic.go` 的
  `OrganizationUpdateLender` 方法
- [ ] T026 [US3] 添加组织不存在时的错误处理

**Checkpoint**: 此时，用户故事 1-3 都应该独立功能

---

## Phase 6: User Story 4 - 切换组织状态 (Priority: P2)

**Goal**: 管理员可以启用/禁用组织，系统自动同步所有组织成员的状态

**Independent Test**: 切换组织状态 → 验证组织状态变更且所有 organization_users 状态同步更新

### Implementation for User Story 4

- [ ] T027 [US4] 实现 `app/admin/rpc/internal/logic/organization/switch_status_logic.go` 的 `SwitchStatus` 方法
- [ ] T028 [US4] 使用事务更新组织状态和批量更新组织用户状态
- [ ] T029 [US4] 调用 `OrganizationUserRepo.UpdateStatusByOrganizationId` 批量更新状态（Single SQL）
- [ ] T030 [US4] 实现 `app/admin/api/internal/logic/organization/organization_switch_status_logic.go` 的
  `OrganizationSwitchStatus` 方法
- [ ] T031 [US4] 添加状态验证，只允许 Active 和 Disabled 之间的切换

**Checkpoint**: 此时，用户故事 1-4 都应该独立功能

---

## Phase 7: User Story 5 - 删除组织 (Priority: P3)

**Goal**: 管理员可以删除组织，系统级联删除关联数据

**Independent Test**: 删除组织 → 验证组织、用户、角色、关系全部被删除

### Implementation for User Story 5

- [ ] T032 [US5] 实现 `app/admin/rpc/internal/logic/organization/delete_logic.go` 的 `Delete` 方法
- [ ] T033 [US5] 使用事务按序删除：organization_users → roles → users → organization
- [ ] T034 [US5] 调用 `OrganizationUserRepo.DeleteByOrganizationIdWithTx` 删除用户关系
- [ ] T035 [US5] 调用 `RoleRepo.DeleteByOrganizationIdWithTx` 删除组织角色
- [ ] T036 [US5] 调用 `UserRepo.DeleteByOrganizationIdWithTx` 删除组织用户
- [ ] T037 [US5] 最后删除组织记录
- [ ] T038 [US5] 实现 `app/admin/api/internal/logic/organization/organization_delete_logic.go` 的 `OrganizationDelete`
  方法
- [ ] T039 [US5] 添加删除确认逻辑（可选，根据需求）

**Checkpoint**: 所有用户故事现在都应该独立功能

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: 影响多个用户故事的改进

- [ ] T040 [P] 运行 golangci-lint 代码检查，修复所有警告
- [ ] T041 [P] 运行 `go test ./...` 执行单元测试，确保测试通过
- [ ] T042 添加或更新 API 文档注释
- [ ] T043 验证所有接口使用 `copierx.Trans[T, U]` 进行类型转换
- [ ] T044 验证所有错误处理使用 `common/errorx/` 包
- [ ] T045 运行 `quickstart.md` 中的验证步骤

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成 - 阻塞所有用户故事
- **User Stories (Phase 3-7)**: 全部依赖 Foundational 完成
    - 用户故事完成后可并行执行（如果有足够人力）
    - 或按优先级顺序执行（P1 → P2 → P3）
- **Polish (Phase 8)**: 依赖所有用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational 完成后可开始 - 无其他故事依赖
- **User Story 2 (P1)**: Foundational 完成后可开始 - 无其他故事依赖
- **User Story 3 (P2)**: Foundational 完成后可开始 - 可与 US1/US2 集成但应独立测试
- **User Story 4 (P2)**: Foundational 完成后可开始 - 需要 US1 创建的组织
- **User Story 5 (P3)**: Foundational 完成后可开始 - 需要 US1 创建的组织

### Within Each User Story

- RPC Logic 实现 → API Logic 实现 → 集成
- 核心实现 → 错误处理 → 日志记录

### Parallel Opportunities

- 所有 Setup 任务标记 [P] 可并行执行
- 所有 Foundational 任务标记 [P] 可并行执行（在 Phase 2 内）
- Foundational 完成后，US1 和 US2 可并行开发（都是 P1）
- US3 和 US4 可并行开发（都是 P2）
- Polish 阶段所有标记 [P] 的任务可并行执行

---

## Parallel Example: User Story 1 & 2 (P1 Stories)

```bash
# Foundational 完成后，US1 和 US2 可并行开发：

# Developer A: User Story 1
Task T011: 实现 create_lender_logic.go
Task T015: 实现 organization_create_lender_logic.go

# Developer B: User Story 2 (并行)
Task T018: 实现 detail_logic.go
Task T020: 实现 organization_detail_logic.go
```

---

## Implementation Strategy

### MVP First (User Stories 1 & 2 Only)

1. 完成 Phase 1: Setup
2. 完成 Phase 2: Foundational (CRITICAL - 阻塞所有故事)
3. 完成 Phase 3: User Story 1
4. 完成 Phase 4: User Story 2
5. **STOP and VALIDATE**: 独立测试 US1 和 US2
6. 如准备就绪则部署/演示

### Incremental Delivery

1. 完成 Setup + Foundational → 基础就绪
2. 添加 User Story 1 → 独立测试 → 部署/演示（MVP 第 1 部分）
3. 添加 User Story 2 → 独立测试 → 部署/演示（MVP 完成）
4. 添加 User Story 3 → 独立测试 → 部署/演示
5. 添加 User Story 4 → 独立测试 → 部署/演示
6. 添加 User Story 5 → 独立测试 → 部署/演示
7. 每个故事都增加价值且不破坏之前的故事

### Parallel Team Strategy

有多名开发人员时：

1. 团队一起完成 Setup + Foundational
2. Foundational 完成后：
    - Developer A: User Story 1
    - Developer B: User Story 2
    - Developer C: User Story 3（如果 US1/US2 完成后）
3. 故事独立完成和集成

---

## Notes

- [P] 任务 = 不同文件，无依赖
- [Story] 标签将任务映射到特定用户故事以便追踪
- 每个用户故事应该独立完成和可测试
- 每个任务或逻辑组后提交代码
- 在任何 checkpoint 停止以独立验证故事
- 避免：模糊任务、同一文件冲突、破坏独立性的跨故事依赖

---

## Task Summary

- **Total Tasks**: 45
- **Setup Phase**: 6 tasks
- **Foundational Phase**: 4 tasks (BLOCKS all user stories)
- **User Story 1**: 7 tasks
- **User Story 2**: 4 tasks
- **User Story 3**: 5 tasks
- **User Story 4**: 5 tasks
- **User Story 5**: 8 tasks
- **Polish Phase**: 6 tasks

### Parallel Opportunities

- **Setup**: 1 parallel task (T006)
- **Foundational**: 4 parallel tasks (T007-T010)
- **User Stories**: US1 & US2 可并行（都是 P1），US3 & US4 可并行（都是 P2）
- **Polish**: 5 parallel tasks (T040-T044)

### Independent Test Criteria

- **US1**: 创建组织 → 验证 Active 状态 + 关联用户
- **US2**: 查看组织详情 → 验证完整信息显示
- **US3**: 编辑组织 → 验证可编辑字段已更新
- **US4**: 切换状态 → 验证组织 + 所有 users 状态同步
- **US5**: 删除组织 → 验证级联删除成功

### MVP Scope

**Recommended MVP**: User Stories 1 & 2 (P1 stories only)

- 完成组织创建和查看功能
- 提供核心价值：管理员可以管理基本的组织生命周期
- 在更高级功能（编辑、状态切换、删除）之前可独立部署
