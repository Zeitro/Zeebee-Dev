# Project Index: Zeebee Backend

Generated: 2026-06-24 | Go 1.23 | go-zero 1.10 | ~172k LoC

## 📁 Project Structure

```
backend/
├── app/                  # 4 microservices
│   ├── borrower/         # 借款方 (439 go, 68k LoC)
│   ├── lender/           # 贷款方 (247 go, 36k LoC)
│   ├── admin/            # 管理后台 (299 go, 43k LoC)
│   └── pdf_generator/    # PDF 生成 (14 go, 1k LoC)
├── domains/zeebee/       # Shared domain (172 go, 20k LoC)
├── common/               # 22 utility packages (4.7k LoC)
├── deploy/               # Docker Compose per service
├── scripts/              # Build/management scripts
├── specs/                # Feature specs
├── third_party/          # Third-party proto (google/api)
├── .claude/docs/         # AI-generated docs + analysis
├── .specify/memory/      # Project constitution
├── go.mod / go.sum       # Dependencies
├── Makefile              # Code generation + build
└── CLAUDE.md             # Dev conventions
```

## 🚀 Entry Points

| Service | Entry | HTTP | gRPC |
|---------|-------|------|------|
| Borrower | `app/borrower/cmd/main.go` | 10000 | 10001 |
| Lender | `app/lender/cmd/main.go` | 11000 | 11001 |
| Admin | `app/admin/cmd/main.go` | 12000 | 12001 |
| PDF Generator | `app/pdf_generator/cmd/main.go` | 13000 | — |

## 📦 Core Modules

### Shared Domain (`domains/zeebee/`)
- **model/** (38 files) — GORM entities: User, Loan, Organization, Borrower*, PricingPlan, etc.
- **repo/** (38 files) — Generic BaseRepo pattern, `provider.go` for DI
- **logic/** — Shared business logic (UserLogic, OrganizationUserLogic, etc.)
- **data/dto/** (26+ files) — DTOs for request/response mapping
- **auth/** — JWT claim, RBAC, Casbin integration
- **constants/** — Business constants + Redis key format strings
- **email/, limiter/, templater/, drive/, lenderdiscover/, zeitro/**

### Common Utilities (`common/`, 22 packages)

| Package | Purpose | Key API |
|---------|---------|---------|
| `copierx` | Generic type conversion | `Trans[T,U](*T) → *U` |
| `errorx` | Error wrapping + stack traces | `New()`, `Wrap()`, `WithStack()` |
| `gormx` | Generic GORM BaseRepo | `BaseRepo[T]`, `Page()`, `Transaction()` |
| `redisx` | Redis client wrapper | `Get()`, `SetExpire()`, `GetTTL()`, `SetNX()` |
| `timex` | Duration ↔ String | `StringToDuration()`, `DurationToString()` |
| `stringx` | String formatting | `Format()` (sprintf wrapper), `FormatFullName()` |
| `validator` | Request validation | Struct tag-based |
| `response` | Unified HTTP response | JSON envelope |
| `cron` | Cron engine | Atomic task scheduling |
| `rbac` | Role-based access control | Casbin adapter |
| `preapproval_engine` | Loan pre-approval rules | Decision engine |
| `aws_s3` | AWS S3 file storage | Upload/download/presigned URL |
| `limiter` | Rate limiting | Redis-backed |
| `uuid` | ID generation | Sonyflake |

### App Layer Pattern (per service)
```
app/{service}/
├── cmd/main.go              # Entry point
├── api/internal/            # HTTP layer
│   ├── handler/{module}/    # Route handlers (generated)
│   ├── logic/{module}/      # API logic
│   ├── middleware/           # Auth, CORS, etc.
│   └── svc/                 # ServiceContext (DI container)
└── rpc/internal/            # gRPC layer
    ├── logic/{module}/      # Core business logic
    ├── server/              # gRPC server impl
    └── svc/                 # ServiceContext
```

## 🔧 Configuration

| Type | Path | Purpose |
|------|------|---------|
| Main Makefile | `Makefile` | Code gen (`make create-service`), build, lint |
| Go module | `go.mod` | zeebee-backend, Go 1.23, 197 deps |
| Service configs | `app/{svc}/etc/config-*.yaml` | Per-environment (local/staging/production) |
| Deploy | `deploy/{svc}/docker-compose.yml` | Docker Compose per service |
| Linting | `.golangci.yml` | golangci-lint rules |
| Constitution | `.specify/memory/constitution.md` | Git format, tech stack, code standards |

## 📚 Documentation

| Document | Purpose |
|----------|---------|
| `CLAUDE.md` | **Read first** — dev conventions, data conversion, error handling |
| `.claude/docs/analysis/architecture.md` | Architecture: layers, data flow, directory structure |
| `.claude/docs/analysis/analysis.md` | Deep audit (2026-06-24): 90 findings, 6.2/10 |
| `.specify/memory/constitution.md` | Git commit format, tech stack, governance |
| `.claude/docs/INDEX.md` | Knowledge base navigation |
| `.claude/docs/guides/` | API index, project structure, domain architecture |

## 🧪 Test Coverage

- Unit tests: 27 files (10 in timex, remaining across domains/common)
- Integration/DB: minimal
- Coverage: < 10%

## 🔗 Key Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| go-zero | 1.10.1 | Microservice framework (REST + gRPC) |
| gorm.io/gorm | v2 | ORM (PostgreSQL, MySQL) |
| redis/go-redis | v9.19.0 | Redis client |
| google.golang.org/grpc | v1.81.1 | gRPC |
| chromedp | v0.15.1 | Headless Chrome for PDF |
| golang-jwt/jwt | v4.5.2 | JWT authentication |
| duke-git/lancet | v2.3.9 | Go utility library |
| casbin | v3 | RBAC authorization |

## 📝 Quick Start

1. Read `CLAUDE.md` + `.claude/docs/analysis/architecture.md`
2. `make create-service <service>` — generate code from .api/.proto
3. `go build ./app/<service>/...` — build
4. `go test ./common/... ./domains/...` — test

## 🏷️ Key Stats

| Metric | Value |
|--------|-------|
| Go source files | 1,205 |
| Test files | 27 |
| Proto definitions | 20 |
| API definitions | 82 |
| Generated pb.go | 6 |
| YAML configs | 24 |
| Docker files | 8 |
| Markdown docs | 81 |
| Common packages | 22 |
| Domain models | 38 |
| Domain repos | 38 |
| Lines of code | ~172k |
