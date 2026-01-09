# Quickstart Guide: Borrower Zeitro API Proxy

**Feature**: 004-borrower-zeitro-proxy
**Date**: 2025-01-09
**Phase**: Phase 1 - Design & Contracts

## 概述

本指南将帮助开发者快速开始使用 Borrower Zeitro API Proxy 功能。该功能提供了从 Zeitro API 搜索房产列表的代理服务，包括自动
token 管理、数据转换和错误处理。

## 前置条件

- Go 1.23+
- PostgreSQL 14+
- Redis 7+
- 有效的 Zeitro API 凭据

## 安装步骤

### 1. 配置 Zeitro API 凭据

在 `app/borrower/etc/config-{env}.yaml` 中添加 Zeitro 配置：

```yaml
Zeitro:
  Url: https://lo.zeitro.com
  Username: cli_hu4tKiltzCfQfheOa3
  Password: sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd
  TokenExpireTime: 3600  # 秒
  Timeout: 10           # 秒
  MaxRetries: 3
```

### 2. 生成 API 代码

```bash
# 生成 borrower API 代码
cd app/borrower/api
goctl api go -api desc/zeitro.api -dir . --style=go_zero
```

### 3. 实现 RPC Logic

在 `app/borrower/rpc/internal/logic/` 中创建 Zeitro 相关逻辑：

```go
// zeitro_listing_search_logic.go
package logic

import (
    "context"
    "yourusername/app/borrower/rpc/internal/svc"
    "yourusername/app/borrower/rpc/pb"
    "yourusername/domains/zeebee/zeitro"
)

type ZeitroListingSearchLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewZeitroListingSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ZeitroListingSearchLogic {
    return &ZeitroListingSearchLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *ZeitroListingSearchLogic) ZeitroListingSearch(in *pb.ZeitroListingSearchReq) (*pb.ZeitroListingSearchResp, error) {
    // 调用领域服务
    client := l.svcCtx.ZeitroClient
    listings, err := client.SearchListings(l.ctx, &model.SearchRequest{
        Latitude:  in.Latitude,
        Longitude: in.Longitude,
        Distance:  in.Distance,
        Address:   in.Address,
    })
    if err != nil {
        return nil, err
    }

    // 转换为 protobuf 格式
    var result []*pb.ZeitroPropertyEntity
    for _, listing := range listings {
        result = append(result, convertToProto(listing))
    }

    return &pb.ZeitroListingSearchResp{
        List: result,
    }, nil
}
```

### 4. 注册 RPC 服务

在 ServiceContext 中添加 Zeitro 客户端：

```go
// app/borrower/rpc/internal/svc/servicecontext.go
type ServiceContext struct {
    Config      config.Config
    ZeitroClient *zeitro.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config:      c,
        ZeitroClient: zeitro.NewClient(c.Zeitro),
    }
}
```

## 使用示例

### 基本搜索

```bash
curl -X POST http://localhost:10000/borrower/v1/zeitro/listings/search \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lat": 37.33874,
    "lon": -121.8852525,
    "dist": 5,
    "addr": "San Jose, CA, USA"
  }'
```

### 响应示例

```json
{
  "list": [
    {
      "full_address": "523 Union Ave, Campbell, CA 95008",
      "state": "CA",
      "county": "Santa Clara",
      "city": "Campbell",
      "zip_code": "95008",
      "property_type": "Townhouse",
      "price": "1699000",
      "income_indicator": "4",
      "minority_percent": "93.91",
      "total_population": "3662",
      "hispanic": "3302",
      "black": "9",
      "asian": "75",
      "pacific": "10",
      "msa_code": "31084",
      "days_on_market": "31",
      "agent_name": "Jina Lee",
      "agent_phone": "3109855462",
      "agent_email": "jina@metisre.com",
      "distance": "4"
    }
  ]
}
```

## 代码结构

```
domains/zeebee/zeitro/
├── model/                   # 数据模型
│   ├── config.go           # 配置结构
│   ├── token.go            # Token 结构
│   ├── property.go         # 房产列表结构
│   └── request.go          # 请求结构
├── client/                  # HTTP 客户端
│   └── client.go           # Zeitro API 客户端
├── token/                   # Token 管理
│   └── manager.go          # Token 管理器
├── converter/               # 数据转换
│   └── converter.go        # Zeitro → Internal 转换
└── errors/                  # 错误定义
    └── errors.go           # 错误类型

app/borrower/
├── api/
│   ├── desc/
│   │   └── zeitro.api      # API 定义
│   └── internal/
│       └── handler/
│           └── zeitro_listing_search_handler.go
└── rpc/
    ├── internal/
    │   └── logic/
    │       └── zeitro_listing_search_logic.go
    └── pb/
        └── borrower.proto  # Protobuf 定义
```

## 开发工作流

### 1. 添加新的 API 端点

1. 在 `app/borrower/api/desc/zeitro.api` 中定义
2. 运行 `goctl api go` 生成代码
3. 实现 Handler 和 Logic
4. 在 RPC 层实现业务逻辑

### 2. 添加新的数据字段

1. 更新 `domains/zeebee/zeitro/model/` 中的结构
2. 更新转换器 `converter/converter.go`
3. 更新 API 定义 `zeitro.api`
4. 重新生成代码

### 3. 测试

```bash
# 单元测试
go test ./domains/zeebee/zeitro/...

# 集成测试
go test ./app/borrower/api/...

# 运行服务
cd app/borrower/api
go run borrower.go -f etc/borrower-api.yaml

# 运行 RPC 服务
cd app/borrower/rpc
go run borrower.go -f etc/borrower.yaml
```

## 调试

### 启用详细日志

在配置文件中设置：

```yaml
Log:
  Level: debug
  Mode: console
```

### 查看 Token 状态

```go
import "yourusername/domains/zeebee/zeitro/token"

manager := token.NewManager(config)
token, err := manager.GetToken(ctx)
logx.Infof("Token expires at: %s", token.ExpiresAt)
```

### 测试 Zeitro API 连接

```bash
# 测试 Token 获取
curl -X POST https://lo.zeitro.com/auth/access-token \
  -u "cli_hu4tKiltzCfQfheOa3:sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials"

# 测试 Listing 搜索
curl -X POST https://lo.zeitro.com/los/openhousesearch \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"Lat":37.33874,"Lon":-121.8852525,"Dist":5}'
```

## 常见问题

### Q: Token 刷新失败怎么办？

A: 系统会保留旧 token 并记录错误日志。检查：

1. Zeitro API 凭据是否正确
2. 网络连接是否正常
3. Zeitro API 服务是否可用

### Q: 如何处理 Zeitro API 超时？

A: 系统已配置 10 秒超时。如果需要调整：

```yaml
Zeitro:
  Timeout: 15  # 增加到 15 秒
```

### Q: 如何缓存搜索结果？

A: 当前不实现本地缓存。前端或 API Gateway 可以添加缓存层。

### Q: 支持哪些房产类型？

A: 取决于 Zeitro API 返回的类型，包括但不限于：

- Townhouse
- Single Family
- Condo
- Multi-Family

## 性能优化建议

1. **并发控制**: 使用 `singleflight` 防止重复 token 刷新
2. **连接池**: HTTP 客户端已配置连接池参数
3. **超时设置**: 合理设置超时避免阻塞
4. **日志级别**: 生产环境使用 `info` 级别

## 监控指标

- Token 缓存命中率
- API 响应时间
- Token 刷新频率
- 错误率

## 相关文档

- [API 契约](contracts/openapi.yaml)
- [数据模型](data-model.md)
- [研究文档](research.md)
- [实施计划](plan.md)

## 支持

如有问题，请联系开发团队或提交 Issue。
