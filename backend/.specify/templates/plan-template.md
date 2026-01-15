# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the
execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

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

- [ ] **Shared Domain-First**: New models/repositories in `domains/zeebee/`? Or service-specific?
- [ ] **API/RPC Separation**: Does feature respect API (validation) vs RPC (business logic) boundary?
- [ ] **Code Generation**: Will `.api`/`.proto` definitions be used with `goctl`?
- [ ] **Data Conversion**: All type conversions using `copierx.Trans[T, U]`?
- [ ] **Type Safety**: Are unit tests included for `common/` utilities?
- [ ] **Test Coverage**: Does implementation plan include adequate tests (target: 80%+)?

**Violations Requiring Justification**:

| Principle           | Potential Violation | Justification Required?                          |
|---------------------|---------------------|--------------------------------------------------|
| Shared Domain-First |                     | If NOT using `domains/zeebee/` for shared models |
| API/RPC Separation  |                     | If blurring API/RPC responsibilities             |
| Code Generation     |                     | If NOT using goctl code generation               |
| Data Conversion     |                     | If NOT using copierx.Trans for conversions       |

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# Zeebee Backend - Microservices Architecture
app/
├── borrower/           # Borrower service
│   ├── api/           # HTTP API layer
│   │   ├── desc/      # API definitions (.api files)
│   │   ├── internal/
│   │   │   ├── handler/   # HTTP handlers
│   │   │   ├── logic/     # API logic (validation, routing)
│   │   │   └── svc/       # Service context
│   │   └── etc/       # Configuration
│   └── rpc/           # gRPC RPC layer
│       ├── pb/        # Protocol buffers (.proto files)
│       ├── internal/
│       │   └── logic/ # RPC logic (business logic)
│       └── etc/       # Configuration
├── lender/            # Lender service (same structure as borrower)
└── admin/             # Admin service (same structure as borrower)

domains/zeebee/         # Shared domain layer
├── model/             # Unified data models
├── repo/              # Shared repositories
└── constants/         # Business constants

common/                # Shared utilities
├── errorx/            # Error handling
├── copierx/           # Type conversion
├── validator/         # Data validation
├── gormx/             # Database operations
└── utils/             # Utilities (JWT, password, etc.)

tests/                 # Test files (alongside source)
```

**Structure Decision**: Microservices with API/RPC dual-layer architecture. Shared domain models in `domains/zeebee/`,
shared utilities in `common/`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation                                               | Why Needed         | Simpler Alternative Rejected Because   |
|---------------------------------------------------------|--------------------|----------------------------------------|
| [e.g., Service-specific model instead of shared domain] | [current need]     | [why `domains/zeebee/` insufficient]   |
| [e.g., Direct data access instead of copierx.Trans]     | [specific problem] | [why standard conversion insufficient] |
