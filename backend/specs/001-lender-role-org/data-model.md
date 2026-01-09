# Data Model: Lender Role and Organization User Management

**Feature**: 001-lender-role-org
**Date**: 2025-12-29
**Status**: ✅ Completed

## Overview

This document describes the data model for role and organization user management. The implementation uses shared domain
models from `domains/zeebee/model`, which are already used by the borrower service.

**Validation Result**: All shared domain models (Role, OrganizationUser, Resource, RoleAndResource,
OrganizationUserRole) work correctly for lender service. No schema changes required.

## Entity Definitions

### Role

**Purpose**: Represents a collection of permissions that can be assigned to users within an organization.

**Attributes**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Id | int64 | Primary key | auto-generated |
| OrganizationId | int64 | Foreign key to organization | required, > 0 |
| Name | string | Role name | required, 2-100 chars, unique per org |
| Description | string | Role description | optional, max 255 chars |
| Type | string | Role type (e.g., "custom", "system") | default: "custom" |
| Status | string | Active/Inactive | "Active" or "Inactive" |
| CreatedAt | timestamp | Creation timestamp | auto-generated |
| UpdatedAt | timestamp | Last update timestamp | auto-generated |
| DeletedAt | timestamp | Soft delete timestamp | nullable |

**Relationships**:

- Belongs to: Organization (OrganizationId)
- Has many: OrganizationUser (via role assignments)
- Has many: Resource (via RoleAndResource join table)

**State Transitions**:

```
[Created] → [Active] ↔ [Inactive] → [Deleted]
                    ↓
                 (Cannot delete if users assigned)
```

**Validation Rules**:

- Name cannot be "admin" (system reserved)
- Name must be unique within organization
- Cannot delete if assigned to any users

### OrganizationUser

**Purpose**: Represents a user's membership in an organization with assigned roles.

**Attributes**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Id | int64 | Primary key | auto-generated |
| OrganizationId | int64 | Foreign key to organization | required, > 0 |
| UserId | int64 | Foreign key to user | required, > 0 |
| Status | string | Membership status | "Active", "Disabled", "Pending" |
| FirstName | string | User's first name | required, 2-100 chars |
| MidName | string | User's middle name | optional |
| LastName | string | User's last name | required, 2-100 chars |
| Phone | string | Contact phone | optional |
| Remark | string | Additional notes | optional, max 255 chars |
| CreatedAt | timestamp | Creation timestamp | auto-generated |
| UpdatedAt | timestamp | Last update timestamp | auto-generated |
| DeletedAt | timestamp | Soft delete timestamp | nullable |

**Relationships**:

- Belongs to: Organization (OrganizationId)
- Belongs to: User (UserId)
- Has many: Role (via OrganizationUserRole join table)

**State Transitions**:

```
[Pending] → [Active] ↔ [Disabled]
    ↓           ↓
[Invitation   (Can access
 Expires]    organization)
```

**Validation Rules**:

- UserId + OrganizationId combination must be unique
- User must be "Active" status to access organization resources
- Pending users cannot access resources until invitation accepted

### Resource (API)

**Purpose**: Represents a specific capability or permission in the system (e.g., "create:loan", "view:reports").

**Attributes**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Id | int64 | Primary key | auto-generated |
| Code | string | Resource code (e.g., "user:create") | required, unique |
| Name | string | Resource display name | required |
| Type | string | Resource type | "menu", "button", "api" |
| ParentId | int64 | Parent resource ID | nullable (for tree structure) |
| Path | string | URL path | optional |
| Method | string | HTTP method | optional |
| Sort | int | Display order | default: 0 |
| Status | string | Active/Inactive | "Active" or "Inactive" |
| CreatedAt | timestamp | Creation timestamp | auto-generated |
| UpdatedAt | timestamp | Last update timestamp | auto-generated |
| DeletedAt | timestamp | Soft delete timestamp | nullable |

**Relationships**:

- Has many: Role (via RoleAndResource join table)
- Self-referencing: ParentId → Id (tree structure)

### RoleAndResource (Join Table)

**Purpose**: Many-to-many relationship between roles and resources.

**Attributes**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Id | int64 | Primary key | auto-generated |
| RoleId | int64 | Foreign key to role | required, > 0 |
| ResourceId | int64 | Foreign key to resource | required, > 0 |
| CreatedAt | timestamp | Creation timestamp | auto-generated |

**Constraints**:

- RoleId + ResourceId combination must be unique

### OrganizationUserRole (Join Table)

**Purpose**: Many-to-many relationship between organization users and roles.

**Attributes**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Id | int64 | Primary key | auto-generated |
| OrganizationUserId | int64 | Foreign key to org user | required, > 0 |
| RoleId | int64 | Foreign key to role | required, > 0 |
| CreatedAt | timestamp | Creation timestamp | auto-generated |

**Constraints**:

- OrganizationUserId + RoleId combination must be unique

### Invitation (Derived)

**Purpose**: Represents a pending invitation for a user to join an organization. This may be stored as part of
OrganizationUser with status "Pending" and an invitation code.

**Attributes**:
| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| Code | string | Unique invitation code | UUID-based |
| Email | string | Invitee email address | required, valid email |
| OrganizationId | int64 | Target organization | required, > 0 |
| RoleIds | []int64 | Roles to assign upon acceptance | optional |
| ExpiresAt | timestamp | Invitation expiration | default: +7 days |
| Status | string | Invitation status | "Pending", "Accepted", "Expired" |
| CreatedAt | timestamp | Creation timestamp | auto-generated |

**State Transitions**:

```
[Pending] → [Accepted]
    ↓
[Expired] (after 7 days or when used)
```

## ER Diagram

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│ Organization│ 1     N │ OrganizationUser │ N     1 │    User     │
└─────────────┘─────────└──────────────────┘─────────└─────────────┘
       │ 1                    │ N
       │                      │
       │ N                    │ N
┌──────▼──────┐         ┌────▼────────┐
│    Role     │         │OrgUserRole  │
└──────┬──────┘         └─────────────┘
       │ N
       │
       │ N
┌──────▼──────┐         ┌──────────────┐
│RoleAndResource│ N   N │  Resource    │
└─────────────┘─────────└──────────────┘
```

## Indexes

For optimal query performance, the following indexes are recommended:

**Role**:

- `idx_org_id` (OrganizationId)
- `idx_org_status` (OrganizationId, Status)
- `idx_org_name` (OrganizationId, Name) - unique

**OrganizationUser**:

- `idx_org_id` (OrganizationId)
- `idx_user_id` (UserId)
- `idx_org_status` (OrganizationId, Status)
- `idx_org_user` (OrganizationId, UserId) - unique

**Resource**:

- `idx_code` (Code) - unique
- `idx_type_status` (Type, Status)
- `idx_parent_id` (ParentId)

**RoleAndResource**:

- `idx_role_resource` (RoleId, ResourceId) - unique

**OrganizationUserRole**:

- `idx_orguser_role` (OrganizationUserId, RoleId) - unique

## Notes

1. **Shared Models**: All models defined above exist in `domains/zeebee/model` and are shared with borrower service.

2. **Soft Deletes**: All main entities use soft deletes via `DeletedAt` timestamp. Queries must filter out deleted
   records.

3. **Organization Scope**: All role and permission data is scoped to organizations. There is no global role system.

4. **Caching Consideration**: Role-resource mappings may be cached for performance, with invalidation on updates.

5. **Audit Trail**: Consider adding audit tables (role_audit, organization_user_audit) if detailed change tracking is
   required.
