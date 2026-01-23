# Implementation Plan: Admin 模板 CRUD

**Branch**: `6-admin-template-crud` | **Date**: 2026-01-23 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/006-admin-template-crud/spec.md`

## Summary

为 Admin 后端实现模板 CRUD 管理功能，包括基础模板和组织模板的完整生命周期管理。核心功能包括：
- 基础模板的创建、读取、更新、删除（CRUD）
- 模板状态切换（Active/Disabled）
- 应用层实现的级联操作（状态更改和删除级联到组织模板）
- 复杂的 JSONB 变量架构验证
- 分页查询和筛选功能

**技术方法**：
- 使用应用层代码实现级联操作（非数据库触发器）
- 遵循 admin 端现有的分页结构和代码风格
- 通过数据库事务保证级联操作的原子性

## Technical Context

**Language/Version**: Go 1.23+
**Primary Dependencies**: go-zero 1.9.3, GORM, Redis, AWS S3
**Storage**: PostgreSQL 14+
**Testing**: go test, golangci-lint
**Target Platform**: Linux server (microservices architecture)
**Project Type**: microservices (API/RPC separation)
**Performance Goals**: <200ms API response (p95), 1000+ req/s
**Constraints**: API/RPC dual-layer, shared domain models, copierx.Trans for conversions, application-layer cascade
**Scale/Scope**: Admin service only, uses shared domain layer

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Based on `.specify/memory/constitution.md`:

- [x] **Shared Domain-First**: New models/repositories in `domains/zeebee/`
- [x] **API/RPC Separation**: Does feature respect API (validation) vs RPC (business logic) boundary
- [x] **Code Generation**: `.api`/`.proto` definitions used with `goctl`
- [x] **Data Conversion**: All type conversions using `copierx.Trans[T, U]`
- [ ] **Type Safety**: Unit tests for `common/` utilities (N/A - no new common utilities)
- [x] **Test Coverage**: Implementation plan includes adequate tests (target: 80%+)

**Violations Requiring Justification**:

| Principle | Potential Violation | Justification Required? |
|-----------|---------------------|------------------------|
| Shared Domain-First | | No - using `domains/zeebee/` for models/repositories |
| API/RPC Separation | | No - clear API validation → RPC business logic separation |
| Code Generation | | No - using goctl for code generation |
| Data Conversion | | No - using copierx.Trans for all conversions |

**No violations** - All constitutional principles are followed.

## Project Structure

### Documentation (this feature)

```text
specs/006-admin-template-crud/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output - Technical research
├── data-model.md        # Phase 1 output - Data models
├── quickstart.md        # Phase 1 output - Quick start guide
├── contracts/           # Phase 1 output - API contracts
│   └── template.api     # API definition
├── checklists/
│   └── requirements.md  # Quality checklist
└── spec.md              # Feature specification
```

### Source Code (repository root)

```text
app/
├── admin/
│   ├── api/
│   │   ├── desc/
│   │   │   └── template.api          # API definition (from contracts/)
│   │   └── internal/
│   │       ├── handler/template/     # HTTP handlers
│   │       │   ├── create_base_template_handler.go
│   │       │   ├── update_base_template_handler.go
│   │       │   ├── base_template_detail_handler.go
│   │       │   ├── base_template_page_handler.go
│   │       │   ├── base_template_switch_status_handler.go
│   │       │   └── delete_base_template_handler.go
│   │       └── logic/template/       # API logic (validation, routing)
│   │           ├── create_base_template_logic.go
│   │           ├── update_base_template_logic.go
│   │           ├── base_template_detail_logic.go
│   │           ├── base_template_page_logic.go
│   │           ├── base_template_switch_status_logic.go
│   │           └── delete_base_template_logic.go
│   └── rpc/
│       ├── client/template/          # RPC client
│       │   └── template.go
│       └── internal/
│           ├── logic/template/       # RPC logic (business logic)
│           │   ├── create_logic.go
│           │   ├── update_logic.go
│           │   ├── detail_logic.go
│           │   ├── page_logic.go
│           │   ├── switch_status_logic.go
│           │   └── delete_logic.go
│           └── server/template/      # RPC server
│               └── template_server.go
│
domains/zeebee/                         # Shared domain layer
├── constants/
│   └── template.go                   # Template constants (updated)
├── model/
│   ├── base_template.go              # BaseTemplate model
│   └── organization_template.go      # OrganizationTemplate model
├── repo/
│   ├── base_template.go              # BaseTemplate repository
│   └── organization_template.go      # OrganizationTemplate repository
└── logic/
    └── template.go                   # Template business logic (cascade operations)
```

**Structure Decision**: Microservices with API/RPC dual-layer architecture. Shared domain models in `domains/zeebee/`, cascade operations implemented in application layer.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - this section is not applicable.

## Implementation Phases

### Phase 0: Research ✓ Complete

**Output**: [research.md](./research.md)

**Key Decisions**:
- 级联操作使用应用层代码实现（用户明确要求）
- JSONB 验证使用应用层 + 数据库约束
- 标签唯一性使用数据库唯一索引 + 应用层验证
- 分页使用 admin 端现有的 BasePageRequest/BasePageResponse 结构

### Phase 1: Design ✓ Complete

**Outputs**:
- [data-model.md](./data-model.md) - Complete data models
- [contracts/template.api](./contracts/template.api) - API definition
- [quickstart.md](./quickstart.md) - Quick start guide

**Design Artifacts**:
1. Database schema with proper indexes and constraints
2. Go struct definitions for Model, BO, and API layers
3. API contract following admin endpoint patterns
4. Cascade operation implementation patterns

### Phase 2: Implementation (Next - via /speckit.tasks)

The following phases will be detailed in the tasks generated by `/speckit.tasks`:

1. **Domain Layer** (`domains/zeebee/`)
   - Create model definitions
   - Create repository interfaces
   - Implement cascade logic

2. **RPC Layer** (`app/admin/rpc/`)
   - Create protobuf definitions
   - Implement RPC logic
   - Register RPC server

3. **API Layer** (`app/admin/api/`)
   - Generate API code from .api file
   - Implement handlers and logic
   - Wire up dependencies

4. **Testing & Validation**
   - Unit tests for repositories
   - Integration tests for cascade operations
   - E2E tests for API endpoints

## Dependencies

### Internal Dependencies
- `domains/zeebee/constants` - User types, common statuses
- `domains/zeebee/model` - Base models (Base, etc.)
- `common/errorx` - Error handling
- `common/copierx` - Type conversions
- `common/validator` - Data validation

### External Dependencies
- `github.com/zeromicro/go-zero` v1.9.3
- `gorm.io/gorm` v1.25.12
- `gorm.io/driver/postgres` v1.5.9
- `github.com/go-playground/validate/v10` v10.22.0

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Cascade operation data inconsistency | Use database transactions with rollback |
| Large content performance impact | Validate content size (10MB limit) |
| JSONB validation errors | Application layer validation with detailed error messages |
| Concurrent update conflicts | Optimistic locking with updated_at version check |
| Tag uniqueness violations | Unique database index + application pre-check |

## Success Criteria

From [spec.md](./spec.md):

- **SC-001**: Admin users can create a new base template with all required fields in under 60 seconds
- **SC-002**: Template list page loads and displays up to 100 templates in under 2 seconds
- **SC-003**: Search and filter operations return results in under 1 second for databases with up to 10,000 templates
- **SC-004**: Status change operations (enable/disable) complete in under 3 seconds including cascading to all organization templates
- **SC-005**: 100% of base template status changes correctly cascade to all derived organization templates
- **SC-006**: 100% of base template deletions correctly remove all derived organization templates
- **SC-007**: Variable schema validation rejects 100% of invalid JSON structures before saving
- **SC-008**: Admin users report successful template creation on first attempt 95% of the time

---

**Next Step**: Run `/speckit.tasks` to generate detailed implementation task list.
