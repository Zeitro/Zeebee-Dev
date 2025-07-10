# Zeebee Backend

基于 go-zero 框架的微服务架构贷款管理系统后端。

## 项目结构

# 数据库搭建 To do

```
backend/
├── common/                 # 共享组件
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

