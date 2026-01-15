# Implementation Plan: Lender Organization Management

**Branch**: `003-lender-organization` | **Date**: 2025-01-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/003-lender-organization/spec.md`

## Summary

实现 lender 端组织管理功能，包括三个核心能力：查看当前组织信息、更新组织信息、创建新组织。该功能遵循 go-zero 微服务架构，API 层负责请求验证，RPC 层处理业务逻辑，使用共享领域层进行数据操作。

## Technical Context

**Language/Version**: Go 1.23+
**Primary Dependencies**: go-zero 1.9.3, GORM, Redis, PostgreSQL
**Storage**: PostgreSQL 14+ with existing organizations and organization_users tables
**Testing**: go test, golangci-lint
**Target Platform**: Linux server (microservices architecture)
**Project Type**: microservices (API/RPC separation)
**Performance Goals**: <200ms API response (p95), 1000+ req/s
**Constraints**: API/RPC dual-layer, shared domain models, copierx.Trans for conversions
**Scale/Scope**: Lender service only, leverages existing shared domain layer

### Known Patterns

- **User Context Retrieval**: Use `l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)` to get JWT claims from gRPC metadata
- **Organization Query**: Use `l.baseLogic.OrganizationLogic.Detail(ctx, organizationId, userType)` to get organization with preloaded associations
- **Data Conversion**: Use `copierx.Trans[SourceType, TargetType](source, copierx.TimeToRFC3339Str())` for all type conversions
- **Error Handling**: Use `common/errorx` for consistent error wrapping and messages
- **Repository Operations**: Use `OrganizationRepo` methods (GetOrganization, UpdateColumnsById, CreateWithTx) for database operations

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Based on `.specify/memory/constitution.md`:

- [x] **Shared Domain-First**: Uses existing models in `domains/zeebee/model/` (Organization, OrganizationUser)
- [x] **API/RPC Separation**: API layer validates requests, RPC layer contains business logic
- [x] **Code Generation**: Uses existing `.api` and `.proto` definitions with goctl
- [x] **Data Conversion**: All conversions use `copierx.Trans[T, U]` with time converters
- [x] **Type Safety**: Implementation will follow Go best practices with proper typing
- [x] **Test Coverage**: Unit tests for logic layer, integration tests for RPC endpoints

**Violations Requiring Justification**:

| Principle           | Potential Violation | Justification Required?                          |
|---------------------|---------------------|--------------------------------------------------|
| Shared Domain-First |                     | ✅ Using existing `domains/zeebee/` models       |
| API/RPC Separation  |                     | ✅ Clear separation maintained                   |
| Code Generation     |                     | ✅ Using goctl-generated code                    |
| Data Conversion     |                     | ✅ All conversions use copierx.Trans             |

**Status**: ✅ All checks passed - No violations

## Project Structure

### Documentation (this feature)

```text
specs/003-lender-organization/
├── spec.md              # Feature specification (completed)
├── plan.md              # This file
├── research.md          # Technical research and decisions
├── data-model.md        # Entity definitions and relationships
├── quickstart.md        # Developer quick start guide
├── contracts/           # API contracts (already exist)
│   ├── organization.api
│   └── user.api
└── tasks.md             # Implementation tasks (created by /speckit.tasks)
```

### Source Code (repository root)

```text
app/lender/                      # Lender service
├── api/                         # HTTP API layer
│   ├── desc/                    # API definitions (.api files)
│   │   ├── organization.api     # Organization API definitions (exist)
│   │   └── user.api             # User API definitions (exist)
│   └── internal/
│       ├── handler/             # HTTP handlers (already generated)
│       │   └── organization/
│       │       ├── organization_profile_handler.go
│       │       ├── organization_update_profile_handler.go
│       │       └── organization_create_handler.go
│       ├── logic/               # API logic (to be implemented)
│       │   └── organization/
│       │       ├── organization_profile_logic.go
│       │       ├── organization_update_profile_logic.go
│       │       └── organization_create_logic.go
│       └── types/               # Type definitions (auto-generated)
└── rpc/                         # gRPC RPC layer
    ├── pb/                      # Protocol buffers (already defined)
    │   └── lender/
    │       ├── lender.proto      # gRPC service definitions
    │       └── lender.pb.go      # Generated protobuf code
    ├── internal/
    │   └── logic/               # RPC logic (to be implemented)
    │       └── organization/
    │           ├── profile_logic.go
    │           ├── update_profile_logic.go
    │           └── create_logic.go
    └── client/                  # RPC clients (already generated)
        └── organization/
            └── organization.go

domains/zeebee/                  # Shared domain layer (existing)
├── model/                       # Unified data models (existing)
│   ├── organization.go          # Organization model
│   └── organization_user.go     # OrganizationUser model
├── repo/                        # Shared repositories (existing)
│   └── organization.go          # Organization repository
├── logic/                       # Business logic (existing)
│   ├── organization.go          # Organization business logic
│   └── user_logic.go            # User business logic
└── constants/                   # Business constants (existing)
    └── organization.go          # Organization constants

common/                          # Shared utilities (existing)
├── errorx/                      # Error handling
├── copierx/                     # Type conversion
├── validator/                   # Data validation
└── utils/                       # Utilities (JWT, context, etc.)
```

**Structure Decision**: Leveraging existing microservices architecture with API/RPC dual-layer. All new code follows established patterns, uses existing shared domain models, and maintains consistency with codebase standards.

## Implementation Phases

### Phase 0: Research & Technical Decisions

**Status**: ✅ Complete

Research completed during specification clarification phase. Key findings documented in spec.md Clarifications section:

1. **User Context Access Pattern**: Confirmed use of `baseLogic.UserLogic.GetUserClaimFromCtx()` for gRPC metadata
2. **Organization Query Pattern**: Confirmed use of `baseLogic.OrganizationLogic.Detail()` for data retrieval
3. **Data Conversion Pattern**: Confirmed use of `copierx.Trans[T, U]()` for type conversions

**Output**: See spec.md - Clarifications section (2025-01-09 session)

### Phase 1: Design & Contracts

**Status**: In Progress

This phase will produce:
- `research.md` - Consolidated technical research and decisions
- `data-model.md` - Entity definitions and relationships
- `contracts/` - API contract definitions (already exist in codebase)
- `quickstart.md` - Developer quick start guide

#### 1.1 Data Model

Entities to implement:

1. **Organization** (existing model in `domains/zeebee/model/organization.go`)
   - Fields: id, name, status, owner_id, address, contact_email, contact_phone, contact_name, owner_type, type, logo, website, tagline, metadata
   - Relationships: Owner (User), OrganizationEmails (hasMany)

2. **OrganizationUser** (existing model in `domains/zeebee/model/organization_user.go`)
   - Fields: id, organization_id, user_id, status, first_name, mid_name, last_name, full_name, phone, remark
   - Relationships: User (belongsTo), Organization (belongsTo), Roles (manyToMany)

3. **User** (existing model, referenced via relationships)
   - Used for owner verification and user-context operations

#### 1.2 API Contracts

API endpoints already defined in:
- `app/lender/api/desc/organization.api`

Endpoints to implement:
1. `GET /organization/profile` - Get current user's organization
2. `PUT /organization/profile` - Update organization information
3. `POST /organization` - Create new organization

#### 1.3 RPC Services

RPC services already defined in:
- `app/lender/rpc/pb/lender/lender.proto`

Services to implement:
1. `rpc Profile(EmptyRequest) returns (OrganizationEntity)` - Get organization by context
2. `rpc UpdateProfile(OrganizationUpdateProfileRequest) returns (OrganizationEntity)` - Update organization
3. `rpc Create(OrganizationCreateRequest) returns (OrganizationCreateResponse)` - Create organization

### Phase 2: Implementation Strategy

**Status**: Pending (will be executed via `/speckit.tasks`)

#### 2.1 RPC Layer Implementation

**Priority: P1 - Core Business Logic**

Files to implement:
1. `app/lender/rpc/internal/logic/organization/profile_logic.go`
   - Get user claim from context using `baseLogic.UserLogic.GetUserClaimFromCtx()`
   - Validate organization_id is not zero
   - Call `baseLogic.OrganizationLogic.Detail()` to get organization
   - Convert BO to proto entity using `copierx.Trans[]`
   - Return organization entity

2. `app/lender/rpc/internal/logic/organization/update_profile_logic.go`
   - Get user claim from context
   - Retrieve organization and verify user is owner
   - Validate input data (email format, phone length, required fields)
   - Build update map with allowed fields only
   - Call `OrganizationRepo.UpdateColumnsById()` or `UpdateOneById()`
   - Return updated organization

3. `app/lender/rpc/internal/logic/organization/create_logic.go`
   - Get user claim from context (creator will be owner)
   - Validate required fields
   - Create organization model with:
     - Status = Inactive
     - OwnerType = ThirdParty
     - Type = UserTypeLender
   - Use transaction to:
     - Create organization
     - Update owner_id
     - Create OrganizationUser association with admin role
   - Return created organization

#### 2.2 API Layer Implementation

**Priority: P2 - Request Validation & RPC Orchestration**

Files to implement:
1. `app/lender/api/internal/logic/organization/organization_profile_logic.go`
   - Call RPC client's Profile method
   - Convert proto response to API types using `copierx.Trans[]`
   - Return API response

2. `app/lender/api/internal/logic/organization/organization_update_profile_logic.go`
   - Validate request structure (handled by go-zero)
   - Build RPC request from API request
   - Call RPC client's UpdateProfile method
   - Convert proto response to API types
   - Return API response

3. `app/lender/api/internal/logic/organization/organization_create_logic.go`
   - Validate request structure (handled by go-zero)
   - Build RPC request from API request
   - Call RPC client's Create method
   - Convert proto response to API types
   - Return API response

#### 2.3 Testing Strategy

**Unit Tests** (for RPC logic layer):
- Test Profile: happy path, no organization context, invalid organization
- Test UpdateProfile: owner success, non-owner rejection, validation errors
- Test Create: success with correct defaults, validation failures

**Integration Tests** (for API endpoints):
- Test each endpoint with authenticated user context
- Test permission checks
- Test data validation
- Test error scenarios

### Phase 3: Quality Assurance

**Status**: Pending

1. **Code Quality**:
   - Run `golangci-lint` on all modified files
   - Ensure error handling uses `common/errorx`
   - Verify all conversions use `copierx.Trans[]`

2. **Testing**:
   - Unit tests for RPC logic (target: >80% coverage)
   - Integration tests for API endpoints
   - Manual testing with real JWT tokens

3. **Documentation**:
   - Code comments in English
   - Update quickstart.md with setup instructions
   - Document any new patterns or conventions

## Dependencies & Risks

### External Dependencies
1. **PostgreSQL**: Organizations and organization_users tables must exist
2. **Redis**: For JWT token caching (if applicable)
3. **Shared Domain**: Models and repositories in `domains/zeebee/` must be up-to-date

### Internal Dependencies
1. **User Authentication**: JWT middleware must populate context correctly
2. **BaseLogic**: Must have access to UserLogic and OrganizationLogic
3. **RPC Client**: API layer must have initialized RPC client

### Risks & Mitigations

| Risk                                    | Impact | Mitigation                                               |
|-----------------------------------------|--------|---------------------------------------------------------|
| Organization not found                 | Medium | Return clear error, guide user to create organization   |
| User not owner of organization          | Medium | Validate ownership before update, return permission error |
| Concurrent organization updates        | Low    | Database row-level locking via GORM transactions         |
| Invalid data in create/update           | Low    | Use go-zero validation tags, return validation errors   |
| Context missing user claim              | Low    | Check early, return authentication error                 |

## Success Metrics

From spec.md Success Criteria:
- **SC-001**: Profile response < 1 second (p95)
- **SC-002**: UpdateProfile response < 2 seconds (p95)
- **SC-003**: Create response < 3 seconds (p95)
- **SC-004**: 100% permission check accuracy
- **SC-005**: 100% data validation coverage
- **SC-006**: 100% compliance with organization creation defaults
- **SC-007**: Support 1000+ concurrent profile queries
- **SC-008**: 99.9%+ transaction success rate for updates

## Next Steps

1. ✅ Complete specification (spec.md)
2. ✅ Add technical clarifications
3. 🔄 Generate this plan (plan.md) - **CURRENT**
4. ⏭️ Run `/speckit.tasks` to generate detailed task list
5. ⏭️ Implement RPC layer logic
6. ⏭️ Implement API layer logic
7. ⏭️ Write tests
8. ⏭️ Quality assurance and validation

---

**Plan Version**: 1.0
**Last Updated**: 2025-01-09
**Status**: Ready for task generation via `/speckit.tasks`
