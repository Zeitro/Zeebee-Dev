# Admin 模板 CRUD - 技术研究

**功能**: Admin 模板 CRUD 管理系统
**日期**: 2026-01-23
**状态**: 已完成

## 研究概述

本研究为 Admin 模板 CRUD 功能提供技术决策和最佳实践建议。功能涉及基础模板和组织模板的完整生命周期管理，包括级联状态更改和删除操作。

---

## 1. 级联操作实现策略

### 决策：使用应用层代码实现级联操作

**理由**：
- **业务逻辑集中**：所有级联逻辑在应用层统一管理，便于维护和测试
- **灵活控制**：可以根据业务需求定制级联行为，支持复杂条件判断
- **事务一致性**：使用数据库事务确保操作的原子性
- **可观测性**：可以在应用层添加日志和监控，追踪级联操作影响
- **易于调试**：逻辑清晰，问题定位更容易

**实现方案**：

```go
// CascadeStatus 级联状态更改到组织模板
func (s *BaseTemplateService) CascadeStatus(
    ctx context.Context,
    baseTemplateID int64,
    status constants.CommonStatus,
) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 更新所有相关的组织模板状态
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

// CascadeSoftDelete 级联软删除到组织模板
func (s *BaseTemplateService) CascadeSoftDelete(
    ctx context.Context,
    baseTemplateID int64,
) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        result := tx.Model(&model.OrganizationTemplate{}).
            Where("base_template_id = ? AND is_del = 0", baseTemplateID).
            Updates(map[string]interface{}{
                "is_del":     1,
                "updated_at": time.Now(),
            })

        if result.Error != nil {
            return fmt.Errorf("failed to cascade soft delete: %w", result.Error)
        }

        logx.WithContext(ctx).Infof(
            "Cascaded soft delete to %d organization templates for base template %d",
            result.RowsAffected, baseTemplateID,
        )

        return nil
    })
}

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

// UpdateStatusWithCascade 更新基础模板状态并级联
func (l *UpdateBaseTemplateStatusLogic) UpdateStatusWithCascade(
    ctx context.Context,
    req *types.UpdateBaseTemplateStatusReq,
) error {
    // 获取当前模板状态
    template, err := l.svcCtx.TemplateRepo.GetByID(ctx, req.Id)
    if err != nil {
        return err
    }

    // 如果状态没有变化，直接返回
    if template.Status == req.Status {
        return nil
    }

    // 在事务中执行更新和级联
    return l.svcCtx.TemplateRepo.Transaction(ctx, func(ctx context.Context) error {
        // 更新基础模板状态
        if err := l.svcCtx.TemplateRepo.UpdateStatus(ctx, req.Id, req.Status); err != nil {
            return err
        }

        // 级联到组织模板
        if err := l.svcCtx.TemplateService.CascadeStatus(ctx, req.Id, req.Status); err != nil {
            return err
        }

        return nil
    })
}

// DeleteWithCascade 删除基础模板并级联删除组织模板
func (l *DeleteBaseTemplateLogic) DeleteWithCascade(
    ctx context.Context,
    id int64,
) error {
    return l.svcCtx.TemplateRepo.Transaction(ctx, func(ctx context.Context) error {
        // 先执行级联硬删除
        if err := l.svcCtx.TemplateService.CascadeHardDelete(ctx, id); err != nil {
            return err
        }

        // 再删除基础模板
        if err := l.svcCtx.TemplateRepo.HardDelete(ctx, id); err != nil {
            return err
        }

        return nil
    })
}
```

**替代方案评估**：

| 方案 | 优点 | 缺点 | 未选择原因 |
|------|------|------|-----------|
| 数据库触发器 | 性能较好、自动执行 | 逻辑分散、难以调试、维护成本高 | 用户明确要求应用层实现 |
| GORM 回调 | 集成到 ORM 层 | 生命周期管理复杂、事务边界不清晰 | 代码实现不如显式调用清晰 |

---

## 2. JSONB 变量架构验证

### 决策：应用层验证 + 数据库约束

**理由**：
- **Go 类型安全**：在应用层使用强类型结构体确保架构正确性
- **友好的错误消息**：应用层可以提供详细的验证错误信息
- **数据库作为最终防线**：使用 CHECK 约束确保 JSONB 格式有效性

**实现方案**：

```go
// VariableGroup 变量组定义
type VariableGroup struct {
    Name      string          `json:"name" validate:"required,min=1,max=100"`
    Variables []*TemplateVar   `json:"variables" validate:"required,min=1"`
}

// TemplateVar 模板变量定义
type TemplateVar struct {
    Label        string                 `json:"label" validate:"required"`
    Type         string                 `json:"type" validate:"required,oneof=string number boolean array object"`
    Key          string                 `json:"key" validate:"required,min=1,max=100"`
    IsEditable   bool                   `json:"is_editable"`
    Source       string                 `json:"source" validate:"oneof=Fixed Dynamic"`
    FixedValue   any                    `json:"fixed_value"`
    ExampleValue any                    `json:"example_value,omitempty"`
    Properties   map[string]*TemplateVar `json:"properties,omitempty"` // for object type
    Items        []*TemplateVar          `json:"items,omitempty"`      // for array type
}

// Validate 验证变量结构
func (v *TemplateVar) Validate() error {
    // 类型特定验证
    switch v.Type {
    case "object":
        if len(v.Properties) == 0 {
            return errors.New("object type must have properties")
        }
        for key, prop := range v.Properties {
            if err := prop.Validate(); err != nil {
                return fmt.Errorf("property %s: %w", key, err)
            }
        }
    case "array":
        if len(v.Items) == 0 {
            return errors.New("array type must have items")
        }
        for i, item := range v.Items {
            if err := item.Validate(); err != nil {
                return fmt.Errorf("item %d: %w", i, err)
            }
        }
    }
    return nil
}
```

**数据库约束**：

```sql
-- 确保 variable_group_schema 是有效 JSON
ALTER TABLE base_templates
ADD CONSTRAINT valid_variable_schema
CHECK (jsonb_typeof(variable_group_schema) = 'object' OR variable_group_schema = '{}');
```

---

## 3. 模板标签唯一性约束

### 决策：数据库唯一约束 + 应用层验证

**理由**：
- **数据库层面强制**：唯一索引防止并发插入重复标签
- **友好错误处理**：应用层捕获违反约束错误并转换为可读消息
- **支持 FR-016**：标签创建后不可更改，确保引用稳定性

**实现方案**：

```sql
-- 唯一约束：类型 + 类别 + 标签组合
CREATE UNIQUE INDEX idx_base_template_tag_unique
ON base_templates (type, category, tag)
WHERE is_del = 0;

COMMENT ON INDEX idx_base_template_tag_unique IS
'Ensure tag uniqueness within type and category for active templates';
```

**应用层处理**：

```go
// 创建前验证
func (l *CreateBaseTemplateLogic) checkTagUnique(ctx context.Context, tag, templateType, category string) error {
    exists, err := l.svcCtx.TemplateRepo.ExistsByTag(ctx, tag, templateType, category)
    if err != nil {
        return errorx.NewDefaultError("Failed to check tag uniqueness")
    }
    if exists {
        return errorx.NewBusinessError("Tag already exists in this type and category")
    }
    return nil
}

// 更新时禁止更改标签
func (l *UpdateBaseTemplateLogic) validateTagImmutable(ctx context.Context, id int64, newTag string) error {
    template, err := l.svcCtx.TemplateRepo.GetByID(ctx, id)
    if err != nil {
        return err
    }
    if template.Tag != newTag {
        return errorx.NewBusinessError("Tag cannot be changed after creation")
    }
    return nil
}
```

---

## 4. 分页和搜索性能优化

### 决策：数据库索引 + 游标分页

**理由**：
- **索引优化**：为常用查询字段建立索引
- **游标分页**：对于大数据集，游标分页比偏移量分页更高效
- **全文搜索**：使用 PostgreSQL 全文搜索功能处理名称和标签搜索

**实现方案**：

```sql
-- 复合索引：支持类型、类别、状态筛选
CREATE INDEX idx_base_template_filter
ON base_templates (type, category, status, is_del)
WHERE is_del = 0;

-- GIN 索引：支持 JSONB 查询
CREATE INDEX idx_base_template_var_schema
ON base_templates USING GIN (variable_group_schema);

-- 全文搜索索引
CREATE INDEX idx_base_template_search
ON base_templates USING GIN (to_tsvector('simple', name || ' ' || tag));
```

**Go 实现**：

```go
// PaginationRequest 分页请求
type PaginationRequest struct {
    Page      int    `json:"page" validate:"min=1"`
    PageSize  int    `json:"page_size" validate:"min=1,max=100"`
    Type      string `json:"type"`
    Category  string `json:"category"`
    Status    string `json:"status"`
    Search    string `json:"search"`
}

// PaginatedResult 分页结果
type PaginatedResult[T any] struct {
    Items       []T    `json:"items"`
    Total       int64  `json:"total"`
    Page        int    `json:"page"`
    PageSize    int    `json:"page_size"`
    TotalPages  int    `json:"total_pages"`
    HasNext     bool   `json:"has_next"`
    HasPrevious bool   `json:"has_previous"`
}
```

---

## 5. 并发控制策略

### 决策：乐观锁 + 版本号

**理由**：
- **无锁竞争**：乐观锁避免数据库行锁，提高并发性能
- **版本控制**：使用 updated_at 作为版本指示器
- **冲突检测**：检测并发修改并返回友好错误

**实现方案**：

```go
// UpdateWithConcurrencyCheck 带并发检查的更新
func (r *BaseTemplateRepo) UpdateWithConcurrencyCheck(
    ctx context.Context,
    template *model.BaseTemplate,
    previousUpdatedAt time.Time,
) error {
    result := r.db.WithContext(ctx).
        Model(&model.BaseTemplate{}).
        Where("id = ? AND updated_at = ?", template.ID, previousUpdatedAt).
        Updates(template)

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return errorx.NewConcurrencyError(
            "Template was modified by another user. Please refresh and try again.",
        )
    }

    return nil
}
```

---

## 6. 大内容处理策略

### 决策：TOAST 存储 + 分块传输

**理由**：
- **PostgreSQL TOAST**：自动处理大字段存储
- **流式传输**：避免一次性加载大内容到内存
- **大小限制**：应用层验证内容大小（最大 10MB）

**实现方案**：

```go
const maxTemplateContentSize = 10 * 1024 * 1024 // 10MB

func validateContentSize(content string) error {
    if len(content) > maxTemplateContentSize {
        return errorx.NewBusinessError(
            fmt.Sprintf("Template content exceeds maximum size of %d MB", maxTemplateContentSize/(1024*1024)),
        )
    }
    return nil
}
```

---

## 技术决策总结

| 决策点 | 选择 | 理由 |
|--------|------|------|
| 级联操作 | 应用层代码 + 事务 | 业务逻辑集中、灵活控制、易于调试维护 |
| JSONB 验证 | 应用层 + 数据库约束 | 类型安全 + 最终防线 |
| 标签唯一性 | 唯一索引 + 应用验证 | 强制约束 + 友好错误 |
| 分页策略 | 游标分页 + 索引 | 性能优化 |
| 并发控制 | 乐观锁 | 无锁竞争 |
| 大内容 | TOAST + 大小验证 | PostgreSQL 原生支持 |

---

## 依赖库版本

```go
require (
    github.com/zeromicro/go-zero v1.9.3
    gorm.io/gorm v1.25.12
    gorm.io/driver/postgres v1.5.9
    github.com/go-playground/validate/v10 v10.22.0
)
```

---

## 未解决的问题

无 - 所有技术决策已明确。

---

**下一步**: Phase 1 - 数据模型和 API 契约设计
