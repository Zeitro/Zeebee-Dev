# Admin 模板 CRUD - 数据模型

**功能**: Admin 模板 CRUD 管理系统
**日期**: 2026-01-23
**状态**: 已完成

## 数据模型概述

本文档定义了 Admin 模板 CRUD 功能的完整数据模型，包括数据库表结构、Go 结构体定义和关系映射。**注意：SQL 迁移已手动完成，本文档仅作参考。**

---

## 1. 数据库表结构（已手动创建）

### 1.1 base_templates 表

```sql
CREATE TABLE "public"."base_templates" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" smallint NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "type" varchar(20) NOT NULL DEFAULT '',
    "category" varchar(20) NOT NULL DEFAULT '',
    "tag" varchar(100) NOT NULL DEFAULT '',
    "content" text NOT NULL DEFAULT '',
    "variable_group_schema" jsonb NOT NULL DEFAULT '{}',
    "status" varchar(20) NOT NULL DEFAULT 'Active',
    "description" varchar(255) NOT NULL DEFAULT ''
);

-- 注释
COMMENT ON COLUMN "public"."base_templates"."type" IS 'Borrower/Lender';
COMMENT ON COLUMN "public"."base_templates"."category" IS 'Email/Contract';
COMMENT ON COLUMN "public"."base_templates"."status" IS 'Active/Disabled';
COMMENT ON COLUMN "public"."base_templates"."tag" IS 'Unique identifier within type and category, immutable after creation';
```

### 1.2 organization_templates 表

```sql
CREATE TABLE "public"."organization_templates" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" smallint NOT NULL DEFAULT 0,
    "organization_id" integer NOT NULL DEFAULT 0,
    "base_template_id" integer NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "type" varchar(20) NOT NULL DEFAULT '',
    "category" varchar(20) NOT NULL DEFAULT '',
    "tag" varchar(100) NOT NULL DEFAULT '',
    "content" text NOT NULL DEFAULT '',
    "variable_group_schema" jsonb NOT NULL DEFAULT '{}',
    "status" varchar(20) NOT NULL DEFAULT 'Active',
    "description" varchar(255) NOT NULL DEFAULT ''
);

-- 外键
ALTER TABLE "public"."organization_templates"
ADD CONSTRAINT "fk_org_tpl_org"
FOREIGN KEY ("organization_id") REFERENCES "public"."organizations"("id") ON DELETE CASCADE;

ALTER TABLE "public"."organization_templates"
ADD CONSTRAINT "fk_org_tpl_base"
FOREIGN KEY ("base_template_id") REFERENCES "public"."base_templates"("id") ON DELETE CASCADE;

-- 注释
COMMENT ON COLUMN "public"."organization_templates"."type" IS 'Borrower/Lender';
COMMENT ON COLUMN "public"."organization_templates"."category" IS 'Email/Contract';
COMMENT ON COLUMN "public"."organization_templates"."status" IS 'Active/Disabled';
COMMENT ON COLUMN "public"."organization_templates"."base_template_id" IS 'References base_templates table, cascade operations handled in application layer';
```

---

## 2. Go 结构体定义

### 2.1 Model 层 (domains/zeebee/model/)

```go
package model

import (
    "time"

    "github.com/yourorg/zeebee/domains/zeebee/constants"
)

// BaseTemplate represents the base_templates table
type BaseTemplate struct {
    ID                 int64                       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
    CreatedAt          time.Time                   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt          time.Time                   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
    IsDel              int16                       `json:"is_del" gorm:"column:is_del;default:0;index"`
    Name               string                      `json:"name" gorm:"column:name;type:varchar(100);not null;default:''"`
    Type               constants.UserType          `json:"type" gorm:"column:type;type:varchar(20);not null;default:''"`
    Category           constants.TemplateCategory  `json:"category" gorm:"column:category;type:varchar(20);not null;default:''"`
    Tag                string                      `json:"tag" gorm:"column:tag;type:varchar(100);not null;default:''"`
    Content            string                      `json:"content" gorm:"column:content;type:text;not null;default:''"`
    VariableGroupSchema *VariableGroupSchema       `json:"variable_group_schema" gorm:"column:variable_group_schema;type:jsonb;not null;default:'{}';serializer:json"`
    Status             constants.CommonStatus       `json:"status" gorm:"column:status;type:varchar(20);not null;default:'Active'"`
    Description        string                      `json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`
}

// TableName returns the table name for BaseTemplate
func (BaseTemplate) TableName() string {
    return "base_templates"
}

// VariableGroupSchema represents the variable_group_schema JSONB field
type VariableGroupSchema []*TemplateVariableGroup

// TemplateVariableGroup represents a group of template variables
type TemplateVariableGroup struct {
    Name      string                  `json:"name" validate:"required,min=1,max=100"`
    Variables []*TemplateVariable     `json:"variables" validate:"required,min=1"`
}

// TemplateVariable represents a single template variable definition
type TemplateVariable struct {
    Label        string                      `json:"label" validate:"required"`
    Type         string                      `json:"type" validate:"required,oneof=string number boolean array object"`
    Key          string                      `json:"key" validate:"required,min=1,max=100"`
    IsEditable   bool                        `json:"is_editable"`
    Source       string                      `json:"source" validate:"oneof=Fixed Dynamic"`
    FixedValue   any                         `json:"fixed_value"`
    ExampleValue any                         `json:"example_value,omitempty"`
    Properties   map[string]*TemplateVariable `json:"properties,omitempty"` // for object type
    Items        []*TemplateVariable         `json:"items,omitempty"`      // for array type
}

// Validate validates the template variable structure
func (v *TemplateVariable) Validate() error {
    switch v.Type {
    case "object":
        if len(v.Properties) == 0 {
            return fmt.Errorf("object type must have properties")
        }
        for key, prop := range v.Properties {
            if err := prop.Validate(); err != nil {
                return fmt.Errorf("property %s: %w", key, err)
            }
        }
    case "array":
        if len(v.Items) == 0 {
            return fmt.Errorf("array type must have items")
        }
        for i, item := range v.Items {
            if err := item.Validate(); err != nil {
                return fmt.Errorf("item %d: %w", i, err)
            }
        }
    }
    return nil
}

// OrganizationTemplate represents the organization_templates table
type OrganizationTemplate struct {
    ID                 int64                       `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
    CreatedAt          time.Time                   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt          time.Time                   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
    IsDel              int16                       `json:"is_del" gorm:"column:is_del;default:0;index"`
    OrganizationID     int64                       `json:"organization_id" gorm:"column:organization_id;default:0;index:idx_org_tpl_org_id"`
    BaseTemplateID     int64                       `json:"base_template_id" gorm:"column:base_template_id;default:0;index:idx_org_tpl_base_id"`
    Name               string                      `json:"name" gorm:"column:name;type:varchar(100);not null;default:''"`
    Type               constants.UserType          `json:"type" gorm:"column:type;type:varchar(20);not null;default:''"`
    Category           constants.TemplateCategory  `json:"category" gorm:"column:category;type:varchar(20);not null;default:''"`
    Tag                string                      `json:"tag" gorm:"column:tag;type:varchar(100);not null;default:''"`
    Content            string                      `json:"content" gorm:"column:content;type:text;not null;default:''"`
    VariableGroupSchema *VariableGroupSchema       `json:"variable_group_schema" gorm:"column:variable_group_schema;type:jsonb;not null;default:'{}';serializer:json"`
    Status             constants.CommonStatus       `json:"status" gorm:"column:status;type:varchar(20);not null;default:'Active'"`
    Description        string                      `json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`

    // Relations
    BaseTemplate       *BaseTemplate               `json:"base_template,omitempty" gorm:"foreignKey:BaseTemplateID"`
    Organization       *Organization               `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
}

// TableName returns the table name for OrganizationTemplate
func (OrganizationTemplate) TableName() string {
    return "organization_templates"
}
```

### 2.2 BO 层 (domains/zeebee/data/bo/)

```go
package bo

import (
    "time"

    "github.com/yourorg/zeebee/domains/zeebee/constants"
    "github.com/yourorg/zeebee/domains/zeebee/model"
)

// BaseTemplateBo represents a base template business object
type BaseTemplateBo struct {
    ID                 int64                       `json:"id"`
    CreatedAt          *time.Time                  `json:"created_at,omitempty"`
    UpdatedAt          *time.Time                  `json:"updated_at,omitempty"`
    Name               string                      `json:"name"`
    Type               constants.UserType          `json:"type"`
    Category           constants.TemplateCategory  `json:"category"`
    Tag                string                      `json:"tag"`
    Content            string                      `json:"content"`
    VariableGroupSchema *model.VariableGroupSchema `json:"variable_group_schema"`
    Status             constants.CommonStatus      `json:"status"`
    Description        string                      `json:"description"`
}

// BaseTemplatePageBo represents base template for page display
type BaseTemplatePageBo struct {
    ID          int64                      `json:"id"`
    Name        string                     `json:"name"`
    Type        constants.UserType         `json:"type"`
    Category    constants.TemplateCategory `json:"category"`
    Tag         string                     `json:"tag"`
    Status      constants.CommonStatus     `json:"status"`
    Description string                     `json:"description"`
    CreatedAt   *time.Time                 `json:"created_at"`
    UpdatedAt   *time.Time                 `json:"updated_at"`
}
```

### 2.3 API Types 层 (app/admin/api/internal/types/)

遵循 admin 端现有风格，使用 `BasePageRequest` 和 `BasePageResponse`：

```go
package types

// BaseTemplateEntity represents a base template
type BaseTemplateEntity struct {
    Id                 int64                       `json:"id"`
    Name               string                      `json:"name"`
    Type               string                      `json:"type"`
    Category           string                      `json:"category"`
    Tag                string                      `json:"tag"`
    Content            string                      `json:"content"`
    VariableGroupSchema *model.VariableGroupSchema `json:"variable_group_schema"`
    Status             string                      `json:"status"`
    Description        string                      `json:"description"`
    CreatedAt          string                      `json:"created_at"`
    UpdatedAt          string                      `json:"updated_at"`
}

// BaseTemplatePageEntity represents base template in page list
type BaseTemplatePageEntity struct {
    Id          int64  `json:"id"`
    Name        string `json:"name"`
    Type        string `json:"type"`
    Category    string `json:"category"`
    Tag         string `json:"tag"`
    Status      string `json:"status"`
    Description string `json:"description"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

// BaseTemplateFilter represents filter criteria
type BaseTemplateFilter struct {
    Keyword string `json:"keyword,optional"`
    Type     string `json:"type,optional" validate:"omitempty,oneof=Borrower Lender"`
    Category string `json:"category,optional" validate:"omitempty,oneof=Email Contract"`
    Status   string `json:"status,optional" validate:"omitempty,oneof=Active Disabled"`
}

// BaseTemplatePageRequest - 使用 BasePageRequest
type BaseTemplatePageRequest struct {
    BasePageRequest
    Filter *BaseTemplateFilter `json:"filter,optional"`
    Sort   *BasePageSort        `json:"sort,optional"`
}

// BaseTemplatePageResponse - 使用 BasePageResponse
type BaseTemplatePageResponse struct {
    BasePageResponse
    List []*BaseTemplatePageEntity `json:"list"`
}

// CreateBaseTemplateRequest represents create request
type CreateBaseTemplateRequest struct {
    Name                string                       `json:"name" validate:"required,min=1,max=100"`
    Type                string                       `json:"type" validate:"required,oneof=Borrower Lender"`
    Category            string                       `json:"category" validate:"required,oneof=Email Contract"`
    Tag                 string                       `json:"tag" validate:"required,min=1,max=100"`
    Content             string                       `json:"content" validate:"required,max=10485760"`
    VariableGroupSchema *model.VariableGroupSchema   `json:"variable_group_schema" validate:"required"`
    Description         string                       `json:"description" validate:"max=255"`
}

// UpdateBaseTemplateRequest represents update request
type UpdateBaseTemplateRequest struct {
    Id                  int64                        `path:"id" validate:"required,min=1"`
    Name                string                       `json:"name" validate:"required,min=1,max=100"`
    Type                string                       `json:"type" validate:"required,oneof=Borrower Lender"`
    Category            string                       `json:"category" validate:"required,oneof=Email Contract"`
    Tag                 string                       `json:"tag" validate:"required,min=1,max=100"`
    Content             string                       `json:"content" validate:"required,max=10485760"`
    VariableGroupSchema *model.VariableGroupSchema   `json:"variable_group_schema" validate:"required"`
    Description         string                       `json:"description" validate:"max=255"`
}

// BaseTemplateSwitchStatusRequest represents switch status request
type BaseTemplateSwitchStatusRequest struct {
    Id     int64  `path:"id" validate:"required,min=1"`
    Status string `json:"status" validate:"required,oneof=Active Disabled"`
}

// BaseTemplateDetailRequest represents detail request
type BaseTemplateDetailRequest struct {
    Id int64 `path:"id" validate:"required,min=1"`
}

// BaseTemplateDetailResponse represents detail response
type BaseTemplateDetailResponse struct {
    BaseTemplateEntity
}

// DeleteBaseTemplateRequest represents delete request
type DeleteBaseTemplateRequest struct {
    Id int64 `path:"id" validate:"required,min=1"`
}
```

---

## 3. 数据验证规则

### 3.1 字段验证

| 字段 | 验证规则 | 说明 |
|------|----------|------|
| name | required, min=1, max=100 | 必填，最大100字符 |
| type | required, oneof=Borrower Lender | 必填，只能是借款方或贷款方 |
| category | required, oneof=Email Contract | 必填，只能是邮件或合同 |
| tag | required, min=1, max=100, immutable | 必填，创建后不可更改 |
| status | required, oneof=Active Disabled | 必填，只能是启用或禁用 |
| content | required, max=10485760 (10MB) | 必填，最大10MB |
| variable_group_schema | required, json | 必填，有效JSON |
| type + category + tag | unique | 组合唯一约束 |

### 3.2 JSONB 验证函数

```go
// ValidateVariableGroupSchema validates the variable group schema
func ValidateVariableGroupSchema(schema *model.VariableGroupSchema) error {
    if schema == nil || len(*schema) == 0 {
        return nil // Empty schema is valid
    }

    for _, group := range *schema {
        if err := validateVariableGroup(group); err != nil {
            return fmt.Errorf("variable group %s: %w", group.Name, err)
        }
    }
    return nil
}

func validateVariableGroup(group *model.TemplateVariableGroup) error {
    if group.Name == "" {
        return errors.New("group name is required")
    }
    if len(group.Variables) == 0 {
        return errors.New("group must have at least one variable")
    }

    for i, variable := range group.Variables {
        if err := variable.Validate(); err != nil {
            return fmt.Errorf("variable %d: %w", i, err)
        }
    }
    return nil
}
```

---

## 4. 状态转换

### 4.1 基础模板状态转换

```
     ┌─────────┐
     │ Created │
     └────┬────┘
          │
          ▼
     ┌─────────┐  disable   ┌──────────┐
     │ Active  │────────────→│ Disabled │
     └─────────┘            └──────────┘
          │                       │
          │ enable               │ enable
          ▼                       ▼
     ┌─────────┐            ┌──────────┐
     │ Active  │←───────────│ Disabled │
     └─────────┘            └──────────┘
```

**级联规则**（应用层代码实现）：
- Active → Disabled: 所有派生的组织模板也变为 Disabled
- Disabled → Active: 所有派生的组织模板也变为 Active

### 4.2 删除操作（应用层代码级联）

```
     ┌─────────┐
     │ Any     │
     └────┬────┘
          │ delete (application layer cascade)
          ▼
     ┌─────────────────┐
     │ 1. Cascade to   │
     │    org templates│
     └────────┬────────┘
             │
             ▼
     ┌─────────────────┐
     │ 2. Delete base  │
     │    template     │
     └─────────────────┘
```

---

## 5. 关系图

```
┌─────────────────┐
│  organizations │
└────────┬────────┘
         │ 1
         │
         │ N
         ▼
┌─────────────────────┐         1:N         ┌─────────────────┐
│organization_templates│◄────────────────────│  base_templates │
└─────────────────────┘                      └─────────────────┘
         │                                               │
         │ organization_id                                │
         │                                               │
         └─── BaseTemplateID (FK) ───────────────────────┘
```

---

## 6. 分页结构（遵循 admin 端风格）

### 6.1 基础分页结构（已在 common.api 定义）

```go
// BasePageRequest 基础分页请求
type BasePageRequest struct {
    Limit int `json:"limit,optional"`
    Page  int `json:"page, optional"`
}

// BasePageResponse 基础分页响应
type BasePageResponse struct {
    Page      int64 `json:"page"`
    Limit     int64 `json:"limit"`
    Total     int64 `json:"total"`
    TotalPage int64 `json:"total_page"`
}

// BasePageFilter 基础筛选
type BasePageFilter struct {
    Keyword string `json:"keyword,optional"`
}

// BasePageSort 基础排序
type BasePageSort struct {
    Column string `json:"column,optional"`
    Sort   string `json:"sort,optional"`
}
```

### 6.2 模板分页请求/响应

```go
// 请求
type BaseTemplatePageRequest struct {
    BasePageRequest
    Filter *BaseTemplateFilter `json:"filter,optional"`
    Sort   *BasePageSort        `json:"sort,optional"`
}

// 响应
type BaseTemplatePageResponse struct {
    BasePageResponse
    List []*BaseTemplatePageEntity `json:"list"`
}
```

---

## 7. 级联操作实现（应用层代码）

### 7.1 级联状态更改

```go
// CascadeStatus 级联状态更改到组织模板
func (s *BaseTemplateService) CascadeStatus(
    ctx context.Context,
    baseTemplateID int64,
    status constants.CommonStatus,
) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        result := tx.Model(&model.OrganizationTemplate{}).
            Where("base_template_id = ? AND status != ?", baseTemplateID, status).
            Updates(map[string]interface{}{
                "status":     status,
                "updated_at": time.Now(),
            })

        if result.Error != nil {
            return fmt.Errorf("failed to cascade status: %w", result.Error)
        }

        logx.WithContext(ctx).Infof(
            "Cascaded status %s to %d organization templates for base template %d",
            status, result.RowsAffected, baseTemplateID,
        )

        return nil
    })
}
```

### 7.2 级联删除

```go
// CascadeHardDelete 级联硬删除到组织模板
func (s *BaseTemplateService) CascadeHardDelete(
    ctx context.Context,
    baseTemplateID int64,
) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        result := tx.Where("base_template_id = ?", baseTemplateID).
            Delete(&model.OrganizationTemplate{})

        if result.Error != nil {
            return fmt.Errorf("failed to cascade hard delete: %w", result.Error)
        }

        logx.WithContext(ctx).Infof(
            "Cascaded hard delete to %d organization templates for base template %d",
            result.RowsAffected, baseTemplateID,
        )

        return nil
    })
}
```

---

**下一步**: Phase 1 - API 契约定义 (contracts/)
