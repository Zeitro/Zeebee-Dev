# Admin 模板 CRUD 功能需求

**创建日期**: 2026-01-23
**目标**: 为 Admin 后端实现模板管理功能，支持基础模板和组织模板的完整生命周期管理

---

## 1. 业务背景

### 问题场景
当前系统需要支持邮件和合同模板的管理功能。管理员可以创建基础模板（如"借款方默认邮件模板"），然后组织可以基于这些基础模板创建自己的定制版本。

### 核心概念
- **基础模板 (BaseTemplate)**: 由管理员创建的原始模板，作为所有组织模板的源头
- **组织模板 (OrganizationTemplate)**: 组织基于基础模板创建的定制版本
- **级联关系**: 组织模板依赖于基础模板，基础模板的变更会影响所有派生的组织模板

---

## 2. 功能需求

### 2.1 基础模板 CRUD (P0 - 核心功能)

管理员需要能够对基础模板进行完整的增删改查操作：

| 操作 | 描述 | 优先级 |
|------|------|--------|
| 创建 | 创建新的基础模板，包含名称、类型、类别、标签、内容、变量架构 | P0 |
| 查询 | 支持单个查询和分页列表查询，支持按类型/类别/状态筛选 | P0 |
| 更新 | 更新模板内容、变量架构、描述等信息 | P0 |
| 删除 | 软删除基础模板 | P0 |
| 状态切换 | 启用/禁用模板 | P0 |

**关键字段**:
- `type`: Borrower（借款方）/ Lender（贷款方）
- `category`: Email（邮件）/ Contract（合同）
- `tag`: 模板唯一标识符（在 type+category 组合内唯一），**创建后不可更改**
- `status`: Active（启用）/ Disabled（禁用）
- `content`: 模板内容（纯文本或 HTML，最大 10MB）
- `variable_group_schema`: 变量架构（JSONB 格式，支持嵌套对象和数组）

### 2.2 级联操作 (P0 - 数据完整性)

**关键决策**: 使用应用层代码实现级联，不使用数据库触发器

| 触发操作 | 级联行为 | 实现方式 |
|----------|----------|----------|
| 禁用基础模板 | 所有派生的组织模板也被禁用 | 应用层代码 |
| 启用基础模板 | 所有派生的组织模板也被启用 | 应用层代码 |
| 删除基础模板 | 所有派生的组织模板也被删除 | 应用层代码 |

### 2.3 变量架构管理 (P1 - 高级功能)

模板支持复杂的变量结构定义，用于动态内容生成：

- 支持的数据类型: string, number, boolean, array, object
- object 类型可以包含嵌套属性
- array 类型可以定义项目架构
- 每个变量包含: label, key, type, is_editable, source, fixed_value, example_value

---

## 3. 数据模型设计

### 3.1 表结构

```sql
-- 基础模板表
CREATE TABLE base_templates (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    is_del smallint NOT NULL DEFAULT 0,
    name varchar(100) NOT NULL,
    type varchar(20) NOT NULL,              -- Borrower/Lender
    category varchar(20) NOT NULL,          -- Email/Contract
    tag varchar(100) NOT NULL,              -- type+category 内唯一，不可变
    content text NOT NULL,                  -- 最大 10MB
    variable_group_schema jsonb NOT NULL,   -- 变量架构
    status varchar(20) NOT NULL DEFAULT 'Active',
    description varchar(255)
);

-- 组织模板表
CREATE TABLE organization_templates (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    is_del smallint NOT NULL DEFAULT 0,
    organization_id int NOT NULL,
    base_template_id int NOT NULL,           -- 关联基础模板
    name varchar(100) NOT NULL,
    type varchar(20) NOT NULL,
    category varchar(20) NOT NULL,
    tag varchar(100) NOT NULL,
    content text NOT NULL,
    variable_group_schema jsonb NOT NULL,
    status varchar(20) NOT NULL DEFAULT 'Active',
    description varchar(255),

    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (base_template_id) REFERENCES base_templates(id)
);
```

### 3.2 GORM Model

```go
// BaseTemplate 基础模板
type BaseTemplate struct {
    Base
    Name                string                     `json:"name" gorm:"type:varchar(100)"`
    Type                constants.UserType        `json:"type" gorm:"type:varchar(20)"`
    Category            constants.TemplateCategory `json:"category" gorm:"type:varchar(20)"`
    Tag                 string                     `json:"tag" gorm:"type:varchar(100)"`
    Content             string                     `json:"content" gorm:"type:text"`
    VariableGroupSchema []*TemplateVariableGroup  `json:"variable_group_schema" gorm:"type:jsonb"`
    Status              constants.CommonStatus    `json:"status" gorm:"type:varchar(20)"`
    Description         string                     `json:"description" gorm:"type:varchar(255)"`
}

// OrganizationTemplate 组织模板
type OrganizationTemplate struct {
    Base
    OrganizationId      int64                      `json:"organization_id" gorm:"index"`
    BaseTemplateId      int64                      `json:"base_template_id" gorm:"index"`
    Name                string                     `json:"name" gorm:"type:varchar(100)"`
    Type                constants.UserType        `json:"type" gorm:"type:varchar(20)"`
    Category            constants.TemplateCategory `json:"category" gorm:"type:varchar(20)"`
    Tag                 string                     `json:"tag" gorm:"type:varchar(100)"`
    Content             string                     `json:"content" gorm:"type:text"`
    VariableGroupSchema []*TemplateVariableGroup  `json:"variable_group_schema" gorm:"type:jsonb"`
    Status              constants.CommonStatus    `json:"status" gorm:"type:varchar(20)"`
    Description         string                     `json:"description" gorm:"type:varchar(255)"`
}
```

---

## 4. API 接口设计

### 4.1 端点列表

遵循 admin 端现有风格，使用 `BasePageRequest` 和 `BasePageResponse`：

| 方法 | 路径 | 功能 |
|------|------|------|
| POST | /admin/v1/base_template | 创建基础模板 |
| PUT | /admin/v1/base_template/:id | 更新基础模板 |
| GET | /admin/v1/base_template/:id | 获取模板详情 |
| POST | /admin/v1/base_template/page | 分页查询（支持筛选） |
| PUT | /admin/v1/base_template/:id/status | 切换状态（级联） |
| DELETE | /admin/v1/base_template/:id | 删除模板（级联） |

### 4.2 筛选条件

- `keyword`: 模糊匹配名称或标签
- `type`: Borrower/Lender
- `category`: Email/Contract
- `status`: Active/Disabled

---

## 5. 技术约束与要求

### 5.1 架构约束

1. **微服务架构**: API 层 + RPC 层分离
2. **共享代码**: RPC 业务逻辑可放在 `domains/zeebee` 中，供 admin 和 lender 服务复用
3. **代码风格**: 遵循 admin 端现有代码风格
4. **数据转换**: 使用 `copierx.Trans[T, U]` 进行类型转换

### 5.2 级联实现

**必须使用应用层代码实现级联**，不使用数据库触发器：

```go
// 级联状态切换示例
func CascadeStatus(ctx context.Context, baseTemplateID int64, status string) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 1. 更新基础模板状态
        // 2. 批量更新所有派生的组织模板状态
        return nil
    })
}
```

### 5.3 验证规则

| 字段 | 验证规则 |
|------|----------|
| name | 必填，1-100 字符 |
| type | 必填，Borrower/Lender |
| category | 必填，Email/Contract |
| tag | 必填，1-100 字符，type+category 内唯一，**创建后不可更改** |
| content | 必填，最大 10MB |
| status | Active/Disabled |
| variable_group_schema | 必填，有效 JSONB |

---

## 6. 开发分支

**要求**: 从 `feat-backend-template` 分支切出新分支进行开发

```bash
git checkout feat-backend-template
git checkout -b 6-admin-template-crud
```

---

## 7. 不在范围内

- 模板版本控制和历史
- 模板预览和渲染
- 组织模板的分发流程（分配给组织）
- 组织级模板定制界面
- 模板导入/导出
- 富文本编辑器集成

---

## 8. 参考信息

- **数据库迁移文件**: `deploy/zeebee/sql/20260123144148_create_template_table.sql`
- **常量定义**: `domains/zeebee/constants/template.go`
- **参考模块**: admin 端组织管理模块（类似的 CRUD 结构）
