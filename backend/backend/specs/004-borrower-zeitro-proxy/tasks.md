# Tasks: Borrower Zeitro API Proxy

**Input**: Design documents from `backend/specs/004-borrower-zeitro-proxy/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/openapi.yaml ✅

**Tests**: 根据项目宪章，本功能需要包含单元测试（目标覆盖率 80%+），集成测试使用 httptest mock Zeitro API。

**Organization**: 任务按用户故事分组，支持独立实施和测试。根据 plan.md，使用 `make create-service borrower` 命令生成模板代码。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件，无依赖关系）
- **[Story]**: 任务所属用户故事（US1, US2, US3）
- **包含精确文件路径**

## Path Conventions

**Zeebee Backend (Microservices)**:

- API 层: `app/borrower/api/`
- RPC 层: `app/borrower/rpc/`
- 领域层: `domains/zeebee/zeitro/`
- 工具包: `common/errorx/`, `common/copierx/`, `common/validator/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 项目初始化和基础结构，Proto 定义已添加，现在生成代码

- [ ] T001 验证 Proto 定义已在 `app/borrower/rpc/pb/borrower.proto` (第 2002-2045 行)
- [ ] T002 [P] 创建领域层目录结构 `domains/zeebee/zeitro/{client,token,converter,errors,model}`
- [ ] T003 [P] 创建 API 定义文件 `app/borrower/api/desc/zeitro.api` 包含 ZeitroListingSearchRequest/Response 和
  ZeitroPropertyEntity
- [ ] T004 [P] 在 `app/borrower/api/desc/borrower.api` 中导入 `zeitro.api`
- [ ] T005 [P] 配置 golangci-lint 和代码格式化工具（如果尚未配置）

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 所有用户故事依赖的核心基础设施 - 必须先完成此阶段

**⚠️ CRITICAL**: 在此阶段完成前，不能开始任何用户故事的开发

### 领域层模型 (domains/zeebee/zeitro/model/)

- [ ] T006 [P] 实现 `domains/zeebee/zeitro/model/config.go` - ZeitroConfig 结构体和 Validate() 方法
- [ ] T007 [P] 实现 `domains/zeebee/zeitro/model/token.go` - ZeitroToken 结构体和 IsExpired()/IsValid() 方法
- [ ] T008 [P] 实现 `domains/zeebee/zeitro/model/request.go` - SearchRequest 结构体和 Validate() 方法
- [ ] T009 [P] 实现 `domains/zeebee/zeitro/model/property.go` - PropertyListing 和 ZeitroAPIResponse 结构体

### 错误定义 (domains/zeebee/zeitro/errors/)

- [ ] T010 实现 `domains/zeebee/zeitro/errors/errors.go` - 使用 common/errorx 定义所有错误类型（ErrConfigInvalid,
  ErrAuthFailed, ErrTokenExpired 等）

### Token 管理器 (domains/zeebee/zeitro/token/)

- [ ] T011 实现 `domains/zeebee/zeitro/token/manager.go` - TokenManager 结构体包含 GetToken() 和 refreshToken() 方法，使用
  sync.RWMutex + singleflight 实现并发控制

### HTTP 客户端 (domains/zeebee/zeitro/client/)

- [ ] T012 实现 `domains/zeebee/zeitro/client/client.go` - NewClient() 创建客户端，包含 fetchToken() 和 SearchListings()
  方法，配置 HTTP 客户端超时和连接池

### 数据转换器 (domains/zeebee/zeitro/converter/)

- [ ] T013 [P] 实现 `domains/zeebee/zeitro/converter/converter.go` - FromZeitroResponse() 转换 Zeitro API 响应到内部模型

### 单元测试

- [ ] T014 [P] 为 `domains/zeebee/zeitro/model/config.go` 编写单元测试，测试 Validate() 方法
- [ ] T015 [P] 为 `domains/zeebee/zeitro/model/token.go` 编写单元测试，测试 IsExpired() 和 IsValid() 方法
- [ ] T016 [P] 为 `domains/zeebee/zeitro/token/manager.go` 编写单元测试，使用 httptest mock Zeitro API
- [ ] T017 [P] 为 `domains/zeebee/zeitro/client/client.go` 编写单元测试，使用 httptest mock HTTP 响应
- [ ] T018 [P] 为 `domains/zeebee/zeitro/converter/converter.go` 编写单元测试，验证字段映射正确性

**Checkpoint**: 基础设施就绪 - 用户故事实现现在可以并行开始

---

## Phase 3: User Story 1 - 房产列表搜索 (Priority: P1) 🎯 MVP

**Goal**: 用户可以通过 borrower API 搜索 Zeitro 房产列表，系统自动处理 token 认证和数据转换

**Independent Test**: 调用 `POST /borrower/v1/zeitro/listings/search` 接口，传入有效经纬度和半径，验证返回房产列表数据格式符合
borrower 规范

### 实现任务

#### RPC 层实现

- [ ] T019 生成 RPC 服务代码：在 `app/borrower` 目录运行 `make create-service borrower`，这将从 borrower.proto 生成 RPC
  模板代码
- [ ] T020 在 `app/borrower/rpc/internal/config/config.go` 中添加 `Zeitro zeitro.ZeitroConfig` 配置字段
- [ ] T021 在 `app/borrower/rpc/internal/svc/servicecontext.go` 中初始化 ZeitroClient:
  `ZeitroClient: zeitro.NewClient(c.Zeitro)`
- [ ] T022 实现 RPC Logic `app/borrower/rpc/internal/logic/zeitrolistingsearchlogic.go` - 调用
  ZeitroClient.SearchListings()，转换数据为 pb.ZeitroPropertyEntity

#### API 层实现

- [ ] T023 生成 API 服务代码：在 `app/borrower` 目录运行 `make create-service borrower`，这将从 zeitro.api 生成 API 模板代码
- [ ] T024 在 `app/borrower/api/internal/config/config.go` 中添加 `Zeitro zeitro.ZeitroConfig` 配置字段
- [ ] T025 实现 API Logic `app/borrower/api/internal/logic/zeitrolistingsearchlogic.go` - 调用 RPC ZeitroListingSearch()
  ，使用 copierx.Trans 转换响应数据
- [ ] T026 在 `app/borrower/api/internal/handler/zeitrolistingsearchhandler.go` 中添加请求验证（如需要）

#### 配置集成

- [ ] T027 [P] 更新 `app/borrower/api/etc/borrower-api.yaml` 添加 Zeitro 配置段（Url, Username, Password,
  TokenExpireTime, Timeout, MaxRetries）
- [ ] T028 [P] 更新 `app/borrower/rpc/etc/borrower.yaml` 添加相同的 Zeitro 配置段
- [ ] T029 [P] 更新 `app/borrower/etc/config-local.yaml` 添加 Zeitro 配置
- [ ] T030 [P] 更新 `app/borrower/etc/config-production.yaml` 添加 Zeitro 配置

#### 错误处理和日志

- [ ] T031 在 RPC Logic 中添加错误处理，使用 common/errorx 包返回友好错误信息
- [ ] T032 在 API Logic 中添加错误处理，捕获 RPC 错误并转换为 HTTP 响应
- [ ] T033 在所有与 Zeitro API 交互处添加结构化日志（使用 go-zero logx），记录请求、响应和错误

**Checkpoint**: 此时 User Story 1 应完全可用并可独立测试

---

## Phase 4: User Story 2 - 配置管理 (Priority: P2)

**Goal**: 系统管理员可以在配置文件中设置 Zeitro API 参数，系统启动时自动加载并验证配置

**Independent Test**: 修改配置文件中的 Zeitro URL 或认证凭据，重启服务，验证系统能正确加载配置并返回明确的错误信息

### 实现任务

- [ ] T034 [P] 验证 `domains/zeebee/zeitro/model/config.go` 中的 Validate() 方法能检测所有无效配置（空 URL、空用户名、空密码）
- [ ] T035 [P] 在服务启动时调用 ZeitroConfig.Validate()，配置无效时返回明确错误并阻止服务启动
- [ ] T036 [P] 添加配置热重载支持（如需要），监听配置文件变化并在 10 秒内自动重载（FR-007）

**Checkpoint**: 此时 User Stories 1 和 2 都应独立可用

---

## Phase 5: User Story 3 - Token 自动管理 (Priority: P2)

**Goal**: 系统自动管理 Zeitro API token，包括首次获取、缓存、过期检测和自动刷新

**Independent Test**: 监控 token 生命周期，首次请求自动获取 token，token 即将过期时自动刷新，刷新失败时返回明确错误

### 实现任务

- [ ] T037 [P] 验证 TokenManager.GetToken() 能自动获取首次 token 并缓存
- [ ] T038 [P] 验证 TokenManager.GetToken() 能检测 token 过期（剩余 < 10%）并自动刷新
- [ ] T039 [P] 验证 singleflight 机制防止并发重复刷新（模拟多个并发请求）
- [ ] T040 在 token 刷新失败时保留旧 token 并记录错误日志

**Checkpoint**: 所有用户故事现在都应独立可用

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 跨多个用户故事的改进和质量保证

### 测试和验证

- [ ] T041 [P] 编写集成测试，使用 httptest mock Zeitro API 服务器，测试完整的房产搜索流程
- [ ] T042 [P] 运行 `go test ./domains/zeebee/zeitro/...` 验证测试覆盖率达到 80%+
- [ ] T043 [P] 运行 `golangci-lint run` 检查代码质量，修复所有问题
- [ ] T044 [P] 运行 `go fmt` 格式化所有代码
- [ ] T045 [P] 按照成功标准验证：平均响应时间 < 2s, Token 刷新成功率 99.9%, 数据转换准确率 100%

### 性能优化

- [ ] T046 [P] 验证并发处理能力：使用负载测试工具测试 100 并发房产搜索请求
- [ ] T047 [P] 验证缓存命中率：确保缓存命中率达到 95% 以上（除首次请求和 token 过期外）
- [ ] T048 [P] 优化 HTTP 客户端连接池配置（MaxIdleConns, MaxIdleConnsPerHost, IdleConnTimeout）

### 文档

- [ ] T049 [P] 更新 `.claude/docs/INDEX.md` 添加 Zeitro API Proxy 功能索引
- [ ] T050 [P] 运行 `backend/specs/004-borrower-zeitro-proxy/quickstart.md` 中的快速入门步骤验证

### 错误处理完善

- [ ] T051 验证所有边界情况处理：Zeitro API 不可用、超时、数据格式变化、配置参数缺失或格式错误

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: 无依赖 - 可立即开始
- **Foundational (Phase 2)**: 依赖 Setup 完成 - 阻塞所有用户故事
- **User Stories (Phase 3-5)**: 都依赖 Foundational 完成
    - User Story 1 (P1): 核心搜索功能，可独立实现和测试
    - User Story 2 (P2): 配置管理，可独立测试
    - User Story 3 (P2): Token 自动管理，可独立测试
- **Polish (Phase 6)**: 依赖所有用户故事完成

### User Story Dependencies

- **User Story 1 (P1)**: Foundational 完成后可开始 - 无其他故事依赖
- **User Story 2 (P2)**: Foundational 完成后可开始 - 与 US1 集成但可独立测试
- **User Story 3 (P2)**: Foundational 完成后可开始 - 内部功能，独立测试

### Within Each Phase

- **Foundational Phase**: 模型 → 错误 → Token管理 → 客户端 → 转换器 → 测试
- **User Story 1**: RPC 代码生成 → RPC Logic → API 代码生成 → API Logic → 配置 → 错误处理 → 日志
- **User Story 2 & 3**: 独立验证任务，无依赖关系

### Parallel Opportunities

- **Setup Phase**: T002, T003, T004, T005 可并行
- **Foundational Phase**: T006-T009 可并行，T010-T013 可并行，T014-T018 可并行
- **User Story 1**: T027-T030 可并行
- **User Story 2**: T034-T036 可并行
- **User Story 3**: T037-T040 可并行
- **Polish Phase**: T041-T050 可并行

---

## Parallel Example: Foundational Phase

```bash
# 可并行：创建所有领域模型
Task: "实现 config.go" # T006
Task: "实现 token.go" # T007
Task: "实现 request.go" # T008
Task: "实现 property.go" # T009

# 可并行：编写所有单元测试
Task: "测试 config.go" # T014
Task: "测试 token.go" # T015
Task: "测试 token manager" # T016
Task: "测试 client" # T017
Task: "测试 converter" # T018
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. ✅ 完成 Proto 定义追加
2. 完成 Phase 1: Setup
3. **完成 Phase 2: Foundational (CRITICAL - 阻塞所有故事)**
4. 完成 Phase 3: User Story 1
5. **STOP and VALIDATE**: 独立测试 User Story 1
6. 如果就绪，部署/演示 MVP

### Incremental Delivery

1. Setup + Foundational → 基础设施就绪
2. User Story 1 → 独立测试 → 部署/演示（MVP!）
3. User Story 2 → 独立测试 → 部署/演示
4. User Story 3 → 独立测试 → 部署/演示
5. Polish → 最终质量保证

### Parallel Team Strategy

有多位开发者时：

1. 团队共同完成 Setup + Foundational
2. Foundational 完成后：
    - Developer A: User Story 1（核心功能）
    - Developer B: User Story 2（配置管理）
    - Developer C: User Story 3（Token 管理）
3. 故事独立完成并集成

---

## Notes

- **[P] 任务** = 不同文件，无依赖，可并行执行
- **[Story] 标签** = 将任务映射到特定用户故事，便于追溯
- 每个用户故事应独立可完成和可测试
- 在每个检查点停止以独立验证故事
- **开发流程**: 根据 plan.md，使用 `make create-service borrower` 生成模板，然后填充逻辑
- **避免**: 模糊的任务、同文件冲突、破坏独立性的跨故事依赖

## 总计任务数

- **Setup**: 5 个任务
- **Foundational**: 13 个任务
- **User Story 1**: 13 个任务
- **User Story 2**: 3 个任务
- **User Story 3**: 4 个任务
- **Polish**: 11 个任务

**总计**: 49 个任务

## Suggested MVP Scope

**MVP = Phase 1 + Phase 2 + Phase 3 (User Story 1)**

总计 21 个任务，提供核心的房产列表搜索功能。
