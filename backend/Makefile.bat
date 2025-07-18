@echo off
setlocal enabledelayedexpansion

:: 检查是否传入了参数
if "%~1"=="" (
    call :help
    exit /b
)

:: 设置服务名（第二个参数）
set "SERVICE=%~2"

:: 判断调用哪个命令
call :%1
if errorlevel 1 (
    echo Unknown command: %1
    call :help
)
exit /b

rem 显示帮助信息
:help
echo Usage: make.bat ^<target^> ^<service^>
echo Targets:
echo   create-api ^<service^>        Create API file
echo   create-swagger ^<service^>    Create API Swagger
echo   create-api-service ^<service^> Create API Service
echo   create-proto ^<service^>      Create proto file
echo   create-rpc-service ^<service^> Create RPC service
echo   run-api-service ^<service^>   Run API service
echo   run-rpc-service ^<service^>   Run RPC service
exit /b

rem 创建 API 文件
:create-api
if "%SERVICE%"=="" (
    echo Error: Service name is required for create-api.
    call :help
    exit /b 1
)
echo create api file: %SERVICE%
if not exist "app\%SERVICE%\api" (
    mkdir app\%SERVICE%\api
)
goctl api -o=app\%SERVICE%\api\%SERVICE%.api --home third_party\goctl_template
exit /b 0

rem 创建 API Swagger
:create-swagger
if "%SERVICE%"=="" (
    echo Error: Service name is required for create-swagger.
    call :help
    exit /b 1
)
echo create api swagger: %SERVICE%
goctl api swagger --api app\%SERVICE%\api\%SERVICE%.api --dir app\%SERVICE%\api\swagger
exit /b 0

rem 创建 API 服务
:create-api-service
if "%SERVICE%"=="" (
    echo Error: Service name is required for create-api-service.
    call :help
    exit /b 1
)
echo create api server: %SERVICE%
goctl api go --api app\%SERVICE%\api\%SERVICE%.api --dir app\%SERVICE%\api --style go_zero --home third_party\goctl_template
goctl api plugin -p goctl-validate="validate --custom" --api app\%SERVICE%\api\%SERVICE%.api --dir app\%SERVICE%\api
exit /b 0

rem 创建 proto 文件
:create-proto
if "%SERVICE%"=="" (
    echo Error: Service name is required for create-proto.
    call :help
    exit /b 1
)
echo create proto file: %SERVICE%
goctl rpc -o app\%SERVICE%\rpc\pb\%SERVICE%.proto
exit /b 0

rem 创建 RPC 服务
:create-rpc-service
if "%SERVICE%"=="" (
    echo Error: Service name is required for create-rpc-service.
    call :help
    exit /b 1
)
echo create rpc service: %SERVICE%
goctl rpc protoc app\%SERVICE%\rpc\pb\%SERVICE%.proto --go_out=./app\%SERVICE%\rpc\pb --go-grpc_out=./app\%SERVICE%\rpc\pb --zrpc_out=./app\%SERVICE%\rpc -m --home third_party\goctl_template --style go_zero
exit /b 0

rem 运行 API 服务
:run-api-service
if "%SERVICE%"=="" (
    echo Error: Service name is required for run-api-service.
    call :help
    exit /b 1
)
echo run api service: %SERVICE%
go run app\%SERVICE%\api\%SERVICE%.go -f app\%SERVICE%\api\etc\%SERVICE%-api.yaml
exit /b 0

rem 运行 RPC 服务
:run-rpc-service
if "%SERVICE%"=="" (
    echo Error: Service name is required for run-rpc-service.
    call :help
    exit /b 1
)
echo run rpc service: %SERVICE%
go run app\%SERVICE%\rpc\%SERVICE%.go -f app\%SERVICE%\rpc\etc\%SERVICE%.yaml
exit /b 0