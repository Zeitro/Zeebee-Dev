# Data Model: Lender Organization Management

**Feature**: 003-lender-organization
**Date**: 2025-01-09
**Status**: Stable

## Overview

本文档定义了 lender 端组织管理功能涉及的数据模型和关系。

## Entity Definitions

### 1. Organization (组织)

**Table**: `organizations`
**Model**: `domains/zeebee/model/organization.go`

**Fields**:

| Field          | Type               | Description                              | Constraints            |
|----------------|--------------------|------------------------------------------|------------------------|
| Id             | int64              | Primary key                              | PK, auto-increment     |
| Name           | string             | Organization name                        | NOT NULL, max 100      |
| Status         | OrganizationStatus | Organization status                      | NOT NULL               |
| OwnerId        | int64              | ID of the owner user                     | NOT NULL, default 0    |
| Address        | Address           | Organization address (JSONB)             | nullable               |
| ContactEmail   | string             | Contact email                            | NOT NULL               |
| ContactPhone   | string             | Contact phone                            | max 20                 |
| ContactName    | string             | Contact person name                      | max 100                |
| OwnerType      | OrganizationOwnerType | Owner type (Self/ThirdParty/Partner)   | NOT NULL               |
| Type           | UserType           | User type (borrower/lender)              | NOT NULL               |
| Logo           | string             | Logo URL                                 | max 255                |
| Website        | string             | Website URL                              | max 255                |
| Tagline        | string             | Tagline                                  | max 255                |
| Metadata       | OrganizationMetadata | Additional metadata (JSONB)            | default '{}'           |
| CreatedAt      | time.Time          | Creation timestamp                       | auto                   |
| UpdatedAt      | time.Time          | Last update timestamp                    | auto                   |

**Relationships**:
- `Owner` → User (belongsTo, foreign_key: owner_id)
- `OrganizationEmails` → OrganizationEmail (hasMany)
- `Users` → User (manyToMany through organization_users)

**Enums**:
- `OrganizationStatus`: Active, Inactive
- `OrganizationOwnerType`: None, Self, ThirdParty, Partner
- `UserType`: Borrower, Lender

### 2. OrganizationUser (组织用户关联)

**Table**: `organization_users`
**Model**: `domains/zeebee/model/organization_user.go`

**Fields**:

| Field           | Type                         | Description                        | Constraints            |
|-----------------|------------------------------|------------------------------------|------------------------|
| Id              | int64                        | Primary key                        | PK, auto-increment     |
| OrganizationId  | int64                        | Organization ID                    | FK, NOT NULL           |
| UserId          | int64                        | User ID                            | FK, NOT NULL           |
| Status          | OrganizationUserStatus       | Association status                 | NOT NULL               |
| FirstName       | string                       | First name                         | max 50                 |
| MidName         | string                       | Middle name                        |                        |
| LastName        | string                       | Last name                          | max 50                 |
| FullName        | string                       | Full name                          |                        |
| Phone           | string                       | Phone number                       |                        |
| Remark          | string                       | Remarks                            |                        |
| CreatedAt       | time.Time                    | Creation timestamp                 | auto                   |
| UpdatedAt       | time.Time                    | Last update timestamp              | auto                   |

**Relationships**:
- `User` → User (belongsTo, foreign_key: user_id)
- `Organization` → Organization (belongsTo, foreign_key: organization_id)
- `Roles` → Role (manyToMany through organization_user_and_roles)

**Enums**:
- `OrganizationUserStatus`: Not Activated, Active, Inactive, Suspended

### 3. User (用户)

**Table**: `users`
**Model**: `domains/zeebee/model/user.go`

**Relevant Fields**:

| Field                | Type    | Description                    | Constraints            |
|----------------------|---------|--------------------------------|------------------------|
| Id                   | int64   | Primary key                   | PK, auto-increment     |
| DefaultOrganizationId| int64   | Current active organization ID | default 0              |
| FirstName            | string  | First name                     | NOT NULL               |
| LastName             | string  | Last name                      | NOT NULL               |
| Email                | string  | Email address                  | NOT NULL, unique       |
| Type                 | UserType | User type (borrower/lender)    | NOT NULL               |
| Status               | UserStatus| User status                    | NOT NULL               |

**Relationships**:
- `Organizations` → Organization (manyToMany through organization_users)

**Enums**:
- `UserStatus`: Unverified, Active, Inactive, Suspended
- `UserType`: Borrower, Lender

### 4. OrganizationEmail (组织邮箱)

**Table**: `organization_emails`
**Model**: `domains/zeebee/model/organization_email.go`

**Fields**:

| Field            | Type          | Description              | Constraints            |
|------------------|---------------|--------------------------|------------------------|
| Id               | int64         | Primary key              | PK, auto-increment     |
| OrganizationId   | int64         | Organization ID         | FK, NOT NULL           |
| Email            | string        | Email address            | NOT NULL               |
| AuthStrategy     | string        | Auth strategy            |                        |
| Purpose          | string        | Email purpose            |                        |
| IsAuthorized     | bool          | Authorization status     | default false          |
| CreatedAt        | time.Time     | Creation timestamp       | auto                   |
| UpdatedAt        | time.Time     | Last update timestamp    | auto                   |

**Relationships**:
- `Organization` → Organization (belongsTo, foreign_key: organization_id)

## Entity Relationships (ER Diagram)

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│    User     │         │ OrganizationUser │         │Organization │
├─────────────┤         ├──────────────────┤         ├─────────────┤
│ Id          │───┬────│ UserId           │─────────││ Id          │
│ Email       │   │    │ OrganizationId   │─────────┘│ Name        │
│ Type        │   │    │ Status           │         │ OwnerId     │──┐
│ Status      │   │    │ FullName         │         │ Status      │  │
│ ...         │   │    │ ...              │         │ Type        │  │
└─────────────┘   │    └──────────────────┘         │ ...         │  │
                  │                                   └─────────────┘  │
                  │                                                      │
                  │         ┌──────────────────┐                            │
                  └─────────│  OrganizationRole│                            │
                            │ (many-to-many)   │                            │
                            ├──────────────────┤                            │
                            │ OrganizationUserId│                            │
                            │ RoleId            │                            │
                            └──────────────────┘                            │
                                                                              │
                            ┌──────────────────┐                            │
                            │ OrganizationEmail│                            │
                            ├──────────────────┤                            │
                            │ OrganizationId    │────────────────────────────────┘
                            │ Email             │
                            │ IsAuthorized      │
                            │ ...              │
                            └──────────────────┘
```

## State Transitions

### Organization Status

```
┌───────────┐
│ Inactive  │  ← Initial state for newly created organizations
└─────┬─────┘
      │
      │ (activation process)
      ▼
┌───────────┐
│  Active   │  ← Activated organization
└───────────┘
```

### OrganizationUser Status

```
┌─────────────────┐
│ Not Activated   │  ← Initial state
└───────┬─────────┘
        │
        │ (activation)
        ▼
┌─────────────────┐
│     Active      │  ← Active user
└───────┬─────────┘
        │
        │ (suspension)
        ▼
┌─────────────────┐
│    Suspended    │  ← Suspended user
└─────────────────┘
```

## Validation Rules

### Organization Creation

**Required Fields**:
- Name: NOT NULL, max 100 chars
- ContactEmail: valid email format
- ContactPhone: NOT NULL, min 8 chars, max 20 chars
- ContactName: NOT NULL, min 2 chars, max 50 chars

**Auto-Set Fields**:
- Status: Inactive
- OwnerType: ThirdParty
- Type: Lender (for this feature)
- OwnerId: Current user's ID

### Organization Update

**Allowed Fields**:
- ContactEmail
- ContactPhone
- ContactName
- Logo
- Website
- Tagline
- Metadata (lender_metadata)

**Validation**:
- Only owner can update
- Email format validation
- Phone length validation

## Indexes

### organizations
- PRIMARY KEY (id)
- INDEX idx_type_status (type, status)
- INDEX idx_owner_id (owner_id)
- UNIQUE INDEX idx_name_type (name, type)

### organization_users
- PRIMARY KEY (id)
- INDEX idx_user_id (user_id)
- INDEX idx_organization_id (organization_id)
- INDEX idx_status (status)
- UNIQUE INDEX idx_user_org (user_id, organization_id)

### organization_emails
- PRIMARY KEY (id)
- INDEX idx_organization_id (organization_id)
- INDEX idx_email (email)

## Data Integrity

### Constraints

1. **Organization Owner**:
   - `organizations.owner_id` must reference a valid user
   - Owner must be a member of the organization

2. **OrganizationUser Association**:
   - A user can only have one active association per organization
   - At least one user must be the owner

3. **Referential Integrity**:
   - All foreign keys are enforced
   - CASCADE delete is NOT used (soft delete only)

## JSONB Schemas

### Organization Metadata

```json
{
  "lender": {
    "loan_amount_range": {
      "min": 1000.0,
      "max": 100000.0
    },
    "loan_term_range": {
      "min": 12,
      "max": 360
    },
    "loan_rate_range": {
      "min": 3.5,
      "max": 25.0
    },
    "loan_point_range": {
      "min": 0.0,
      "max": 5.0
    }
  }
}
```

### Address

```json
{
  "street": "123 Main St",
  "city": "Springfield",
  "state": "IL",
  "zip_code": "62701",
  "country": "USA"
}
```

## Migration Notes

All tables and relationships already exist in the database. This feature does not require schema changes, only implements business logic for CRUD operations.

## References

- `domains/zeebee/model/organization.go` - Organization model definition
- `domains/zeebee/model/organization_user.go` - OrganizationUser model definition
- `domains/zeebee/model/user.go` - User model definition
- `domains/zeebee/constants/organization.go` - Organization constants and enums
