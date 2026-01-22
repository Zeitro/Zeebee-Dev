# Quickstart Guide: Borrower Zeitro API Proxy

**Feature**: 004-borrower-zeitro-proxy
**Date**: 2025-01-09
**Status**: ✅ Implementation Complete

## 概述

本指南将帮助开发者快速开始使用 Borrower Zeitro API Proxy 功能。该功能提供了从 Zeitro API 搜索房产列表的代理服务，包括自动 token 管理、数据转换和错误处理。

**架构方案**: 方案C（简化架构）- 单文件客户端实现
**测试覆盖率**: 81.9%
**状态**: 生产就绪

## 前置条件

- Go 1.23+
- 有效的 Zeitro API 凭据

## 配置说明

### 1. Zeitro API 凭据配置

在 `app/borrower/etc/config-{env}.yaml` 中配置 Zeitro API：

```yaml
RpcConf:
  ZeitroConf:
    Url: https://lo.zeitro.com
    Username: cli_hu4tKiltzCfQfheOa3
    Password: sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd
```

**注意**: 配置已在以下环境完成：
- `config.yaml` - 开发环境
- `config-local.yaml` - 本地环境
- `config-dev.yaml` - 开发环境
- `config-staging.yaml` - 预发布环境
- `config-production.yaml` - 生产环境

## 使用示例

### API 调用

```bash
curl -X POST http://localhost:10000/borrower/v1/zeitro/listings/search \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "lat": 37.7749,
    "lon": -122.4194,
    "dist": 10.0,
    "addr": "San Francisco, CA"
  }'
```

### 响应格式

```json
{
  "data": [
    {
      "full_address": "123 Main St",
      "state": "CA",
      "county": "San Francisco",
      "city": "San Francisco",
      "zip_code": "94102",
      "property_type": "Single Family",
      "price": "1250000",
      "income_indicator": "High",
      "minority_percent": "25",
      "total_population": "50000",
      "hispanic": "15",
      "black": "5",
      "asian": "30",
      "pacific": "2",
      "msa_code": "41860",
      "days_on_market": "10",
      "agent_name": "John Doe",
      "agent_phone": "(415) 555-1234",
      "agent_email": "john@example.com",
      "distance": "5.2"
    }
  ]
}
```

## 代码结构

### 简化架构（方案C）

```
domains/zeebee/zeitro/
├── client.go           # 核心客户端实现（331行）
│   ├── Config         # 配置结构
│   ├── Token          # Token 管理
│   ├── Client         # HTTP 客户端
│   └── SearchListings # 房产搜索
└── client_test.go     # 测试套件（540行，81.9%覆盖）

app/borrower/
├── api/
│   ├── desc/zeitro.api              # API 定义
│   └── internal/logic/zeitro/
│       └── zeitro_listing_search_logic.go  # API 层逻辑
└── rpc/
    ├── internal/logic/zeitro/
    │   └── zeitro_listing_search_logic.go  # RPC 层逻辑
    └── internal/svc/
        └── service_context.go       # 服务上下文（集成 ZeitroClient）
```

## 核心功能

### 1. 自动 Token 管理

```go
// Token 会自动获取并缓存
// 在过期前 10% 时间自动刷新
// 使用 singleflight 防止并发重复获取

client, _ := zeitro.NewClient(config, httpClient)
listings, err := client.SearchListings(ctx, searchRequest)
```

**Token 管理特性**:
- ✅ 首次请求自动获取
- ✅ 10% 阈值过期检测
- ✅ singleflight 并发控制
- ✅ 失败时回退到旧 token

### 2. HTTP 客户端注入

```go
// 在 ServiceContext 中初始化
httpClient := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

zeitroClient, err := zeitro.NewClient(zeitro.Config{
    URL:      c.ZeitroConf.Url,
    Username: c.ZeitroConf.Username,
    Password: c.ZeitroConf.Password,
}, httpClient)
```

### 3. 数据转换

使用 `copierx.Trans` 进行类型安全的自动转换：

```go
// RPC 层：领域模型 → Protobuf
result, err := copierx.Trans[[]zeitro.PropertyListing, []*borrower.ZeitroPropertyEntity](&listings)

// API 层：Protobuf → API 类型
resp, err := copierx.Trans[borrower.ZeitroListingSearchResponse, types.ZeitroListingSearchResponse](rpcResp)
```

## 测试

### 运行测试

```bash
# 单元测试
go test ./domains/zeebee/zeitro/... -v

# 查看覆盖率
go test ./domains/zeebee/zeitro/... -cover

# 编译验证
go build -o /dev/null ./domains/zeebee/zeitro/... ./app/borrower/api/... ./app/borrower/rpc/...
```

**测试结果**: 11个测试函数，20+子测试，覆盖率 81.9%

### 测试覆盖范围

- ✅ 配置验证（5个测试场景）
- ✅ Token 过期检测（4个测试场景）
- ✅ Token 有效性检查（4个测试场景）
- ✅ Token 获取（使用 httptest mock）
- ✅ Token 缓存验证
- ✅ 房产搜索（完整集成测试）
- ✅ 错误处理
- ✅ 电话号码格式化

## API 端点

### POST /borrower/v1/zeitro/listings/search

搜索指定位置附近的房产列表。

**请求参数**:

| 字段 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| lat | float64 | ✅ | 纬度（-90 到 90） | 37.7749 |
| lon | float64 | ✅ | 经度（-180 到 180） | -122.4194 |
| dist | float64 | ✅ | 搜索半径（0-100英里） | 10.0 |
| addr | string | ❌ | 地址字符串 | "San Francisco, CA" |

**响应字段**:

| 字段 | 类型 | 说明 |
|------|------|------|
| data | []ZeitroPropertyEntity | 房产列表数组 |

**ZeitroPropertyEntity 字段**:

- `full_address`: 完整地址
- `state`, `county`, `city`, `zip_code`: 地理信息
- `property_type`: 房产类型
- `price`: 价格
- `income_indicator`, `minority_percent`, `total_population`: 人口统计
- `hispanic`, `black`, `asian`, `pacific`: 种族分布
- `days_on_market`: 上市天数
- `agent_name`, `agent_phone`, `agent_email`: 经纪人信息
- `distance`: 距离搜索中心的距离

## 错误处理

### 错误响应格式

```json
{
  "code": 500,
  "msg": "Failed to search listings: ...",
  "details": "..."
}
```

### 常见错误

1. **配置错误**: Zeitro URL 或凭据未配置
2. **Token 获取失败**: API 凭据无效或网络问题
3. **搜索超时**: Zeitro API 响应时间超过 10 秒
4. **无效参数**: 纬度/经度超出范围

## 性能特性

- **连接池**: 100 个最大空闲连接，每主机 10 个
- **超时控制**: 10 秒 HTTP 超时
- **Token 缓存**: 自动缓存，减少重复请求
- **并发控制**: singleflight 防止并发 token 刷新

## 调试

### 启用详细日志

在配置文件中设置：

```yaml
Log:
  Level: debug
  Mode: console
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
  -d '{"Lat":37.7749,"Lon":-122.4194,"Dist":10}'
```

## 常见问题

### Q: Token 刷新失败怎么办？

A: 系统会保留旧 token 并记录错误日志。检查：
1. Zeitro API 凭据是否正确
2. 网络连接是否正常
3. Zeitro API 服务是否可用

### Q: 如何调整超时时间？

A: 修改 `service_context.go` 中的 HTTP 客户端配置：

```go
httpClient := &http.Client{
    Timeout: 15 * time.Second,  // 调整为 15 秒
}
```

### Q: 如何添加新的搜索参数？

A: 更新以下文件：
1. `client.go` 中的 `SearchRequest` 结构
2. `zeitro.api` 中的请求定义
3. 重新生成代码：`make api-borrower`

### Q: 支持哪些房产类型？

A: 取决于 Zeitro API 返回的类型，包括但不限于：
- Townhouse
- Single Family
- Condo
- Multi-Family

## 相关文档

- [功能规格](spec.md) - 完整需求说明
- [实施计划](plan.md) - 技术设计和架构决策
- [任务清单](tasks.md) - 开发任务和进度跟踪
- [数据模型](data-model.md) - 数据结构设计
- [研究文档](research.md) - 技术调研结果

## 项目索引

已添加到项目文档索引：`.claude/docs/INDEX.md` → **外部集成** 章节

## 支持

如有问题，请查阅相关文档或联系开发团队。

