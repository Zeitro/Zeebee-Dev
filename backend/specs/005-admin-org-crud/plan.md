# Implementation Plan: 管理后台组织管理 CRUD

**Branch**: `005-admin-org-crud` | **Date**: 2025-01-12 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `backend/specs/005-admin-org-crud/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the
execution workflow.

## Summary

实现管理后台对 Lender（贷款方）组织的完整 CRUD 功能，包括创建组织、查看详情、编辑信息、切换状态和删除组织。技术方案基于 go-zero
微服务架构，遵循 API/RPC 双层设计模式，使用共享领域层进行数据访问，所有数据类型转换使用 `copierx.Trans[T, U]`。

## Technical Context

**Language/Version**: Go 1.23+
**Primary Dependencies**: go-zero 1.9.3, GORM, Redis, AWS S3
**Storage**: PostgreSQL 14+
**Testing**: go test, golangci-lint
**Target Platform**: Linux server (microservices architecture)
**Project Type**: microservices (API/RPC separation)
**Performance Goals**: <200ms API response (p95), 1000+ req/s
**Constraints**: API/RPC dual-layer, shared domain models, copierx.Trans for conversions
**Scale/Scope**: Multiple services (borrower, lender, admin), shared domain layer

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Based on `.specify/memory/constitution.md`:

- [x] **Shared Domain-First**: 新模型/仓储使用 `domains/zeebee/` 中现有的
  OrganizationLogic、UserLogic、OrganizationUserRepo
- [x] **API/RPC Separation**: 功能遵循 API（验证）vs RPC（业务逻辑）边界
- [x] **Code Generation**: API 定义使用 `.api` 文件，已通过 `goctl` 生成模板代码
- [x] **Data Conversion**: 所有类型转换使用 `copierx.Trans[T, U]`
- [x] **Type Safety**: 单元测试覆盖新功能（目标：80%+）
- [x] **Test Coverage**: 实施计划包含单元测试和集成测试

**Violations Requiring Justification**:

| Principle           | Potential Violation | Justification Required?         |
|---------------------|---------------------|---------------------------------|
| Shared Domain-First |                     | ✅ 使用现有 `domains/zeebee/` 中的共享逻辑 |
| API/RPC Separation  |                     | ✅ 遵循现有 API/RPC 双层架构模式           |
| Code Generation     |                     | ✅ 已生成 API/RPC 模板代码              |
| Data Conversion     |                     | ✅ 所有转换使用 `copierx.Trans[T, U]`  |

**结论**: 所有宪章原则都已满足，无需额外理由说明。

## Project Structure

### Documentation (this feature)

```text
backend/specs/005-admin-org-crud/
├── plan.md              # 本文件 (/speckit.plan 命令输出)
├── research.md          # Phase 0 输出
├── data-model.md        # Phase 1 输出
├── quickstart.md        # Phase 1 输出
├── contracts/           # Phase 1 输出 (API/RPC 定义)
└── tasks.md             # Phase 2 输出 (/speckit.tasks 命令)
```

### Source Code (repository root)

```text
# Zeebee Backend - Microservices Architecture
app/
├── admin/              # Admin 服务（本功能所在服务）
│   ├── api/           # HTTP API 层
│   │   ├── desc/      # API 定义 (.api 文件)
│   │   │   └── organization.api  # 组织管理 API 定义
│   │   └── internal/
│   │       ├── handler/   # HTTP 处理器（已生成）
│   │       ├── logic/     # API 逻辑（验证、路由）
│   │       │   └── organization/
│   │       │       ├── organization_page_logic.go         # 已实现
│   │       │       ├── organization_create_lender_logic.go  # 待实现
│   │       │       ├── organization_detail_logic.go       # 待实现
│   │       │       ├── organization_update_lender_logic.go # 待实现
│   │       │       ├── organization_switch_status_logic.go # 待实现
│   │       │       └── organization_delete_logic.go       # 待实现
│   │       ├── types/     # API 类型定义
│   │       └── svc/       # Service context
│   └── rpc/           # gRPC RPC 层
│       └── internal/
│           └── logic/ # RPC 逻辑（业务逻辑）
│               └── organization/
│                   ├── page_logic.go         # 已实现
│                   ├── create_lender_logic.go # 待实现
│                   ├── detail_logic.go       # 待实现
│                   ├── update_lender_logic.go # 待实现
│                   ├── switch_status_logic.go # 待实现
│                   └── delete_logic.go       # 待实现

domains/zeebee/         # 共享领域层
├── logic/             # 业务逻辑
│   ├── organization_logic.go     # 组织业务逻辑（已存在）
│   ├── user_logic.go             # 用户业务逻辑（参考 RegisterWithOrganization）
│   └── organization_user_logic.go # 组织用户关系逻辑（已存在）
├── model/             # 统一数据模型
│   ├── organization.go           # 组织模型
│   ├── user.go                   # 用户模型
│   └── organization_user.go      # 组织用户关系模型
└── repo/              # 共享仓储
    ├── organization_repo.go      # 组织仓储
    ├── user_repo.go              # 用户仓储
    └── organization_user_repo.go # 组织用户仓储

common/                # 共享工具
├── errorx/            # 错误处理
├── copierx/           # 类型转换（copierx.Trans[T, U]）
├── validator/         # 数据验证
├── gormx/             # 数据库操作
└── constants/         # 业务常量
    └── organization.go            # 组织状态常量

tests/                 # 测试文件（与源码并列）
```

**Structure Decision**: 微服务架构 + API/RPC 双层设计。共享领域模型在 `domains/zeebee/`，共享工具在 `common/`。本功能在
admin 服务中实现。

## Complexity Tracking

> **无需填写** - Constitution Check 全部通过，无违规需要说明

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| N/A       | N/A        | N/A                                  |
