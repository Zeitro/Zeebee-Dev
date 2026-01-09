# 共享代码架构设计指南

**文档创建时间**: 2025-11-24
**文档类型**: 架构设计指导
**维护状态**: 活跃

## 概述

本文档详细说明 Zeebee 项目中共享代码层架构的实施，包括共享数据模型的引入、代码复用实现以及服务间的协作模式。此架构设计的核心目标是减少重复代码，提高开发效率，确保borrower和lender服务使用一致的数据访问层。

**重要说明**: 本项目采用共享代码层模式，而非完整的领域驱动设计(DDD)。domains目录包含的是可复用的数据模型和仓储实现，主要为提高代码复用率和数据一致性。

## 架构演进

### 原始架构 (服务特定数据访问)

```
app/borrower/rpc/internal/repo/model/  # Borrower 专用模型
app/borrower/rpc/internal/repo/        # Borrower 专用仓储
app/lender/rpc/internal/repo/model/    # Lender 专用模型
app/lender/rpc/internal/repo/          # Lender 专用仓储
```

### 新架构 (共享代码层)

```
domains/zeebee/model/    # 共享数据模型定义
domains/zeebee/repo/     # 共享数据访问实现
app/borrower/rpc/internal/logic/    # 借款方业务逻辑
app/lender/rpc/internal/logic/      # 贷款方业务逻辑
```

**架构优势**:

- **代码复用**: borrower和lender服务共享相同的数据模型和仓储
- **数据一致性**: 统一的数据访问层确保数据操作的一致性
- **维护效率**: 数据模型变更只需在一处修改
- **开发效率**: 减少重复的数据访问代码编写

## 共享代码层组织结构

### 共享数据模型 (domains/zeebee/model/)

#### 核心业务实体

```go
// 基础模型，提供软删除和通用字段
type Base struct {
    Id        int64                  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    CreatedAt time.Time              `gorm:"column:created_at" json:"created_at"`
    UpdatedAt time.Time              `gorm:"column:updated_at" json:"updated_at"`
    DeletedAt gorm.DeletedAt         `gorm:"column:deleted_at" json:"deleted_at"`
}

// 贷款实体 - 核心业务对象
type Loan struct {
    Base
    Code               string            `gorm:"column:code;size:100;not null" json:"code"`
    BorrowerId         int64             `gorm:"column:borrower_id;not null" json:"borrower_id"`
    LenderCompanyId    int64             `gorm:"column:lender_company_id;not null" json:"lender_company_id"`
    // ... 业务字段
}
```

#### 通用值对象

```go
// 地址值对象 - 被多个实体复用
type Address struct {
    State      string `gorm:"column:state;size:100" json:"state"`
    City       string `gorm:"column:city;size:100" json:"city"`
    Address1   string `gorm:"column:address1;size:255" json:"address1"`
    Address2   string `gorm:"column:address2;size:255" json:"address2"`
    ZipCode    string `gorm:"column:zip_code;size:20" json:"zip_code"`
}
```

### 共享数据访问层 (domains/zeebee/repo/)

#### 泛型仓储模式

```go
// 基础仓储，提供通用 CRUD 操作
type BaseRepo[T schema.Tabler] struct {
    Client func(ctx context.Context) *gorm.DB
}

// 贷款数据访问层特化
type LoanRepo struct {
    *gormx.BaseRepo[*model.Loan]
    m *model.Loan
}

func NewLoanRepo(db *gorm.DB) *LoanRepo {
    return &LoanRepo{
        BaseRepo: gormx.NewBaseRepo[*model.Loan](db),
        m:        new(model.Loan),
    }
}
```

## 服务集成模式

### 借款方服务集成

```go
// app/borrower/rpc/internal/logic/loan_logic.go
import (
    "zeebee-backend/domains/zeebee/repo"
    "zeebee-backend/domains/zeebee/model"
)

type LoanLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    loanRepo *repo.LoanRepo  // 使用共享仓储
}
```

### 贷款方服务集成

```go
// app/lender/rpc/internal/logic/loan_logic.go
import (
    "zeebee-backend/domains/zeebee/repo"
    "zeebee-backend/domains/zeebee/model"
)

type LoanLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    loanRepo *repo.LoanRepo  // 使用相同的共享仓储
}
```

## 代码复用统计

### 当前复用情况

- **共享数据层文件**: 37 个 Go 文件
- **使用共享层的服务**: 45 个文件引用 domains 包
- **复用率**: 约 24% 的业务逻辑通过共享数据层实现

### 复用收益

1. **代码一致性**: 借款方和贷款方使用相同的数据模型
2. **维护成本降低**: 修改业务逻辑只需在一处更新
3. **数据完整性**: 统一的数据验证和业务规则
4. **开发效率**: 新功能可以复用现有的数据模型和仓储

## 架构优势

### 1. 关注点分离

- **应用服务层**: 专注于业务流程编排和用例实现
- **共享数据层**: 专注于数据模型和仓储逻辑
- **基础设施层**: 专注于技术实现细节

### 2. 代码复用

- **模型复用**: 借款方和贷款方共享相同的数据结构
- **数据访问复用**: 数据库操作逻辑统一实现
- **业务规则复用**: 验证逻辑和计算逻辑集中管理

### 3. 可测试性

```go
// 共享数据层可以独立测试
func TestLoanRepo_GetLoan(t *testing.T) {
    db := setupTestDB()
    repo := NewLoanRepo(db)

    loan, err := repo.GetLoan(context.Background(), 1)
    assert.NoError(t, err)
    assert.NotNil(t, loan)
}
```

## 最佳实践

### 1. 数据模型设计原则

- **富模型**: 将业务行为封装在模型中
- **值对象**: 使用不可变值对象表示概念
- **聚合根**: 明确定义聚合边界和一致性规则

### 2. 数据访问层设计原则

- **接口分离**: 定义清晰的数据访问接口
- **泛型复用**: 使用泛型减少重复代码
- **查询优化**: 在共享数据层封装复杂查询逻辑

### 3. 服务集成原则

- **依赖注入**: 通过接口注入共享数据层依赖
- **边界清晰**: 服务只负责自己的业务用例
- **事务管理**: 在服务层管理事务边界

## 扩展指南

### 添加新的数据实体

1. 在 `domains/zeebee/model/` 中定义新模型
2. 在 `domains/zeebee/repo/` 中实现对应数据访问层
3. 在需要使用的服务中引用数据模型和数据访问层

### 服务特化

```go
// 服务可以在共享数据模型基础上添加特定逻辑
type BorrowerLoanService struct {
    *LoanRepo  // 复用数据访问层
    borrowerSpecificLogic  // 借款方特有逻辑
}
```

### 新服务接入

```go
// 新服务可以直接使用现有数据模型
type NewServiceLogic struct {
    loanRepo *repo.LoanRepo    // 复用贷款数据访问
    userRepo *repo.UserRepo     // 复用用户数据访问
    // 新服务特定逻辑
}
```

## 架构权衡

### 优势

- ✅ **代码复用**: 减少重复开发和维护成本
- ✅ **一致性**: 统一的数据模型和业务规则
- ✅ **可测试性**: 共享数据层可以独立测试
- ✅ **扩展性**: 新服务可以快速接入现有共享数据层

### 注意事项

- ⚠️ **耦合度**: 服务间通过共享数据层产生一定耦合
- ⚠️ **变更影响**: 共享数据层变更会影响所有使用服务
- ⚠️ **团队协作**: 需要跨团队协调共享数据层变更

## 未来演进方向

### 短期优化

- 完善数据模型的业务方法
- 增加共享数据层单元测试覆盖率
- 优化数据访问查询性能

### 中期规划

- 引入数据同步机制，实现服务间数据一致性
- 实现读写分离 (CQRS) 模式
- 添加数据服务抽象层

### 长期愿景

- 按业务子域进一步拆分数据模型
- 完善微服务架构的数据层设计
- 引入事件溯源 (Event Sourcing) 模式（如果业务需要）

## 相关文档

- **[系统架构分析](../analysis/architecture_analysis_20250124.md)**: 整体架构评估
- **[项目结构指南](project_structure_guide.md)**: 详细的代码组织说明
- **[API 文档索引](api_documentation_index.md)**: 完整的 API 接口文档

---

*本文档最后更新时间: 2025-11-24*