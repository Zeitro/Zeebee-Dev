# Feature Specification: Borrower Zeitro API Proxy

**Feature Branch**: `004-borrower-zeitro-proxy`
**Created**: 2025-01-09
**Status**: Draft
**Input**: User description: "转发 Zeitro listing 接口到 borrower 端，实现配置管理、token 缓存和数据结构转换"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 房产列表搜索 (Priority: P1)

作为借款人（borrower）端用户，我需要在系统中搜索 Zeitro 提供的房产列表数据，以便查看可用的房产信息。

**Why this priority**: 这是核心功能，用户访问系统的首要目的是获取房产列表数据，没有这个功能系统无法提供基本服务。

**Independent Test**: 可以通过调用 borrower API 的房产列表搜索接口进行独立测试，验证系统能够成功转发请求到 Zeitro 并返回转换后的数据。

**Acceptance Scenarios**:

1. **Given** 用户通过 borrower API 发起房产搜索请求，**When** 请求包含有效的地理位置和搜索半径，**Then** 系统返回匹配的房产列表数据，数据格式符合 borrower 端规范
2. **Given** Zeitro API 认证 token 有效，**When** 用户发起房产搜索请求，**Then** 系统使用缓存的 token 完成请求，无需重新认证
3. **Given** Zeitro API 认证 token 已过期，**When** 用户发起房产搜索请求，**Then** 系统自动刷新 token 并完成请求，用户无感知

---

### User Story 2 - 配置管理 (Priority: P2)

作为系统管理员，我需要在配置文件中设置 Zeitro API 的连接参数，以便系统能够正确连接到 Zeitro 服务。

**Why this priority**: 配置管理是系统运行的基础，但可以通过配置文件预先设置，不需要用户交互，优先级低于核心搜索功能。

**Independent Test**: 可以通过修改配置文件并重启服务，验证系统能够读取正确的配置并成功连接 Zeitro API。

**Acceptance Scenarios**:

1. **Given** 配置文件中包含有效的 Zeitro URL、用户名和密码，**When** 系统启动，**Then** 系统能够成功连接到 Zeitro API 并获取认证 token
2. **Given** 配置文件中的 Zeitro URL 错误，**When** 系统尝试获取 token，**Then** 系统返回明确的连接错误信息
3. **Given** 配置文件中的认证凭据无效，**When** 系统尝试获取 token，**Then** 系统返回认证失败错误信息

---

### User Story 3 - Token 自动管理 (Priority: P2)

作为系统，我需要自动管理 Zeitro API 的认证 token，包括获取、缓存和自动刷新，以避免频繁的手动认证操作。

**Why this priority**: Token 管理是系统内部功能，对用户不可见，但是确保系统稳定运行的关键机制。

**Independent Test**: 可以通过监控 token 生命周期，验证系统能够在 token 过期前自动刷新，确保请求不会因为 token 过期而失败。

**Acceptance Scenarios**:

1. **Given** 系统首次启动或 token 不存在，**When** 首次请求 Zeitro API，**Then** 系统自动获取并缓存 token
2. **Given** 缓存的 token 即将过期（有效期剩余 < 10%），**When** 发起 API 请求，**Then** 系统自动刷新 token 并更新缓存
3. **Given** token 刷新请求失败，**When** 发起 API 请求，**Then** 系统返回明确的错误信息，并记录错误日志

---

### Edge Cases

- 当 Zeitro API 服务不可用时，系统应返回友好的错误信息，而不是暴露底层异常
- 当 Zeitro API 响应超时时，系统应有合理的超时处理和重试机制
- 当 Zeitro API 返回的数据格式发生变化时，系统应能够优雅降级或返回明确的错误信息
- 当同时有多个请求需要刷新 token 时，系统应避免重复刷新，使用并发控制机制
- 当配置文件中的参数缺失或格式错误时，系统应在启动时返回明确的配置错误信息

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须从配置文件中读取 Zeitro API 连接参数（URL、用户名、密码）
- **FR-002**: 系统必须实现 token 获取功能，使用 Basic Auth 和 client_credentials 授权方式
- **FR-003**: 系统必须实现 token 缓存机制，避免每次请求都重新获取 token
- **FR-004**: 系统必须实现 token 自动刷新机制，在 token 过期前自动更新
- **FR-005**: 系统必须提供房产列表搜索接口，转发请求到 Zeitro API
- **FR-006**: 系统必须将 Zeitro API 的响应数据转换为 borrower 端的数据结构格式
- **FR-007**: 系统必须在配置文件变更后支持热重载，无需重启服务
- **FR-008**: 系统必须记录所有与 Zeitro API 交互的日志，包括请求、响应和错误信息
- **FR-009**: 系统必须实现合理的超时机制，避免长时间等待 Zeitro API 响应
- **FR-010**: 系统必须实现错误处理机制，将 Zeitro API 的错误转换为 borrower 端友好的错误信息

### Key Entities

- **Zeitro Config**: Zeitro API 配置信息
  - URL: Zeitro API 基础地址
  - Username: API 用户名
  - Password: API 密码
  - TokenExpireTime: Token 过期时间（秒）

- **Zeitro Token**: 认证令牌
  - AccessToken: 访问令牌
  - TokenType: 令牌类型（Bearer）
  - ExpiresIn: 过期时间（秒）
  - CreatedAt: 创建时间
  - ExpiresAt: 过期时间戳

- **Property Listing**: 房产列表信息
  - FullAddress: 完整地址
  - State: 州/省
  - County: 县
  - City: 城市
  - ZipCode: 邮编
  - PropertyType: 房产类型
  - Price: 价格
  - Location: 地理位置信息（经纬度）
  - AgentInfo: 房产经纪人信息
  - DaysOnMarket: 上架天数

- **Search Request**: 搜索请求参数
  - Latitude: 纬度
  - Longitude: 经度
  - Distance: 搜索半径（英里）
  - Address: 地址描述

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 房产列表搜索请求的平均响应时间 < 2 秒（不包括 Zeitro API 的响应时间）
- **SC-002**: Token 自动刷新成功率达到 99.9%
- **SC-003**: 数据转换准确率达到 100%，所有字段正确映射
- **SC-004**: 系统能够处理至少 100 并发的房产搜索请求而不出现性能下降
- **SC-005**: 当 Zeitro API 服务不可用时，系统应在 5 秒内返回友好的错误信息
- **SC-006**: 缓存命中率达到 95% 以上（除首次请求和 token 过期外）
- **SC-007**: 错误日志记录完整率达到 100%，所有与 Zeitro API 的交互都有日志记录
- **SC-008**: 配置变更后，系统能够在 10 秒内自动重载配置

## Assumptions & Dependencies

### Assumptions

- Zeitro API 服务稳定可用，SLA 达到 99.5% 以上
- Zeitro API 的认证机制和响应格式在短期内不会发生重大变化
- borrower 端已有成熟的数据结构定义，新接口需要遵循现有规范
- 系统部署环境能够访问外部 Zeitro API 服务

### Dependencies

- Zeitro API 服务必须正常运行
- borrower 端的 API Gateway 和 RPC 服务已实现
- 配置管理系统已就绪，支持 YAML 格式的配置文件
- 日志系统已配置，能够记录 API 交互日志

## Out of Scope

- 本功能不包括 Zeitro API 的其他接口（如房产详情、收藏等），仅实现 listing 接口
- 不包括前端页面的开发，仅提供后端 API
- 不包括 Zeitro API 的监控和告警，仅在日志中记录错误
- 不包括数据的本地持久化缓存，仅实现 token 的内存缓存
