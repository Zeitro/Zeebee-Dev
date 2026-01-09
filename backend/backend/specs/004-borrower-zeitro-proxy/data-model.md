# Data Model: Borrower Zeitro API Proxy

**Feature**: 004-borrower-zeitro-proxy
**Date**: 2025-01-09
**Phase**: Phase 1 - Design & Contracts

## 实体关系图

```
┌─────────────────┐
│  ZeitroConfig   │
│  (配置)          │
└─────────────────┘
         │
         │ uses
         ↓
┌─────────────────┐       ┌──────────────────┐
│  TokenManager   │──────→│   ZeitroToken    │
│  (Token 管理器)  │       │   (认证令牌)      │
└─────────────────┘       └──────────────────┘
         │
         │ calls
         ↓
┌─────────────────┐       ┌──────────────────┐
│  ZeitroClient   │──────→│ PropertyListing  │
│  (Zeitro 客户端) │       │  (房产列表)       │
└─────────────────┘       └──────────────────┘
```

## 核心实体定义

### 1. ZeitroConfig（配置实体）

**位置**: `domains/zeebee/zeitro/model/config.go`

```go
package model

// ZeitroConfig represents the configuration for Zeitro API integration
type ZeitroConfig struct {
    URL             string  `json:"url" toml:"url"`                     // Zeitro API base URL
    Username        string  `json:"username" toml:"username"`           // API username
    Password        string  `json:"password" toml:"password"`           // API password
    TokenExpireTime int64   `json:"token_expire_time" toml:"token_expire_time"` // Token expiration in seconds
    Timeout         int64   `json:"timeout" toml:"timeout"`             // Request timeout in seconds
    MaxRetries      int     `json:"max_retries" toml:"max_retries"`     // Maximum retry attempts
}

// Validate checks if the configuration is valid
func (c *ZeitroConfig) Validate() error {
    if c.URL == "" {
        return errors.New("Zeitro URL is required")
    }
    if c.Username == "" {
        return errors.New("Zeitro username is required")
    }
    if c.Password == "" {
        return errors.New("Zeitro password is required")
    }
    if c.Timeout <= 0 {
        c.Timeout = 10 // default 10 seconds
    }
    if c.MaxRetries <= 0 {
        c.MaxRetries = 3 // default 3 retries
    }
    return nil
}
```

**验证规则**:

- URL: 必填，有效的 HTTP(S) URL
- Username: 必填，非空字符串
- Password: 必填，非空字符串
- TokenExpireTime: 默认 3600 秒（1 小时）
- Timeout: 默认 10 秒，最小 1 秒
- MaxRetries: 默认 3 次，最小 0 次

### 2. ZeitroToken（认证令牌）

**位置**: `domains/zeebee/zeitro/model/token.go`

```go
package model

import "time"

// ZeitroToken represents the authentication token for Zeitro API
type ZeitroToken struct {
    AccessToken  string    `json:"access_token"`  // The actual access token
    TokenType    string    `json:"token_type"`    // Token type (e.g., "Bearer")
    ExpiresIn    int64     `json:"expires_in"`    // Expiration time in seconds
    CreatedAt    time.Time `json:"created_at"`    // Token creation timestamp
    ExpiresAt    time.Time `json:"expires_at"`    // Token expiration timestamp
}

// IsExpired checks if the token is expired or will expire soon
func (t *ZeitroToken) IsExpired() bool {
    // Consider token expired if less than 10% of lifetime remains
    threshold := time.Duration(t.ExpiresIn/10) * time.Second
    return time.Until(t.ExpiresAt) < threshold
}

// IsValid checks if the token is valid and not expired
func (t *ZeitroToken) IsValid() bool {
    return t.AccessToken != "" && !t.ExpiresAt.IsZero() && !t.IsExpired()
}
```

**状态转换**:

```
┌──────────┐
│  未创建   │
└────┬─────┘
     │ 获取 token
     ↓
┌──────────┐
│  有效期   │ ←──┐
└────┬─────┘    │
     │ 过期      │ 10% 时间阈值
     │           │
     ↓           │
┌──────────┐    │
│  即将过期 │────┘
└────┬─────┘
     │ 刷新
     ↓
┌──────────┐
│  已刷新   │
└──────────┘
```

### 3. PropertyListing（房产列表）

**位置**: `domains/zeebee/zeitro/model/property.go`

```go
package model

// PropertyListing represents a property listing from Zeitro API
type PropertyListing struct {
    // Location information
    FullAddress   string  `json:"full_address"`   // Complete address
    State         string  `json:"state"`           // State code (e.g., "CA")
    County        string  `json:"county"`          // County name
    City          string  `json:"city"`            // City name
    ZipCode       string  `json:"zip_code"`        // ZIP code

    // Property details
    PropertyType  string  `json:"property_type"`   // Property type (e.g., "Townhouse")
    Price         string  `json:"price"`           // Listing price

    // Demographics
    IncomeIndicator   string  `json:"income_indicator"`   // Income level indicator
    MinorityPercent   string  `json:"minority_percent"`   // Minority percentage
    TotalPopulation   string  `json:"total_population"`  // Total population
    Hispanic          string  `json:"hispanic"`           // Hispanic population count
    Black             string  `json:"black"`              // Black population count
    Asian             string  `json:"asian"`              // Asian population count
    Pacific           string  `json:"pacific"`            // Pacific population count

    // Market information
    MsaCode       string  `json:"msa_code"`       // Metropolitan Statistical Area code
    DaysOnMarket  string  `json:"days_on_market"` // Days listed on market

    // Agent information
    AgentName  string `json:"agent_name"`  // Agent's name
    AgentPhone string `json:"agent_phone"` // Agent's phone number
    AgentEmail string `json:"agent_email"` // Agent's email address

    // Search metadata
    Distance    string  `json:"distance"`     // Distance from search location
}

// PropertyList represents a list of property listings
type PropertyList []*PropertyListing
```

**字段说明**:

- 所有字段使用字符串类型以保持与 Zeitro API 响应一致
- 价格、人口等数值字段需要前端或消费者进行类型转换

### 4. SearchRequest（搜索请求）

**位置**: `domains/zeebee/zeitro/model/request.go`

```go
package model

// SearchRequest represents a search request for property listings
type SearchRequest struct {
    Latitude  float64 `json:"lat" validate:"required,min=-90,max=90"`  // Latitude in decimal degrees
    Longitude float64 `json:"lon" validate:"required,min=-180,max=180"` // Longitude in decimal degrees
    Distance  float64 `json:"dist" validate:"required,min=0,max=100"`  // Search radius in miles
    Address   string  `json:"addr"`                                     // Human-readable address
}

// Validate checks if the search request is valid
func (r *SearchRequest) Validate() error {
    if r.Latitude < -90 || r.Latitude > 90 {
        return errors.New("latitude must be between -90 and 90")
    }
    if r.Longitude < -180 || r.Longitude > 180 {
        return errors.New("longitude must be between -180 and 180")
    }
    if r.Distance <= 0 || r.Distance > 100 {
        return errors.New("distance must be between 0 and 100 miles")
    }
    return nil
}
```

**验证规则**:

- Latitude: -90 到 90 度
- Longitude: -180 到 180 度
- Distance: 0 到 100 英里（可配置最大值）
- Address: 可选，用于日志记录

## 数据转换映射

### Zeitro API 响应 → Borrower Entity

| Zeitro 字段       | Borrower 字段      | 转换规则         |
|-----------------|------------------|--------------|
| FullAddress     | full_address     | 直接映射，保留原格式   |
| State           | state            | 直接映射         |
| County          | county           | 直接映射         |
| City            | city             | 直接映射         |
| ZipCode         | zip_code         | 直接映射，去除空格    |
| PropertyType    | property_type    | 直接映射         |
| Price           | price            | 去除首尾空格       |
| IncomeIndicator | income_indicator | 直接映射         |
| MinorityPercent | minority_percent | 直接映射         |
| TotalPopulation | total_population | 直接映射         |
| Hispanic        | hispanic         | 直接映射         |
| Black           | black            | 直接映射         |
| Asian           | asian            | 直接映射         |
| Pacific         | pacific          | 直接映射         |
| MsaCode         | msa_code         | 直接映射         |
| DaysOnMarket    | days_on_market   | 转换为整数        |
| AgentName       | agent_name       | 直接映射         |
| AgentPhone      | agent_phone      | 格式化为标准格式（可选） |
| AgentEmail      | agent_email      | 验证邮箱格式（可选）   |
| Distance        | distance         | 转换为浮点数       |

### 转换函数示例

**位置**: `domains/zeebee/zeitro/converter/converter.go`

```go
package converter

import (
    "strconv"
    "strings"
    "zeitro/model"
)

// PropertyConverter handles conversion between Zeitro and internal formats
type PropertyConverter struct{}

// NewPropertyConverter creates a new property converter
func NewPropertyConverter() *PropertyConverter {
    return &PropertyConverter{}
}

// FromZeitroResponse converts Zeitro API response to internal model
func (c *PropertyConverter) FromZeitroResponse(z *ZeitroAPIResponse) *model.PropertyListing {
    if z == nil {
        return nil
    }

    return &model.PropertyListing{
        FullAddress:       strings.TrimSpace(z.FullAddress),
        State:            strings.TrimSpace(z.State),
        County:           strings.TrimSpace(z.County),
        City:             strings.TrimSpace(z.City),
        ZipCode:          strings.TrimSpace(z.ZipCode),
        PropertyType:     strings.TrimSpace(z.PropertyType),
        Price:            strings.TrimSpace(z.Price),
        IncomeIndicator:  strings.TrimSpace(z.IncomeIndicator),
        MinorityPercent:  strings.TrimSpace(z.MinorityPercent),
        TotalPopulation:  strings.TrimSpace(z.TotalPopulation),
        Hispanic:         strings.TrimSpace(z.Hispanic),
        Black:            strings.TrimSpace(z.Black),
        Asian:            strings.TrimSpace(z.Asian),
        Pacific:          strings.TrimSpace(z.Pacific),
        MsaCode:          strings.TrimSpace(z.MsaCode),
        DaysOnMarket:     strings.TrimSpace(z.DaysOnMarket),
        AgentName:        strings.TrimSpace(z.AgentName),
        AgentPhone:       formatPhone(z.AgentPhone),
        AgentEmail:       strings.TrimSpace(z.AgentEmail),
        Distance:         strings.TrimSpace(z.Distance),
    }
}

// formatPhone formats phone number to standard format
func formatPhone(phone string) string {
    // Remove all non-numeric characters
    phone = strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, phone)
    return phone
}
```

## 存储策略

### Token 缓存

- **存储方式**: 内存（进程内）
- **数据结构**: `*ZeitroToken` 单例
- **并发控制**: `sync.RWMutex` + `singleflight`
- **生命周期**:
    - 首次获取: 服务启动后第一次 API 调用
    - 刷新时机: 剩余时间 < 10%
    - 失效处理: 刷新失败时保留旧 token，记录错误日志

### 配置存储

- **存储方式**: YAML 配置文件
- **文件位置**:
    - 开发: `app/borrower/etc/config-local.yaml`
    - 生产: `app/borrower/etc/config-production.yaml`
- **热重载**: 使用 go-zero 的配置监听机制

### 数据持久化

- **不持久化**: 房产列表数据不本地存储
- **原因**: 实时数据，时效性强
- **缓存策略**: 前端或 API Gateway 层可添加缓存

## 数据验证

### 输入验证

使用 `common/validator` 进行请求验证：

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateSearchRequest(req *SearchRequest) error {
    return validate.Struct(req)
}
```

### 输出验证

- 检查必要字段是否为空
- 验证价格格式
- 验证地理坐标范围

## 错误处理

### 错误类型

```go
package errors

import "github.com/yourusername/common/errorx"

var (
    // Configuration errors
    ErrConfigInvalid     = errorx.NewDefaultError("Zeitro configuration is invalid")
    ErrConfigMissing     = errorx.NewDefaultError("Zeitro configuration is missing")

    // Authentication errors
    ErrAuthFailed        = errorx.NewDefaultError("Zeitro authentication failed")
    ErrTokenExpired      = errorx.NewDefaultError("Zeitro token has expired")
    ErrTokenRefreshFailed = errorx.NewDefaultError("Failed to refresh token")

    // Network errors
    ErrConnectionFailed  = errorx.NewDefaultError("Failed to connect to Zeitro API")
    ErrTimeout          = errorx.NewDefaultError("Zeitro API request timed out")
    ErrRateLimited      = errorx.NewDefaultError("Zeitro API rate limit exceeded")

    // API errors
    ErrAPIError         = errorx.NewDefaultError("Zeitro API returned an error")
    ErrInvalidResponse  = errorx.NewDefaultError("Invalid response from Zeitro API")

    // Validation errors
    ErrInvalidInput     = errorx.NewDefaultError("Invalid search request parameters")
)
```

## 下一步

数据模型已定义完成，接下来生成 API 契约文档。
