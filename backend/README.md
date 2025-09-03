# Zeebee Backend

基于 go-zero 框架的微服务架构贷款管理系统后端。

## 开发环境
```bash
# 安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 查看是否安装成功
goctl version

# protoc一键安装
goctl env check --install --verbose --force

# goctl依赖检查
goctl env check --verbose

# goctl-validate插件安装
go install github.com/linabellbiu/goctl-validate@latest
```
## 项目快速生成
```bash
# 生成api文件
.\Makefile.bat create-api api-name

# 根据api文件生成http服务
.\Makefile.bat create-api-service api-name

# 运行http服务
.\Makefile.bat run-api-service api-name

# 根据api文件生成swagger文件
.\Makefile.bat create-swagger api-name

# 生成proto文件
.\Makefile.bat create-proto proto-name

# 根据proto文件生成rpc服务
.\Makefile.bat create-rpc-service proto-name

# 运行rpc服务
.\Makefile.bat run-rpc-service proto-name

```
## 项目结构

# 数据库搭建 To do

```
backend/
├── app/                        # 服务
│   ├── app_name/               # 服务名
│   │   ├── api/                # http服务
│   │   │   ├── etc/           # 配置文件
│   │   │   ├── internal/      # 内部实现
│   │   │   ├── swagger/       # swagger文件
│   │   │   ├── app_name.api   # http api定义文件
│   │   │   └── app_name.go    # 服务入口文件
│   │   └── rpc/                # rpc服务
│   │       ├── client/        # rpc客户端
│   │       ├── etc/           # 配置文件 
│   │       ├── internal/      # 内部实现
│   │       ├── pb/            # proto文件
│   │       └── app_name.go    # 服务入口文件  
├── common/                # 共享组件
│   ├── constants/         # 常量定义
│   ├── response/          # 统一响应处理
│   └── utils/             # 工具函数
├── webserver/             # 主要Web服务
│   ├── etc/              # 配置文件
│   ├── internal/         # 内部实现
│   │   ├── config/       # 配置结构
│   │   ├── handler/      # HTTP处理器
│   │   ├── logic/        # 业务逻辑
│   │   ├── middleware/   # 中间件
│   │   ├── model/        # 数据模型
│   │   ├── routes/       # 路由配置
│   │   └── svc/          # 服务上下文
│   └── main.go           # 主程序入口
├── streamline/            # 通信服务 (TODO)
│   ├── etc/              # 配置文件
│   └── main.go           # 主程序入口
├── forms/                 # 文件服务 (TODO)
│   ├── etc/              # 配置文件
│   └── main.go           # 主程序入口
├── scripts/               # 脚本文件
├── go.mod                 # Go模块文件
└── README.md             # 项目文档
```

