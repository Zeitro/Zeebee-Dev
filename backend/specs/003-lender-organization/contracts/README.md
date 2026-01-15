# API Contracts

本目录包含 lender 端组织管理功能的 API 契约定义。

## 现有契约文件

API 契约已经在代码库中定义：

### Organization API
**文件**: `app/lender/api/desc/organization.api`

**定义的端点**:

1. **Profile - 获取当前组织信息**
   - 路径: `GET /organization/profile`
   - 请求: `EmptyRequest`
   - 响应: `OrganizationEntity`

2. **UpdateProfile - 更新组织信息**
   - 路径: `PUT /organization/profile`
   - 请求: `OrganizationUpdateProfileRequest`
   - 响应: `OrganizationEntity`

3. **Create - 创建新组织**
   - 路径: `POST /organization`
   - 请求: `OrganizationCreateRequest`
   - 响应: `OrganizationEntity`

### RPC Service Definition
**文件**: `app/lender/rpc/pb/lender/lender.proto`

**定义的服务**:

```protobuf
service Organization {
    rpc Profile(EmptyRequest) returns (OrganizationEntity);
    rpc UpdateProfile(OrganizationUpdateProfileRequest) returns (OrganizationEntity);
    rpc Create(OrganizationCreateRequest) returns (OrganizationCreateResponse);
}
```

## 数据类型

### OrganizationEntity
```go
type OrganizationEntity {
    Id int64
    Name string
    Status string
    OwnerId int64
    Address AddressEntity
    ContactEmail string
    ContactPhone string
    ContactName string
    OwnerType string
    Type string
    Logo string
    Website string
    Tagline string
    Metadata *OrganizationMetadataEntity
    CreatedAt string
    UpdatedAt string
    OrganizationEmails []*OrganizationEmailEntity
    Owner *UserEntity
}
```

### OrganizationUpdateProfileRequest
```go
type OrganizationUpdateProfileRequest {
    ContactEmail string  // required, email
    ContactPhone string  // required, min=8, max=20
    ContactName string   // required, min=2, max=50
    Logo string         // optional
    Website string      // optional
    Tagline string      // optional
    LenderMetadata *OrganizationLenderMetadataEntity  // optional
}
```

### OrganizationCreateRequest
```go
type OrganizationCreateRequest {
    Name string  // required, min=2, max=50
    ContactEmail string  // required, email
    ContactPhone string  // required, min=8, max=20
    ContactName string   // required, min=2, max=50
    Logo string         // optional
    Website string      // optional
    Tagline string      // optional
    LenderMetadata *OrganizationLenderMetadataEntity  // optional
}
```

## 验证规则

所有输入验证通过 go-zero 的 validate 标签实现：

- `validate:"required"` - 必填字段
- `validate:"email"` - 邮箱格式验证
- `validate:"min=8,max=20"` - 长度范围验证
- `validate:"required,min=2,max=50"` - 必填且长度范围

## 错误响应

所有端点返回统一的错误格式：

```json
{
    "code": 0,
    "msg": "error message",
    "data": null
}
```

常见错误码：
- `401` - 未认证（JWT token 无效或过期）
- `403` - 权限不足（非 owner 尝试更新）
- `400` - 请求参数验证失败
- `404` - 组织不存在
- `500` - 服务器内部错误

## 参考文档

- GoZero API 文档: https://go-zero.dev/docs/tutorials/api/api
- Protocol Buffers: https://protobuf.dev/
- 项目规范: `CLAUDE.md`
- 架构概览: `specs/003-lender-organization/plan.md`
