# Research: 管理后台组织管理 CRUD

**Feature**: [spec.md](./spec.md)
**Date**: 2025-01-12
**Status**: Complete

## Overview

本文档记录了实现管理后台组织 CRUD 功能的技术研究和决策。所有决策基于现有代码模式、项目宪章原则和 go-zero 最佳实践。

## Business Rules

### 组织名称约束

- **不可修改**: 组织名称在创建后不能修改
- **唯一性验证**: 创建时必须检查组织名称的唯一性
- **验证范围**: 同类型（Lender）组织中名称不能重复

### Owner 邮箱约束

- **唯一性验证**: 创建时必须检查 owner 邮箱的唯一性
- **验证范围**: 同类型（Lender）用户中邮箱不能重复

## Technical Decisions

### 1. 组织创建逻辑

**Decision**: 参考 `user_logic.go` 中的 `RegisterWithOrganization` 方法，但不使用激活流程

**Rationale**:

- `RegisterWithOrganization` 提供了完整的创建组织、用户和关联关系的事务处理模式
- Admin 创建的 lender 组织不需要激活，状态直接设为 `Active`
- 现有代码已处理组织名称唯一性验证、用户邮箱验证等

**Implementation Pattern**:

```go
// 参考代码位置: domains/zeebee/logic/user_logic.go:481-576
// 关键差异：
// 1. organization.Status = constants.OrganizationStatusActive（而非 NotActivated）
// 2. user.Status = constants.UserStatusActive（而非 Unverified）
// 3. organizationUser.Status = constants.OrganizationUserStatusActive（而非 NotActivated）

// 唯一性验证：
isOrgExisted, err := l.repos.OrganizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
    "name": organizationInput.Name,
    "type": constants.UserTypeLender,
})
if isOrgExisted {
    return nil, errorx.New("Organization Name Already Exists")
}

isUserExisted, err := l.repos.UserRepo.IsExistedByColumns(ctx, map[string]interface{}{
    "email": userInput.Email,
    "type": constants.UserTypeLender,
})
if isUserExisted {
    return nil, errorx.New("Owner Email Already Exists")
}
```

**Alternatives Considered**:

- 直接使用 `RegisterWithOrganization` + 激活流程：增加了不必要的复杂度
- 重新实现创建逻辑：可能遗漏现有的事务处理和验证逻辑

---

### 2. 编辑组织限制

**Decision**: 编辑接口只允许修改特定字段，组织名称不可修改

**Rationale**:

- 业务规则要求组织名称不可修改
- 组织名称是业务标识符，修改会影响业务追踪和审计

**Editable Fields**:

- ContactEmail（联系邮箱）
- ContactPhone（联系电话）
- ContactName（联系人姓名）
- Logo（组织 logo）
- Website（网站）
- Tagline（标语）
- LenderMetadata（贷款方元数据）

**Non-Editable Fields**:

- Name（组织名称）
- Type（组织类型）
- OwnerId（所有者 ID）
- Status（状态通过单独接口切换）

**Implementation Pattern**:

```go
// API 定义中不包含 name 字段
type OrganizationUpdateLenderRequest {
    Id int64 `path:"id" validate:"required,min=1"`
    // name 字段不包含在内
    ContactEmail string `json:"contact_email" validate:"required,email"`
    ContactPhone string `json:"contact_phone" validate:"required,min=8,max=20"`
    ContactName string `json:"contact_name" validate:"required,min=2,max=100"`
    Logo string `json:"logo,optional"`
    Website string `json:"website,optional"`
    Tagline string `json:"tagline,optional"`
    LenderMetadata *OrganizationLenderMetadataEntity `json:"lender_metadata,optional"`
}
```

---

### 3. 状态切换实现

**Decision**: 在 RPC Logic 层实现状态切换，使用 GORM 事务批量更新 organization_user

**Rationale**:

- 状态切换涉及多表更新（organizations + organization_users）
- 需要保证数据一致性，必须使用事务
- 批量更新使用 GORM 的 `Where().Updates()` 可以高效处理大量用户

**Implementation Pattern**:

```go
// 切换为 Active 时
organization.Status = constants.OrganizationStatusActive
// 批量更新关联用户
repos.OrganizationUserRepo.UpdateStatusByOrganizationId(ctx, organizationId, constants.OrganizationUserStatusActive)

// 切换为 Disabled 时
organization.Status = constants.OrganizationStatusDisabled
// 批量更新关联用户
repos.OrganizationUserRepo.UpdateStatusByOrganizationId(ctx, organizationId, constants.OrganizationUserStatusDisabled)
```

**Alternatives Considered**:

- 逐个更新用户状态：性能问题，当组织有大量用户时（1000+）会超时
- 使用异步任务：增加了复杂度，且状态切换应该是同步操作

---

### 4. 删除组织级联逻辑

**Decision**: 使用事务按序删除：organization_user → roles → users → organization

**Rationale**:

- 数据库外键约束要求特定删除顺序
- 事务保证原子性，任何步骤失败都能回滚
- 现有 Repository 已提供 `DeleteWithTx` 方法

**Implementation Pattern**:

```go
repos.UserRepo.Transaction(ctx, func(tx *gorm.DB) error {
    // 1. 删除 organization_user 关联
    repos.OrganizationUserRepo.DeleteByOrganizationIdWithTx(ctx, tx, organizationId)
    // 2. 删除组织角色
    repos.RoleRepo.DeleteByOrganizationIdWithTx(ctx, tx, organizationId)
    // 3. 删除组织用户
    repos.UserRepo.DeleteByOrganizationIdWithTx(ctx, tx, organizationId)
    // 4. 删除组织
    repos.OrganizationRepo.DeleteWithTx(ctx, tx, organizationId)
    return nil
})
```

**Alternatives Considered**:

- 使用数据库级联删除：降低了应用层控制，可能产生意外的数据删除
- 软删除（标记 is_del）：现有架构使用软删除模式，符合需求

---

### 5. API/RPC 数据流

**Decision**: 遵循现有模式：API Handler → API Logic → RPC Call → RPC Logic → Repository

**Rationale**:

- 项目宪章强制要求 API/RPC 分离
- API 层负责验证和格式转换
- RPC 层负责业务逻辑和数据访问

**Implementation Pattern**:

```go
// API Logic (app/admin/api/internal/logic/organization/xxx_logic.go)
func (l *OrganizationCreateLenderLogic) OrganizationCreateLender(req *types.OrganizationCreateLenderRequest) (*types.OrganizationCreateLenderResponse, error) {
    // 1. 类型转换: types → pb
    rpcReq, err := copierx.Trans[types.OrganizationCreateLenderRequest, admin.OrganizationCreateLenderRequestData](req)
    // 2. 调用 RPC
    rpcResp, err := l.svcCtx.OrganizationRpcClient.CreateLender(l.ctx, rpcReq)
    // 3. 类型转换: pb → types
    resp, err := copierx.Trans[admin.OrganizationCreateLenderResponseData, types.OrganizationCreateLenderResponse](rpcResp.Data)
    return resp, err
}

// RPC Logic (app/admin/rpc/internal/logic/organization/xxx_logic.go)
func (l *CreateLenderLogic) CreateLender(in *admin.OrganizationCreateLenderRequest) (*admin.OrganizationCreateLenderResponse, error) {
    // 业务逻辑实现
}
```

**Alternatives Considered**:

- API 层直接访问数据库：违反宪章原则，破坏架构分层
- RPC 层处理验证：混淆职责边界，API 层应该处理输入验证

---

### 6. 数据转换标准

**Decision**: 所有类型转换使用 `copierx.Trans[T, U]` 泛型方法

**Rationale**:

- 项目宪章强制要求
- 统一转换逻辑，减少错误
- 支持时间字段自动转换（`copierx.TimeToRFC3339Str()`）

**Implementation Pattern**:

```go
// 单个对象转换
data, err := copierx.Trans[bo.OrganizationBo, admin.OrganizationEntity](result, copierx.TimeToRFC3339Str())

// 切片转换
data, err := copierx.Trans[[]*model.Organization, []*admin.OrganizationEntity](&results, copierx.TimeToRFC3339Str())

// 注意：泛型参数 T 和 U 不包括指针符号 *
//       src 参数本身是指针类型 *T
```

---

### 7. 错误处理

**Decision**: 使用 `common/errorx/` 包进行统一错误处理

**Rationale**:

- 项目宪章强制要求
- 统一错误格式和日志记录
- 支持错误堆栈追踪

**Implementation Pattern**:

```go
defer func() {
    if err != nil {
        l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
    }
}()

// 返回错误
if err != nil {
    return nil, errorx.WithStack(err)  // 包装错误并堆栈追踪
}
```

---

## Dependencies Analysis

### 现有可复用组件

| Component            | Location                                        | Usage                               |
|----------------------|-------------------------------------------------|-------------------------------------|
| OrganizationLogic    | `domains/zeebee/logic/organization_logic.go`    | 组织业务逻辑（分页、查询等）                      |
| UserLogic            | `domains/zeebee/logic/user_logic.go`            | 用户业务逻辑（RegisterWithOrganization 参考） |
| OrganizationUserRepo | `domains/zeebee/repo/organization_user_repo.go` | 组织用户关系数据访问                          |
| OrganizationRepo     | `domains/zeebee/repo/organization_repo.go`      | 组织数据访问                              |
| UserRepo             | `domains/zeebee/repo/user_repo.go`              | 用户数据访问                              |
| BaseLogic            | `app/admin/rpc/internal/logic/base_logic.go`    | RPC 基础逻辑                            |

### 需要新增的 Repository 方法

| Method                       | Purpose    | Signature                                                                      |
|------------------------------|------------|--------------------------------------------------------------------------------|
| UpdateStatusByOrganizationId | 批量更新组织用户状态 | `UpdateStatusByOrganizationId(ctx, organizationId int64, status string) error` |
| DeleteByOrganizationIdWithTx | 删除组织所有用户   | `DeleteByOrganizationIdWithTx(ctx, tx, organizationId int64) error`            |
| DeleteByOrganizationIdWithTx | 删除组织所有角色   | `DeleteByOrganizationIdWithTx(ctx, tx, organizationId int64) error`            |
| DeleteByOrganizationIdWithTx | 删除组织所有用户关系 | `DeleteByOrganizationIdWithTx(ctx, tx, organizationId int64) error`            |

---

## Performance Considerations

### 1. 状态切换性能

**Requirement**: 组织状态切换在 3 秒内完成，包括所有关联用户状态更新

**Strategy**:

- 使用批量更新（Single SQL UPDATE with WHERE）而非逐条更新
- 预期 SQL: `UPDATE organization_users SET status = ? WHERE organization_id = ?`
- 对于 1000+ 用户的组织，单条 SQL 更新比循环更新快 100x

### 2. 删除操作性能

**Requirement**: 删除操作在 5 秒内完成

**Strategy**:

- 使用事务保证一致性
- 按 foreign key 顺序删除避免约束检查开销
- 考虑后台软删除（标记 is_del = 1）而非硬删除

### 3. 分页查询性能

**Strategy**:

- 使用现有的分页查询模式
- 添加适当的数据库索引（organization_id, status）

---

## Security Considerations

### 1. 权限验证

**Requirement**: Admin 用户必须有权访问组织管理功能

**Implementation**:

- API 层使用 JWT 中间件验证用户身份
- RPC 层验证用户类型为 Admin

### 2. 输入验证

**Requirement**: 所有输入必须验证

**Implementation**:

- API 层使用 `validate` tag 进行声明式验证
- 使用 `common/validator/` 进行自定义验证逻辑

### 3. 唯一性验证

**Requirement**:

- 组织名称在 Lender 类型中唯一
- Owner 邮箱在 Lender 类型中唯一

**Implementation**:

```go
// 创建组织前验证
isOrgNameExisted, _ := l.repos.OrganizationRepo.IsExistedByColumns(ctx, map[string]interface{}{
    "name": organizationInput.Name,
    "type": constants.UserTypeLender,
})
if isOrgNameExisted {
    return nil, errorx.New("Organization name already exists")
}

isEmailExisted, _ := l.repos.UserRepo.IsExistedByColumns(ctx, map[string]interface{}{
    "email": userInput.Email,
    "type": constants.UserTypeLender,
})
if isEmailExisted {
    return nil, errorx.New("Owner email already exists")
}
```

### 4. SQL 注入防护

**Implementation**:

- 使用 GORM 参数化查询
- 避免拼接 SQL 字符串

---

## Testing Strategy

### 单元测试

- 覆盖所有 RPC Logic 方法
- Mock Repository 层
- 测试正常流程和异常场景
- **重点测试唯一性验证逻辑**

### 集成测试

- 测试完整的数据流：API → RPC → Repository → Database
- 使用测试数据库和事务回滚
- 测试并发场景（状态切换）
- **测试组织名称重复和邮箱重复场景**

### 测试覆盖率目标

- 80%+ 代码覆盖率
- 100% 核心业务逻辑覆盖

---

## Open Questions (Resolved)

### Q1: Admin 创建组织时是否需要发送激活邮件？

**Answer**: 不需要。Admin 创建的组织直接激活（Status = Active），无需邮件验证流程。

**Rationale**: 功能规格明确说明 "admin创建的lender组织不需要激活"。

### Q2: 组织名称可以修改吗？

**Answer**: 不可以。组织名称在创建后不可修改，是组织的固定标识符。

**Rationale**: 业务规则明确要求，组织名称不可修改以保证业务追踪和审计的准确性。

### Q3: 编辑组织时需要验证唯一性吗？

**Answer**: 不需要。因为唯一字段（组织名称）不可编辑，所以编辑时不需要唯一性验证。

### Q4: 删除组织时如何处理关联的业务数据（如贷款记录）？

**Answer**: 规格只要求删除用户、角色和组织关系。业务数据（贷款记录等）不在本次功能范围。

**Future**: 可能需要软删除或归档功能来处理历史业务数据。

---

## Next Steps

1. **Phase 1**: 生成数据模型和 API 契约
2. **Phase 2**: 生成详细的任务清单（`/speckit.tasks`）
3. **Implementation**: 按照 tasks.md 执行开发
