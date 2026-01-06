# 项目结构指南

**文档创建时间**: 2025-11-24
**文档类型**: 开发指导文档
**维护状态**: 活跃

## 概述

本文档详细说明 Zeebee 贷款管理系统的项目结构、目录组织和代码组织方式。系统采用基于 go-zero
框架的微服务架构，遵循领域驱动设计 (DDD) 原则，具有清晰的分层结构和模块化设计。

## 整体项目结构

```
zeebee-backend/
├── 📋 README.md                           # 项目说明文档
├── 📜 CLAUDE.md                          # AI 开发指导文档
├── 📦 go.mod                            # Go 模块依赖管理
├── 📦 go.sum                            # Go 依赖版本锁定
├── 🛠️ Makefile                          # 代码生成和管理命令
├── 📏 .golangci.yml                      # Go 代码检查配置
├── 📄 .gitignore                        # Git 忽略文件配置
├── 🗃️ app/                              # 应用服务层
├── 🔧 common/                           # 通用组件库
├── 🗂️ domains/                          # 领域模型层
├── 🚀 deploy/                           # 部署配置文件
├── 📚 claudedocs/                       # AI 生成文档库
├── 🛠️ third_party/                      # 第三方资源和模板
└── 📜 scripts/                          # 脚本工具文件
```

## 应用服务层 (app/)

### 服务组织结构

```
app/
├── 🏢 borrower/                         # 借款方微服务
│   ├── 🔌 api/                          # HTTP API 服务
│   │   ├── 📄 borrower.api             # 主 API 定义文件
│   │   ├── 📂 desc/                     # 消息体定义目录
│   │   │   ├── 📄 common.api           # 通用消息体
│   │   │   ├── 📄 user.api             # 用户相关消息体
│   │   │   ├── 📄 loan.api             # 贷款相关消息体
│   │   │   ├── 📄 property.api         # 房产相关消息体
│   │   │   ├── 📄 borrower_entity.api  # 借款实体消息体
│   │   │   ├── 📄 borrower_applicant.api # 借款申请人消息体
│   │   │   ├── 📄 loan_applicant.api    # 贷款申请人消息体
│   │   │   ├── 📄 loan_borrower.api     # 贷款借款人消息体
│   │   │   ├── 📄 loan_property.api     # 贷款房产消息体
│   │   │   ├── 📄 loan_demographic.api  # 贷款人口统计消息体
│   │   │   ├── 📄 loan_document.api     # 贷款文档消息体
│   │   │   ├── 📄 loan_experience.api   # 贷款经验消息体
│   │   │   ├── 📄 loan_third_service.api # 第三方服务消息体
│   │   │   ├── 📄 file.api             # 文件相关消息体
│   │   │   ├── 📄 storage.api          # 存储相关消息体
│   │   │   ├── 📄 us_geography.api      # 美国地理消息体
│   │   │   └── 📄 pdf.api               # PDF 生成消息体
│   │   ├── 🏭 internal/                 # 内部实现
│   │   │   ├── 🎯 handler/             # HTTP 请求处理器
│   │   │   │   ├── 👤 user/             # 用户相关处理器
│   │   │   │   │   ├── 📄 info_handler.go
│   │   │   │   │   ├── 📄 login_handler.go
│   │   │   │   │   ├── 📄 register_handler.go
│   │   │   │   │   └── 📄 reset_password_handler.go
│   │   │   │   ├── 💰 loan/             # 贷款相关处理器
│   │   │   │   │   ├── 📄 create_handler.go
│   │   │   │   │   ├── 📄 detail_handler.go
│   │   │   │   │   ├── 📄 page_handler.go
│   │   │   │   │   └── 📄 update_handler.go
│   │   │   │   ├── 🏠 property/         # 房产相关处理器
│   │   │   │   │   ├── 📄 property_create_handler.go
│   │   │   │   │   ├── 📄 property_delete_handler.go
│   │   │   │   │   ├── 📄 property_detail_handler.go
│   │   │   │   │   ├── 📄 property_page_handler.go
│   │   │   │   │   └── 📄 property_update_handler.go
│   │   │   │   └── 🔧 common/           # 通用功能处理器
│   │   │   │       ├── 📄 get_us_cities_handler.go
│   │   │   │       ├── 📄 get_us_states_handler.go
│   │   │   │       └── 📄 select_options_handler.go
│   │   │   ├── 🧠 logic/               # API 层业务逻辑
│   │   │   │   ├── 👤 user/             # 用户业务逻辑
│   │   │   │   ├── 💰 loan/             # 贷款业务逻辑
│   │   │   │   ├── 🏠 property/         # 房产业务逻辑
│   │   │   │   └── 🔧 common/           # 通用业务逻辑
│   │   │   ├── 🛡️ middleware/         # 中间件
│   │   │   │   ├── 📄 auth_middleware.go  # JWT 认证中间件
│   │   │   │   └── 📄 cors_middleware.go  # CORS 跨域中间件
│   │   │   ├── 📋 validate_rules/      # 数据验证规则
│   │   │   │   └── 📄 property.go       # 房产验证规则
│   │   │   ├── 📝 svc/                 # 服务上下文
│   │   │   │   └── 📄 service_context.go
│   │   │   ├── 📄 config.go            # 配置结构定义
│   │   │   ├── 📄 borrower.go          # 服务注册
│   │   │   └── 📄 swagger/              # Swagger API 文档
│   │   ├── ⚙️ etc/                    # 环境配置文件
│   │   │   ├── 📄 borrower-api.yaml    # API 服务配置
│   │   │   ├── 📄 config-dev.yaml      # 开发环境配置
│   │   │   ├── 📄 config-local.yaml    # 本地环境配置
│   │   │   ├── 📄 config-production.yaml # 生产环境配置
│   │   │   └── 📄 config-staging.yaml  # 预发布环境配置
│   │   └── 🚀 cmd/                    # 服务启动入口
│   │       └── 📄 main.go              # 主程序文件
│   └── 🔌 rpc/                          # RPC 业务服务
│       ├── 📄 borrower.proto          # gRPC 服务定义
│       ├── 🏭 client/                  # RPC 客户端封装
│       │   ├── 📄 applicant/           # 申请人服务客户端
│       │   ├── 📄 entity/              # 实体服务客户端
│       │   ├── 📄 property/            # 房产服务客户端
│       │   └── 📄 user/                # 用户服务客户端
│       ├── 🏭 internal/               # RPC 内部实现
│       │   ├── 🧠 logic/               # RPC 业务逻辑
│       │   │   ├── 👤 user/             # 用户 RPC 逻辑
│       │   │   ├── 💰 loan/             # 贷款 RPC 逻辑
│       │   │   ├── 🏠 property/         # 房产 RPC 逻辑
│       │   │   └── 🔧 common/           # 通用 RPC 逻辑
│       │   ├── 🖥️ server/              # gRPC 服务端
│       │   ├── 📝 svc/                 # RPC 服务上下文
│       │   ├── 📄 config.go            # RPC 配置定义
│       │   ├── 📄 borrower.go          # RPC 服务注册
│       │   ├── 🌐 global/              # 全局组件
│       │   │   ├── 📄 config.go        # 全局配置
│       │   │   ├── 📄 email.go         # 邮件服务
│       │   │   └── 📄 global.go         # 全局初始化
│       │   └── 📧 email_tpl/            # 邮件模板
│       │       ├── 📄 send_login_url_email.html
│       │       ├── 📄 send_register_verify_email.html
│       │       ├── 📄 send_reset_pwd_email.html
│       │       ├── 📄 send_self_preapproval_notification.html
│       │       └── 📄 send_third_party_preapproval_notification.html
│       ├── 📂 pb/                      # Protocol Buffers 生成文件
│       │   ├── 📄 borrower.pb.go
│       │   ├── 📄 borrower_grpc.pb.go
│       │   └── 📄 borrower.proto
│       ├── ⚙️ etc/                    # RPC 环境配置
│       │   ├── 📄 borrower.yaml
│       │   ├── 📄 config-dev.yaml
│       │   ├── 📄 config-local.yaml
│       │   ├── 📄 config-production.yaml
│       │   └── 📄 config-staging.yaml
│       └── 🚀 cmd/                    # RPC 服务启动入口
│           └── 📄 main.go
└── 🏢 lender/                          # 贷款方微服务 (结构类似 borrower)
    ├── 🔌 api/                          # HTTP API 服务
    └── 🔌 rpc/                          # RPC 业务服务
```

## 通用组件库 (common/)

### 组件结构详解

```
common/
├── 📦 constants/                        # 业务常量定义
│   ├── 📄 constants.go                 # 系统基础常量
│   ├── 📄 applicant.go                 # 申请人相关常量
│   ├── 📄 borrower.go                  # 借款方相关常量
│   ├── 📄 file.go                      # 文件相关常量
│   ├── 📄 lender_company.go            # 贷款公司相关常量
│   ├── 📄 loan.go                      # 贷款相关常量
│   ├── 📄 loan_document.go             # 贷款文档相关常量
│   ├── 📄 loan_operation.go            # 贷款操作相关常量
│   ├── 📄 property.go                  # 房产相关常量
│   ├── 📄 redis_key.go                 # Redis 键值常量
│   └── 📄 user.go                      # 用户相关常量
├── 🚨 errorx/                           # 自定义错误处理
│   ├── 📄 error.go                      # 自定义错误类型定义
│   ├── 📄 error_group.go                # 错误分组处理
│   ├── 📄 error_test.go                 # 错误处理测试
│   └── 📄 copier_error.go               # 数据复制错误
├── 🗃️ gormx/                            # 数据库操作封装
│   ├── 📄 base_repo.go                  # 基础仓储模式 (泛型实现)
│   ├── 📄 logger.go                     # 数据库日志配置
│   ├── 📄 option.go                    # 查询选项定义
│   └── 📄 page.go                       # 分页组件实现
├── 📤 response/                         # HTTP 响应标准化
│   ├── 📄 http_handler.go              # HTTP 响应处理器
│   └── 📄 response.go                  # 标准响应格式定义
├── 🛠️ utils/                            # 通用工具函数
│   ├── 📄 jwt.go                        # JWT 令牌工具
│   ├── 📄 password.go                  # 密码处理工具
│   ├── 📄 password_test.go             # 密码工具测试
│   ├── 📄 clone.go                      # 对象深拷贝工具
│   └── 📄 stringx.go                    # 字符串处理工具
├── ☁️ aws_s3/                            # AWS S3 文件存储
│   ├── 📄 presigner.go                  # S3 预签名 URL 生成
│   └── 📄 presigner_test.go             # S3 工具测试
├── 📧 email/                            # 邮件服务组件
│   ├── 📄 email.go                      # 邮件发送服务
│   └── 📄 email_test.go                 # 邮件服务测试
├── 🎯 preapproval_engine/                # 贷款预审批引擎
│   ├── 📄 pre_approval_engine.go       # 预审批核心逻辑
│   ├── 📄 pre_approval_engine_test.go  # 预审批引擎测试
│   └── 📄 pre_approval_rule_test.go     # 预审批规则测试
├── 🪪 uuid/                             # UUID 生成服务
│   └── 📄 sonyflake.go                  # Sonyflake ID 生成器
├── ✅ validator/                         # 数据验证工具
│   └── 📄 validator_test.go             # 验证工具测试
├── 🔄 limiter/                           # 限流器组件
│   └── 📄 redis_limiter.go              # Redis 基础限流实现
├── 📏 enums/                            # 枚举类型定义
│   └── 📄 enum.go                       # 业务枚举类型
└── ⏰ timex/                            # 时间处理工具
    └── 📄 datetime.go                   # 日期时间工具
```

### 核心组件说明

#### 1. 错误处理系统 (errorx/)

```go
// 自定义错误类型
type CodeError interface {
    error
    GetCode() int
    GetMsg() string
    GetData() interface{}
}

// 错误分组处理
type ErrorGroup struct {
    errors []error
}
```

#### 2. 数据库抽象层 (gormx/)

```go
// 泛型基础仓储模式
type BaseRepo[T schema.Tabler] struct {
    db *gorm.DB
}

// 统一的 CRUD 操作
func (repo *BaseRepo[T]) Create(ctx context.Context, data T) (T, error)
func (repo *BaseRepo[T]) FindOneById(ctx context.Context, id int64) (T, error)
func (repo *BaseRepo[T]) Page(query *gorm.DB, pageReq PageRequest) (PageResponse[T], error)
```

#### 3. 响应标准化 (response/)

```go
// 标准响应格式
type Response struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data"`
}
```

## 领域模型层 (domains/)

### 领域驱动设计结构

```
domains/
└── 🏢 zeebee/                           # Zeebee 领域模型
    ├── 🗂️ model/                       # 数据模型定义
    │   ├── 📄 base.go                   # 基础模型和通用结构
    │   ├── 👤 user.go                   # 用户模型
    │   ├── 🏢 borrower.go               # 借款方模型
    │   ├── 👥 borrower_applicant.go     # 借款申请人模型
    │   ├── 🏗️ borrower_entity.go        # 借款实体模型
    │   ├── 🏠 borrower_property.go       # 借款房产模型
    │   ├── 💰 loan.go                   # 贷款模型
    │   ├── 📋 loan_applicant.go          # 贷款申请人模型
    │   ├── 🏦 loan_borrower.go           # 贷款借款关系模型
    │   ├── 📊 loan_demographic.go        # 贷款人口统计模型
    │   ├── 📄 loan_document.go           # 贷款文档模型
    │   ├── 🏗️ loan_experience.go         # 贷款经验模型
    │   ├── 📍 loan_operation.go          # 贷款操作模型
    │   ├── 🏠 loan_property.go           # 贷款房产模型
    │   ├── 🤝 loan_third_service.go       # 第三方服务模型
    │   ├── 🏦 lender_company.go          # 贷款公司模型
    │   ├── 👨‍💼 lender_employee.go         # 贷款员工模型
    │   ├── 📁 file.go                   # 文件模型
    │   └── 🌍 us_geography.go            # 美国地理数据模型
    └── 🗂️ repo/                        # 数据访问层
        ├── 📄 borrower.go               # 借款方数据访问
        ├── 📄 lender_company.go         # 贷款公司数据访问
        ├── 📄 loan_applicant.go          # 贷款申请人数据访问
        ├── 📄 loan_borrower.go           # 贷款借款关系数据访问
        ├── 📄 loan_demographic.go        # 人口统计数据访问
        ├── 📄 loan_documents.go          # 贷款文档数据访问
        ├── 📄 loan_experience.go          # 贷款经验数据访问
        ├── 📄 loan_operation.go          # 贷款操作数据访问
        ├── 📄 loan_property.go           # 贷款房产数据访问
        ├── 📄 loan_third_service.go       # 第三方服务数据访问
        ├── 📄 user.go                    # 用户数据访问
        └── 📄 us_geography.go            # 美国地理数据访问
```

### 数据模型设计原则

#### 1. 基础模型结构

```go
// 所有模型的基础结构
type Base struct {
    Id        int64                 `gorm:"primaryKey;autoIncrement;<-:create" json:"id"`
    CreatedAt time.Time             `gorm:"column:created_at;default:now();<-:create" json:"created_at"`
    UpdatedAt time.Time             `gorm:"column:updated_at;default:now()" json:"updated_at"`
    IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-"`
}
```

#### 2. 领域关系映射

```sql
-- 核心业务关系
borrower_entities (1:N) → borrower_applicants
borrower_entities (1:N) → borrower_properties
borrower_entities (1:N) → loans
loans (N:M) → loan_applicants
loans (1:N) → loan_documents
loans (1:N) → loan_properties
```

## 部署配置 (deploy/)

### 容器化部署结构

```
deploy/
├── 🏢 borrower/                         # 借款方服务部署配置
│   ├── 🐳 Dockerfile                   # Docker 镜像构建文件
│   └── 🐳 docker-compose.yml           # 容器编排配置
│       📄 borrower-api/
│       📄 borrower-rpc/
│       📄 postgresql/
│       📄 redis/
└── 🏢 lender/                           # 贷款方服务部署配置
    ├── 🐳 Dockerfile                   # Docker 镜像构建文件
    └── 🐳 docker-compose.yml           # 容器编排配置
        📄 lender-api/
        📄 lender-rpc/
        📄 postgresql/
        📄 redis/
```

## 开发工具和模板 (third_party/)

### go-zero 代码生成模板

```
third_party/
├── 🏗️ goctl_template/                   # go-zero 代码生成模板
│   └── 📦 1.8.4/                       # 模板版本
│       ├── 📝 api/                      # API 代码生成模板
│       │   ├── 📄 config.tpl            # 配置文件模板
│       │   ├── 📄 etc.tpl              # 配置目录模板
│       │   ├── 📄 handler.tpl           # 处理器模板
│       │   ├── 📄 logic.tpl             # 业务逻辑模板
│       │   ├── 📄 main.tpl              # 主程序模板
│       │   ├── 📄 middleware.tpl        # 中间件模板
│       │   ├── 📄 response.tpl          # 响应模板
│       │   ├── 📄 routes.tpl            # 路由模板
│       │   └── 📄 types.tpl             # 类型定义模板
│       ├── 🔄 rpc/                      # RPC 代码生成模板
│       │   ├── 📄 call.tpl              # RPC 调用模板
│       │   ├── 📄 config.tpl            # RPC 配置模板
│       │   ├── 📄 etc.tpl              # RPC 配置目录模板
│       │   ├── 📄 logic.tpl             # RPC 逻辑模板
│       │   ├── 📄 main.tpl              # RPC 主程序模板
│       │   ├── 📄 server-func.tpl       # 服务端函数模板
│       │   ├── 📄 server.tpl            # 服务端模板
│       │   ├── 📄 svc.tpl              # 服务上下文模板
│       │   └── 📄 template.tpl          # RPC 模板
│       ├── 🐳 docker/                    # Docker 容器模板
│       │   └── 📄 docker.tpl            # Dockerfile 模板
│       ├── ☸️ kube/                      # Kubernetes 部署模板
│       │   ├── 📄 deployment.tpl        # 部署模板
│       │   ├── 📄 job.tpl              # 任务模板
│       │   └── 📄 service.tpl           # 服务模板
│       └── 🗄️ model/                     # 数据模型生成模板
│           ├── 📄 delete.tpl            # 删除方法模板
│           ├── 📄 err.tpl               # 错误处理模板
│           ├── 📄 find-one-by-field-extra-method.tpl
│           ├── 📄 find-one-by-field.tpl
│           ├── 📄 find-one.tpl          # 查找单个记录模板
│           ├── 📄 import-no-cache.tpl    # 导入模板
│           ├── 📄 import.tpl            # 导入模板
│           ├── 📄 insert.tpl            # 插入模板
│           ├── 📄 model-gen.tpl         # 模型生成模板
│           ├── 📄 model-new.tpl          # 新模型模板
│           ├── 📄 model.tpl             # 模型模板
│           ├── 📄 table-name.tpl        # 表名模板
│           ├── 📄 types.tpl             # 类型模板
│           ├── 📄 update.tpl            # 更新模板
│           ├── 📄 var.tpl               # 变量模板
│           └── 📄 view.tpl              # 视图模板
└── 🌐 google/                            # Google Protocol Buffers 定义
    ├── 📑 api/                          # Google API 定义
    │   ├── 📄 annotations.proto       # API 注解
    │   ├── 📄 client.proto            # 客户端定义
    │   ├── 📄 field_behavior.proto   # 字段行为定义
    │   ├── 📄 http.proto              # HTTP 映射
    │   └── 📄 httpbody.proto           # HTTP 消息体
    └── 📑 protobuf/                      # Protocol Buffers 定义
        ├── 📄 any.proto                # 任意类型
        ├── 📄 api.proto                # API 定义
        ├── 📄 compiler/                # 编译器相关
        ├── 📄 descriptor.proto         # 描述符定义
        ├── 📄 duration.proto           # 时间间隔
        ├── 📄 empty.proto              # 空消息
        ├── 📄 field_mask.proto         # 字段掩码
        ├── 📄 source_context.proto     # 源码上下文
        ├── 📄 struct.proto             # 结构体
        ├── 📄 timestamp.proto          # 时间戳
        ├── 📄 type.proto                # 类型定义
        └── 📄 wrappers.proto            # 包装类型
```

## 开发工作流程

### 1. 服务创建流程

```bash
# 创建新的 API 定义
make create-api service-name

# 生成 API 服务
make create-api-service service-name

# 创建 RPC 定义
make create-proto service-name

# 生成 RPC 服务
make create-rpc-service service-name

# 生成 Swagger 文档
make create-swagger service-name
```

### 2. 代码生成流程

```
API 定义文件 (.api)
        ↓
goctl 代码生成工具
        ↓
┌─────────────────────────────────────┐
│              生成的代码                │
├─────────────────────────────────────┤
│ 📄 handler/     # HTTP 请求处理器      │
│ 📄 logic/       # API 层业务逻辑       │
│ 📄 types.go     # 类型定义            │
│ 📄 routes.go    # 路由注册            │
│ 📄 main.go      # 服务启动入口        │
└─────────────────────────────────────┘
```

### 3. 数据流向架构

```
HTTP Request
    ↓
API Handler (接收HTTP请求)
    ↓
API Logic (数据校验、参数整合)
    ↓
RPC Call (gRPC Protocol)
    ↓
RPC Logic (核心业务处理)
    ↓
Repository (数据访问层)
    ↓
Model (数据模型)
    ↓
Database (PostgreSQL)
```

## 配置管理

### 环境配置结构

```
配置层次:
├── 🏠 本地开发 (Local)
│   └── 📄 config-local.yaml
├── 🔧 开发环境 (Development)
│   └── 📄 config-dev.yaml
├── 🧪 预发布环境 (Staging)
│   └── 📄 config-staging.yaml
└── 🏭 生产环境 (Production)
    └── 📄 config-production.yaml
```

### 配置文件内容结构

```yaml
# 服务配置示例
Name: borrower-api
Host: 0.0.0.0
Port: 10000
Mode: dev

# 数据库配置
DbConf:
  Type: postgres
  Host: localhost
  Port: 5432
  Name: zeebee_dev
  Username: dev_user
  Password: dev_password
  MaxOpenConns: 10
  MaxIdleConns: 5
  MaxLifetime: 300

# Redis 配置
RedisConf:
  Host: localhost
  Port: 6379
  Type: node
  Pass: ""

# JWT 配置
Auth:
  AccessSecret: "your-access-secret"
  AccessExpire: 86400
  RefreshExpire: 604800

# AWS S3 配置
AWS:
  Region: us-west-2
  Bucket: zeebee-files
  AccessKey: ""
  SecretKey: ""
```

## 代码质量工具

### golangci-lint 配置

```yaml
# .golangci.yml 配置要点
version: 2
run:
  timeout: 5m
  tests: true
  fix: true

linters:
  enable:
    - govet          # 静态分析
    - errcheck       # 错误检查
    - staticcheck    # 静态检查
    - revive         # 代码风格
    - cyclop         # 复杂度检查
    - gocritic       # 代码质量
    - gosec          # 安全检查
    - ineffassign    # 无效赋值检查
    - misspell       # 拼写检查
    - prealloc       # 预分配检查
    - unconvert      # 类型转换检查
    - unparam        # 未使用参数检查
    - whitespace     # 空白字符检查
```

## 测试结构

### 测试文件组织

```
测试分布:
├── 🔬 单元测试 (10个文件)
│   ├── common/aws_s3/presigner_test.go
│   ├── common/email/email_test.go
│   ├── common/errorx/error_test.go
│   ├── common/preapproval_engine/pre_approval_engine_test.go
│   ├── common/utils/password_test.go
│   ├── common/uuid/sonyflake_test.go
│   ├── common/validator/validator_test.go
│   └── app/borrower/api/internal/pdf/apply_form_test.go
├── 🧪 集成测试 (待补充)
├── 🌐 E2E 测试 (待补充)
└── 🚀 性能测试 (待补充)
```

### 测试覆盖率目标

- **当前覆盖率**: < 20%
- **目标覆盖率**: 80%+
- **优先覆盖**: 业务逻辑层、数据访问层、工具组件

## 文档管理

### AI 生成文档结构

```
claudedocs/
├── 📊 analysis/                      # 分析类文档
│   ├── 📄 architecture_analysis_20250124.md
│   └── 📂 project-research/
├── 📋 planning/                      # 规划类文档
│   ├── 📂 task-lists/
│   ├── 📂 progress-tracker/
│   └── 📂 roadmap/
├── 📚 guides/                         # 开发指导文档
│   ├── 📄 api_documentation_index.md
│   ├── 📄 navigation-guide.md
│   ├── 📄 project_structure_guide.md
│   ├── 📂 best-practices/
│   └── 📂 troubleshooting/
├── 📈 reports/                        # 各类分析报告
│   ├── 📂 code-quality/
│   ├── 📂 security-analysis/
│   └── 📂 performance/
├── 🛠️ tools/                          # 开发工具文档
│   ├── 📂 workflows/
│   └── 📂 automation/
└── 📄 INDEX.md                       # 项目文档索引
```

## 最佳实践

### 1. 代码组织原则

- **单一职责**: 每个包和模块职责明确
- **依赖倒置**: 高层模块不依赖低层模块
- **接口隔离**: 定义细粒度接口
- **开闭原则**: 对扩展开放，对修改封闭

### 2. 命名规范

- **文件命名**: 使用 snake_case (示例: `user_service.go`)
- **目录命名**: 使用 snake_case (示例: `user_service/`)
- **函数命名**: 使用驼峰命名 (示例: `CreateUser`)
- **常量命名**: 使用大写下划线 (示例: `MAX_RETRY_COUNT`)

### 3. 错误处理规范

```go
// 统一错误处理模式
defer func() {
    if err != nil {
        logx.WithContext(r.Context()).
            WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
            Error(err)
    }
}()

if err != nil {
    return nil, errorx.NewDefaultError("操作失败")
}
```

### 4. 数据库操作规范

```go
// 使用基础仓储模式
type UserRepo struct {
    *gormx.BaseRepo[*model.User]
}

func (repo *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := repo.Client(ctx).
        Where("email = ?", email).
        First(&user).Error
    return &user, errorx.WithStack(err)
}
```

## 扩展指南

### 1. 添加新服务

1. 在 `app/` 下创建服务目录
2. 定义 API 和 Proto 文件
3. 使用 goctl 生成代码
4. 实现业务逻辑
5. 添加测试用例
6. 更新部署配置

### 2. 添加新的 API 端点

1. 在 `desc/` 目录定义消息体
2. 在主 API 文件中定义路由
3. 重新生成代码
4. 实现 Handler 和 Logic
5. 添加测试

### 3. 数据库模型变更

1. 修改或创建新的 Model 文件
2. 创建数据库迁移脚本
3. 更新 Repository 文件
4. 添加相关业务逻辑
5. 更新 API 定义

## 相关文档

- **API 文档索引**: `claudedocs/guides/api_documentation_index.md`
- **导航指南**: `claudedocs/guides/navigation-guide.md`
- **最佳实践**: `claudedocs/guides/best-practices/`
- **故障排除**: `claudedocs/guides/troubleshooting/`
- **架构分析**: `claudedocs/analysis/architecture_analysis_20250124.md`

---

*本文档最后更新时间: 2025-01-24*