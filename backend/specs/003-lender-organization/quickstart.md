# Developer Quick Start: Lender Organization Management

**Feature**: 003-lender-organization
**Date**: 2025-01-09
**Branch**: `003-lender-organization`

## Overview

本指南帮助开发者快速理解和实施 lender 端组织管理功能。

## 功能概述

实现三个核心 API 端点：

1. **GET /organization/profile** - 获取当前用户所在组织信息
2. **PUT /organization/profile** - 更新组织信息（仅 owner）
3. **POST /organization** - 创建新组织

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                         API Layer                           │
│  (Request Validation, RPC Orchestration, Type Conversion)   │
└──────────────────────────┬──────────────────────────────────┘
                           │ gRPC
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                         RPC Layer                           │
│     (Business Logic, Permission Check, Database Operations)  │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                             │
│        (Models, Repositories, Business Logic)                │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Database Layer                           │
│              (PostgreSQL - organizations, etc.)              │
└─────────────────────────────────────────────────────────────┘
```

## 核心模式

### 1. 获取用户上下文

在 RPC Logic 中获取当前登录用户信息：

```go
// 从 gRPC metadata 中获取 JWT claims
claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
if err != nil {
    return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
}

// claim 包含:
// - claim.UserId: 用户ID
// - claim.OrganizationId: 当前组织ID
// - claim.Username: 用户名
// - claim.UserType: 用户类型
```

### 2. 查询组织信息

使用业务逻辑层获取组织详情：

```go
// 获取组织详情（包含预加载的关联数据）
organization, err := l.baseLogic.OrganizationLogic.Detail(
    l.ctx,
    claim.OrganizationId,
    constants.UserTypeLender,
)
if err != nil {
    return nil, errorx.Wrap(err, "Failed To Get Organization")
}

// organization 是 *bo.OrganizationBo 类型
// 包含 OrganizationEmails 等关联数据
```

### 3. 数据类型转换

使用 `copierx.Trans` 进行类型转换：

```go
import "zeebee-backend/common/copierx"

// 业务对象 → Proto 实体
resp, err = copierx.Trans[*bo.OrganizationBo, *lender.OrganizationEntity](
    organization,
    copierx.TimeToRFC3339Str(),  // 时间字段转换
)
if err != nil {
    return nil, errorx.GetFailedTransRespErr(err)
}

// Proto → API 类型
apiResp, err = copierx.Trans[*lender.OrganizationEntity, *types.OrganizationEntity](
    rpcResp.Data,
    copierx.TimeToRFC3339Str(),
)
```

### 4. 错误处理

统一的错误处理模式：

```go
defer func() {
    if err != nil {
        l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
    }
}()

// 使用 errorx.Wrap 包装错误
if err != nil {
    return nil, errorx.Wrap(err, "操作描述")
}

// 使用 errorx.New 创建验证错误
if validationFailed {
    return nil, errorx.NewDefaultError("清晰的错误消息")
}
```

## 实现步骤

### Phase 1: RPC Layer (业务逻辑层)

#### 1.1 Profile Logic

**文件**: `app/lender/rpc/internal/logic/organization/profile_logic.go`

```go
func (l *ProfileLogic) Profile(in *lender.EmptyRequest) (resp *lender.OrganizationEntity, err error) {
    defer func() {
        if err != nil {
            l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
        }
    }()

    // 1. 获取用户上下文
    claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
    if err != nil {
        return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
    }

    if claim.OrganizationId == 0 {
        return nil, errorx.NewDefaultError("User Not Belong To Any Organization")
    }

    // 2. 获取组织信息
    organization, err := l.baseLogic.OrganizationLogic.Detail(
        l.ctx,
        claim.OrganizationId,
        constants.UserTypeLender,
    )
    if err != nil {
        return nil, errorx.Wrap(err, "Failed To Get Organization")
    }

    // 3. 类型转换
    resp, err = copierx.Trans[*bo.OrganizationBo, *lender.OrganizationEntity](
        organization,
        copierx.TimeToRFC3339Str(),
    )
    if err != nil {
        return nil, errorx.GetFailedTransRespErr(err)
    }

    return resp, nil
}
```

#### 1.2 UpdateProfile Logic

**文件**: `app/lender/rpc/internal/logic/organization/update_profile_logic.go`

```go
func (l *UpdateProfileLogic) UpdateProfile(in *lender.OrganizationUpdateProfileRequest) (resp *lender.OrganizationEntity, err error) {
    defer func() {
        if err != nil {
            l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
        }
    }()

    // 1. 获取用户上下文
    claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
    if err != nil {
        return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
    }

    // 2. 获取组织并验证权限
    organization, err := l.baseLogic.OrganizationLogic.Detail(
        l.ctx,
        claim.OrganizationId,
        constants.UserTypeLender,
    )
    if err != nil {
        return nil, errorx.Wrap(err, "Failed To Get Organization")
    }

    // 3. 验证是否为 owner
    if organization.OwnerId != claim.UserId {
        return nil, errorx.NewDefaultError("Only Owner Can Update Organization Information")
    }

    // 4. 构建更新数据
    updateData := make(map[string]interface{})
    if len(in.Data.ContactEmail) > 0 {
        updateData["contact_email"] = in.Data.ContactEmail
    }
    if len(in.Data.ContactPhone) > 0 {
        updateData["contact_phone"] = in.Data.ContactPhone
    }
    if len(in.Data.ContactName) > 0 {
        updateData["contact_name"] = in.Data.ContactName
    }
    if len(in.Data.Logo) > 0 {
        updateData["logo"] = in.Data.Logo
    }
    if len(in.Data.Website) > 0 {
        updateData["website"] = in.Data.Website
    }
    if len(in.Data.Tagline) > 0 {
        updateData["tagline"] = in.Data.Tagline
    }
    if in.Data.LenderMetadata != nil {
        updateData["metadata"] = in.Data.LenderMetadata
    }

    // 5. 执行更新
    err = l.svcCtx.OrganizationRepo.UpdateColumnsById(l.ctx, claim.OrganizationId, updateData)
    if err != nil {
        return nil, errorx.Wrap(err, "Failed To Update Organization")
    }

    // 6. 返回更新后的组织信息
    return l.Profile(&lender.EmptyRequest{})
}
```

#### 1.3 Create Logic

**文件**: `app/lender/rpc/internal/logic/organization/create_logic.go`

```go
func (l *CreateLogic) Create(in *lender.OrganizationCreateRequest) (resp *lender.OrganizationCreateResponse, err error) {
    defer func() {
        if err != nil {
            l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
        }
    }()

    // 1. 获取用户上下文（创建者将成为 owner）
    claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
    if err != nil {
        return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
    }

    // 2. 转换请求数据为模型
    organization := &model.Organization{}
    organization, err = copierx.Trans[lender.OrganizationCreateRequestData, *model.Organization](
        in.Data,
    )
    if err != nil {
        return nil, errorx.GetFailedTransRespErr(err)
    }

    // 3. 设置默认值
    organization.Status = constants.OrganizationStatusInactive
    organization.OwnerType = constants.OrganizationOwnerTypeThirdParty
    organization.Type = constants.UserTypeLender

    // 4. 使用事务创建组织和关联
    err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
        // 创建组织
        organization, err = l.svcCtx.OrganizationRepo.CreateWithTx(l.ctx, tx, organization, false)
        if err != nil {
            return errorx.Wrap(err, "Failed To Create Organization")
        }

        // 更新 owner
        err = l.svcCtx.OrganizationRepo.UpdateOwnerIdWithTx(l.ctx, tx, organization.Id, claim.UserId)
        if err != nil {
            return errorx.Wrap(err, "Failed To Update Organization Owner")
        }

        // 创建组织用户关联
        adminRole := &model.Role{
            Name:           constants.AdminRoleName,
            Type:           constants.UserTypeLender,
            OrganizationId: organization.Id,
        }

        organizationUser := &model.OrganizationUser{
            OrganizationId: organization.Id,
            UserId:         claim.UserId,
            Status:         constants.OrganizationUserStatusActive,
            FullName:       claim.Username,
            Roles:          []*model.Role{adminRole},
        }

        _, err = l.svcCtx.OrganizationUserRepo.CreateWithTx(l.ctx, tx, organizationUser, true)
        if err != nil {
            return errorx.Wrap(err, "Failed To Create Organization User")
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    // 5. 返回创建的组织信息
    resultOrganization, err := copierx.Trans[*model.Organization, *lender.OrganizationEntity](
        organization,
        copierx.TimeToRFC3339Str(),
    )
    if err != nil {
        return nil, errorx.GetFailedTransRespErr(err)
    }

    return &lender.OrganizationCreateResponse{
        Data: resultOrganization,
    }, nil
}
```

### Phase 2: API Layer (请求处理层)

API 层逻辑相对简单，主要负责调用 RPC 并转换数据类型：

```go
func (l *OrganizationProfileLogic) OrganizationProfile(req *types.EmptyRequest) (resp *types.OrganizationEntity, err error) {
    // 调用 RPC
    rpcResp, err := l.svcCtx.OrganizationRpcClient.Profile(l.ctx, &lender.EmptyRequest{})
    if err != nil {
        return nil, err
    }

    // 转换返回类型
    resp, err = copierx.Trans[*lender.OrganizationEntity, *types.OrganizationEntity](
        rpcResp,
        copierx.TimeToRFC3339Str(),
    )
    if err != nil {
        return nil, errorx.GetFailedTransRespErr(err)
    }

    return resp, nil
}
```

## 测试

### 单元测试 (RPC Layer)

```go
func TestProfileLogic_Profile(t *testing.T) {
    // Setup test context with mock user claim
    // Test scenarios:
    // 1. Success - user with organization
    // 2. Error - user without organization
    // 3. Error - invalid organization ID
}
```

### 集成测试 (API Layer)

```go
func TestOrganizationProfileAPI(t *testing.T) {
    // Setup test server
    // Create authenticated request with valid JWT
    // Verify response structure and data
}
```

## 常见问题

### Q: 如何处理用户不属于任何组织的情况？

A: 在 Profile Logic 中检查 `claim.OrganizationId` 是否为 0：

```go
if claim.OrganizationId == 0 {
    return nil, errorx.NewDefaultError("User Not Belong To Any Organization")
}
```

### Q: 如何验证用户是否为组织 owner？

A: 获取组织后比较 `OwnerId` 与用户的 `UserId`：

```go
if organization.OwnerId != claim.UserId {
    return nil, errorx.NewDefaultError("Only Owner Can Update Organization Information")
}
```

### Q: 创建组织时如何设置正确的默认值？

A: 在创建前手动设置字段：

```go
organization.Status = constants.OrganizationStatusInactive
organization.OwnerType = constants.OrganizationOwnerTypeThirdParty
organization.Type = constants.UserTypeLender
```

### Q: 如何确保数据一致性？

A: 使用事务处理多个数据库操作：

```go
err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
    // Multiple operations here
    // All succeed or all fail
    return nil
})
```

## 相关文件

### Models
- `domains/zeebee/model/organization.go`
- `domains/zeebee/model/organization_user.go`
- `domains/zeebee/model/user.go`

### Repositories
- `domains/zeebee/repo/organization.go`
- `domains/zeebee/repo/organization_user.go`

### Business Logic
- `domains/zeebee/logic/organization.go`
- `domains/zeebee/logic/user_logic.go`

### Constants
- `domains/zeebee/constants/organization.go`

### Utilities
- `common/copierx/trans.go` - 类型转换
- `common/errorx/` - 错误处理
- `common/utils/ctx.go` - 上下文工具

## 参考资料

- Project Constitution: `.specify/memory/constitution.md`
- Development Standards: `CLAUDE.md`
- Feature Specification: `specs/003-lender-organization/spec.md`
- Implementation Plan: `specs/003-lender-organization/plan.md`
