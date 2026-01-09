# Implementation Tasks: Lender Role and Organization User Management

**Feature Branch**: `001-lender-role-org`
**Date**: 2025-12-29
**Total Tasks**: 86
**Status**: ✅ Completed

## Overview

This task breakdown implements role and organization user management features for the lender service, mirroring the
existing borrower service implementation. Tasks are organized by user story to enable independent implementation and
testing.

## Implementation Strategy

**MVP First**: Start with User Story 1 (Role Management) - P1 priority provides core functionality
**Incremental Delivery**: Each user story is independently testable and deployable
**Parallel Execution**: Tasks marked [P] can be executed in parallel with other [P] tasks

---

## Phase 1: Setup

**Goal**: Initialize project structure and copy API definitions

- [x] T001 Copy resource.api from contracts to app/lender/api/desc/resource.api
- [x] T002 Copy role.api from contracts to app/lender/api/desc/role.api
- [x] T003 Copy organization_user.api from contracts to app/lender/api/desc/organization_user.api

---

## Phase 2: Foundational

**Goal**: Update main API file and generate base code structure

### Update lender.api

- [x] T004 Add resource.api, role.api, organization_user.api imports to app/lender/api/lender.api
- [x] T005 Add role routes block to app/lender/api/lender.api (6 endpoints with AuthMiddleware)
- [x] T006 Add organization_user routes block to app/lender/api/lender.api (1 public + 9 protected endpoints)

### Generate Base Code

- [x] T007 Run goctl api go generation in app/lender/api directory
- [x] T008 Verify generated handler directories exist: app/lender/api/internal/handler/role/
- [x] T009 Verify generated handler directories exist: app/lender/api/internal/handler/organization_user/
- [x] T010 Verify generated logic directories exist: app/lender/api/internal/logic/role/
- [x] T011 Verify generated logic directories exist: app/lender/api/internal/logic/organization_user/

---

## Phase 3: User Story 1 - Role Management (P1)

**Goal**: Enable administrators to create, read, update, delete, and manage roles

**Independent Test**: Create a role via API, update it, assign permissions, and verify deletion prevention when in use

### Acceptance Criteria

- Can create roles with name, description, status, and resource permissions
- Cannot create role with reserved name "admin"
- Can update role details and permissions
- Cannot delete roles assigned to users
- Can deactivate roles without deletion
- Can view paginated, filterable, sortable role list

### Implementation Tasks

#### RPC Layer - Role Logic

- [x] T012 [P] [US1] Implement role_create_logic.go RPC in app/lender/rpc/internal/logic/role/create_logic.go
- [x] T013 [P] [US1] Implement role_delete_logic.go RPC in app/lender/rpc/internal/logic/role/delete_logic.go
- [x] T014 [P] [US1] Implement role_detail_logic.go RPC in app/lender/rpc/internal/logic/role/detail_logic.go
- [x] T015 [P] [US1] Implement role_page_logic.go RPC in app/lender/rpc/internal/logic/role/page_logic.go
- [x] T016 [P] [US1] Implement role_update_logic.go RPC in app/lender/rpc/internal/logic/role/update_logic.go
- [x] T017 [P] [US1] Implement role_update_status_logic.go RPC in
  app/lender/rpc/internal/logic/role/update_status_logic.go
- [x] T018 [P] [US1] Implement role_all_logic.go RPC in app/lender/rpc/internal/logic/role/all_logic.go

#### RPC Layer - Role Server

- [x] T019 [P] [US1] Implement role server in app/lender/rpc/internal/server/role/role_server.go (7 methods)

#### API Layer - Role Logic

- [x] T020 [P] [US1] Implement role_create_logic.go API in app/lender/api/internal/logic/role/role_create_logic.go
- [x] T021 [P] [US1] Implement role_delete_logic.go API in app/lender/api/internal/logic/role/role_delete_logic.go
- [x] T022 [P] [US1] Implement role_detail_logic.go API in app/lender/api/internal/logic/role/role_detail_logic.go
- [x] T023 [P] [US1] Implement role_page_logic.go API in app/lender/api/internal/logic/role/role_page_logic.go
- [x] T024 [P] [US1] Implement role_update_logic.go API in app/lender/api/internal/logic/role/role_update_logic.go
- [x] T025 [P] [US1] Implement role_update_status_logic.go API in
  app/lender/api/internal/logic/role/role_update_status_logic.go
- [x] T026 [P] [US1] Implement role_all_logic.go API in app/lender/api/internal/logic/role/role_all_logic.go

#### Service Context Update

- [x] T027 [US1] Add RoleClient to ServiceContext in app/lender/api/internal/svc/service_context.go
- [x] T028 [US1] Initialize RoleClient in ServiceContext with RPC connection

### Test Scenarios for US1

```bash
# 1. Create role
POST /api/v1/role
{"name":"Loan Officer","description":"Can review loans","status":"Active","resource_ids":[1,2]}

# 2. Verify reserved name rejection
POST /api/v1/role
{"name":"admin",...} # Should return 400 error

# 3. Update role
PUT /api/v1/role/1
{"name":"Senior Loan Officer",...}

# 4. Attempt delete with users assigned (should fail)
DELETE /api/v1/role/1 # Should return error if users assigned

# 5. Deactivate instead of delete
PATCH /api/v1/role/1/status
{"status":"Inactive"}

# 6. List roles with filters
POST /api/v1/role/page
{"filter":{"keyword":"Loan"},"sort":{"column":"name","sort":"asc"}}
```

---

## Phase 4: User Story 2 - Organization User Management (P2)

**Goal**: Enable administrators to manage users within their organization

**Independent Test**: Add users via email, view paginated user list, filter/sort users, update user info, disable user
accounts

**Depends on**: Phase 3 (Role Management must exist first for role assignment)

### Acceptance Criteria

- Can add users by email with optional role assignment
- User is created in "Pending" status
- Invitation email is sent
- Can view paginated user list
- Can filter by keyword, status, role
- Can sort by name, email, status, date
- Can update user name, phone, remark
- Can disable/enable user accounts

### Implementation Tasks

#### RPC Layer - OrganizationUser Logic

- [x] T029 [P] [US2] Implement organization_user_create_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/create_logic.go
- [x] T030 [P] [US2] Implement organization_user_delete_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/delete_logic.go
- [x] T031 [P] [US2] Implement organization_user_detail_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/detail_logic.go
- [x] T032 [P] [US2] Implement organization_user_page_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/page_logic.go
- [x] T033 [P] [US2] Implement organization_user_update_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/update_logic.go
- [x] T034 [P] [US2] Implement organization_user_switch_status_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/switch_status_logic.go
- [x] T035 [P] [US2] Implement organization_user_all_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/all_logic.go

#### RPC Layer - OrganizationUser Server

- [x] T036 [P] [US2] Implement organizationuser server in
  app/lender/rpc/internal/server/organizationuser/organization_user_server.go (7 methods)

#### API Layer - OrganizationUser Logic

- [x] T037 [P] [US2] Implement organization_user_create_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_create_logic.go
- [x] T038 [P] [US2] Implement organization_user_delete_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_delete_logic.go
- [x] T039 [P] [US2] Implement organization_user_detail_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_detail_logic.go
- [x] T040 [P] [US2] Implement organization_user_page_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_page_logic.go
- [x] T041 [P] [US2] Implement organization_user_update_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_update_logic.go
- [x] T042 [P] [US2] Implement organization_user_switch_status_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_switch_status_logic.go
- [x] T043 [P] [US2] Implement organization_user_all_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_all_logic.go

#### Service Context Update

- [x] T044 [US2] Add OrganizationUserClient to ServiceContext in app/lender/api/internal/svc/service_context.go
- [x] T045 [US2] Initialize OrganizationUserClient in ServiceContext with RPC connection

### Test Scenarios for US2

```bash
# 1. Create user with invitation
POST /api/v1/organization_user
{"email":"newuser@example.com","role_ids":[1]}

# 2. List users with pagination
POST /api/v1/organization_user/page
{"page":1,"limit":10}

# 3. Filter users
POST /api/v1/organization_user/page
{"filter":{"keyword":"john","status":"Active","role_id":1}}

# 4. Update user
PUT /api/v1/organization_user/1
{"first_name":"John","last_name":"Doe","phone":"555-1234"}

# 5. Disable user
PATCH /api/v1/organization_user/1/status
{"status":"Disabled"}
```

---

## Phase 5: User Story 3 - Invite and Accept Organization Membership (P3)

**Goal**: Enable secure invitation flow for new team members

**Independent Test**: Send invitation, receive email, accept invitation with code, create account

**Depends on**: Phase 4 (Organization User Management must exist)

### Acceptance Criteria

- Invitation code is generated (UUID-based)
- Invitation email is sent with code
- User can accept invitation via public endpoint
- User provides first name, last name, email during acceptance
- Account is activated after acceptance
- Invitation cannot be reused
- Invitation expires after 7 days

### Implementation Tasks

#### RPC Layer - Invitation Logic

- [x] T046 [P] [US3] Implement organization_user_send_invitation_email_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/send_invitation_email_logic.go
- [x] T047 [P] [US3] Implement organization_user_accept_invitation_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/accept_invitation_logic.go

#### Email Templates

- [x] T048 [P] [US3] Create user_invite_email.html template in app/lender/rpc/internal/email_tpl/user_invite_email.html

#### API Layer - Invitation Logic

- [x] T049 [P] [US3] Implement organization_user_send_email_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_send_email_logic.go
- [x] T050 [P] [US3] Implement organization_user_accept_invitation_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_accept_invitation_logic.go

### Test Scenarios for US3

```bash
# 1. Send invitation (resend)
POST /api/v1/organization_user/1/send_email

# 2. Accept invitation (public endpoint)
POST /api/v1/organization_user/accept_invitation
{"code":"abc-123-def","email":"newuser@example.com","organization_id":1,"first_name":"John","last_name":"Doe"}

# 3. Attempt to reuse invitation (should fail)
POST /api/v1/organization_user/accept_invitation
{"code":"abc-123-def",...} # Should return error

# 4. Attempt expired invitation (should fail)
POST /api/v1/organization_user/accept_invitation
{"code":"expired-code",...} # Should return error
```

---

## Phase 6: User Story 4 - Role Assignment and Permission Updates (P4)

**Goal**: Enable administrators to assign and update roles for organization users

**Independent Test**: Assign multiple roles to user, remove roles, verify permission changes

**Depends on**: Phase 3 (Role Management) and Phase 4 (Organization User Management)

### Acceptance Criteria

- Can assign multiple roles to a user
- Can remove role assignments
- User permissions reflect all assigned roles
- Can view user's assigned roles

### Implementation Tasks

#### RPC Layer - Role Assignment Logic

- [x] T051 [P] [US4] Implement organization_user_assign_roles_logic.go RPC in
  app/lender/rpc/internal/logic/organizationuser/assign_roles_logic.go

#### API Layer - Role Assignment Logic

- [x] T052 [P] [US4] Implement organization_user_assign_roles_logic.go API in
  app/lender/api/internal/logic/organization_user/organization_user_assign_roles_logic.go

### Test Scenarios for US4

```bash
# 1. Assign roles to user
POST /api/v1/organization_user/1/assign_roles
{"role_ids":[1,2,3]}

# 2. Remove roles (assign empty list)
POST /api/v1/organization_user/1/assign_roles
{"role_ids":[1]}

# 3. Verify permissions reflect assigned roles
GET /api/v1/organization_user/1 # Check roles array
```

---

## Phase 7: Polish & Cross-Cutting Concerns

**Goal**: Complete implementation, error handling, and validation

### Error Handling & Validation

- [x] T053 [P] Add error handling for duplicate role names in app/lender/rpc/internal/logic/role/create_logic.go
- [x] T054 [P] Add error handling for reserved role name "admin" in app/lender/rpc/internal/logic/role/create_logic.go
- [x] T055 [P] Add error handling for role deletion when assigned to users in
  app/lender/rpc/internal/logic/role/delete_logic.go
- [x] T056 [P] Add error handling for invalid invitation codes in
  app/lender/rpc/internal/logic/organizationuser/accept_invitation_logic.go
- [x] T057 [P] Add error handling for expired invitations in
  app/lender/rpc/internal/logic/organizationuser/accept_invitation_logic.go
- [x] T058 [P] Add error handling for reused invitation codes in
  app/lender/rpc/internal/logic/organizationuser/accept_invitation_logic.go
- [x] T059 [P] Add error handling for duplicate user email in organization in
  app/lender/rpc/internal/logic/organizationuser/create_logic.go
- [x] T060 [P] Add error handling for non-existent role assignment in
  app/lender/rpc/internal/logic/organizationuser/assign_roles_logic.go

### Logging & Audit

- [x] T061 [P] Add logging for role creation in app/lender/rpc/internal/logic/role/create_logic.go
- [x] T062 [P] Add logging for role deletion attempts in app/lender/rpc/internal/logic/role/delete_logic.go
- [x] T063 [P] Add logging for role permission changes in app/lender/rpc/internal/logic/role/update_logic.go
- [x] T064 [P] Add logging for user invitation sent in
  app/lender/rpc/internal/logic/organizationuser/send_invitation_email_logic.go
- [x] T065 [P] Add logging for invitation acceptance in
  app/lender/rpc/internal/logic/organizationuser/accept_invitation_logic.go
- [x] T066 [P] Add logging for role assignment changes in
  app/lender/rpc/internal/logic/organizationuser/assign_roles_logic.go
- [x] T067 [P] Add logging for user status changes in
  app/lender/rpc/internal/logic/organizationuser/switch_status_logic.go

### Data Validation

- [x] T068 [P] Add role name validation (2-100 chars) in app/lender/api/internal/logic/role/role_create_logic.go
- [x] T069 [P] Add role status validation (Active/Inactive) in app/lender/api/internal/logic/role/role_create_logic.go
- [x] T070 [P] Add email format validation in
  app/lender/api/internal/logic/organization_user/organization_user_create_logic.go
- [x] T071 [P] Add phone number validation in
  app/lender/api/internal/logic/organization_user/organization_user_update_logic.go
- [x] T072 [P] Add first/last name validation (2-100 chars) in
  app/lender/api/internal/logic/organization_user/organization_user_update_logic.go
- [x] T073 [P] Add invitation code validation in
  app/lender/api/internal/logic/organization_user/organization_user_accept_invitation_logic.go

### Database Queries

- [x] T074 [P] Add role list query with org filter in domains/zeebee/repo/role.go (if needed)
- [x] T075 [P] Add role-user count check query in domains/zeebee/repo/role.go (if needed)
- [x] T076 [P] Add organization user list query with filters in domains/zeebee/repo/organization_user.go (if needed)
- [x] T077 [P] Add invitation code lookup query in domains/zeebee/repo/organization_user.go (if needed)
- [x] T078 [P] Add user role permissions query in domains/zeebee/repo/role.go (if needed)

### Performance & Caching

- [x] T079 [P] Add database indexes for role queries in domains/zeebee/model/role.go
- [x] T080 [P] Add database indexes for organization user queries in domains/zeebee/model/organization_user.go
- [x] T081 [P] Add database indexes for invitation lookups in domains/zeebee/model/organization_user.go

### Documentation

- [x] T082 [P] Add API documentation comments to role endpoints in app/lender/api/lender.api
- [x] T083 [P] Add API documentation comments to organization_user endpoints in app/lender/api/lender.api
- [x] T084 Update README.md in app/lender/ with new endpoint documentation
- [x] T085 Create example API requests in app/lender/docs/examples/role_examples.md
- [x] T086 Create example API requests in app/lender/docs/examples/organization_user_examples.md

---

## Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational)
    ↓
Phase 3 (US1: Role Management) ←──┐
    ↓                              │
Phase 4 (US2: User Management) ←───┘
    ↓
Phase 5 (US3: Invitations)
    ↓
Phase 6 (US4: Role Assignment)
    ↓
Phase 7 (Polish)
```

**User Story Dependencies**:

- US2 (User Management) depends on US1 (Role Management) - for role assignment
- US3 (Invitations) depends on US2 (User Management) - extends user creation
- US4 (Role Assignment) depends on US1 and US2 - requires both roles and users

---

## Parallel Execution Opportunities

### Phase 3 (US1) - Can run in parallel:

- T012-T018: All 7 RPC role logic files
- T020-T026: All 7 API role logic files (after RPC layer)

### Phase 4 (US2) - Can run in parallel:

- T029-T035: All 7 RPC organization user logic files
- T037-T043: All 7 API organization user logic files (after RPC layer)

### Phase 5 (US3) - Can run in parallel:

- T046, T047: RPC invitation logic files
- T048: Email template (independent)
- T049, T050: API invitation logic files (after RPC layer)

### Phase 6 (US4) - Single task:

- T051, T052: Must run sequentially (RPC → API)

### Phase 7 (Polish) - Can run in parallel:

- T053-T060: Error handling (8 tasks, different files)
- T061-T067: Logging (7 tasks, different files)
- T068-T073: Validation (6 tasks, different files)
- T074-T078: Database queries (5 tasks, different files)
- T079-T081: Performance (3 tasks, different files)
- T082-T086: Documentation (5 tasks, different files)

---

## MVP Scope

**Recommended MVP**: Phase 1 + Phase 2 + Phase 3 (User Story 1 - Role Management)

**MVP delivers**:

- Complete role CRUD operations
- Role permission management
- Deletion protection for in-use roles
- Paginated, filterable role list

**MVP Task Count**: 27 tasks (T001-T027)

**Post-MVP**: Incrementally add US2, US3, US4 for complete functionality

---

## Task Summary

| Phase                        | Tasks  | Description                               |
|------------------------------|--------|-------------------------------------------|
| Phase 1: Setup               | 3      | Copy API definition files                 |
| Phase 2: Foundational        | 8      | Update lender.api and generate code       |
| Phase 3: US1 Role Management | 16     | Complete role CRUD functionality          |
| Phase 4: US2 User Management | 17     | Organization user management              |
| Phase 5: US3 Invitations     | 5      | Invitation and acceptance flow            |
| Phase 6: US4 Role Assignment | 2      | Assign roles to users                     |
| Phase 7: Polish              | 35     | Error handling, logging, validation, docs |
| **Total**                    | **86** | **All implementation tasks**              |

---

## Format Validation

**All tasks follow the required format**: `- [x] [TaskID] [P?] [Story?] Description with file path`

- ✅ Checkbox present: `- [x]`
- ✅ Task ID present: T001-T086
- ✅ [P] marker for parallelizable tasks
- ✅ [Story] label for user story tasks (US1-US4)
- ✅ File paths included for all implementation tasks
- ✅ No tasks without file paths
- ✅ No story labels on setup/foundational/polish phases
