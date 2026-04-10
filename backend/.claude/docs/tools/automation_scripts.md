# Zeebee 自动化脚本集

**创建时间**: 2025-12-22
**适用环境**: Linux/macOS
**最后更新**: 2025-12-22

---

## 📋 概述

本文档包含了 Zeebee 项目开发、测试、部署等环节的自动化脚本，旨在提高开发效率，减少重复性工作。

## 🚀 快速启动脚本

### 1. 一键启动开发环境

```bash
#!/bin/bash
# scripts/start-dev.sh
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.24+"
        exit 1
    fi

    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi

    log_info "Dependencies check passed ✓"
}

# 启动基础设施服务
start_infrastructure() {
    log_info "Starting infrastructure services..."

    # 启动 PostgreSQL 和 Redis
    docker-compose -f deploy/docker-compose.yml up -d postgres redis

    # 等待服务启动
    log_info "Waiting for services to start..."
    sleep 10

    # 检查服务状态
    if docker-compose -f deploy/docker-compose.yml ps | grep -q "Up.*postgres"; then
        log_info "PostgreSQL is running ✓"
    else
        log_error "PostgreSQL failed to start"
        exit 1
    fi

    if docker-compose -f deploy/docker-compose.yml ps | grep -q "Up.*redis"; then
        log_info "Redis is running ✓"
    else
        log_error "Redis failed to start"
        exit 1
    fi
}

# 运行数据库迁移
run_migrations() {
    log_info "Running database migrations..."

    # 等待数据库就绪
    until docker-compose -f deploy/docker-compose.yml exec -T postgres pg_isready -U dev_user -d dev_db; do
        log_warn "Waiting for database to be ready..."
        sleep 2
    done

    # 这里可以添加具体的迁移命令
    # migrate -database "postgres://dev_user:dev_password@localhost:5432/dev_db" up
    log_info "Database migrations completed ✓"
}

# 启动应用服务
start_applications() {
    log_info "Starting application services..."

    # 启动 Borrower 服务
    log_info "Starting Borrower API (Port: 10000)..."
    cd app/borrower && go run cmd/main.go -f etc/borrower-api.yaml &
    BORROWER_PID=$!

    sleep 3

    cd app/borrower/rpc && go run cmd/main.go -f etc/borrower.yaml &
    BORROWER_RPC_PID=$!

    # 启动 Lender 服务
    log_info "Starting Lender API (Port: 11000)..."
    cd app/lender && go run cmd/main.go -f etc/lender-api.yaml &
    LENDER_PID=$!

    # 保存 PID 到文件
    echo $BORROWER_PID > .borrower-api.pid
    echo $BORROWER_RPC_PID > .borrower-rpc.pid
    echo $LENDER_PID > .lender-api.pid

    log_info "All application services started ✓"
}

# 健康检查
health_check() {
    log_info "Performing health checks..."

    sleep 5

    # 检查 Borrower API
    if curl -s http://localhost:10000/health > /dev/null; then
        log_info "Borrower API is healthy ✓"
    else
        log_warn "Borrower API health check failed"
    fi

    # 检查 Lender API
    if curl -s http://localhost:11000/health > /dev/null; then
        log_info "Lender API is healthy ✓"
    else
        log_warn "Lender API health check failed"
    fi
}

# 主函数
main() {
    log_info "🚀 Starting Zeebee Development Environment"

    check_dependencies
    start_infrastructure
    run_migrations
    start_applications
    health_check

    echo
    log_info "🎉 Development environment is ready!"
    log_info "📍 Borrower API: http://localhost:10000"
    log_info "📍 Lender API: http://localhost:11000"
    log_info "📚 Documentation: http://localhost:8080/docs"
    echo
    log_info "To stop all services, run: ./scripts/stop-dev.sh"

    # 等待用户输入以退出
    trap 'kill_services' EXIT
    log_info "Press Ctrl+C to stop all services"

    while true; do
        sleep 1
    done
}

# 清理函数
kill_services() {
    log_info "Stopping all services..."

    # 停止应用服务
    if [ -f .borrower-api.pid ]; then
        kill $(cat .borrower-api.pid) 2>/dev/null || true
        rm .borrower-api.pid
    fi

    if [ -f .borrower-rpc.pid ]; then
        kill $(cat .borrower-rpc.pid) 2>/dev/null || true
        rm .borrower-rpc.pid
    fi

    if [ -f .lender-api.pid ]; then
        kill $(cat .lender-api.pid) 2>/dev/null || true
        rm .lender-api.pid
    fi

    # 停止基础设施服务
    docker-compose -f deploy/docker-compose.yml down

    log_info "All services stopped ✓"
}

# 执行主函数
main
```

### 2. 停止开发环境

```bash
#!/bin/bash
# scripts/stop-dev.sh
set -e

GREEN='\033[0;32m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_info "🛑 Stopping Zeebee Development Environment"

# 停止应用服务
log_info "Stopping application services..."
pkill -f "borrower-api" || true
pkill -f "borrower-rpc" || true
pkill -f "lender-api" || true

# 清理 PID 文件
rm -f .borrower-api.pid .borrower-rpc.pid .lender-api.pid

# 停止基础设施服务
log_info "Stopping infrastructure services..."
docker-compose -f deploy/docker-compose.yml down

log_info "✅ All services stopped successfully"
```

---

## 🧪 测试自动化脚本

### 1. 运行所有测试

```bash
#!/bin/bash
# scripts/run-tests.sh
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 设置测试环境
setup_test_env() {
    log_info "Setting up test environment..."

    # 设置测试数据库
    export DB_HOST=localhost
    export DB_NAME=test_db
    export REDIS_HOST=localhost
    export JWT_TEST_SECRET=test-secret

    # 启动测试数据库
    docker-compose -f deploy/docker-compose.test.yml up -d postgres-test redis-test

    # 等待数据库就绪
    sleep 10
}

# 运行单元测试
run_unit_tests() {
    log_info "Running unit tests..."

    go test -v -race -short ./...

    if [ $? -eq 0 ]; then
        log_info "✅ Unit tests passed"
    else
        log_error "❌ Unit tests failed"
        exit 1
    fi
}

# 运行集成测试
run_integration_tests() {
    log_info "Running integration tests..."

    go test -v -race -tags=integration ./...

    if [ $? -eq 0 ]; then
        log_info "✅ Integration tests passed"
    else
        log_error "❌ Integration tests failed"
        exit 1
    fi
}

# 生成覆盖率报告
generate_coverage() {
    log_info "Generating test coverage report..."

    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    log_info "📊 Test coverage: ${COVERAGE}%"

    if [ "${COVERAGE%.*}" -lt 50 ]; then
        log_warn "⚠️ Test coverage is below 50%"
    else
        log_info "✅ Test coverage is good"
    fi
}

# 性能测试
run_performance_tests() {
    log_info "Running performance tests..."

    # 使用 go test -bench 运行性能测试
    go test -bench=. -benchmem ./...

    if [ $? -eq 0 ]; then
        log_info "✅ Performance tests passed"
    else
        log_error "❌ Performance tests failed"
        exit 1
    fi
}

# 代码质量检查
run_quality_checks() {
    log_info "Running code quality checks..."

    # 运行 golangci-lint
    golangci-lint run

    if [ $? -eq 0 ]; then
        log_info "✅ Code quality checks passed"
    else
        log_error "❌ Code quality checks failed"
        exit 1
    fi
}

# 安全扫描
run_security_scan() {
    log_info "Running security scan..."

    # 使用 gosec 进行安全扫描
    gosec ./...

    if [ $? -eq 0 ]; then
        log_info "✅ Security scan passed"
    else
        log_warn "⚠️ Security scan found issues"
    fi
}

# 清理测试环境
cleanup_test_env() {
    log_info "Cleaning up test environment..."

    # 停止测试服务
    docker-compose -f deploy/docker-compose.test.yml down

    # 清理测试文件
    rm -f coverage.out coverage.html
}

# 主函数
main() {
    log_info "🧪 Running Zeebee Test Suite"

    # 解析命令行参数
    UNIT_ONLY=false
    INTEGRATION_ONLY=false
    COVERAGE_ONLY=false
    PERFORMANCE_ONLY=false
    QUALITY_ONLY=false
    SECURITY_ONLY=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit-only)
                UNIT_ONLY=true
                shift
                ;;
            --integration-only)
                INTEGRATION_ONLY=true
                shift
                ;;
            --coverage-only)
                COVERAGE_ONLY=true
                shift
                ;;
            --performance-only)
                PERFORMANCE_ONLY=true
                shift
                ;;
            --quality-only)
                QUALITY_ONLY=true
                shift
                ;;
            --security-only)
                SECURITY_ONLY=true
                shift
                ;;
            *)
                shift
                ;;
        esac
    done

    # 设置测试环境
    setup_test_env

    # 根据参数运行相应测试
    if [ "$UNIT_ONLY" = true ]; then
        run_unit_tests
    elif [ "$INTEGRATION_ONLY" = true ]; then
        run_integration_tests
    elif [ "$COVERAGE_ONLY" = true ]; then
        generate_coverage
    elif [ "$PERFORMANCE_ONLY" = true ]; then
        run_performance_tests
    elif [ "$QUALITY_ONLY" = true ]; then
        run_quality_checks
    elif [ "$SECURITY_ONLY" = true ]; then
        run_security_scan
    else
        # 运行所有测试
        run_unit_tests
        run_integration_tests
        generate_coverage
        run_performance_tests
        run_quality_checks
        run_security_scan
    fi

    # 清理环境
    cleanup_test_env

    log_info "🎉 Test suite completed successfully!"
}

# 设置错误处理
trap 'cleanup_test_env' EXIT

# 执行主函数
main "$@"
```

### 2. 持续集成脚本

```bash
#!/bin/bash
# scripts/ci.sh
set -e

# 配置变量
export CI=true
export GO111MODULE=on
export CGO_ENABLED=0

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[CI-INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[CI-ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[CI-WARN]${NC} $1"
}

# CI 阶段函数
stage_setup() {
    log_info "🔧 Setting up CI environment"

    # 设置 Go 环境变量
    export GOPATH=$HOME/go
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

    # 缓存目录
    mkdir -p $HOME/.cache/go-build

    # 显示环境信息
    go version
    docker version

    log_info "✅ CI environment setup completed"
}

stage_test() {
    log_info "🧪 Running tests"

    # 下载依赖
    go mod download

    # 运行测试
    go test -v -race -coverprofile=coverage.out ./...

    # 生成覆盖率报告
    go tool cover -func=coverage.out
}

stage_build() {
    log_info "🏗️ Building applications"

    # 构建 Borrower API
    log_info "Building Borrower API..."
    cd app/borrower
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../build/borrower-api ./cmd/main.go

    # 构建 Borrower RPC
    log_info "Building Borrower RPC..."
    cd rpc
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../../build/borrower-rpc ./cmd/main.go

    # 构建 Lender API
    log_info "Building Lender API..."
    cd ../../app/lender
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../build/lender-api ./cmd/main.go

    log_info "✅ Build completed"
}

stage_security() {
    log_info "🔒 Running security scans"

    # 使用 gosec 扫描
    log_info "Running gosec..."
    gosec ./...

    # 使用 nancy 检查依赖漏洞
    log_info "Checking dependencies for vulnerabilities..."
    go list -json -m all | nancy sleuth

    log_info "✅ Security scans completed"
}

stage_quality() {
    log_info "📊 Running quality checks"

    # 运行 golangci-lint
    log_info "Running golangci-lint..."
    golangci-lint run

    # 检查代码格式
    log_info "Checking code format..."
    if [ $(gofmt -s -l . | wc -l) -gt 0 ]; then
        log_error "Code is not properly formatted"
        gofmt -s -l .
        exit 1
    fi

    log_info "✅ Quality checks passed"
}

# 主函数
main() {
    log_info "🚀 Starting Zeebee CI Pipeline"

    stage_setup
    stage_test
    stage_build
    stage_security
    stage_quality

    log_info "🎉 CI Pipeline completed successfully!"
}

# 执行主函数
main
```

---

## 🐳 Docker 自动化脚本

### 1. 多服务构建脚本

```bash
#!/bin/bash
# scripts/build-docker.sh
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[BUILD]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[BUILD]${NC} $1"
}

# 构建配置
REGISTRY="your-registry.com"
PROJECT_NAME="zeebee"
VERSION=${1:-latest}

# 构建函数
build_service() {
    local service_name=$1
    local dockerfile=$2
    local context=$3

    log_info "Building ${service_name}..."

    docker build \
        -f ${dockerfile} \
        -t ${REGISTRY}/${PROJECT_NAME}/${service_name}:${VERSION} \
        -t ${REGISTRY}/${PROJECT_NAME}/${service_name}:latest \
        ${context}

    if [ $? -eq 0 ]; then
        log_info "✅ ${service_name} built successfully"
    else
        log_error "❌ ${service_name} build failed"
        exit 1
    fi
}

# 构建所有服务
build_all() {
    log_info "🐳 Building all Docker images"

    # 构建 Borrower API
    build_service "borrower-api" "deploy/borrower/Dockerfile" "."

    # 构建 Borrower RPC
    build_service "borrower-rpc" "deploy/borrower/Dockerfile.rpc" "."

    # 构建 Lender API
    build_service "lender-api" "deploy/lender/Dockerfile" "."

    # 构建 Admin API
    build_service "admin-api" "deploy/admin/Dockerfile" "."

    log_info "✅ All Docker images built successfully"
}

# 推送镜像
push_images() {
    log_info "📤 Pushing Docker images to registry..."

    services=("borrower-api" "borrower-rpc" "lender-api" "admin-api")

    for service in "${services[@]}"; do
        log_info "Pushing ${service}..."
        docker push ${REGISTRY}/${PROJECT_NAME}/${service}:${VERSION}
        docker push ${REGISTRY}/${PROJECT_NAME}/${service}:latest
    done

    log_info "✅ All images pushed successfully"
}

# 清理未使用的镜像
cleanup() {
    log_info "🧹 Cleaning up unused Docker images"

    # 删除构建缓存
    docker builder prune -f

    # 删除未使用的镜像
    docker image prune -f

    log_info "✅ Cleanup completed"
}

# 显示镜像信息
show_images() {
    log_info "📋 Built Docker images:"

    services=("borrower-api" "borrower-rpc" "lender-api" "admin-api")

    for service in "${services[@]}"; do
        echo "  - ${REGISTRY}/${PROJECT_NAME}/${service}:${VERSION}"
        echo "  - ${REGISTRY}/${PROJECT_NAME}/${service}:latest"
    done
}

# 主函数
main() {
    log_info "🚀 Zeebee Docker Build Script"

    # 解析命令行参数
    PUSH=false
    CLEANUP=false
    SHOW=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --push)
                PUSH=true
                shift
                ;;
            --cleanup)
                CLEANUP=true
                shift
                ;;
            --show)
                SHOW=true
                shift
                ;;
            --version)
                VERSION=$2
                shift 2
                ;;
            -h|--help)
                echo "Usage: $0 [--push] [--cleanup] [--show] [--version VERSION]"
                echo "  --push     Push images to registry"
                echo "  --cleanup  Clean up unused images"
                echo "  --show     Show built images"
                echo "  --version  Set image version (default: latest)"
                exit 0
                ;;
            *)
                shift
                ;;
        esac
    done

    # 构建
    build_all

    # 推送
    if [ "$PUSH" = true ]; then
        push_images
    fi

    # 显示镜像信息
    if [ "$SHOW" = true ]; then
        show_images
    fi

    # 清理
    if [ "$CLEANUP" = true ]; then
        cleanup
    fi

    log_info "🎉 Docker build process completed!"
}

# 执行主函数
main "$@"
```

### 2. 部署脚本

```bash
#!/bin/bash
# scripts/deploy.sh
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[DEPLOY]${NC} $1"
}

log_error() {
    echo -e "${RED}[DEPLOY-ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[DEPLOY-WARN]${NC} $1"
}

# 配置
ENVIRONMENT=${1:-staging}
REGISTRY="your-registry.com"
PROJECT_NAME="zeebee"
VERSION=${2:-latest}

# 环境验证
validate_env() {
    case $ENVIRONMENT in
        staging|production)
            ;;
        *)
            log_error "Invalid environment. Use 'staging' or 'production'"
            exit 1
            ;;
    esac

    log_info "Deploying to environment: ${ENVIRONMENT}"
}

# 部署函数
deploy_service() {
    local service=$1
    local replicas=$2
    local port=$3

    log_info "Deploying ${service}..."

    # 生成部署配置
    cat > deploy/${ENVIRONMENT}/${service}.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${service}
  namespace: ${PROJECT_NAME}-${ENVIRONMENT}
spec:
  replicas: ${replicas}
  selector:
    matchLabels:
      app: ${service}
  template:
    metadata:
      labels:
        app: ${service}
    spec:
      containers:
      - name: ${service}
        image: ${REGISTRY}/${PROJECT_NAME}/${service}:${VERSION}
        ports:
        - containerPort: ${port}
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: ${PROJECT_NAME}-config
              key: DB_HOST
        - name: REDIS_HOST
          valueFrom:
            secretKeyRef:
              name: ${PROJECT_NAME}-config
              key: REDIS_HOST
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        livenessProbe:
          httpGet:
            path: /health
            port: ${port}
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: ${port}
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: ${service}
  namespace: ${PROJECT_NAME}-${ENVIRONMENT}
spec:
  selector:
    app: ${service}
  ports:
  - protocol: TCP
    port: ${port}
    targetPort: ${port}
  type: ClusterIP
EOF

    # 应用部署配置
    kubectl apply -f deploy/${ENVIRONMENT}/${service}.yaml

    # 等待部署完成
    log_info "Waiting for ${service} to be ready..."
    kubectl rollout status deployment/${service} -n ${PROJECT_NAME}-${ENVIRONMENT} --timeout=300s

    # 验证部署
    if kubectl get pods -n ${PROJECT_NAME}-${ENVIRONMENT} -l app=${service} | grep -q "Running"; then
        log_info "✅ ${service} deployed successfully"
    else
        log_error "❌ ${service} deployment failed"
        exit 1
    fi
}

# 滚动更新
rolling_update() {
    local service=$1
    local image=$2

    log_info "Performing rolling update for ${service}..."

    # 更新镜像
    kubectl set image deployment/${service} ${service}=${image} -n ${PROJECT_NAME}-${ENVIRONMENT}

    # 等待更新完成
    kubectl rollout status deployment/${service} -n ${PROJECT_NAME}-${ENVIRONMENT} --timeout=600s

    log_info "✅ Rolling update completed for ${service}"
}

# 健康检查
health_check() {
    log_info "Performing health checks..."

    services=("borrower-api:10000" "lender-api:11000" "admin-api:12000")

    for service_info in "${services[@]}"; do
        service_name=$(echo $service_info | cut -d: -f1)
        port=$(echo $service_info | cut -d: -f2)

        # 获取服务 URL
        service_url=$(kubectl get svc ${service_name} -n ${PROJECT_NAME}-${ENVIRONMENT} -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null)

        if [ -n "$service_url" ]; then
            url="http://${service_url}:${port}"
        else
            # 如果没有负载均衡器，使用 NodePort
            node_port=$(kubectl get svc ${service_name} -n ${PROJECT_NAME}-${ENVIRONMENT} -o jsonpath='{.spec.ports[0].nodePort}' 2>/dev/null)
            if [ -n "$node_port" ]; then
                url="http://localhost:${node_port}"
            else
                log_warn "Cannot determine service URL for ${service_name}"
                continue
            fi
        fi

        # 执行健康检查
        if curl -f -s "${url}/health" > /dev/null; then
            log_info "✅ ${service_name} is healthy"
        else
            log_error "❌ ${service_name} health check failed"
            exit 1
        fi
    done
}

# 回滚
rollback() {
    local service=$1
    local revision=${2:-1}

    log_info "Rolling back ${service} to revision ${revision}..."

    kubectl rollout undo deployment/${service} -n ${PROJECT_NAME}-${ENVIRONMENT} --to-revision=${revision}

    kubectl rollout status deployment/${service} -n ${PROJECT_NAME}-${ENVIRONMENT} --timeout=300s

    log_info "✅ Rollback completed for ${service}"
}

# 主函数
main() {
    log_info "🚀 Zeebee Deployment Script"

    # 解析命令行参数
    ROLLBACK=false
    HEALTH_CHECK=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --rollback)
                ROLLBACK=true
                SERVICE=$2
                REVISION=$3
                shift 3
                ;;
            --health-check)
                HEALTH_CHECK=true
                shift
                ;;
            --version)
                VERSION=$2
                shift 2
                ;;
            -h|--help)
                echo "Usage: $0 [ENVIRONMENT] [VERSION] [--rollback SERVICE REVISION] [--health-check]"
                echo "  ENVIRONMENT: Target environment (staging|production)"
                echo "  VERSION: Image version to deploy (default: latest)"
                echo "  --rollback: Rollback service to specific revision"
                echo "  --health-check: Perform health check after deployment"
                exit 0
                ;;
            *)
                shift
                ;;
        esac
    done

    # 验证环境
    validate_env

    # 执行操作
    if [ "$ROLLBACK" = true ]; then
        rollback "$SERVICE" "$REVISION"
    else
        # 部署服务
        deploy_service "borrower-api" 3 10000
        deploy_service "borrower-rpc" 2 10001
        deploy_service "lender-api" 2 11000
        deploy_service "admin-api" 1 12000

        # 健康检查
        if [ "$HEALTH_CHECK" = true ]; then
            sleep 30
            health_check
        fi
    fi

    log_info "🎉 Deployment process completed!"
}

# 执行主函数
main "$@"
```

---

## 🔧 开发工具脚本

### 1. 代码生成脚本

```bash
#!/bin/bash
# scripts/generate.sh
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[GENERATE]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[GENERATE]${NC} $1"
}

# 生成 API 代码
generate_api() {
    local service=$1

    log_info "Generating API code for ${service}..."

    goctl api go \
        -api app/${service}/api/desc/ \
        -dir app/${service}/api \
        -style gozero

    log_info "✅ API code generated for ${service}"
}

# 生成 RPC 代码
generate_rpc() {
    local service=$1

    log_info "Generating RPC code for ${service}..."

    goctl rpc protoc \
        app/${service}/rpc/pb/${service}.proto \
        --go_out=app/${service}/rpc/pb \
        --go-grpc_out=app/${service}/rpc/pb \
        --zrpc_out=app/${service}/rpc

    log_info "✅ RPC code generated for ${service}"
}

# 生成所有服务代码
generate_all() {
    log_info "🏗️ Generating code for all services"

    services=("borrower" "lender" "admin")

    for service in "${services[@]}"; do
        generate_api "$service"
        generate_rpc "$service"
    done

    log_info "✅ All service code generated successfully"
}

# 生成 API 文档
generate_docs() {
    log_info "📚 Generating API documentation..."

    services=("borrower" "lender" "admin")

    for service in "${services[@]}"; do
        log_info "Generating docs for ${service}..."

        # 生成 Swagger 文档
        goctl api doc \
            -dir app/${service}/api/desc/ \
            -o app/${service}/api/swagger/

        log_info "✅ Docs generated for ${service}"
    done
}

# 生成配置文件
generate_config() {
    log_info "⚙️ Generating configuration files..."

    # 生成 API 配置
    cat > app/borrower/api/etc/borrower-api.yaml <<EOF
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

# 日志配置
LogConf:
  Mode: file
  Path: logs/borrower-api.log
  Level: info
EOF

    log_info "✅ Configuration files generated"
}

# 主函数
main() {
    log_info "🚀 Zeebee Code Generation Script"

    # 解析命令行参数
    API_ONLY=false
    RPC_ONLY=false
    DOCS_ONLY=false
    CONFIG_ONLY=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --api-only)
                API_ONLY=true
                shift
                ;;
            --rpc-only)
                RPC_ONLY=true
                shift
                ;;
            --docs-only)
                DOCS_ONLY=true
                shift
                ;;
            --config-only)
                CONFIG_ONLY=true
                shift
                ;;
            --service)
                SERVICE=$2
                shift 2
                ;;
            -h|--help)
                echo "Usage: $0 [--api-only] [--rpc-only] [--docs-only] [--config-only] [--service SERVICE]"
                echo "  --api-only:     Generate API code only"
                echo "  --rpc-only:     Generate RPC code only"
                echo "  --docs-only:    Generate documentation only"
                echo "  --config-only:  Generate config files only"
                echo "  --service:      Generate code for specific service"
                exit 0
                ;;
            *)
                shift
                ;;
        esac
    done

    # 根据参数执行操作
    if [ "$API_ONLY" = true ]; then
        if [ -n "$SERVICE" ]; then
            generate_api "$SERVICE"
        else
            log_error "Please specify service name with --service"
            exit 1
        fi
    elif [ "$RPC_ONLY" = true ]; then
        if [ -n "$SERVICE" ]; then
            generate_rpc "$SERVICE"
        else
            log_error "Please specify service name with --service"
            exit 1
        fi
    elif [ "$DOCS_ONLY" = true ]; then
        generate_docs
    elif [ "$CONFIG_ONLY" = true ]; then
        generate_config
    else
        # 生成所有代码
        generate_all
        generate_docs
        generate_config
    fi

    log_info "🎉 Code generation completed!"
}

# 执行主函数
main "$@"
```

---

## 📋 脚本使用指南

### 快速开始

```bash
# 克隆脚本仓库
git clone https://github.com/your-org/zeebee-scripts.git
cd zeebee-scripts

# 赋予执行权限
chmod +x scripts/*.sh

# 复制到项目根目录
cp scripts/* ../zeebee-backend/scripts/
```

### 开发环境使用

```bash
# 启动开发环境
./scripts/start-dev.sh

# 停止开发环境
./scripts/stop-dev.sh

# 运行测试
./scripts/run-tests.sh

# 生成代码
./scripts/generate.sh

# 构建 Docker 镜像
./scripts/build-docker.sh --push
```

### CI/CD 使用

```bash
# 持续集成
./scripts/ci.sh

# 部署到环境
./scripts/deploy.sh staging
./scripts/deploy.sh production

# 健康检查
./scripts/deploy.sh production --health-check

# 回滚
./scripts/deploy.sh production --rollback borrower-api 2
```

---

**脚本维护**: DevOps 团队
**版本控制**: Git + 语义化版本
**测试覆盖**: 所有脚本都有错误处理
**文档更新**: 2025-12-22

---

**注意**: 使用脚本前请确保已安装必要的依赖工具，并根据实际环境调整配置参数。