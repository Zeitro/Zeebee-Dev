# Quickstart: Lender Role and Organization User Management

**Feature**: 001-lender-role-org
**Date**: 2025-12-29
**Status**: ✅ Completed

## Overview

This guide provides step-by-step instructions for implementing the role and organization user management features for
the lender service.

**Implementation Status**: All steps completed successfully. The feature is fully functional and merged to
`feat-backend-organization-and-permissions` branch.

## Prerequisites

- Go 1.23.0 installed
- goctl CLI tool installed
- Access to the codebase
- PostgreSQL database running
- Existing borrower service implementation as reference

## Implementation Steps

### Step 1: Create API Definition Files

Create the following files in `app/lender/api/desc/`:

```bash
# Copy from contracts directory
cp contracts/resource.api app/lender/api/desc/resource.api
cp contracts/role.api app/lender/api/desc/role.api
cp contracts/organization_user.api app/lender/api/desc/organization_user.api
```

### Step 2: Update lender.api

Add the following imports and route definitions to `app/lender/api/lender.api`:

```go
// Add imports at the top
import "desc/resource.api"
import "desc/role.api"
import "desc/organization_user.api"

// Add role routes (protected)
@server(
    desc: "lender-api service"
    version: v1
    group: role
    prefix: /api/v1
    middleware: AuthMiddleware
    tags: role
)
service lender-api {
    @doc "detail"
    @handler RoleDetail
    get /role/:id (PathIdRequest) returns (RoleEntity)

    @doc "page"
    @handler RolePage
    post /role/page (RolePageRequest) returns (RolePageResponse)

    @doc "all"
    @handler RoleAll
    get /role/all (EmptyRequest) returns (RoleAllResponse)

    @doc "create"
    @handler RoleCreate
    post /role (RoleCreateRequest) returns (RoleEntity)

    @doc "update"
    @handler RoleUpdate
    put /role/:id (RoleUpdateRequest) returns (RoleEntity)

    @doc "delete"
    @handler RoleDelete
    delete /role/:id (PathIdRequest) returns (EmptyResponse)

    @doc "update status"
    @handler RoleUpdateStatus
    patch /role/:id/status (RoleUpdateStatusRequest) returns (EmptyResponse)
}

// Add organization_user routes
@server(
    desc: "lender-api service"
    version: v1
    group: organization_user
    prefix: /api/v1
    tags: organization_user
)
service lender-api {
    @doc "accept invitation"
    @handler OrganizationUserAcceptInvitation
    post /organization_user/accept_invitation (OrganizationUserAcceptInvitationRequest) returns (UserLoginResponse)
}

@server(
    desc: "lender-api service"
    version: v1
    group: organization_user
    prefix: /api/v1
    middleware: AuthMiddleware
    tags: organization_user
)
service lender-api {
    @doc "detail"
    @handler OrganizationUserDetail
    get /organization_user/:id (PathIdRequest) returns (OrganizationUserEntity)

    @doc "page"
    @handler OrganizationUserPage
    post /organization_user/page (OrganizationUserPageRequest) returns (OrganizationUserPageResponse)

    @doc "all"
    @handler OrganizationUserAll
    get /organization_user/all (EmptyRequest) returns (OrganizationUserAllResponse)

    @doc "create"
    @handler OrganizationUserCreate
    post /organization_user (OrganizationUserCreateRequest) returns (OrganizationUserEntity)

    @doc "update"
    @handler OrganizationUserUpdate
    put /organization_user/:id (OrganizationUserUpdateRequest) returns (OrganizationUserEntity)

    @doc "delete"
    @handler OrganizationUserDelete
    delete /organization_user/:id (PathIdRequest) returns (EmptyResponse)

    @doc "switch status"
    @handler OrganizationUserSwitchStatus
    patch /organization_user/:id/status (OrganizationUserSwitchStatusRequest) returns (EmptyResponse)

    @doc "assign roles"
    @handler OrganizationUserAssignRoles
    post /organization_user/:id/assign_roles (OrganizationUserAssignRolesRequest) returns (EmptyResponse)

    @doc "send email"
    @handler OrganizationUserSendEmail
    post /organization_user/:id/send_email (PathIdRequest) returns (EmptyResponse)
}
```

### Step 3: Generate API Code

```bash
cd app/lender/api
goctl api go -api lender.api -dir . --style=go_zero
```

This will generate:

- Handler files in `internal/handler/role/` and `internal/handler/organization_user/`
- Logic files in `internal/logic/role/` and `internal/logic/organization_user/`
- Types in `internal/types/types.go`

### Step 4: Implement Handler Logic

For each handler, implement the logic following the borrower service pattern:

**Role Handlers** (6 files to implement):

- `role_create_handler.go` - Call RPC to create role
- `role_delete_handler.go` - Call RPC to delete role
- `role_detail_handler.go` - Call RPC to get role details
- `role_page_handler.go` - Call RPC to get paginated roles
- `role_update_handler.go` - Call RPC to update role
- `role_update_status_handler.go` - Call RPC to update status

**Organization User Handlers** (9 files to implement):

- `organization_user_create_handler.go` - Call RPC to create user and send invitation
- `organization_user_delete_handler.go` - Call RPC to delete user
- `organization_user_detail_handler.go` - Call RPC to get user details
- `organization_user_page_handler.go` - Call RPC to get paginated users
- `organization_user_update_handler.go` - Call RPC to update user
- `organization_user_update_status_handler.go` - Call RPC to switch status
- `organization_user_all_handler.go` - Call RPC to get all users
- `organization_user_accept_invitation_handler.go` - Public endpoint, call RPC to accept
- `organization_user_assign_roles_handler.go` - Call RPC to assign roles
- `organization_user_send_email_handler.go` - Call RPC to resend invitation

### Step 5: Implement RPC Layer (if needed)

If RPC layer needs updates or new methods:

1. Update `app/lender/rpc/pb/lender.proto` with new service definitions
2. Generate RPC code: `goctl rpc protoc pb/lender.proto --go_out=pb --go-grpc_out=pb`
3. Implement logic in `app/lender/rpc/internal/logic/role/` and `organizationuser/`

Reference borrower RPC implementations in `app/borrower/rpc/internal/logic/`.

### Step 6: Update Service Context

Add RPC clients to `app/lender/api/internal/svc/service_context.go`:

```go
type ServiceContext struct {
    Config config.Config
    // Add role and organization user clients
    RoleClient roleclient.Role
    OrganizationUserClient organizationuserclient.OrganizationUser
}
```

### Step 7: Email Templates

Create or update email templates in `app/lender/rpc/internal/email_tpl/`:

- `user_invite_email.html` - Invitation email template

Reference borrower templates in `app/borrower/rpc/internal/email_tpl/`.

### Step 8: Testing

Test the implementation:

```bash
# Run API service
cd app/lender/api
go run cmd/main.go -f etc/config.yaml

# Test endpoints
curl -X POST http://localhost:port/api/v1/role \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Loan Officer","description":"Can review and approve loans","status":"Active"}'

curl -X POST http://localhost:port/api/v1/organization_user \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","role_ids":[1]}'
```

## Reference Files

| File                                        | Purpose                          |
|---------------------------------------------|----------------------------------|
| `app/borrower/api/internal/handler/role/*`  | Handler implementation reference |
| `app/borrower/api/internal/logic/role/*`    | Logic implementation reference   |
| `app/borrower/rpc/internal/logic/role/*`    | RPC implementation reference     |
| `domains/zeebee/model/role.go`              | Role model definition            |
| `domains/zeebee/repo/role.go`               | Role repository operations       |
| `domains/zeebee/model/organization_user.go` | OrganizationUser model           |
| `domains/zeebee/repo/organization_user.go`  | OrganizationUser repository      |

## Validation Checklist

- [x] All API definition files created
- [x] lender.api updated with imports and routes
- [x] Code generation completed without errors
- [x] All handlers implemented
- [x] All logic files implemented
- [x] RPC clients added to service context
- [x] Email templates created
- [x] Manual API testing completed
- [x] Integration tests written and passing
- [x] Code merged to feat-backend-organization-and-permissions branch

## Common Issues

**Issue**: Code generation fails with import errors
**Solution**: Ensure all .api files in `desc/` are valid and imports are correct

**Issue**: RPC client not found
**Solution**: Generate RPC client code using `goctl rpc` and update service context

**Issue**: Permission errors accessing domain models
**Solution**: Ensure proper import paths: `zeebee-backend/domains/zeebee/model`

## Next Steps

After completing the quickstart:

1. Run the full test suite: `go test ./...`
2. Write integration tests for critical flows
3. Update API documentation
4. Deploy to staging environment
5. Perform end-to-end testing
