# Quickstart: 管理后台组织管理 CRUD

**Feature**: [spec.md](./spec.md)
**Branch**: `005-admin-org-crud`
**Date**: 2025-01-12

## 快速开始

本文档提供实现组织管理 CRUD 功能的快速入门指南。

## 前置条件

1. 阅读项目架构文档：`.claude/docs/analysis/architecture.md`
2. 了解后端开发规范：`backend/CLAUDE.md`
3. 熟悉 go-zero 框架基础
4. 已创建功能分支：`005-admin-org-crud`

## 文件结构

### 需要实现的文件

```
app/admin/api/internal/logic/organization/
├── organization_create_lender_logic.go    # 创建组织 API 逻辑
├── organization_detail_logic.go           # 查看详情 API 逻辑
├── organization_update_lender_logic.go    # 编辑组织 API 逻辑
├── organization_switch_status_logic.go    # 切换状态 API 逻辑
└── organization_delete_logic.go           # 删除组织 API 逻辑

app/admin/rpc/internal/logic/organization/
├── create_lender_logic.go                 # 创建组织 RPC 逻辑
├── detail_logic.go                        # 查看详情 RPC 逻辑
├── update_lender_logic.go                 # 编辑组织 RPC 逻辑
├── switch_status_logic.go                 # 切换状态 RPC 逻辑
└── delete_logic.go                        # 删除组织 RPC 逻辑
```

### 参考文件

```
app/admin/api/internal/logic/organization/
└── organization_page_logic.go             # 已实现的分页查询（参考模式）

app/admin/rpc/internal/logic/organization/
└── page_logic.go                          # 已实现的分页查询（参考模式）

domains/zeebee/logic/
└── user_logic.go                          # RegisterWithOrganization 方法（参考创建逻辑）
```

## 实施步骤

### 步骤 1: 创建 Lender 组织

**文件**: `app/admin/rpc/internal/logic/organization/create_lender_logic.go`

**核心逻辑**:

```go
// 参考 user_logic.go:481-576 RegisterWithOrganization 方法

// 关键差异：
// 1. 状态直接设为 Active（不需要激活流程）
// 2. 验证组织名称唯一性（Lender 类型）
// 3. 验证 owner 邮箱唯一性（Lender 类型）
```

**验证清单**:

- [ ] 组织名称在 Lender 类型中唯一
- [ ] Owner 邮箱在 Lender 类型中唯一
- [ ] 所有必填字段已验证
- [ ] 组织状态 = Active
- [ ] 用户状态 = Active
- [ ] OrganizationUser 状态 = Active

---

### 步骤 2: 查看组织详情

**文件**: `app/admin/rpc/internal/logic/organization/detail_logic.go`

**核心逻辑**:

```go
// 使用 OrganizationLogic.GetOrganizationById
// 预加载关联数据（Owner, OrganizationEmails）
```

---

### 步骤 3: 编辑组织

**文件**: `app/admin/rpc/internal/logic/organization/update_lender_logic.go`

**核心逻辑**:

```go
// 注意：组织名称不可修改
// 可编辑字段：
// - ContactEmail, ContactPhone, ContactName
// - Logo, Website, Tagline
// - LenderMetadata
```

**验证清单**:

- [ ] 组织名称不允许修改
- [ ] 只更新允许的字段

---

### 步骤 4: 切换状态

**文件**: `app/admin/rpc/internal/logic/organization/switch_status_logic.go`

**核心逻辑**:

```go
// 使用事务
// 1. 更新 organization.status
// 2. 批量更新 organization_users 状态

// SQL: UPDATE organization_users
//      SET status = ?
//      WHERE organization_id = ?
```

**验证清单**:

- [ ] 状态同步到所有 organization_users
- [ ] 使用事务保证一致性

---

### 步骤 5: 删除组织

**文件**: `app/admin/rpc/internal/logic/organization/delete_logic.go`

**核心逻辑**:

```go
// 使用事务按序删除：
// 1. organization_users
// 2. roles
// 3. users
// 4. organization
```

**验证清单**:

- [ ] 按正确顺序删除（外键约束）
- [ ] 使用事务保证原子性

---

## 编码规范

### 1. 数据转换

所有类型转换使用 `copierx.Trans[T, U]`：

```go
// API Layer
rpcReq, err := copierx.Trans[types.OrganizationCreateLenderRequest, admin.OrganizationCreateLenderRequestData](req)
if err != nil {
    return nil, errorx.GetFailedTransRespErr(err)
}

// RPC Layer
organization, err := copierx.Trans[dto.OrganizationCreateLenderDto, model.Organization](&in.Data)
if err != nil {
    return nil, errorx.GetFailedTransRespErr(err)
}

// 时间字段转换
data, err := copierx.Trans[bo.OrganizationBo, admin.OrganizationEntity](result, copierx.TimeToRFC3339Str())
```

### 2. 错误处理

```go
defer func() {
    if err != nil {
        l.WithContext(l.ctx).WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).Error(err)
    }
}()

if err != nil {
    return nil, errorx.WithStack(err)
}
```

### 3. 事务处理

```go
err = l.svcCtx.Repos.UserRepo.Transaction(l.ctx, func(tx *gorm.DB) error {
    // 业务逻辑
    return nil
})
if err != nil {
    return nil, err
}
```

### 4. 唯一性验证

```go
// 组织名称唯一性
isExisted, err := l.svcCtx.Repos.OrganizationRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
    "name": input.Name,
    "type": constants.UserTypeLender,
})
if isExisted {
    return nil, errorx.New("Organization name already exists")
}

// 邮箱唯一性
isExisted, err := l.svcCtx.Repos.UserRepo.IsExistedByColumns(l.ctx, map[string]interface{}{
    "email": input.Email,
    "type": constants.UserTypeLender,
})
if isExisted {
    return nil, errorx.New("Email already exists")
}
```

---

## 测试

### 单元测试示例

```go
// app/admin/rpc/internal/logic/organization/create_lender_logic_test.go

func TestCreateLenderLogic_CreateLender_Success(t *testing.T) {
    // Arrange
    mockRepo := &MockOrganizationRepo{}
    logic := NewCreateLenderLogic(context.Background(), &svc.ServiceContext{
        Repos: &repo.AllRepos{OrganizationRepo: mockRepo},
    })

    // Act
    resp, err := logic.CreateLender(&admin.OrganizationCreateLenderRequest{
        Data: &admin.OrganizationCreateLenderRequestData{
            Organization: &admin.OrganizationCreateLenderData{
                Name: "Test Lender Org",
            },
        },
    })

    // Assert
    assert.Nil(t, err)
    assert.NotNil(t, resp)
}
```

---

## 常见问题

### Q: 为什么组织名称不可修改？

A: 组织名称是业务的唯一标识符，修改后会破坏审计追踪和业务一致性。

### Q: 状态切换为什么需要批量更新？

A: 组织状态变化时，所有成员在该组织中的状态应该同步，批量更新可以保证性能（1000+ 用户 < 3 秒）。

### Q: 删除顺序为什么很重要？

A: 需要遵循数据库外键约束。如果先删除组织，关联的 organization_users 会因为外键约束而删除失败。

---

## 下一步

1. 运行 `/speckit.tasks` 生成详细的任务清单
2. 按照任务清单依次实现每个功能
3. 编写单元测试和集成测试
4. 运行测试验证功能正确性
