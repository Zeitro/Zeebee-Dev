# Technical Research: Lender Organization Management

**Feature**: 003-lender-organization
**Date**: 2025-01-09
**Status**: Complete

## Overview

本文档记录了 lender 端组织管理功能开发过程中的技术调研和决策。

## User Context Access Pattern

### Problem
如何在 gRPC RPC 服务中获取当前登录用户的信息？

### Research
调研了现有代码库中的实现模式：
- `app/lender/rpc/internal/logic/user/user_info_logic.go` - 使用 `baseLogic.UserLogic.GetUserClaimFromCtx()`
- `domains/zeebee/logic/user_logic.go` - 实现了从 gRPC metadata 中提取 JWT token 的逻辑

### Decision
**选择**: 使用 `l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)`

**实现原理**:
1. 从 gRPC metadata 中提取 token（key: `authorization`）
2. 使用 `utils.ParseToken()` 解析 JWT
3. 返回 `utils.JwtUser` 结构体，包含：
   - UserId: 用户ID
   - OrganizationId: 当前组织ID
   - Username: 用户名
   - UserType: 用户类型

**代码示例**:
```go
claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
if err != nil {
    return nil, errorx.Wrap(err, "Failed To Get User Claim From Context")
}
// Use claim.UserId, claim.OrganizationId, etc.
```

### Alternatives Considered
- 直接从 context 中提取：❌ 重复实现，不符合 DRY 原则
- 通过 API 层传递用户ID：❌ 破坏 RPC 层独立性

## Organization Query Pattern

### Problem
如何获取组织详情信息，包括关联数据？

### Research
调研了 OrganizationLogic 和 OrganizationRepo 的实现：
- `domains/zeebee/logic/organization.go` - 提供业务逻辑层
- `domains/zeebee/repo/organization.go` - 提供数据访问层
- 参考 borrower 端的实现模式

### Decision
**选择**: 使用 `l.baseLogic.OrganizationLogic.Detail(ctx, organizationId, userType)`

**实现原理**:
1. 调用 `OrganizationRepo.GetOrganization()` 获取组织数据
2. 自动预加载关联数据（OrganizationEmails）
3. 返回业务对象 `bo.OrganizationBo`

**代码示例**:
```go
organization, err := l.baseLogic.OrganizationLogic.Detail(
    l.ctx,
    claim.OrganizationId,
    constants.UserTypeLender,
)
if err != nil {
    return nil, errorx.Wrap(err, "Failed To Get Organization")
}
```

### Alternatives Considered
- 直接使用 Repository：❌ 绕过业务逻辑层，不符合分层架构
- 使用 GORM 原生查询：❌ 重复实现，不利用现有代码

## Data Conversion Pattern

### Problem
如何在不同数据类型之间进行转换（Model → BO → Proto/Entity）？

### Research
项目使用自定义的 `copierx` 包进行类型转换，它封装了 `copier` 库并提供了类型安全支持。

### Decision
**选择**: 使用 `copierx.Trans[SourceType, TargetType](source, copierx.TimeToRFC3339Str())`

**类型转换规则**:
- BO (Business Object) → Proto Entity
- Model → Proto Entity
- 时间字段使用 `copierx.TimeToRFC3339Str()` 转换为 RFC3339 字符串

**代码示例**:
```go
// 单个对象转换
resp, err = copierx.Trans[*bo.OrganizationBo, *lender.OrganizationEntity](
    organization,
    copierx.TimeToRFC3339Str(),
)
if err != nil {
    return nil, errorx.GetFailedTransRespErr(err)
}

// 切片转换
data, err = copierx.Trans[[]*model.OrganizationUser, []*lender.OrganizationUserEntity](
    &results,
    copierx.TimeToRFC3339Str(),
)
```

### Alternatives Considered
- 手动映射字段：❌ 容易出错，代码冗长
- 使用标准库 copier：❌ 缺少类型安全和便利方法

## Organization Creation Pattern

### Problem
如何创建新组织并自动关联用户作为 owner？

### Research
调研了用户注册流程中的组织创建逻辑：
- `domains/zeebee/logic/user_logic.go` - `RegisterWithOrganization()` 方法
- 展示了完整的事务处理流程

### Decision
**选择**: 参考注册流程，使用事务确保数据一致性

**实现要点**:
1. 创建组织时设置：
   - Status = `constants.OrganizationStatusInactive`
   - OwnerType = `constants.OrganizationOwnerTypeThirdParty`
   - Type = `constants.UserTypeLender`
2. 使用事务确保：
   - 组织创建成功
   - OwnerId 正确更新
   - OrganizationUser 关联创建
   - 用户角色正确分配

**代码框架**:
```go
err = l.svcCtx.OrganizationRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
    // Create organization
    organization.Status = constants.OrganizationStatusInactive
    organization.OwnerType = constants.OrganizationOwnerTypeThirdParty
    organization, err = l.svcCtx.OrganizationRepo.CreateWithTx(l.ctx, tx, organization, false)
    if err != nil {
        return err
    }

    // Update owner
    err = l.svcCtx.OrganizationRepo.UpdateOwnerIdWithTx(l.ctx, tx, organization.Id, claim.UserId)
    if err != nil {
        return err
    }

    // Create organization user association
    organizationUser := &model.OrganizationUser{
        OrganizationId: organization.Id,
        UserId:         claim.UserId,
        Status:         constants.OrganizationUserStatusActive,
        // ... other fields
        Roles:          adminRoles,
    }
    _, err = l.svcCtx.OrganizationUserRepo.CreateWithTx(l.ctx, tx, organizationUser, true)
    return err
})
```

## Permission Validation Pattern

### Problem
如何验证用户是否有权限更新组织信息？

### Decision
**选择**: 在更新前验证用户是否为组织的 owner

**实现要点**:
1. 从 context 获取用户信息
2. 获取组织详情
3. 比较 `organization.OwnerId` 与 `claim.UserId`
4. 如果不匹配，返回权限错误

**代码示例**:
```go
claim, err := l.baseLogic.UserLogic.GetUserClaimFromCtx(l.ctx)
if err != nil {
    return nil, err
}

organization, err := l.baseLogic.OrganizationLogic.Detail(l.ctx, claim.OrganizationId, constants.UserTypeLender)
if err != nil {
    return nil, err
}

if organization.OwnerId != claim.UserId {
    return nil, errorx.NewDefaultError("Only Owner Can Update Organization Information")
}
```

## Error Handling Pattern

### Decision
遵循项目错误处理规范：

**最佳实践**:
```go
defer func() {
    if err != nil {
        l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
    }
}()

// Use errorx for consistent error wrapping
if err != nil {
    return nil, errorx.Wrap(err, "Descriptive Context")
}

// Use errorx.New for validation errors
if invalidCondition {
    return nil, errorx.NewDefaultError("Clear Error Message")
}
```

## Summary

| Aspect            | Pattern/Tool                          | Source Location                         |
|-------------------|---------------------------------------|-----------------------------------------|
| User Context      | `baseLogic.UserLogic.GetUserClaimFromCtx()` | `domains/zeebee/logic/user_logic.go:49`  |
| Organization Query | `baseLogic.OrganizationLogic.Detail()` | `domains/zeebee/logic/organization.go:27` |
| Data Conversion   | `copierx.Trans[T, U]()`               | `common/copierx/trans.go`                |
| Transaction       | `Repo.Transaction()`                   | `common/gormx/base_repo.go:53`           |
| Error Handling    | `errorx.Wrap()`, `errorx.New()`        | `common/errorx/`                         |

## References

- `app/lender/rpc/internal/logic/user/user_info_logic.go` - User context retrieval example
- `domains/zeebee/logic/user_logic.go` - UserLogic.GetUserClaimFromCtx implementation
- `domains/zeebee/logic/organization.go` - OrganizationLogic.Detail implementation
- `app/borrower/rpc/internal/logic/organization/info_logic.go` - Similar implementation pattern
