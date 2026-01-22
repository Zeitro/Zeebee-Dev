# Data Model: 管理后台组织管理 CRUD

**Feature**: [spec.md](./spec.md)
**Date**: 2025-01-12

## Overview

本文档定义了组织管理 CRUD 功能涉及的数据模型、实体关系和验证规则。

## Core Entities

### 1. Organization（组织）

**Table**: `organizations`

**Purpose**: 代表一个 Lender（贷款方）组织实体

**Fields**:

| Field         | Type      | Constraints                               | Description       |
|---------------|-----------|-------------------------------------------|-------------------|
| id            | int64     | PK, Auto                                  | 组织唯一标识            |
| name          | string    | Unique, Not Empty                         | **组织名称**（创建后不可修改） |
| type          | string    | "Lender" \| "Borrower"                    | 组织类型              |
| status        | string    | "Active" \| "Disabled" \| "Not activated" | 组织状态              |
| owner_id      | int64     | FK → users.id                             | 所有者用户 ID          |
| contact_email | string    | Email format                              | 联系邮箱              |
| contact_phone | string    | Min 8, Max 20                             | 联系电话              |
| contact_name  | string    | Min 2, Max 100                            | 联系人姓名             |
| owner_type    | string    | -                                         | 所有者类型             |
| logo          | string    | Optional                                  | 组织 logo URL       |
| website       | string    | Optional                                  | 网站 URL            |
| tagline       | string    | Optional                                  | 标语                |
| metadata      | jsonb     | Optional                                  | 元数据（Lender 贷款范围等） |
| created_at    | timestamp | Auto                                      | 创建时间              |
| updated_at    | timestamp | Auto                                      | 更新时间              |
| is_del        | int       | Default 0                                 | 软删除标记             |

**Validation Rules**:

- `name`: 创建时必须验证在 Lender 类型中唯一
- `contact_email`: 必须符合 email 格式
- `contact_phone`: 长度 8-20 字符

**State Transitions**:

```
┌──────────────┐  Switch Status  ┌──────────┐
│ Not Activated│ ────────────────►│ Active   │
└──────────────┘                   └──────────┘
     ▲                                   │
     │                                   │ Switch Status
     └───────────────────────────────────┘
                    │
                    ▼
              ┌──────────┐
              │Disabled  │
              └──────────┘
```

**Note**: Admin 创建的组织状态直接为 `Active`，跳过 `Not Activated` 状态。

---

### 2. User（用户）

**Table**: `users`

**Purpose**: 组织的所有者用户账号

**Fields** (Relevant for this feature):

| Field                   | Type      | Constraints                            | Description            |
|-------------------------|-----------|----------------------------------------|------------------------|
| id                      | int64     | PK, Auto                               | 用户唯一标识                 |
| email                   | string    | Unique, Email                          | 用户邮箱（作为 owner 时需验证唯一性） |
| password_hash           | string    | Not Empty                              | 密码哈希                   |
| first_name              | string    | Min 2, Max 100                         | 名                      |
| last_name               | string    | Min 2, Max 100                         | 姓                      |
| type                    | string    | "Lender" \| "Borrower"                 | 用户类型                   |
| status                  | string    | "Active" \| "Inactive" \| "Unverified" | 用户状态                   |
| default_organization_id | int64     | FK → organizations.id                  | 默认组织 ID                |
| created_at              | timestamp | Auto                                   | 创建时间                   |
| updated_at              | timestamp | Auto                                   | 更新时间                   |
| is_del                  | int       | Default 0                              | 软删除标记                  |

**Validation Rules**:

- `email`: 作为 Lender 组织 owner 时必须验证唯一性

---

### 3. OrganizationUser（组织用户关系）

**Table**: `organization_users`

**Purpose**: 用户与组织的多对多关联关系

**Fields**:

| Field           | Type      | Constraints                               | Description |
|-----------------|-----------|-------------------------------------------|-------------|
| id              | int64     | PK, Auto                                  | 关系唯一标识      |
| organization_id | int64     | FK → organizations.id                     | 组织 ID       |
| user_id         | int64     | FK → users.id                             | 用户 ID       |
| first_name      | string    | Min 2, Max 100                            | 名（在组织中的显示）  |
| last_name       | string    | Min 2, Max 100                            | 姓（在组织中的显示）  |
| full_name       | string    | Computed                                  | 全名          |
| status          | string    | "Active" \| "Disabled" \| "Not activated" | 在组织中的状态     |
| created_at      | timestamp | Auto                                      | 创建时间        |
| updated_at      | timestamp | Auto                                      | 更新时间        |
| is_del          | int       | Default 0                                 | 软删除标记       |

**State Synchronization**:

当组织状态切换时，该组织所有 `organization_users` 记录的状态同步更新：

```
Organization.Active   → OrganizationUser.Active
Organization.Disabled → OrganizationUser.Disabled
```

---

### 4. Role（角色）

**Table**: `roles`

**Purpose**: 组织下的角色定义

**Fields** (Relevant for this feature):

| Field           | Type      | Constraints            | Description     |
|-----------------|-----------|------------------------|-----------------|
| id              | int64     | PK, Auto               | 角色唯一标识          |
| name            | string    | Not Empty              | 角色名称（如 "Admin"） |
| type            | string    | "Lender" \| "Borrower" | 角色类型            |
| organization_id | int64     | FK → organizations.id  | 所属组织 ID         |
| created_at      | timestamp | Auto                   | 创建时间            |
| updated_at      | timestamp | Auto                   | 更新时间            |
| is_del          | int       | Default 0              | 软删除标记           |

**Note**: 创建组织时自动创建 "Admin" 角色。

---

## Entity Relationships

```
┌──────────────┐         1:N         ┌──────────────┐
│ Organization │ ──────────────────► │ OrganizationUser│
└──────────────┘                     └──────────────┘
       │                                    │
       │ 1:1 (owner)                       │ N:1
       ▼                                    ▼
┌──────────────┐                    ┌──────────────┐
│    User      │                    │   User       │
└──────────────┘                    └──────────────┘

┌──────────────┐         1:N         ┌──────────────┐
│ Organization │ ──────────────────► │    Role      │
└──────────────┘                     └──────────────┘
```

**Relationship Rules**:

1. 一个 Organization 有一个 Owner（User）
2. 一个 Organization 有多个 OrganizationUser
3. 一个 User 可以属于多个 Organization（通过 OrganizationUser）
4. 一个 Organization 有多个 Role

---

## CRUD Operations Data Flow

### Create Organization

```
Input: OrganizationCreateLenderRequest
  │
  ├─ Validate: organization.name unique (Lender type)
  ├─ Validate: user.email unique (Lender type)
  ├─ Validate: all required fields
  │
  ▼
Transaction Start
  │
  ├─ Create Organization (status=Active)
  ├─ Create User (status=Active)
  ├─ Update organization.owner_id = user.id
  ├─ Create OrganizationUser (status=Active)
  └─ Create Role (name="Admin")
  │
Transaction End
  │
  ▼
Output: OrganizationEntity
```

### Update Organization

```
Input: OrganizationUpdateLenderRequest
  │
  ├─ Validate: organization exists
  ├─ Validate: editable fields only
  │  (name, type, owner_id, status NOT allowed)
  │
  ▼
Update Organization
  │
  ▼
Output: OrganizationEntity
```

### Switch Status

```
Input: OrganizationSwitchStatusRequest
  │
  ├─ Validate: organization exists
  │
  ▼
Transaction Start
  │
  ├─ Update organization.status = new_status
  └─ Batch UPDATE organization_users
        SET status = new_status
        WHERE organization_id = ?
  │
Transaction End
  │
  ▼
Output: Success
```

### Delete Organization

```
Input: OrganizationDeleteRequest (id)
  │
  ├─ Validate: organization exists
  │
  ▼
Transaction Start
  │
  ├─ Delete organization_users (WHERE organization_id = ?)
  ├─ Delete roles (WHERE organization_id = ?)
  ├─ Delete users (WHERE default_organization_id = ?)
  └─ Delete organization
  │
Transaction End
  │
  ▼
Output: Success
```

---

## Metadata Structure

### Lender Metadata

```json
{
  "lender": {
    "loan_amount_range": {
      "min": 100000.0,
      "max": 5000000.0
    },
    "loan_term_range": {
      "min": 12,
      "max": 360
    },
    "loan_rate_range": {
      "min": 3.5,
      "max": 15.0
    },
    "loan_point_range": {
      "min": 0.0,
      "max": 5.0
    }
  }
}
```

---

## Database Indexes

**Recommended Indexes** (for performance):

| Table              | Columns                   | Purpose                                 |
|--------------------|---------------------------|-----------------------------------------|
| organizations      | (name, type)              | Unique constraint for organization name |
| organizations      | status                    | Filter by status                        |
| users              | (email, type)             | Unique constraint for owner email       |
| organization_users | (organization_id, status) | Status switching batch update           |
| organization_users | user_id                   | User lookup                             |
