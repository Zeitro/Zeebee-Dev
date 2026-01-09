# Research: Lender Role and Organization User Management

**Feature**: 001-lender-role-org
**Date**: 2025-12-29
**Status**: ✅ Completed

## Overview

This document consolidates research findings for implementing role and organization user management in the lender
service. All decisions align with existing borrower service patterns and leverage shared domain models.

**Research Outcome**: All technical decisions validated through implementation. The approach of mirroring borrower
service architecture proved successful.

## Technical Decisions

### 1. Architecture Pattern

**Decision**: Mirror borrower service architecture with API/RPC separation

**Rationale**:

- Consistency across services reduces cognitive load
- Existing patterns are proven and tested
- Shared domain models already support this architecture

**Alternatives Considered**:

- Monolithic API-only approach: Rejected because it would break existing service boundaries
- Unified microservice: Rejected because it would require extensive refactoring

### 2. Data Layer Strategy

**Decision**: Use existing shared domain models from `domains/zeebee`

**Rationale**:

- Models (role, organization_user, resource) already exist and are tested
- Repositories (role.go, organization_user.go) provide CRUD operations
- Borrower service already uses these shared models successfully

**Alternatives Considered**:

- Lender-specific models: Rejected because role/permission logic is domain-agnostic
- Duplicate models in lender service: Rejected because it violates DRY principle

### 3. API Definition Approach

**Decision**: Create separate .api files in `app/lender/api/desc/`

**Files to Create**:

- `role.api` - Role CRUD request/response types
- `organization_user.api` - Organization user management types
- `resource.api` - Resource/permission types (if not already present)

**Rationale**:

- Follows borrower service organization
- Enables goctl code generation
- Clear separation of concerns

**Alternatives Considered**:

- Single api file: Rejected because it would become unwieldy
- Inline types in lender.api: Rejected because it reduces reusability

### 4. Code Generation Strategy

**Decision**: Use goctl for generating handlers and logic from .api definitions

**Commands to Use**:

```bash
make create-api-service lender  # For API layer
make create-rpc-service lender   # For RPC layer (if needed)
```

**Rationale**:

- Standard tool for go-zero projects
- Maintains code consistency
- Reduces boilerplate

**Alternatives Considered**:

- Manual implementation: Rejected because it's time-consuming and error-prone
- Custom generation scripts: Rejected because goctl is already configured

### 5. Invitation Flow Implementation

**Decision**: Implement invitation code generation with email delivery

**Key Components**:

- Unique invitation code generation (UUID-based)
- Email service integration (existing common/email package)
- Invitation expiration handling (7-day default)
- One-time use validation

**Rationale**:

- Matches borrower service invitation pattern
- Email service already exists in common layer
- UUID generation utility available (common/uuid)

**Alternatives Considered**:

- Magic links (email token): Rejected because invitation code provides better tracking
- No expiration: Rejected for security reasons

### 6. Authentication & Authorization

**Decision**: Use existing AuthMiddleware for protected endpoints

**Rationale**:

- AuthMiddleware already implemented in lender service
- JWT-based authentication already established
- Role-based permission checking can be added to middleware

**Implementation Notes**:

- Role CRUD endpoints: AuthMiddleware + ownership check
- Organization user management: AuthMiddleware + organization membership check
- Public endpoints: register, login, accept_invitation

**Alternatives Considered**:

- New permission middleware: Rejected because existing middleware can be extended
- Gateway-level auth: Rejected because it adds infrastructure complexity

### 7. RPC vs Direct Database Access

**Decision**: API layer calls RPC layer for business logic

**Rationale**:

- Follows established service pattern
- Enables potential service-to-service communication
- Clear separation of concerns (API = validation/routing, RPC = business logic)

**Alternatives Considered**:

- Direct database access from API: Rejected because it bypasses business logic layer
- API-only (no RPC): Rejected because it breaks microservice pattern

## Integration Points

### Existing Dependencies

1. **Email Service** (`common/email`)
    - Used for invitation emails
    - Templates may need lender-specific customization

2. **UUID Generator** (`common/uuid`)
    - Used for invitation code generation
    - Sonyflake-based implementation

3. **Domain Models** (`domains/zeebee`)
    - Role, OrganizationUser, Resource models
    - Repositories for CRUD operations

4. **Error Handling** (`common/errorx`)
    - Standardized error responses
    - Validation error handling

## Best Practices

### Go-Zero Framework

1. **API Definition**:
    - Use `validate` tags for input validation
    - Define clear request/response types
    - Group related endpoints with `@server` directive

2. **Handler Implementation**:
    - Keep handlers thin - delegate to logic layer
    - Use consistent error handling patterns
    - Log important operations

3. **RPC Implementation**:
    - Implement business logic in RPC layer
    - Use transactions for multi-step operations
    - Validate permissions before data modifications

### Database Operations

1. **Use Repository Pattern**:
    - Leverage existing `domains/zeebee/repo` implementations
    - Extend repos only if lender-specific queries needed

2. **Transaction Management**:
    - Use transactions for role assignment updates
    - Handle rollback on errors

3. **Soft Deletes**:
    - Respect GORM soft delete functionality
    - Use `DeletedAt` field for filtering

### Security Considerations

1. **Ownership Validation**:
    - Verify organization membership before allowing operations
    - Check role ownership before role modifications

2. **Invitation Security**:
    - One-time use invitation codes
    - Expiration checking
    - Email verification during acceptance

3. **Permission Enforcement**:
    - Validate permissions at RPC layer
    - Log authorization failures

## Unresolved Questions

None - all technical decisions have been made based on existing patterns and requirements.

## References

- Borrower service implementation: `app/borrower/api/`, `app/borrower/rpc/`
- Shared domain models: `domains/zeebee/model/`, `domains/zeebee/repo/`
- Go-zero documentation: https://go-zero.dev/
- Project CLAUDE.md: `/backend/CLAUDE.md`
