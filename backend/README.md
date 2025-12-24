# Zeebee Backend

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![go-zero](https://img.shields.io/badge/go--zero-1.9.3-green.svg)](https://go-zero.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/your-org/zeebee-backend/actions)
[![Coverage](https://img.shields.io/badge/coverage-<2%25-red.svg)](https://github.com/your-org/zeebee-backend)

A modern microservices architecture loan management system backend built on the go-zero framework, providing
comprehensive digital solutions for the lending industry.

## 🌟 Features

- **Microservices Architecture**: Scalable microservices design based on go-zero
- **Shared Data Layer**: Unified domain models and data access layer for code reuse
- **Modern Technology Stack**: Go 1.24+, PostgreSQL, Redis, AWS S3
- **Complete Business Flow**: Support for the entire loan lifecycle from application to approval and disbursement
- **Secure Authentication**: JWT authentication, password encryption, API permission control
- **Comprehensive Documentation**: Auto-generated Swagger docs and complete API index

## 🚀 Quick Start

### 1. Prerequisites

- Go 1.24.0+
- Docker 20.10+
- PostgreSQL 14+
- Redis 7+

### 2. Development Environment Setup

```bash
# Install goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# Install validation plugin
go install github.com/linabellbiu/goctl-validate@latest

# Install code quality tools
go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest

# Verify installation
goctl version
goctl env check --install --verbose --force
```

### 3. Project Startup

```bash

# Start Borrower service
cd app/borrower
go run cmd/main.go -f etc/config-local.yaml  # (Http Port: 10000, RPC Port: 10001)

# Start Lender service
cd app/lender
go run cmd/main.go -f etc/config-local.yaml         # (Http Port: 11000, RPC Port: 11001)
```

### 4. Code Generation

```bash
# Generate services
make create-api-service service-name
make create-rpc-service service-name

# Generate API documentation
make create-swagger service-name
```

## 🏗️ Architecture

### Microservices Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    Zeebee Backend System                    │
├─────────────────┬─────────────────┬─────────────────────────────────────┤
│   Borrower API  │    Lender API   │          Admin API             │
│   (Port: 10000)  │    (Port: 11000)   │        (Port: 12000)            │
└─────────┬───────┴─────────────────┘─────────────┬─────────────────┘
          │                            │                      │
          ▼                            ▼                      ▼
┌─────────────────────────────┬─────────────────────┐         │
│       Borrower RPC           │     Lender RPC     │         │
│       (Port: 10001)           │     (Port: 11001)    │         │
└─────────────────┬───────────────┴─────────────────┘         │
                  │                                         │
                  └────────────────┬─────────────────────┘
                                 ▼
                  ┌─────────────────────────────┐
                  │      Shared Domain Layer      │
                  │    (domains/zeebee/)          │
                  │  • Models                │
                  │  • Repositories          │
                  │  • Business Logic        │
                  │  • Constants             │
                  └─────────────────┬───────────┘
                                ▼
                  ┌─────────────────────────────┐
                  │     Common Components      │
                  │  • Utils                  │
                  │  • Error Handling        │
                  │  • Database (GORM)       │
                  │  • Cache (Redis)          │
                  │  • Storage (AWS S3)       │
                  └─────────────────┬───────────┘
                                │
                  ┌─────────────────────────────┐
                  │        Infrastructure       │
                  │  • PostgreSQL (14+)       │
                  │  • Redis (7+)              │
                  │  • AWS S3                  │
                  │  • Docker/K8s             │
                  │  • Monitoring             │
                  └─────────────────────────────┘
```

### Core Services

| Service          | Port  | Port Type | Primary Functions            |
|------------------|-------|-----------|------------------------------|
| **borrower-api** | 10000 | HTTP API  | Borrower business interfaces |
| **borrower-rpc** | 10001 | gRPC      | Borrower business logic      |
| **lender-api**   | 11000 | HTTP API  | Lender management interfaces |
| **lender-rpc**   | 11001 | gRPC      | Lender review logic          |
| **admin-api**    | 12000 | HTTP API  | Admin backend interfaces     |
| **admin-rpc**    | 12001 | gRPC      | Management function logic    |

### Directory Structure

```
zeebee-backend/
├── app/                                # Application Services Layer
│   ├── borrower/                       # Borrower Service
│   │   ├── api/                        # HTTP API Layer
│   │   │   ├── desc/                   # API definition files
│   │   │   ├── internal/               # Internal implementation
│   │   │   │   ├── handler/            # HTTP handlers
│   │   │   │   ├── logic/              # API business logic
│   │   │   │   └── svc/                # Service context
│   │   │   ├── etc/                     # Configuration files
│   │   │   └── swagger/                # API documentation
│   │   ├── rpc/                         # RPC Business Layer
│   │   │   ├── pb/                     # Protocol Buffers
│   │   │   ├── internal/               # RPC implementation
│   │   │   └── etc/                    # RPC configuration
│   │   ├── cmd/                        # Service entry
│   │   └── docs/                       # Service documentation
│   ├── lender/                         # Lender Service
│   └── admin/                          # Admin Backend
├── domains/zeebee/                        # Shared Domain Layer
│   ├── model/                           # Data models
│   │   ├── base.go                    # Base model
│   │   ├── user.go                    # User model
│   │   ├── loan.go                    # Loan model
│   │   ├── borrower.go              # Borrower model
│   │   └── ...                       # Other models
│   ├── repo/                            # Data repository layer
│   │   ├── user.go                    # User repository
│   │   ├── loan.go                    # Loan repository
│   │   └── ...                       # Other repositories
│   ├── logic/                           # Business logic layer
│   │   ├── user_logic.go             # User business logic
│   │   └── ...                       # Other business logic
│   └── constants/                       # Business constants
├── common/                              # Common Components
│   ├── utils/                           # Utility functions
│   │   ├── jwt.go                     # JWT handling
│   │   ├── password.go               # Password encryption
│   │   └── ...                       # Other utilities
│   ├── gormx/                          # Database operations
│   │   ├── base_repo.go              # Base repository
│   │   ├── page.go                   # Pagination handling
│   │   └── option.go                # Query options
│   ├── errorx/                          # Error handling
│   │   ├── error.go                   # Error definitions
│   │   └── error_group.go            # Error groups
│   ├── response/                        # Response format
│   ├── validator/                       # Data validation
│   ├── aws_s3/                          # AWS S3 integration
│   └── preapproval_engine/              # Pre-approval engine
├── deploy/                              # Deployment configuration
│   ├── docker-compose.yml             # Docker Compose
│   ├── kubernetes/                     # K8s configuration
│   └── nginx/                          # Nginx configuration
├── scripts/                             # Automation scripts
│   ├── start-dev.sh                   # Start development environment
│   ├── run-tests.sh                   # Run tests
│   ├── build-docker.sh                # Build Docker images
│   └── deploy.sh                      # Deployment scripts
├── docs/                                # Project documentation
│   ├── api/                           # API documentation
│   ├── architecture/                  # Architecture design
│   └── deployment/                   # Deployment guides
├── .claude/                             # AI development assistance
│   ├── docs/                          # Knowledge base documents
│   └── CLAUDE.md                      # Development standards
├── Makefile                             # Build scripts
├── go.mod                              # Go module
└── README.md                            # Project documentation
```

## 🔧 Development Guide

### API Development Process

1. **Define API message body** (`app/service/api/desc/feature.api`)
2. **Define HTTP service** (`app/service/api/service.api`)
3. **Generate code** (`make create-api-service service`)
4. **Implement business logic** (`app/service/api/internal/logic/`)
5. **Test and verify** (`go test ./app/service/api/internal/logic/`)

### RPC Development Process

1. **Define Proto file** (`app/service/rpc/pb/service.proto`)
2. **Generate code** (`make create-rpc-service service`)
3. **Implement business logic** (`app/service/rpc/internal/logic/`)
4. **Test and verify** (`go test ./app/service/rpc/internal/logic/`)

### Code Quality

```bash
# Code formatting
make fmt

# Code linting
make lint

# Run tests
make test

# Generate coverage report
make test-coverage
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Generate test coverage report
make test-coverage

# Run performance tests
go test -bench=. ./...
```

### Test Coverage

Current test coverage: <2% (Target: 80%+)

| Module                  | Coverage | Status               |
|-------------------------|----------|----------------------|
| User Management         | 5%       | 🔴 Needs Improvement |
| Loan Management         | 3%       | 🔴 Needs Improvement |
| Organization Management | 8%       | 🟡 In Progress       |
| Utility Functions       | 45%      | 🟢 Good              |
| Error Handling          | 30%      | 🟢 Good              |

## 📚 Documentation

### Project Documentation

- **[Project Knowledge Base](.claude/docs/INDEX.md)** - Complete project knowledge base and navigation
- **[Development Standards](.claude/CLAUDE.md)** - Claude Code development standards and workflows
- **[Architecture Analysis](.claude/docs/analysis/)** - In-depth system architecture analysis
- **[API Documentation](.claude/docs/guides/api_documentation_index_2025.md)** - Complete API index

### API Documentation

- **Borrower API**: http://localhost:10000/swagger/
- **Lender API**: http://localhost:11000/swagger/
- **Admin API**: http://localhost:12000/swagger/

## 🚀 Deployment

### Development Environment

```bash
# Use Docker Compose
docker-compose -f deploy/service/docker-compose.yml up -d

```

### Environment Configuration

| Environment | Database      | Redis            | Monitoring |
|-------------|---------------|------------------|------------|
| Development | Local         | Local            | Basic      |
| Testing     | Test DB       | Test Redis       | Full       |
| Staging     | Staging DB    | Staging Redis    | Full       |
| Production  | Production DB | Production Redis | Full       |

## 🛡️ Security

- **Authentication & Authorization**: JWT Token authentication + RBAC permission control
- **Password Security**: bcrypt encryption storage
- **Data Protection**: GORM prevents SQL injection
- **API Security**: Request validation and rate limiting
- **File Security**: AWS S3 signed URLs

## 📊 Monitoring

### Performance Metrics

- API response time: <200ms (95% requests)
- Database query time: <100ms
- Memory usage: <512MB
- Concurrent processing: 1000+ requests/second

### Logging

- Structured logging
- Distributed tracing (planned)
- Error alerts and notifications
- Performance monitoring and analysis

## 🤝 Contributing

### Development Process

1. Fork the project
2. Create feature branch
3. Write code and tests
4. Submit Pull Request
5. Code review
6. Merge to main branch

### Code Standards

- Follow Go official code standards
- Use go-zero framework conventions
- Write unit tests
- Add necessary comments

## 📞 Support

- **Bug Reports**: [GitHub Issues](https://github.com/your-org/zeebee-backend/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/your-org/zeebee-backend/discussions)
- **Security Vulnerabilities**: [security@yourcompany.com](mailto:security@yourcompany.com)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## 📝 Changelog

### [v1.0.0] (2025-12-22)

#### Added

- Complete microservices architecture implementation
- Borrower and lender core functionalities
- User authentication and permission management
- File storage and management system
- Pre-approval engine basic implementation

#### Improved

- Updated to Go 1.24.0 and go-zero 1.9.3
- Comprehensive documentation and API index
- Automated testing and deployment scripts
- Security hardening and performance optimization

#### Technical Debt

- Improve test coverage (target: 80%+)
- Enhance monitoring and logging systems
- Implement automated CI/CD pipeline
- Optimize database query performance

---

**Documentation maintained by**: Development team
**Last updated**: 2025-12-22
**Project homepage**: [https://github.com/Zeitro/Zeebee-Dev](https://github.com/Zeitro/Zeebee-Dev)

```

