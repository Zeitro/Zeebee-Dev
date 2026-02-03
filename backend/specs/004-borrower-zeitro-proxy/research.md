# Research: Borrower Zeitro API Proxy

**Feature**: 004-borrower-zeitro-proxy
**Date**: 2025-01-09
**Phase**: Phase 0 - Outline & Research

## 研究目标

本文档记录了 Borrower Zeitro API Proxy 功能的技术研究和设计决策，解决所有技术上下文中的未知问题。

## 技术决策

### 1. 架构模式选择

**决策**: 采用 API/RPC 双层架构 + 领域层封装模式

**理由**:
- 项目使用微服务架构，borrower 服务已有 API/RPC 分层
- 需求明确要求在 `domains/zeebee/zeitro` 中封装 Zeitro 接口
- 符合项目宪章的"共享领域优先"原则
- 便于其他服务（lender、admin）复用 Zeitro 集成

**替代方案**:
- 直接在 borrower RPC 层实现：会导致代码重复，不符合 DRY 原则
- 创建独立的 zeitro 服务：过度设计，当前仅 borrower 端需要

**实施策略**:
```
app/borrower/api           → HTTP 接口层（验证、路由）
    ↓ RPC 调用
app/borrower/rpc          → 业务逻辑层（调用领域服务）
    ↓ 调用
domains/zeebee/zeitro     → 领域服务层（Zeitro API 封装）
    ↓ HTTP 请求
外部 Zeitro API
```

### 2. Token 缓存策略

**决策**: 使用内存缓存 + 单例模式 + 互斥锁

**理由**:
- Token 数据量小，无需持久化
- 需要支持并发访问，使用 `sync.RWMutex` 保护
- 使用 `sync.Once` 确保只初始化一次
- Token 刷新需要原子性操作

**实施细节**:
```go
type TokenManager struct {
    mu          sync.RWMutex
    token       *ZeitroToken
    config      *ZeitroConfig
    httpClient  *http.Client
}

func (tm *TokenManager) GetToken(ctx context.Context) (string, error) {
    // 1. 尝试从缓存读取
    // 2. 检查是否即将过期（< 10% 剩余时间）
    // 3. 如需要，使用互斥锁刷新
}
```

**替代方案**:
- Redis 缓存：增加复杂度，Token 量小无需分布式缓存
- 定时刷新：无法应对 token 意外失效，被动刷新更可靠

### 3. HTTP 客户端选择

**决策**: 使用标准库 `net/http` + `http.Client` 自定义配置

**理由**:
- Go 标准库足够强大，无需额外依赖
- 可配置超时、重试策略
- 与 go-zero 框架集成良好

**配置参数**:
```go
&http.Client{
    Timeout: 10 * time.Second,  // 响应超时
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

**替代方案**:
- `resty`: 功能更多但增加依赖，当前需求标准库足够
- `httpretry`: 可自行实现简单的重试逻辑

### 4. 配置管理

**决策**: 使用 go-zero 的配置结构 + YAML 文件

**理由**:
- 项目已有成熟的配置管理机制
- borrower 服务使用 `borrower-api.yaml` 和 `config-{env}.yaml`
- 支持环境变量和多环境配置

**配置结构**:
```yaml
# app/borrower/api/etc/borrower-api.yaml
Zeitro:
  Url: http://localhost:8080
  Username: cli_hu4tKiltzCfQfheOa3
  Password: sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd
  TokenExpireTime: 3600  # 秒
  Timeout: 10  # 秒
  MaxRetries: 3
```

### 5. 数据转换策略

**决策**: 使用项目标准的 `copierx.Trans[T, U]` + 手动映射

**理由**:
- 项目宪章强制要求使用 `copierx.Trans`
- Zeitro 响应格式与 borrower 数据结构差异较大，需要手动映射
- 分两步：先解析 Zeitro 响应，再转换为目标格式

**实施流程**:
```go
// 1. 定义 Zeitro API 响应结构
type ZeitroListingResponse struct {
    FullAddress     string  `json:"FullAddress"`
    State          string  `json:"State"`
    // ...
}

// 2. 定义 borrower 端数据结构（已有或在 borrower.api 中定义）
type ZeitroPropertyEntity {
    FullAddress     string  `json:"full_address"`
    State          string  `json:"state"`
    // ...
}

// 3. 转换函数
func convertZeitroToEntity(z *ZeitroListingResponse) *ZeitroPropertyEntity {
    // 手动映射字段名（大小写、下划线等差异）
}
```

### 6. 错误处理策略

**决策**: 统一使用 `common/errorx` 包 + 结构化日志

**理由**:
- 项目标准错误处理机制
- 需要区分不同的错误类型：配置错误、网络错误、认证错误、API 错误

**错误分类**:
```go
// 配置错误
var ErrZeitroConfigInvalid = errorx.NewDefaultError("Zeitro configuration is invalid")

// 网络错误
var ErrZeitroConnectionFailed = errorx.NewDefaultError("Failed to connect to Zeitro API")

// 认证错误
var ErrZeitroAuthFailed = errorx.NewDefaultError("Zeitro authentication failed")

// API 错误
var ErrZeitroAPIError = errorx.NewDefaultError("Zeitro API returned error")
```

### 7. API 定义位置

**决策**: 在 `app/borrower/api/desc/` 创建 `zeitro.api`

**理由**:
- borrower 服务已有多个 API 定义文件（property、loan 等）
- 保持代码组织一致性
- 便于 goctl 生成代码

**API 文件结构**:
```go
// app/borrower/api/desc/zeitro.api
syntax = "v1"

info(
    title: "zeitro-listing-api"
    desc: "zeitro listing proxy http service"
    version: "v1"
)

type (
    // 请求结构
    ZeitroListingSearchRequest {
        Latitude   float64  `json:"lat"`
        Longitude  float64  `json:"lon"`
        Distance   float64  `json:"dist"`
        Address    string   `json:"addr"`
    }

    // 响应结构（转换后的格式）
    ZeitroPropertyEntity {
        FullAddress     string  `json:"full_address"`
        State          string  `json:"state"`
        City           string  `json:"city"`
        ZipCode        string  `json:"zip_code"`
        PropertyType   string  `json:"property_type"`
        Price          string  `json:"price"`
        // ...
    }

    ZeitroListingSearchResponse {
        List []ZeitroPropertyEntity `json:"list"`
    }
)

@server(
    prefix: /borrower/v1
    group: zeitro
    middleware: Authority
)
service borrower-api {
    @doc "Search Zeitro property listings"
    @handler zeitroListingSearch
    post /zeitro/listings/search (ZeitroListingSearchRequest) returns (ZeitroListingSearchResponse)
}
```

### 8. 并发控制

**决策**: 使用 `sync/singleflight` 防止重复 token 刷新

**理由**:
- 多个并发请求可能同时检测到 token 过期
- 避免对 Zeitro API 造成压力
- 标准库提供的优雅解决方案

**实施**:
```go
import "golang.org/x/sync/singleflight"

type TokenManager struct {
    sf   singleflight.Group
    // ...
}

func (tm *TokenManager) refreshToken(ctx context.Context) (string, error) {
    result, err, _ := tm.sf.Do("refresh-token", func() (interface{}, error) {
        // 实际刷新逻辑
    })
    // ...
}
```

### 9. 测试策略

**决策**: 单元测试 + 集成测试 + Mock

**理由**:
- 项目宪章要求 80%+ 测试覆盖率
- 需要测试 Token 管理、HTTP 客户端、数据转换等
- 使用 `gomock` 或 `testify/mock` mock Zeitro API

**测试重点**:
1. **TokenManager 单元测试**:
   - Token 获取和缓存
   - Token 过期检测
   - 并发刷新控制

2. **ZeitroClient 单元测试**:
   - HTTP 请求构建
   - 响应解析
   - 错误处理

3. **集成测试**:
   - 使用 httptest mock Zeitro API 服务器
   - 端到端测试完整流程

4. **基准测试**:
   - Token 并发获取性能
   - 数据转换性能

### 10. 日志和监控

**决策**: 使用 go-zero logx 结构化日志 + 自定义字段

**理由**:
- 项目已有日志基础设施
- 需要记录所有 Zeitro API 交互（FR-008）

**日志记录点**:
```go
// Token 刷新
logx.WithContext(ctx).Infow("Zeitro token refreshed",
    logx.Field("expires_at", token.ExpiresAt))

// API 调用
logx.WithContext(ctx).Infow("Zeitro listing search",
    logx.Field("lat", req.Latitude),
    logx.Field("lon", req.Longitude),
    logx.Field("distance", req.Distance))

// 错误记录
logx.WithContext(ctx).Errorw("Zeitro API error",
    logx.Field("status_code", resp.StatusCode),
    logx.Field("response", string(body)))
```

## 技术栈总结

| 层级 | 技术 | 理由 |
|------|------|------|
| HTTP 客户端 | net/http | Go 标准库，无额外依赖 |
| Token 缓存 | sync.RWMutex + singleflight | 并发安全，防止重复刷新 |
| 配置管理 | go-zero Config | 项目标准 |
| 数据转换 | copierx.Trans + 手动映射 | 项目标准 + 灵活性 |
| 错误处理 | common/errorx | 项目标准 |
| 日志 | go-zero logx | 项目标准 |
| 测试 | go test + gomock | Go 标准库 + 项目标准 |

## 未解决的问题

✅ **无未解决问题** - 所有技术决策已明确

## 依赖和风险

### 依赖项
- ✅ Zeitro API 服务稳定可用
- ✅ borrower 服务的 API/RPC 架构已就绪
- ✅ common/errorx、common/copierx 等工具包已实现

### 风险评估
- 🟡 **中等风险**: Zeitro API 响应格式可能变化
  - 缓解措施：版本检查、优雅降级
- 🟢 **低风险**: Token 并发刷新
  - 缓解措施：singleflight 控制
- 🟢 **低风险**: 网络超时
  - 缓解措施：合理超时配置 + 重试机制

## 下一步

Phase 0 完成，所有技术决策已明确。进入 Phase 1：设计数据模型和 API 契约。
