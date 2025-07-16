#!/bin/bash

# 帮助函数，显示可用命令
help() {
    echo "Usage: ./makefile.sh <target>"
    echo "Targets:"
    echo "  create-api <service>        Create API file"
    echo "  create-swagger <service>    Create API Swagger"
    echo "  create-api-service <service> Create API Service"
    echo "  create-proto <service>      Create proto file"
    echo "  create-rpc-service <service> Create RPC service"
    echo "  run-api-service <service>   Run API service"
    echo "  run-rpc-service <service>   Run RPC service"
}

# 创建 API 文件
create-api() {
    SERVICE="$1"
    echo "create api file: $SERVICE"
    if [ ! -d "app/$SERVICE/api" ]; then
        mkdir -p app/$SERVICE/api
    fi
    goctl api -o=app/$SERVICE/api/$SERVICE.api --home third_party/goctl_template
}

# 创建 API Swagger
create-swagger() {
    SERVICE="$1"
    echo "create api swagger: $SERVICE"
    goctl api swagger --api app/"$SERVICE"/api/"$SERVICE".api --dir app/"$SERVICE"/api/swagger
}

# 创建 API 服务
create-api-service() {
    SERVICE="$1"
    echo "create api server: $SERVICE"
    goctl api go --api app/"$SERVICE"/api/"$SERVICE".api --dir app/"$SERVICE"/api --style go_zero --home third_party/goctl_template --type-group
}

# 创建 proto 文件
create-proto() {
    SERVICE="$1"
    goctl rpc -o app/"$SERVICE"/rpc/pb/"$SERVICE".proto
}

# 创建 RPC 服务
create-rpc-service() {
    SERVICE="$1"
    goctl rpc protoc app/"$SERVICE"/rpc/pb/"$SERVICE".proto --go_out=./app/"$SERVICE"/rpc/pb --go-grpc_out=./app/"$SERVICE"/rpc/pb --zrpc_out=./app/"$SERVICE"/rpc -m --home third_party/goctl_template --style go_zero
}

# 运行 API 服务
run-api-service() {
    SERVICE="$1"
    go run app/"$SERVICE"/api/"$SERVICE".go -f app/"$SERVICE"/api/etc/"$SERVICE"-api.yaml
}

# 运行 RPC 服务
run-rpc-service() {
    SERVICE="$1"
    go run app/"$SERVICE"/rpc/"$SERVICE".go -f app/"$SERVICE"/rpc/etc/"$SERVICE".yaml
}

if [ $# -eq 0 ]; then
    help
else
    case "$1" in
        help)
            help
            ;;
        create-api)
            create-api "$2"
            ;;
        create-swagger)
            create-swagger "$2"
            ;;
        create-api-service)
            create-api-service "$2"
            ;;
        create-proto)
            create-proto "$2"
            ;;
        create-rpc-service)
            create-rpc-service "$2"
            ;;
        run-api-service)
            run-api-service "$2"
            ;;
        run-rpc-service)
            run-rpc-service "$2"
            ;;
        *)
            echo "Unknown target: $1"
            help
            ;;
    esac
fi
