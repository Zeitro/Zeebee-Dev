# Zeebee Backend

[![Go](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org)
[![go-zero](https://img.shields.io/badge/go--zero-1.10.1-green.svg)](https://go-zero.dev/)

Microservices loan management system backend built on go-zero.

## Features

- **Microservices**: Borrower, Lender, Admin + standalone PDF Generator
- **Single binary per service**: one `cmd/main.go` boots both API (HTTP) and RPC (gRPC)
- **Shared Domain Layer**: Unified models, repos, logic in `domains/zeebee/`
- **Tech Stack**: Go 1.26, go-zero 1.10.1, PostgreSQL, Redis, AWS S3, gRPC
- **Code Generation**: go-zero based API/RPC scaffolding via Makefile

## Quick Start

### Prerequisites

- Go 1.26+
- PostgreSQL 14+, Redis 7+
- goctl CLI

### Setup

```bash
# Install tools
go install github.com/zeromicro/go-zero/tools/goctl@latest
goctl env check --install --verbose --force
go install github.com/linabellbiu/goctl-validate@latest

# Start services (each boots API + RPC in one process)
cd app/borrower && go run cmd/main.go -f etc/config-local.yaml  # :10000 / :10001
cd app/lender   && go run cmd/main.go -f etc/config-local.yaml  # :11000 / :11001
cd app/admin    && go run cmd/main.go -f etc/config-local.yaml  # :12000 / :12001
cd app/pdf_generator && go run cmd/main.go -f etc/config.yaml   # standalone
```

### Code Generation

```bash
make create-service <name>       # Generate API + RPC + Swagger
make create-api-service <name>   # Generate API service only
make create-rpc-service <name>   # Generate RPC service only
make create-swagger <name>       # Generate Swagger docs
make create-sql <service> <name> # Create timestamped SQL file in deploy/<service>/sql/
```

## Architecture

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Borrower API │  │  Lender API  │  │  Admin API   │
│   :10000     │  │   :11000     │  │   :12000     │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Borrower RPC │  │  Lender RPC  │  │  Admin RPC   │
│   :10001     │  │   :11001     │  │   :12001     │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       └────────────┬─────┴───────────────┘
                    ▼
          ┌─────────────────┐
          │  domains/zeebee │  Shared Domain Layer
          │  model/ repo/   │  (38 models, 38 repos, 34 DTOs)
          │  logic/ dto/    │
          └────────┬────────┘
                   ▼
     PostgreSQL + Redis + AWS S3
```

**Data flow**: `HTTP → Handler → API Logic → RPC → RPC Logic → Repo → Model → DB`

**Layer responsibilities**:
- **API Layer**: Request validation, parameter formatting, RPC forwarding
- **RPC Layer**: Core business logic, database operations, transaction management
- **Domain Layer**: Shared models, repositories, business rules
- **Common Layer**: Utilities, error handling, validation, storage

## Project Structure

```
backend/
├── app/                    # Microservices
│   ├── borrower/           # Borrower service (20 API modules, 21 RPC modules)
│   ├── lender/             # Lender service (13 API modules, 14 RPC modules)
│   ├── admin/              # Admin service (13 API modules, 14 RPC modules)
│   └── pdf_generator/      # PDF generation (standalone binary)
├── domains/zeebee/         # Shared domain layer
│   ├── model/              # 38 GORM models
│   ├── repo/               # 38 repositories + provider.go (DI)
│   ├── data/dto/           # 34 DTOs
│   ├── logic/              # Shared business logic
│   ├── auth/               # Authentication
│   ├── templater/          # Template engine
│   └── zeitro/             # Zeitro API integration
├── common/                 # 22 shared utility packages
│   ├── copierx/            # Generic data conversion (copierx.Trans)
│   ├── errorx/             # Unified error handling
│   ├── gormx/              # GORM extensions (BaseRepo, pagination)
│   ├── validator/          # Input validation
│   ├── preapproval_engine/ # Pre-approval engine
│   └── ...                 # aws_s3, email, limiter, rbac, uuid, swagger, etc.
├── deploy/                 # Dockerfile + docker-compose per service
├── scripts/                # Build and management scripts
├── third_party/            # Proto definitions, goctl templates
└── Makefile                # Code generation commands
```

## Development Workflow

### API Service

1. Define messages in `app/{service}/api/desc/*.api`
2. Define routes in `app/{service}/api/{service}.api`
3. `make create-api-service {service}` to generate code
4. Implement logic in `app/{service}/api/internal/logic/`

### RPC Service

1. Define proto in `app/{service}/rpc/pb/*.proto`
2. `make create-rpc-service {service}` to generate code
3. Implement logic in `app/{service}/rpc/internal/logic/`

### Domain Model

1. Add model in `domains/zeebee/model/`
2. Add repo in `domains/zeebee/repo/`
3. Add shared logic in `domains/zeebee/logic/` (if needed)
4. Add DTO in `domains/zeebee/data/dto/` (if needed)

## Testing

```bash
go test ./...                                    # All tests
go test ./common/utils/...                       # Specific package
go test ./... -coverprofile=coverage.out         # With coverage
go tool cover -html=coverage.out                 # Coverage report
golangci-lint run                                # Lint
```

## Deployment

```bash
docker-compose -f deploy/{service}/docker-compose.yml up -d
```

| Environment | Config |
|-------------|--------|
| Local | `app/{service}/etc/config-local.yaml` |
| Staging | `app/{service}/etc/config-staging.yaml` |
| Production | `app/{service}/etc/config-production.yaml` |

## Documentation

- **Architecture**: `.claude/docs/analysis/architecture.md`
- **API Index**: `.claude/docs/guides/api_documentation_index.md`
- **Knowledge Base**: `.claude/docs/INDEX.md`
- **Dev Standards**: `CLAUDE.md`
- **Swagger**: `app/{service}/api/swagger/`

---

*Last updated: 2026-06-16*


