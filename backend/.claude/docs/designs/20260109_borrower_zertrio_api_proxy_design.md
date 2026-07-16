# zeitro list 接口转发 方案设计

> 来源：tasks/20260109_borrower_zertrio_api_proxy.md（重组为规范设计）

## 1. 背景与目标

### 1.1 业务背景

后端 borrower 端需要展示 zeitro 提供的房源（listing）相关数据。zeitro 对外暴露一套 HTTP 接口，数据结构与本项目的领域模型不同。若由各服务直接对接 zeitro，会出现重复开发与多份鉴权/转换逻辑。

### 1.2 目标

- 统一在 `domains/zeebee/zeitro` 封装 zeitro 接口的调用与数据转换，避免重复开发。
- borrower 通过既有 `web -> api -> rpc -> zeitro` 调用链获取数据，保持代码风格与 borrower 既有模块一致。
- 处理 zeitro access token 过期问题，避免每次请求重新换取 token。

### 1.3 范围边界

- 本设计仅覆盖 borrower 端对 zeitro 接口的转发与转换。
- 涉及分支：从 `feat-backend-lender-register` 切出新分支开发（请勿从主分支开发）。
- 接口文档中给出的数据结构不一定准确，需根据实际情况修改。

## 2. 总体方案

### 2.1 转发架构

沿用项目既有双层架构，zeitro 封装下沉到共享领域层：

```
web (前端)
  └─> borrower API 层（HTTP，校验/路由）
        └─> borrower RPC 层（gRPC，业务编排）
              └─> domains/zeebee/zeitro（接口封装 + 数据转换）
                    └─> zeitro HTTP 接口
```

- zeitro 调用与转换逻辑集中在 `domains/zeebee/zeitro`，作为共享领域层独立可复用。
- borrower 侧 RPC logic 调用该封装，按 borrower 返回给前端的数据结构组织响应。

### 2.2 配置获取

在服务 config 中声明 zeitro 相关配置，运行时从配置注入：

```yaml
Zeitro:
  Url: http://localhost:8080
  Username: admin
  Password: <PASSWORD>
```

### 2.3 接口封装与 Token 缓存

由于 zeitro access token 会过期，在 `domains/zeebee/zeitro` 封装层添加 token 缓存：

- 首次调用：使用 Basic Auth 换取 access token 并缓存。
- token 过期：重新获取并刷新缓存。
- 后续业务接口（如 `/los/openhousesearch`）统一在封装层注入 `Authorization: Bearer <token>`。

### 2.4 接口转换

zeitro 返回的数据结构需转换成 borrower 返回给前端的数据结构。转换逻辑收敛在 `domains/zeebee/zeitro`，保持代码风格与 borrower 既有模块一致。

### 2.5 核心决策

- **封装位置**：`domains/zeebee/zeitro`（共享领域层），非 borrower 服务内部，便于其他服务复用。
- **token 管理**：封装层内缓存 + 过期重取，对调用方透明。
- **调用链不变**：沿用 `web -> api-rpc -> zeitro`，不新增传输层。

## 3. 数据模型

（本设计未涉及）

> 说明：zeitro 返回字段以接口文档为准（见 §4），转换目标为 borrower 前端数据结构；本设计未在文档中定义新的 model / DTO 字段清单，需在实现阶段根据实际响应与 borrower 返回结构补充。

## 4. 接口设计

> 文档中给出的数据结构不一定准确，需要根据实际情况修改。

### 4.1 获取 access token

```
接口：POST  https://lo.zeitro.com/auth/access-token
  - Authorization: Basic Auth
    - Username - "cli_hu4tKiltzCfQfheOa3"
    - Password - "sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd"
  - Headers
    - Content-Type - "application/x-www-form-urlencoded"
  - Body： x-www-form-urlencoded
    - "grant_type"  - "client_credentials"
  - 返回数据参考:
    - {
    -    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2Fub2ZmaWNlcmlkIjoiIiwiY3VzdG9tZXJpZCI6InJhdGVkbmEiLCJyb2xlcyI6bnVsbCwiaXNoYXJkbW9uZXlsb2FuY3VzdG9tZXIiOmZhbHNlLCJleHAiOjE3Njc5MzEwMTgsImlhdCI6MTc2NzkyNzQxOH0.BIF86xnxxVFBW1Z1xfwVpVqzk9OdBAc6HmP6jwbqIuc",
    -    "token_type": "Bearer",
    -    "expires_in": %!d(float64=3600)
    - }
```

### 4.2 房源搜索（openhousesearch）

```
接口：POST  https://lo.zeitro.com/los/openhousesearch
  - Authorization: Bear Token
    - 使用/auth/access-token接口返回的“access_token”值
  - Body: raw
    - 请求数据参考：
    - {
    -    "Lat":37.33874,
    -    "Lon":-121.8852525,
    -    "Dist":5,
    -    "Addr":"San Jose, CA, USA"
    - }
  - 返回数据参考：
    - [
    -    {
    -        "FullAddress": "523 Union Ave, Campbell, CA 95008",
    -        "State": "CA",
    -        "County": "Santa Clara",
    -        "City": "Campbell",
    -        "ZipCode": "95008",
    -        "PropertyType": "Townhouse",
    -        "Price": "1699000                             ",
    -        "IncomeIndicator": "4",
    -        "MinorityPercent": "93.91",
    -        "TotalPopulation": "3662",
    -        "Hispanic": "3302",
    -        "Black": "9",
    -        "Asian": "75",
    -        "Pacific": "10",
    -        "MsaCode": "31084",
    -        "DaysOnMarket": "31",
    -        "AgentName": "Jina Lee",
    -        "AgentPhone": "3109855462",
    -        "AgentEmail": "jina@metisre.com",
    -        "Distance": "4"
    -    },
    - ]
           const response: Response = await fetch("/los/openhousesearch", {
               method: "POST",
               headers: {
                   "Content-Type": "application/json",
                   Authorization: "Bearer " + token,
               },
               body: JSON.stringify({
                   Lat: location.lat,
                   Lon: location.lng,
                   Dist: parseFloat(form.getFieldValue("radius")),
                   Addr: location.address,
               })
           });
```

### 4.3 字段映射

（本设计未涉及）

> 说明：源文档未给出 zeitro 响应字段到 borrower 前端结构的映射表，需在实现阶段根据 borrower 返回结构补充。

## 5. 核心流程

### 5.1 配置加载

1. 服务启动时加载 `Zeitro` 配置（`Url` / `Username` / `Password`）。
2. 将配置注入 `domains/zeebee/zeitro` 封装层。

### 5.2 接口调用与 token 管理

1. borrower RPC logic 调用 `domains/zeebee/zeitro` 封装方法。
2. 封装层检查本地 token 缓存：
   - 无缓存 / 已过期 → 调用 `/auth/access-token`（Basic Auth + `grant_type=client_credentials`）换取新 token，写入缓存。
   - 缓存有效 → 直接复用。
3. 使用 `Authorization: Bearer <token>` 调用 `/los/openhousesearch`，按经纬度 + 半径 + 地址搜索房源。

### 5.3 数据转换与返回

1. 接收 zeitro 响应（字段见 §4.2）。
2. 在 `domains/zeebee/zeitro` 将响应字段转换为 borrower 返回给前端的数据结构。
3. 沿 `rpc -> api -> web` 链路回传，保持 borrower 既有代码风格。

## 6. 风险与权衡

- **鉴权**：Basic Auth 凭证（`cli_hu4tKiltzCfQfheOa3` / `sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd`）须走配置注入，禁止硬编码。
- **token 过期**：需在封装层处理 token 过期重取；`expires_in` 为 3600 秒，缓存策略须考虑时钟漂移与并发刷新（源文档未细化并发场景，实现时需补充）。
- **超时**：源文档未规定 zeitro 接口超时阈值，实现时需为 HTTP 客户端设置合理超时与重试。
- **错误处理**：zeitro 接口失败须按项目规范经 `common/errorx/` 返回，禁止裸 `errors.New`。
- **数据结构偏差**：接口文档明确标注"数据结构不一定准确"，实现阶段须以 zeitro 实际响应为准，必要时回写映射关系。
- **边界**：本设计仅覆盖 borrower 端转发；其他服务接入 zeitro 时复用 `domains/zeebee/zeitro` 封装即可，不在本设计范围内扩展。

## 7. 实现要点 / 测试

（本设计未涉及）

> 说明：源文档未给出实现要点与测试要求。建议按项目宪法（覆盖率目标 80%+）为 `domains/zeebee/zeitro` 封装层补充单元测试，覆盖 token 缓存命中/过期、接口调用失败、字段转换等路径。
