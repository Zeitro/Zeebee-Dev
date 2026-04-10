# Zeebee 开发工作流程

**创建时间**: 2025-12-22
**适用版本**: Go 1.24+, go-zero 1.9.3+
**最后更新**: 2025-12-22

---

## 📋 概述

本文档定义了 Zeebee 项目的标准化开发工作流程，包括环境设置、代码开发、测试、部署等各个环节的最佳实践。

## 🔧 开发环境设置

### 系统要求

- **操作系统**: Linux/macOS/Windows (推荐 WSL2)
- **Go 版本**: 1.24.0+
- **Docker**: 20.10+
- **Git**: 2.30+

### 必需工具安装

#### 1. Go 开发环境

```bash
# 安装 Go (使用官方安装包)
# 下载地址: https://golang.org/dl/

# 验证安装
go version
go env GOPATH
go env GOROOT
```

#### 2. go-zero 工具链

```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 安装验证插件
go install github.com/linabellbiu/goctl-validate@latest

# 检查环境
goctl env check --install --verbose --force
```

#### 3. 代码质量工具

```bash
# 安装 golangci-lint
go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest

# 安装 staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest

# 安装 goimports
go install golang.org/x/tools/cmd/goimports@latest
```

### IDE 配置

#### VS Code 推荐

```json
// .vscode/settings.json
{
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  "go.gopath": "",
  "go.goroot": "",
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "go.buildOnSave": "workspace",
  "go.testOnSave": true,
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll": true,
    "source.organizeImports": true
  }
}
```

#### 推荐扩展

- Go (Google 官方)
- Go Test Explorer
- Better Comments
- GitLens
- Docker

---

## 🏗️ 项目结构

### 标准目录结构

```
zeebee-backend/
├── app/                       # 应用服务
│   ├── borrower/             # 借款方服务
│   │   ├── api/              # HTTP API 层
│   │   │   ├── desc/        # API 定义文件
│   │   │   ├── internal/    # 内部实现
│   │   │   │   ├── handler/ # HTTP 处理器
│   │   │   │   ├── logic/   # API 业务逻辑
│   │   │   │   └── svc/     # 服务上下文
│   │   │   └── etc/         # 配置文件
│   │   ├── rpc/              # RPC 业务层
│   │   │   ├── pb/           # Protocol Buffers
│   │   │   └── internal/    # RPC 实现
│   │   └── cmd/              # 服务入口
│   └── lender/               # 贷款方服务 (类似结构)
├── domains/                   # 共享领域层
│   └── zeebee/               # 共享数据模型和业务逻辑
│       ├── model/            # 数据模型
│       ├── repo/             # 数据仓储
│       ├── logic/            # 业务逻辑
│       └── constants/        # 业务常量
├── common/                    # 通用组件
│   ├── utils/                # 工具函数
│   ├── gormx/                # 数据库操作
│   ├── errorx/               # 错误处理
│   ├── response/             # 响应格式
│   ├── validator/            # 数据验证
│   └── aws_s3/               # S3 集成
├── third_party/               # 第三方模板
│   └── goctl_template/       # go-zero 代码生成模板
├── .claude/                   # AI 辅助开发配置
├── docs/                      # 项目文档
├── scripts/                   # 脚本工具
├── deploy/                    # 部署配置
├── Makefile                   # 构建脚本
├── go.mod                     # Go 模块定义
└── README.md                  # 项目说明
```

---

## 🔄 开发工作流程

### 1. 功能开发流程

#### API 服务开发

```bash
# 1. 定义 API 消息体
# app/borrower/api/desc/new_feature.api
syntax = "v1"

type NewFeatureRequest {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

type NewFeatureResponse {
    Id   int64  `json:"id"`
    Name string `json:"name"`
}
```

```bash
# 2. 定义 HTTP 服务
# app/borrower/api/borrower.api
import "new_feature.api"

@server(
    prefix: /api
    group:  newfeature
    middleware: AuthMiddleware
)
service borrower-api {
    @handler createNewFeature
    post /new-feature (NewFeatureRequest) returns (NewFeatureResponse)
}
```

```bash
# 3. 生成代码
make create-api-service borrower

# 4. 实现业务逻辑
# app/borrower/api/internal/logic/new_feature_logic.go
func (l *NewFeatureLogic) CreateNewFeature(req *types.NewFeatureRequest) (resp *types.NewFeatureResponse, err error) {
    // 调用 RPC 服务
    rpcResp, err := l.svcCtx.BorrowerClient.CreateNewFeature(l.ctx, &borrower.CreateNewFeatureRequest{
        Name: req.Name,
        Age:  int64(req.Age),
    })
    if err != nil {
        return nil, err
    }

    // 响应转换
    return &types.NewFeatureResponse{
        Id:   rpcResp.Id,
        Name: rpcResp.Name,
    }, nil
}
```

#### RPC 服务开发

```bash
# 1. 定义 Proto
# app/borrower/rpc/pb/borrower.proto
syntax = "proto3";

package borrower;

option go_package = "./pb";

service Borrower {
    rpc CreateNewFeature (CreateNewFeatureRequest) returns (CreateNewFeatureResponse);
}

message CreateNewFeatureRequest {
    string name = 1;
    int64  age = 2;
}

message CreateNewFeatureResponse {
    int64  id = 1;
    string name = 2;
}
```

```bash
# 2. 生成代码
make create-rpc-service borrower

# 3. 实现业务逻辑
# app/borrower/rpc/internal/logic/new_feature_logic.go
func (l *CreateNewFeatureLogic) CreateNewFeature(in *borrower.CreateNewFeatureRequest) (*borrower.CreateNewFeatureResponse, error) {
    // 数据验证
    if in.Name == "" {
        return nil, errorx.NewDefaultError("Name is required")
    }

    // 创建数据模型
    feature := &model.NewFeature{
        Name: in.Name,
        Age:  in.Age,
    }

    // 数据库操作
    result, err := l.svcCtx.NewFeatureRepo.Create(l.ctx, feature)
    if err != nil {
        return nil, errorx.WithStack(err)
    }

    return &borrower.CreateNewFeatureResponse{
        Id:   result.Id,
        Name: result.Name,
    }, nil
}
```

### 2. 数据库操作流程

#### Model 定义

```go
// domains/zeebee/model/new_feature.go
package model

import (
    "zeebee-backend/common/gormx"
    "gorm.io/gorm/schema"
)

type NewFeature struct {
    gormx.Base
    Name string `gorm:"column:name;size:100;not null" json:"name"`
    Age  int64  `gorm:"column:age;not null" json:"age"`
}

func (NewFeature) TableName() string {
    return "new_features"
}
```

#### Repository 操作

```go
// domains/zeebee/repo/new_feature.go
package repo

import (
    "context"

    "zeebee-backend/domains/zeebee/model"
    "zeebee-backend/common/gormx"
)

type NewFeatureRepo struct {
    *gormx.BaseRepo[*model.NewFeature]
}

func NewNewFeatureRepo(db *gorm.DB) *NewFeatureRepo {
    return &NewFeatureRepo{
        BaseRepo: gormx.NewBaseRepo[*model.NewFeature](db),
    }
}

func (repo *NewFeatureRepo) Create(ctx context.Context, feature *model.NewFeature) (*model.NewFeature, error) {
    err := repo.Client(ctx).Create(feature).Error
    return feature, err
}

func (repo *NewFeatureRepo) GetByID(ctx context.Context, id int64) (*model.NewFeature, error) {
    var feature model.NewFeature
    err := repo.Client(ctx).Where("id = ?", id).First(&feature).Error
    return &feature, err
}
```

---

## 🧪 测试流程

### 1. 单元测试

#### 测试文件命名

```
功能文件: new_feature_logic.go
测试文件: new_feature_logic_test.go
```

#### 测试示例

```go
// app/borrower/api/internal/logic/new_feature_logic_test.go
package logic

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "zeebee-backend/app/borrower/api/internal/svc"
    "zeebee-backend/app/borrower/rpc/pb/borrower"
)

// Mock RPC client
type MockBorrowerClient struct {
    mock.Mock
}

func (m *MockBorrowerClient) CreateNewFeature(ctx context.Context, in *borrower.CreateNewFeatureRequest, opts ...grpc.CallOption) (*borrower.CreateNewFeatureResponse, error) {
    args := m.Called(ctx, in, opts)
    return args.Get(0).(*borrower.CreateNewFeatureResponse), args.Error(1)
}

func TestNewFeatureLogic_CreateNewFeature(t *testing.T) {
    // 准备测试数据
    req := &types.NewFeatureRequest{
        Name: "test",
        Age:  25,
    }

    // 创建 mock
    mockClient := new(MockBorrowerClient)
    mockResp := &borrower.CreateNewFeatureResponse{
        Id:   1,
        Name: "test",
    }
    mockClient.On("CreateNewFeature", mock.Anything, mock.Anything, mock.Anything).Return(mockResp, nil)

    // 创建服务上下文
    svcCtx := &svc.ServiceContext{
        BorrowerClient: mockClient,
    }

    // 创建 logic
    logic := NewCreateNewFeatureLogic(context.Background(), svcCtx)

    // 执行测试
    resp, err := logic.CreateNewFeature(req)

    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, resp)
    assert.Equal(t, int64(1), resp.Id)
    assert.Equal(t, "test", resp.Name)

    // 验证 mock 调用
    mockClient.AssertExpectations(t)
}
```

#### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./app/borrower/api/internal/logic/...

# 运行测试并生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 运行性能测试
go test -bench=. ./...
```

### 2. 集成测试

#### 数据库测试

```go
// 测试数据库连接和基本操作
func TestNewFeatureRepo_Integration(t *testing.T) {
    // 设置测试数据库
    db := setupTestDB(t)
    defer cleanupTestDB(db)

    // 创建 repository
    repo := NewNewFeatureRepo(db)

    // 测试创建
    feature := &model.NewFeature{
        Name: "test",
        Age:  25,
    }

    result, err := repo.Create(context.Background(), feature)
    assert.NoError(t, err)
    assert.NotZero(t, result.Id)

    // 测试查询
    found, err := repo.GetByID(context.Background(), result.Id)
    assert.NoError(t, err)
    assert.Equal(t, "test", found.Name)
    assert.Equal(t, int64(25), found.Age)
}
```

---

## 🔍 代码质量检查

### 1. 静态代码分析

#### golangci-lint 配置

```yaml
# .golangci.yml
run:
  timeout: 5m
  tests: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0.8
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2
  misspell:
    locale: US
  lll:
    line-length: 140
  goimports:
    local-prefixes: zeebee-backend

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - exhaustive
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - maligned
    - misspell
    - nakedret
    - rowserrcheck
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

issues:
  exclude-rules:
    # Exclude some linters from running on tests files
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
```

#### 运行代码检查

```bash
# 运行所有检查
golangci-lint run

# 运行特定包
golangci-lint run ./app/borrower/api/internal/logic/

# 生成配置文件
golangci-lint config path
```

### 2. 代码格式化

#### 自动格式化

```bash
# 格式化所有代码
go fmt ./...
goimports -w .

# 检查格式问题
go fmt -d ./...
```

#### Git hooks 配置

```bash
# 安装 pre-commit hooks
go install github.com/pre-commit/pre-commit@latest

# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: gofmt
        name: gofmt
        entry: gofmt
        language: system
        args: [-w]
        files: \.go$

      - id: goimports
        name: goimports
        entry: goimports
        language: system
        args: [-w]
        files: \.go$

      - id: golangci-lint
        name: golangci-lint
        entry: golangci-lint run
        language: system
        files: \.go$
```

---

## 🚀 构建和部署

### 1. 本地开发

#### 启动服务

```bash
# 启动 Borrower API
cd app/borrower
go run cmd/main.go -f etc/borrower-api.yaml

# 启动 Borrower RPC
cd app/borrower/rpc
go run cmd/main.go -f etc/borrower.yaml

# 启动 Lender API
cd app/lender
go run cmd/main.go -f etc/lender-api.yaml
```

#### 环境配置

```yaml
# app/borrower/etc/borrower-api.yaml
Name: borrower-api
Host: 0.0.0.0
Port: 10000

# 数据库配置
DBConf:
  Host: localhost
  Port: 5432
  User: dev_user
  Password: dev_password
  DBName: dev_db
  MaxOpenConns: 100
  MaxIdleConns: 10

# Redis 配置
RedisConf:
  Host: localhost
  Port: 6379
  Type: node
  Pass: ""

# JWT 配置
Auth:
  AccessSecret: your-access-secret
  AccessExpire: 7200
  RefreshSecret: your-refresh-secret
  RefreshExpire: 604800
```

### 2. Docker 构建

#### Dockerfile

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o borrower-api ./app/borrower/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/borrower-api .
COPY --from=builder /app/app/borrower/etc ./etc

EXPOSE 10000
CMD ["./borrower-api", "-f", "etc/borrower-api.yaml"]
```

#### docker-compose.yml

```yaml
version: '3.8'

services:
  borrower-api:
    build: .
    ports:
      - "10000:10000"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: dev_db
      POSTGRES_USER: dev_user
      POSTGRES_PASSWORD: dev_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

---

## 📚 文档编写规范

### 1. 代码注释

#### 函数注释

```go
// CreateUser creates a new user in the system
// It validates the input data, hashes the password, and saves the user to database
// Returns the created user with generated ID and timestamps
func CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
    // Implementation here
}
```

#### 包注释

```go
// Package logic provides business logic implementations for the borrower API.
// It handles the core business operations including user management,
// loan processing, and data validation.
package logic
```

### 2. API 文档

#### Swagger 注释

```go
// CreateUser
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags user
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation request"
// @Success 200 {object} CreateUserResponse "User created successfully"
// @Failure 400 {object} error.Response "Invalid input parameters"
// @Failure 500 {object} error.Response "Internal server error"
// @Router /api/users [post]
func CreateUserHandler(ctx *svc.ServiceContext) http.HandlerFunc {
    // Implementation here
}
```

---

## 🔧 常用工具和脚本

### 1. Makefile

```makefile
# 代码生成
.PHONY: gen-api
gen-api:
	goctl api go -api app/borrower/api/desc/ -dir app/borrower/api

.PHONY: gen-rpc
gen-rpc:
	goctl rpc protoc app/borrower/rpc/pb -dir app/borrower/rpc

# 测试
.PHONY: test
test:
	go test ./...

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# 代码质量
.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./...
	goimports -w .

# 构建
.PHONY: build
build:
	go build -o bin/borrower-api ./app/borrower/cmd/main.go
	go build -o bin/borrower-rpc ./app/borrower/rpc/cmd/main.go

# Docker
.PHONY: docker-build
docker-build:
	docker build -t zeebee/borrower-api .

.PHONY: docker-run
docker-run:
	docker-compose up -d
```

### 2. 实用脚本

#### 数据库迁移脚本

```bash
#!/bin/bash
# scripts/migrate.sh
set -e

echo "Running database migrations..."

# 应用迁移
migrate -database "postgres://dev_user:dev_password@localhost:5432/dev_db?sslmode=disable" -path migrations up

echo "Migration completed successfully!"
```

#### 环境启动脚本

```bash
#!/bin/bash
# scripts/start-dev.sh
set -e

echo "Starting development environment..."

# 启动依赖服务
docker-compose up -d postgres redis

# 等待服务启动
sleep 5

# 运行数据库迁移
./scripts/migrate.sh

# 启动应用
echo "Starting application services..."
go run app/borrower/cmd/main.go -f app/borrower/etc/borrower-api.yaml &
go run app/borrower/rpc/cmd/main.go -f app/borrower/rpc/etc/borrower.yaml &

echo "Development environment started successfully!"
```

---

## 📋 检查清单

### 开发前检查

- [ ] Go 环境已安装并配置
- [ ] go-zero 工具链已安装
- [ ] 数据库和 Redis 已启动
- [ ] 代码仓库已克隆
- [ ] 依赖包已下载

### 代码提交前检查

- [ ] 代码已格式化 (`go fmt`)
- [ ] 代码已通过 lint 检查 (`golangci-lint`)
- [ ] 单元测试已编写并通过
- [ ] API 文档已更新
- [ ] 提交信息符合规范

### 部署前检查

- [ ] 所有测试已通过
- [ ] 性能测试已完成
- [ ] 安全扫描已通过
- [ ] 配置文件已更新
- [ ] 回滚方案已准备

---

## 🔍 故障排除

### 常见问题

#### 1. goctl 命令找不到

```bash
# 重新安装
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 检查 PATH
echo $PATH | grep go
```

#### 2. 数据库连接失败

```bash
# 检查数据库服务
docker-compose ps

# 测试连接
psql "host=localhost user=dev_user password=dev_password dbname=dev_db port=5432"
```

#### 3. 端口冲突

```bash
# 查看端口占用
netstat -tlnp | grep :10000

# 修改配置文件中的端口
vim app/borrower/etc/borrower-api.yaml
```

#### 4. 依赖包问题

```bash
# 清理并重新下载
go clean -modcache
go mod download
go mod tidy
```

---

## 📞 支持和帮助

### 文档资源

- [go-zero 官方文档](https://go-zero.dev/)
- [Go 官方文档](https://golang.org/doc/)
- [项目知识库](../guides/navigation-guide.md)

### 团队沟通

- 技术问题: 项目 Issue
- 开发讨论: 团队会议
- 紧急问题: 项目负责人

---

**文档维护**: 开发团队 + AI 助手
**最后更新**: 2025-12-22
**下次更新**: 工具变更时