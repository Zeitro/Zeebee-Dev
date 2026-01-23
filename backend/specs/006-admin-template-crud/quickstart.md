# Admin 模板 CRUD - 快速入门

**功能**: Admin 模板 CRUD 管理系统
**日期**: 2026-01-23
**分支**: 6-admin-template-crud

## 快速开始

### 1. 数据库迁移

SQL 迁移文件已手动创建，位置：
```
deploy/zeebee/sql/20260123144148_create_template_table.sql
```

### 2. 新增文件列表

```
domains/zeebee/
├── constants/template.go           # 模板相关常量（已存在，需更新）
├── model/
│   ├── base_template.go            # 基础模板模型
│   └── organization_template.go    # 组织模板模型
├── repo/
│   ├── base_template.go            # 基础模板仓储
│   └── organization_template.go    # 组织模板仓储
└── logic/
    └── template.go                 # 模板业务逻辑（级联操作）

app/admin/api/
├── desc/template.api               # API 定义
└── internal/
    ├── handler/template/           # HTTP handlers
    └── logic/template/             # API logic

app/admin/rpc/
├── client/template/                # RPC client
└── internal/
    ├── logic/template/             # RPC logic
    └── server/template/            # RPC server
```

### 3. 核心代码片段

#### 3.1 创建基础模板

**API Handler** (`app/admin/api/internal/handler/template/create_base_template_handler.go`):
```go
package template

import (
    "net/http"

    "github.com/zeromicro/go-zero/rest/httpx"
    "yourorg/app/admin/api/internal/logic/template"
    "yourorg/app/admin/api/internal/svc"
    "yourorg/app/admin/api/internal/types"
)

func CreateBaseTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.CreateBaseTemplateRequest
        if err := httpx.Parse(r, &req); err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
            return
        }

        l := template.NewCreateBaseTemplateLogic(r.Context(), svcCtx)
        resp, err := l.CreateBaseTemplate(&req)
        if err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
        } else {
            httpx.OkJsonCtx(r.Context(), w, resp)
        }
    }
}
```

**API Logic** (`app/admin/api/internal/logic/template/create_base_template_logic.go`):
```go
package template

import (
    "context"
    "fmt"

    "yourorg/common/errorx"
    "yourorg/app/admin/api/internal/svc"
    "yourorg/app/admin/api/internal/types"
    "yourorg/app/admin/rpc/client/template"
)

type CreateBaseTemplateLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCreateBaseTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBaseTemplateLogic {
    return &CreateBaseTemplateLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *CreateBaseTemplateLogic) CreateBaseTemplate(req *types.CreateBaseTemplateRequest) (resp *types.CreateBaseTemplateResponse, err error) {
    // 验证标签唯一性
    exists, err := l.svcCtx.TemplateRPC.ExistsByTag(l.ctx, &template.ExistsByTagReq{
        Tag:     req.Tag,
        Type:     req.Type,
        Category: req.Category,
    })
    if err != nil {
        return nil, errorx.NewDefaultError("Failed to check tag uniqueness")
    }
    if exists.Exists {
        return nil, errorx.NewBusinessError("Tag already exists in this type and category")
    }

    // 调用 RPC 创建模板
    result, err := l.svcCtx.TemplateRPC.Create(l.ctx, &template.CreateReq{
        Name:                req.Name,
        Type:                req.Type,
        Category:            req.Category,
        Tag:                 req.Tag,
        Content:             req.Content,
        VariableGroupSchema: req.VariableGroupSchema,
        Description:         req.Description,
    })

    if err != nil {
        return nil, err
    }

    return &types.CreateBaseTemplateResponse{
        Id:                 result.Id,
        Name:               result.Name,
        Type:               result.Type,
        Category:           result.Category,
        Tag:                result.Tag,
        Content:            result.Content,
        VariableGroupSchema: result.VariableGroupSchema,
        Status:             result.Status,
        Description:        result.Description,
        CreatedAt:          result.CreatedAt,
        UpdatedAt:          result.UpdatedAt,
    }, nil
}
```

#### 3.2 级联状态切换（RPC Logic）

```go
package template

import (
    "context"
    "fmt"

    "yourorg/common/errorx"
    "yourorg/domains/zeebee/constants"
    "yourorg/domains/zeebee/logic"
    "yourorg/domains/zeebee/model"
    "yourorg/domains/zeebee/repo"
    "yourorg/app/admin/rpc/internal/svc"
)

type SwitchStatusLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (l *SwitchStatusLogic) SwitchStatus(in *template.SwitchStatusReq) (*template.EmptyResp, error) {
    // 获取当前模板
    tmpl, err := l.svcCtx.TemplateRepo.GetByID(l.ctx, in.Id)
    if err != nil {
        return nil, errorx.NewDefaultError("Template not found")
    }

    // 状态未变化，直接返回
    if tmpl.Status == constants.CommonStatus(in.Status) {
        return &template.EmptyResp{}, nil
    }

    // 在事务中执行更新和级联
    err = l.svcCtx.TemplateRepo.Transaction(l.ctx, func(ctx context.Context) error {
        // 更新基础模板状态
        if err := l.svcCtx.TemplateRepo.UpdateStatus(ctx, in.Id, constants.CommonStatus(in.Status)); err != nil {
            return fmt.Errorf("failed to update base template status: %w", err)
        }

        // 应用层级联到组织模板
        if err := l.svcCtx.TemplateService.CascadeStatus(ctx, in.Id, constants.CommonStatus(in.Status)); err != nil {
            return fmt.Errorf("failed to cascade status: %w", err)
        }

        return nil
    })

    if err != nil {
        return nil, err
    }

    return &template.EmptyResp{}, nil
}
```

#### 3.3 级联删除

```go
// CascadeHardDelete 级联硬删除
func (s *TemplateService) CascadeHardDelete(
    ctx context.Context,
    baseTemplateID int64,
) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 先删除所有组织模板
        result := tx.Where("base_template_id = ?", baseTemplateID).
            Delete(&model.OrganizationTemplate{})

        if result.Error != nil {
            return fmt.Errorf("failed to cascade delete org templates: %w", result.Error)
        }

        logx.WithContext(ctx).Infof(
            "Cascaded delete to %d organization templates",
            result.RowsAffected,
        )

        return nil
    })
}
```

### 4. API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /admin/v1/template | 创建基础模板 |
| PUT | /admin/v1/template/:id | 更新基础模板 |
| GET | /admin/v1/template/:id | 获取模板详情 |
| POST | /admin/v1/template/page | 分页查询模板 |
| PUT | /admin/v1/template/:id/status | 切换模板状态 |
| DELETE | /admin/v1/template/:id | 删除模板 |

### 5. 关键注意事项

1. **标签不可变性**: 模板创建后，标签字段不可更改
2. **应用层级联**: 所有级联操作（状态切换、删除）通过应用层代码实现
3. **事务安全**: 级联操作必须在数据库事务中执行
4. **数据验证**: JSONB 字段在应用层进行结构验证
5. **内容大小限制**: 模板内容最大 10MB

### 6. 开发顺序建议

1. Model 层（数据模型）
2. Repo 层（数据访问）
3. RPC 层（业务逻辑 + 级联操作）
4. API 层（HTTP 接口）

---

**下一步**: 运行 `/speckit.tasks` 生成详细的任务清单
