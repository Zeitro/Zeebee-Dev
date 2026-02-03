# Implementation Plan: Borrower Zeitro API Proxy

**Branch**: `004-borrower-zeitro-proxy` | **Date**: 2025-01-09 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `backend/specs/004-borrower-zeitro-proxy/spec.md`

## Summary

本功能实现了 borrower 端对 Zeitro 房产列表 API 的代理服务，包括：
1. **配置管理**: 从 YAML 文件读取 Zeitro API 连接参数
2. **Token 自动管理**: 自动获取、缓存和刷新 Zeitro API 认证 token
3. **API 代理**: 转发房产列表搜索请求到 Zeitro API
4. **数据转换**: 将 Zeitro API 响应转换为 borrower 端数据结构
5. **错误处理**: 统一的错误处理和日志记录

**技术方法**: 在 `domains/zeebee/zeitro` 中封装 Zeitro 客户端，通过 borrower API/RPC 双层架构暴露功能，使用 `copierx.Trans` 进行数据转换。

**开发流程**: 按照项目规范，使用 `make create-service borrower` 命令生成模板，然后填充逻辑。

## Technical Context

**Language/Version**: Go 1.23+
**Primary Dependencies**: go-zero 1.9.3, GORM, Redis, AWS S3
**Storage**: PostgreSQL 14+ (不用于此功能)
**Testing**: go test, golangci-lint
**Target Platform**: Linux server (microservices architecture)
**Project Type**: microservices (API/RPC separation)
**Performance Goals**: <2s API response (p95), 100+ req/s
**Constraints**: API/RPC dual-layer, shared domain models, copierx.Trans for conversions
**Scale/Scope**: borrower service only, shared domain layer in `domains/zeebee/`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Based on `.specify/memory/constitution.md`:

- [x] **Shared Domain-First**: ✅ 新模型和客户端放在 `domains/zeebee/zeitro/`
- [x] **API/RPC Separation**: ✅ API 层处理验证，RPC 层调用领域服务，领域服务处理业务逻辑
- [x] **Code Generation**: ✅ 使用 `.api` 和 `.proto` 定义并通过 `goctl` 生成代码
- [x] **Data Conversion**: ✅ 使用 `copierx.Trans[T, U]` 进行数据类型转换
- [x] **Type Safety**: ✅ 所有工具函数包含单元测试
- [x] **Test Coverage**: ✅ 实施计划包含单元测试和集成测试，目标 80%+ 覆盖率

**Violations Requiring Justification**: 无违规项

## Project Structure

### Documentation (this feature)

```text
backend/specs/004-borrower-zeitro-proxy/
├── spec.md              # 功能规格文档
├── plan.md              # 本文件 (/speckit.plan command output)
├── research.md          # 技术研究和决策
├── data-model.md        # 数据模型定义
├── quickstart.md        # 快速入门指南
├── contracts/           # API 契约
│   └── openapi.yaml     # OpenAPI 3.0 规范
└── checklists/          # 质量检查清单
    └── requirements.md  # 规格质量检查
```

### Source Code (repository root)

```text
app/borrower/             # Borrower service
├── api/                 # HTTP API layer
│   ├── desc/
│   │   └── zeitro.api  # Zeitro API 定义（需创建）
│   └── internal/
│       ├── handler/     # HTTP handlers
│       │   └── zeitrolistingsearchhandler.go
│       ├── logic/       # API logic (validation)
│       │   └── zeitrolistingsearchlogic.go
│       ├── svc/        # Service context
│       └── types/      # Generated types
└── rpc/                 # gRPC RPC layer
    ├── pb/
    │   └── borrower.proto  # Protobuf 定义（✅ 已追加）
    └── internal/
        └── logic/       # RPC logic (business logic)
            └── zeitrolistingsearchlogic.go

domains/zeebee/          # Shared domain layer
├── model/               # 数据模型（已有）
│   ├── config.go       # ZeitroConfig（需创建）
│   ├── token.go        # ZeitroToken（需创建）
│   ├── property.go     # PropertyListing（需创建）
│   └── request.go      # SearchRequest（需创建）
├── repo/                # (此功能不需要数据库仓储)
└── zeitro/              # Zeitro 集成（需创建）
    ├── client/         # HTTP 客户端
    │   └── client.go
    ├── token/          # Token 管理
    │   └── manager.go
    ├── converter/      # 数据转换
    │   └── converter.go
    └── errors/         # 错误定义
        └── errors.go

common/                  # Shared utilities（已有）
├── errorx/             # Error handling
├── copierx/            # Type conversion
└── validator/          # Data validation
```

**Structure Decision**: 微服务架构 + API/RPC 双层 + 共享领域层。Zeitro 集成放在 `domains/zeebee/zeitro/` 以便其他服务复用。

## Phase 0: Research Summary

**Status**: ✅ 完成

**Output**: [research.md](research.md)

**技术决策**:
1. **架构模式**: API/RPC 双层 + 领域层封装
2. **Token 缓存**: 内存 + `sync.RWMutex` + `singleflight`
3. **HTTP 客户端**: 标准 `net/http`
4. **配置管理**: go-zero Config + YAML
5. **数据转换**: `copierx.Trans` + 手动映射
6. **错误处理**: `common/errorx` + 结构化日志
7. **API 定义**: `zeitro.api` + `goctl` 生成
8. **并发控制**: `golang.org/x/sync/singleflight`
9. **测试策略**: 单元测试 + 集成测试 + Mock
10. **日志监控**: go-zero logx 结构化日志

**未解决的问题**: ✅ 无

## Phase 1: Design Summary

**Status**: ✅ 完成

**Outputs**:
- [data-model.md](data-model.md) - 数据模型定义
- [contracts/openapi.yaml](contracts/openapi.yaml) - API 契约
- [quickstart.md](quickstart.md) - 快速入门指南
- **✅ borrower.proto 已追加 Zeitro 定义**

**核心实体**:
1. **ZeitroConfig**: 配置实体（URL、用户名、密码）
2. **ZeitroToken**: 认证令牌（access_token、过期时间）
3. **PropertyListing**: 房产列表（地址、价格、经纪人等）
4. **SearchRequest**: 搜索请求（经纬度、半径）

**API 端点**:
- `POST /borrower/v1/zeitro/listings/search` - 搜索房产列表

**RPC 方法**:
- `rpc ZeitroListingSearch(ZeitroListingSearchReq) returns(ZeitroListingSearchResp)` - ✅ 已添加到 borrower.proto

**数据转换**: Zeitro API 响应字段映射到 borrower 实体

## Phase 2: Implementation Plan

### 实施阶段

#### 阶段 1: Proto 文件代码生成 ✅ (已完成)

**目标**: 从更新的 borrower.proto 生成 RPC 代码

**任务**:
1. ✅ 追加 Zeitro 定义到 `app/borrower/rpc/pb/borrower.proto`
2. 生成 RPC 代码
   ```bash
   cd app/borrower/rpc
   goctl rpc protoc pb/borrower.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --style=go_zero
   ```

**验收标准**:
- [x] Proto 定义已追加
- [ ] RPC 代码生成成功
- [ ] 编译无错误

#### 阶段 2: API 定义和生成 (1 天)

**目标**: 定义并生成 borrower API 代码

**任务**:
1. 创建 API 定义 (`app/borrower/api/desc/zeitro.api`)

   ```go
   syntax = "v1"

   info(
       title: "zeitro-listing-api"
       desc: "Zeitro listing proxy API for borrower service"
       version: "v1"
   )

   import "common.api"

   type (
       ZeitroListingSearchRequest {
           Latitude  float64  `json:"lat" validate:"required,min=-90,max=90"`
           Longitude float64  `json:"lon" validate:"required,min=-180,max=180"`
           Distance  float64  `json:"dist" validate:"required,min=0,max=100"`
           Address   string   `json:"addr,optional"`
       }

       ZeitroPropertyEntity {
           FullAddress     string  `json:"full_address"`
           State           string  `json:"state"`
           County          string  `json:"county"`
           City            string  `json:"city"`
           ZipCode         string  `json:"zip_code"`
           PropertyType    string  `json:"property_type"`
           Price           string  `json:"price"`
           IncomeIndicator string  `json:"income_indicator,optional"`
           MinorityPercent string  `json:"minority_percent,optional"`
           TotalPopulation string  `json:"total_population,optional"`
           Hispanic        string  `json:"hispanic,optional"`
           Black           string  `json:"black,optional"`
           Asian           string  `json:"asian,optional"`
           Pacific         string  `json:"pacific,optional"`
           MsaCode         string  `json:"msa_code,optional"`
           DaysOnMarket    string  `json:"days_on_market,optional"`
           AgentName       string  `json:"agent_name,optional"`
           AgentPhone      string  `json:"agent_phone,optional"`
           AgentEmail      string  `json:"agent_email,optional"`
           Distance        string  `json:"distance,optional"`
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
       @doc "Search property listings from Zeitro API"
       @handler zeitroListingSearch
       post /zeitro/listings/search (ZeitroListingSearchRequest) returns (ZeitroListingSearchResponse)
   }
   ```

2. 生成 API 代码
   ```bash
   cd app/borrower
   make create-service borrower
   ```

3. 更新 borrower.api 主文件
   - 导入 zeitro.api

**验收标准**:
- [ ] API 定义正确
- [ ] 代码生成成功
- [ ] 编译无错误

#### 阶段 3: 领域层实现 (3-4 天)

**目标**: 实现 `domains/zeebee/zeitro` 中的核心功能

**任务**:
1. 创建目录结构
   ```bash
   mkdir -p domains/zeebee/zeitro/{client,token,converter,errors,model}
   ```

2. 实现配置模型 (`domains/zeebee/zeitro/model/config.go`)
   - 定义 `ZeitroConfig` 结构
   - 实现 `Validate()` 方法
   - 添加单元测试

3. 实现 Token 模型 (`domains/zeebee/zeitro/model/token.go`)
   - 定义 `ZeitroToken` 结构
   - 实现 `IsExpired()` 和 `IsValid()` 方法
   - 添加单元测试

4. 实现请求模型 (`domains/zeebee/zeitro/model/request.go`)
   - 定义 `SearchRequest` 结构
   - 实现 `Validate()` 方法
   - 添加单元测试

5. 实现房产模型 (`domains/zeebee/zeitro/model/property.go`)
   - 定义 `PropertyListing` 结构
   - 定义 Zeitro API 响应结构
   - 添加单元测试

6. 实现错误定义 (`domains/zeebee/zeitro/errors/errors.go`)
   - 定义所有错误类型
   - 使用 `common/errorx`

7. 实现 Token Manager (`domains/zeebee/zeitro/token/manager.go`)
   - `GetToken()` - 获取有效 token
   - `refreshToken()` - 刷新 token
   - 并发控制（`sync.RWMutex` + `singleflight`）
   - 单元测试

8. 实现 Zeitro 客户端 (`domains/zeebee/zeitro/client/client.go`)
   - `NewClient()` - 创建客户端
   - `fetchToken()` - 获取新 token
   - `SearchListings()` - 搜索房产列表
   - HTTP 客户端配置
   - 单元测试（使用 httptest）

9. 实现数据转换器 (`domains/zeebee/zeitro/converter/converter.go`)
   - `FromZeitroResponse()` - 转换 Zeitro 响应
   - 字段清理（TrimSpace、格式化）
   - 单元测试

**验收标准**:
- [ ] 所有模型有单元测试
- [ ] 测试覆盖率 > 80%
- [ ] `golangci-lint` 检查通过
- [ ] Token 自动获取和缓存
- [ ] Token 过期自动刷新
- [ ] 并发请求不会重复刷新
- [ ] 成功调用 Zeitro Listing API
- [ ] 数据正确转换
- [ ] 错误处理完整

#### 阶段 4: RPC Logic 实现 (2-3 天)

**目标**: 实现 borrower RPC 业务逻辑

0. 生成 RPC 代码
   ```bash
   cd app/borrower
   make create-service borrower
   ```

**任务**:
1. 实现 RPC Logic (`app/borrower/rpc/internal/logic/zeitrolistingsearchlogic.go`)
   ```go
   package logic

   import (
       "context"

       "yourusername/app/borrower/rpc/internal/svc"
       "yourusername/app/borrower/rpc/pb"
       "yourusername/domains/zeebee/zeitro"
       "yourusername/domains/zeebee/zeitro/model"
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
       // 转换请求
       req := &model.SearchRequest{
           Latitude:  in.Lat,
           Longitude: in.Lon,
           Distance:  in.Dist,
           Address:   in.Addr,
       }

       // 验证请求
       if err := req.Validate(); err != nil {
           return nil, errorx.NewDefaultError("Invalid search parameters: " + err.Error())
       }

       // 调用 Zeitro 客户端
       listings, err := l.svcCtx.ZeitroClient.SearchListings(l.ctx, req)
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

   func convertToProto(listing *model.PropertyListing) *pb.ZeitroPropertyEntity {
       return &pb.ZeitroPropertyEntity{
           FullAddress:     listing.FullAddress,
           State:           listing.State,
           County:          listing.County,
           City:            listing.City,
           ZipCode:         listing.ZipCode,
           PropertyType:    listing.PropertyType,
           Price:           listing.Price,
           IncomeIndicator: listing.IncomeIndicator,
           MinorityPercent: listing.MinorityPercent,
           TotalPopulation: listing.TotalPopulation,
           Hispanic:        listing.Hispanic,
           Black:           listing.Black,
           Asian:           listing.Asian,
           Pacific:         listing.Pacific,
           MsaCode:         listing.MsaCode,
           DaysOnMarket:    listing.DaysOnMarket,
           AgentName:       listing.AgentName,
           AgentPhone:      listing.AgentPhone,
           AgentEmail:      listing.AgentEmail,
           Distance:        listing.Distance,
       }
   }
   ```

2. 注册 RPC 服务 (`app/borrower/rpc/internal/svc/servicecontext.go`)
   ```go
   type ServiceContext struct {
       Config      config.Config
       ZeitroClient *zeitro.Client  // 添加这一行
   }

   func NewServiceContext(c config.Config) *ServiceContext {
       return &ServiceContext{
           Config:      c,
           ZeitroClient: zeitro.NewClient(c.Zeitro),  // 添加这一行
       }
   }
   ```

3. 更新配置结构 (`app/borrower/rpc/internal/config/config.go`)
   ```go
   type Config struct {
       // ... 现有配置
       Zeitro zeitro.ZeitroConfig  // 添加这一行
   }
   ```

**验收标准**:
- [ ] RPC Logic 实现完整
- [ ] 调用 Zeitro 客户端成功
- [ ] 数据转换正确
- [ ] 错误处理完整

#### 阶段 5: API 层实现 (1-2 天)

**目标**: 实现 borrower API Handler 和 Logic

**任务**:
1. 实现 API Handler (`app/borrower/api/internal/handler/zeitrolistingsearchhandler.go`)
   - 已由 goctl 生成，只需调整

2. 实现 API Logic (`app/borrower/api/internal/logic/zeitrolistingsearchlogic.go`)
   - 已由 make create-service borrower 生成，需要填充 RPC 调用逻辑
   ```go
   func (l *ZeitroListingSearchLogic) ZeitroListingSearch(req *types.ZeitroListingSearchRequest) (resp *types.ZeitroListingSearchResponse, err error) {
       // 调用 RPC
       rpcResp, err := l.svcCtx.ZeitroRpc.ZeitroListingSearch(l.ctx, &borrower.ZeitroListingSearchReq{
           Lat:   req.Latitude,
           Lon:   req.Longitude,
           Dist:  req.Distance,
           Addr:  req.Address,
       })
       if err != nil {
           return nil, err
       }

       // 转换响应
       var list []types.ZeitroPropertyEntity
       for _, item := range rpcResp.List {
           list = append(list, types.ZeitroPropertyEntity{
               FullAddress:     item.FullAddress,
               State:           item.State,
               // ... 其他字段
           })
       }

       return &types.ZeitroListingSearchResponse{
           List: list,
       }, nil
   }
   ```

3. 更新 API ServiceContext (`app/borrower/api/internal/svc/servicecontext.go`)
   - 添加 ZeitroRpc 客户端（如果需要）

**验收标准**:
- [ ] API 端点可访问
- [ ] 请求验证正确
- [ ] RPC 调用成功
- [ ] 响应格式正确

#### 阶段 6: 配置集成 (1 天)

**目标**: 集成配置管理

**任务**:
1. 更新 API 配置结构 (`app/borrower/api/internal/config/config.go`)
   ```go
   type Config struct {
       // ... 现有配置
       Zeitro zeitro.ZeitroConfig  // 添加这一行
   }
   ```

2. 更新配置文件 (`app/borrower/etc/config-{env}.yaml`)
   ```yaml
   # borrower-api.yaml

   Zeitro:
     Url: https://lo.zeitro.com
     Username: cli_hu4tKiltzCfQfheOa3
     Password: sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd
     TokenExpireTime: 3600
     Timeout: 10
     MaxRetries: 3
   ```

3. 更新 RPC 配置文件 (`app/borrower/rpc/etc/borrower.yaml`)
   - 添加相同的 Zeitro 配置

**验收标准**:
- [ ] 配置正确加载
- [ ] 不同环境配置生效
- [ ] Zeitro 客户端初始化成功

#### 阶段 7: 测试和优化 (2-3 天)

**目标**: 完善测试和性能优化

**任务**:
1. 单元测试
   - Token Manager 测试
   - HTTP 客户端测试
   - 数据转换测试
   - 目标覆盖率 > 80%

2. 集成测试
   - 端到端测试
   - 使用 httptest mock Zeitro API
   - 测试错误场景

3. 性能测试
   - 并发测试（100 QPS）
   - Token 缓存命中率测试
   - 响应时间测试

4. 优化
   - 并发控制优化
   - 日志优化
   - 错误处理优化

**验收标准**:
- [ ] 测试覆盖率 > 80%
- [ ] 性能指标达标
- [ ] `golangci-lint` 检查通过

#### 阶段 8: 文档和部署 (1 天)

**目标**: 完善文档和准备部署

**任务**:
1. 更新 API 文档
2. 更新 README
3. 编写部署指南
4. 准备监控和告警

**验收标准**:
- [ ] 文档完整
- [ ] 部署指南清晰
- [ ] 监控配置就绪

### 任务清单总结

| 阶段 | 任务 | 预估时间 | 优先级 |
|-----|------|---------|-------|
| 1 | Proto 代码生成 | 0.5 天 | P0 |
| 2 | API 定义和生成 | 1 天 | P0 |
| 3 | 领域层实现 | 3-4 天 | P0 |
| 4 | RPC Logic 实现 | 2-3 天 | P0 |
| 5 | API 层实现 | 1-2 天 | P0 |
| 6 | 配置集成 | 1 天 | P1 |
| 7 | 测试和优化 | 2-3 天 | P0 |
| 8 | 文档和部署 | 1 天 | P1 |

**总预估时间**: 11.5-15.5 天

## Dependencies

### 外部依赖
- Zeitro API 服务（必须可用）
- borrower API/RPC 服务（已有）
- 配置管理系统（已有）

### 内部依赖
- `common/errorx` - 错误处理
- `common/copierx` - 数据转换
- `common/validator` - 数据验证
- `common/uuid` - UUID 生成（如需要）

## Risk Mitigation

| 风险 | 影响 | 概率 | 缓解措施 |
|-----|------|------|---------|
| Zeitro API 格式变化 | 高 | 中 | 版本检查、优雅降级 |
| Token 刷新失败 | 高 | 低 | 重试机制、保留旧 token |
| 并发请求冲突 | 中 | 低 | singleflight 控制 |
| 性能不达标 | 中 | 低 | 连接池优化、缓存优化 |

## Success Criteria

根据 [spec.md](spec.md) 的成功标准：

- [ ] **SC-001**: 房产列表搜索请求的平均响应时间 < 2 秒
- [ ] **SC-002**: Token 自动刷新成功率达到 99.9%
- [ ] **SC-003**: 数据转换准确率达到 100%
- [ ] **SC-004**: 系统能够处理至少 100 并发的房产搜索请求
- [ ] **SC-005**: 当 Zeitro API 服务不可用时，系统应在 5 秒内返回友好的错误信息
- [ ] **SC-006**: 缓存命中率达到 95% 以上
- [ ] **SC-007**: 错误日志记录完整率达到 100%
- [ ] **SC-008**: 配置变更后，系统能够在 10 秒内自动重载配置

## Next Steps

1. ✅ **Proto 定义已添加** - 运行生成命令创建 RPC 代码
2. 开始实施：按照阶段顺序执行
3. 进度跟踪：使用任务清单跟踪进度

## References

- [功能规格](spec.md)
- [技术研究](research.md)
- [数据模型](data-model.md)
- [API 契约](contracts/openapi.yaml)
- [快速入门](quickstart.md)
- [项目宪章](.specify/memory/constitution.md)
