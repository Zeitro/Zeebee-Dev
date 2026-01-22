# Implementation Plan: Lender Role and Organization User Management

**Branch**: `001-lender-role-org` | **Date**: 2025-12-29 | **Spec**: [spec.md](./spec.md)
**Status**: ✅ Completed | **Completed**: 2025-12-29
**Input**: Feature specification from `/specs/001-lender-role-org/spec.md`

## Summary

Implement role and organization user management features for the lender service, mirroring the existing borrower service
implementation. This includes creating role CRUD operations, organization user management with invitation flow, and role
assignment capabilities. The feature leverages existing shared domain models (role, organization_user, resource) from
the domains layer.

**Implementation Result**: All 86 tasks completed successfully. The feature has been merged to
`feat-backend-organization-and-permissions` branch.

## Technical Context

**Language/Version**: Go 1.23.0
**Primary Dependencies**: go-zero framework, GORM, PostgreSQL
**Storage**: PostgreSQL with existing shared schema (domains/zeebee/model)
**Testing**: Go testing with `go test` command
**Target Platform**: Linux server (microservice architecture)
**Project Type**: Microservice (web backend)
**Performance Goals**: API responses < 200ms p95, support 1000 concurrent users
**Constraints**: Must follow existing borrower service patterns, use shared domain models
**Scale/Scope**: ~15 new API endpoints, ~20 RPC methods, ~40 handler/logic files

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The project constitution is currently a template without defined principles. Based on project context, the following
gates apply:

| Principle                | Status | Notes                                                    |
|--------------------------|--------|----------------------------------------------------------|
| Follow Existing Patterns | PASS   | Implementation will mirror borrower service architecture |
| Use Shared Domain Models | PASS   | Leveraging existing domains/zeebee layer                 |
| Go-zero Framework        | PASS   | Using goctl for code generation                          |
| API-First Design         | PASS   | API definitions drive implementation                     |

## Project Structure

### Documentation (this feature)

```text
backend/specs/001-lender-role-org/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   ├── role.api         # Role API definitions
│   ├── organization_user.api  # Organization User API definitions
│   └── resource.api     # Resource API definitions (if needed)
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
app/lender/
├── api/
│   ├── desc/            # API definitions (NEW)
│   │   ├── role.api
│   │   ├── organization_user.api
│   │   └── resource.api
│   ├── internal/
│   │   ├── handler/
│   │   │   ├── role/    # NEW - 6 handlers
│   │   │   └── organization_user/  # NEW - 9 handlers
│   │   ├── logic/
│   │   │   ├── role/    # NEW - 6 logic files
│   │   │   └── organization_user/  # NEW - 9 logic files
│   │   └── types/       # Generated types
│   └── lender.api       # Updated with new routes
├── rpc/
│   ├── internal/
│   │   ├── logic/
│   │   │   ├── role/    # NEW - 6 RPC logic files
│   │   │   └── organizationuser/  # NEW - 9 RPC logic files
│   │   └── svc/
│   └── client/          # RPC clients

domains/zeebee/          # SHARED - Already exists, may need extensions
├── model/
│   ├── role.go          # EXISTS
│   ├── organization_user.go  # EXISTS
│   ├── role_and_resource.go   # EXISTS
│   └── resource.go      # EXISTS
└── repo/
    ├── role.go          # EXISTS
    ├── organization_user.go  # EXISTS
    └── role_and_resource.go   # EXISTS
```

**Structure Decision**: Following microservice pattern with clear separation between API layer (HTTP) and RPC layer (
gRPC). Shared domain models from `domains/zeebee` are reused. The implementation mirrors the borrower service structure
for consistency.

## Complexity Tracking

> **No violations to justify** - Implementation follows existing patterns and shared domain models.
