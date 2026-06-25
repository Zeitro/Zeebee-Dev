# Zeebee 后端项目架构分析报告

**分析目标**: 全量后端代码库（~172k LoC，1211 个 Go 源文件）
**分析重点**: 架构、安全、可靠性、性能、可维护性
**分析方式**: 7 维度多 agent 并行 + 对抗性验证
**分析时间**: 2026年6月24日（上次分析: 2025年11月6日）
**分支**: `feat/backend/admin/pricing-plan`

---

## 📋 执行摘要

Zeebee 后端整体架构设计合理——清晰的 handler/logic/repo 分层、微服务边界明确、统一的数据访问层——但本轮深度审计发现 **5 个 Critical、23 个 High、44 个 Medium** 的已确认问题。两个认证后 SQL 注入（9 个 repo 文件）、3 个服务暴露 pprof 端口、pdf_generator 数据竞争，以及系统性的静默错误吞没模式均需立即治理。

架构评估得分: **6.2/10**（上次: 8.2/10）。下降主要来自新发现的安全漏洞和可靠性问题——不是架构退化，而是审计深度增加后的真实暴露。

### 上次行动项完成情况

| 计划项 | 状态 |
|--------|------|
| Prometheus + Grafana 监控集成 | ❌ 未实施 |
| 结构化日志和集中收集 | ⚠️ CloudWatch 日志驱动已配，缺链路追踪 |
| 服务健康检查端点 | ❌ 未实施 |
| Etcd 服务发现和配置管理 | ❌ 未实施 |
| API 网关统一入口 | ❌ 未实施 |
| 熔断器 | ❌ 未实施 |

### 新增问题

- 🔴 **认证后 SQL 注入**：9 个 `Page()` 方法将未校验的 sort 参数直接拼接进 SQL
- 🔴 **pprof 端口公开**：3 个服务在公网端口暴露 heap/goroutine/CPU profile
- 🔴 **数据竞争**：pdf_generator Chrome 池 `allocCtx` 无锁读写
- 🔴 **静默错误吞没**：`_ =` 模式在 6 处丢弃关键错误，导致隐匿数据损坏
- 🟠 **RPC 层缺认证**：5 个 loan/invite 端点绕过 token 校验
- 🟠 **邀请流程断裂**：`accept_invitation`/`check_invitation` 在 ApiAuthMiddleware 内，未认证用户 401

---

## 🏗️ 整体架构概览

### 1.1 当前架构图

```
┌─────────────────────────────────────────────────────────────┐
│                    External Systems                        │
├─────────────────────────────────────────────────────────────┤
│  Web Frontend  │  Mobile App  │  Third-party Services      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Layer                                │
│  CORS │ AuthMiddleware (JWT) │ Rate Limiting               │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────────┐
        ▼                     ▼                          ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────────────┐
│   Borrower    │   │    Lender     │   │       Admin           │
│ (HTTP:10000)  │   │ (HTTP:11000)  │   │    (HTTP:12000)       │
│ (gRPC:10001)  │   │ (gRPC:11001)  │   │    (gRPC:12001)       │
└───────────────┘   └───────────────┘   └───────────────────────┘
                              │
        ┌─────────────────────┤
        ▼                     ▼
┌───────────────┐   ┌───────────────────────────────────────────┐
│ pdf_generator │   │            domains/zeebee/                │
│ (HTTP:13000)  │   │  logic/ │ repo/ │ data/dto/ │ model/     │
│ (Chromedp)    │   │  auth/  │ limiter/ │ constants/           │
└───────────────┘   └───────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Data Layer                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │   PostgreSQL  │  Redis  │  AWS S3  │  Google Drive   │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 服务矩阵

| 服务 | HTTP | gRPC | 职责 | LoC |
|------|------|------|------|-----|
| borrower | 10000 | 10001 | 借款方操作（实体、房产、贷款、申请） | 68,633 |
| lender | 11000 | 11001 | 贷款方操作（审核、文档、批准） | 36,168 |
| admin | 12000 | 12001 | 管理后台（组织、模板、权限、定价计划） | 42,611 |
| pdf_generator | — | — | PDF 生成（Chromedp） | 919 |

### 1.3 架构分层

| 层次 | 职责 | 技术实现 | 评估 |
|------|------|----------|------|
| **表现层** | HTTP API 处理、路由、认证 | go-zero REST framework | ✅ 良好 |
| **业务逻辑层** | 业务规则、流程控制 | Logic 层 + RPC 服务 | ⚠️ 部分端点缺认证 |
| **数据访问层** | 数据持久化、缓存 | GORM + Repository 模式 | 🔴 SQL 注入风险 |
| **基础设施层** | 配置、日志、监控 | go-zero 内置组件 | ⚠️ pprof 公开 + 无链路追踪 |

---

## 🔴 安全审计

### SQL 注入（Critical）

9 个 repo 文件的 `Page()` 方法将未校验的 `sorter.Column` / `sorter.Sort` 直接通过 `fmt.Sprintf` 拼入 GORM `Order()`：

```
影响文件:
  domains/zeebee/repo/organization.go:288
  domains/zeebee/repo/role.go:126
  domains/zeebee/repo/organization_user.go:88
  domains/zeebee/repo/organization_subscription.go:171
  domains/zeebee/repo/base_template.go:111
  domains/zeebee/repo/organization_template.go:89
  domains/zeebee/repo/borrower_entity.go:64
  domains/zeebee/repo/borrower_applicant.go:53
  domains/zeebee/repo/user.go:104
```

修复模式已在 `loan.go:208-221` 和 `borrower_property.go:163-175` 实现（switch 白名单校验）。将所有 `Page()` 统一对齐。

### pprof 端点公开（Critical）

`pdf_generator/api/`、`lender/api/`、`admin/api/` 在公网端口挂载 pprof handler（heap/goroutine/CPU profile）。无环境检查、无认证门控。修复：加 env guard 或 auth middleware。

---

## 🟠 可靠性问题

### pdf_generator 并发缺陷

| 问题 | 位置 | 严重度 |
|------|------|--------|
| `allocCtx` 数据竞争（`URLToPDF`/`HTMLToPDF` 读 vs `tryRestart` 写） | `chrome.go` | High |
| `HTMLToPDF` 挂页面时的无限忙等循环，耗尽信号量槽 | `chrome.go` | High |
| `tryRestart` 持有池互斥锁最长 60s，阻塞所有健康检查和关闭 | `chrome.go` | Medium |
| `chromedp.NewContext(p.allocCtx)` 每请求新启 Chrome 进程 | `chrome.go` | High |

### 静默错误吞没

| 模式 | 文件 | 影响 |
|------|------|------|
| `_ = SyncTotalAmountAndLoanCountWithTx(...)` | `update_borrower_loan_logic.go` | 贷款聚合数据过期 |
| 查询已有用户时忽略 DB 错误 | `register_logic.go` (borrower + lender) | 重复注册 |
| DB 失败视为"未找到" | `sync_data_logic.go` | 幂等性失效 |
| DB 错误与"无默认计划"混同 | `subscribe_default_plan` | 订阅状态异常 |

### N+1 查询（Medium, 5 处）

- `OrganizationUserLogic.Page`: 每用户 2 次额外查询（角色）
- `UserLogic.Info`: 每组织 1 次额外查询（org-user）
- `PricingPlanRepo.Page`: 每计划 1 次额外查询（策略）
- `fillNodeFiles` in DriveLogic: 每文档节点 1 次额外查询（文件）
- `loadOrganizationPlanResources` in auth.go: 每活跃订阅 2 次额外查询

---

## 🟡 代码复用

Borrower 和 Lender 服务间存在大量重复逻辑，`domains/zeebee/logic` 资产生成后未被充分复用：

- 注册邮件发送: `register_logic.go` 在 borrower/lender 中各一份
- 密码验证正则编译: `user_logic.go` 与 `organization_user_logic.go` 重复
- 贷款建设函数: `apply_loan` / `create_borrower_loan` / `update_borrower_loan` / `save_draft_loan` 4 处重复

---

## 📊 架构质量指标（更新）

| 维度 | 评分 (1-10) | 变化 | 说明 |
|------|------------|------|------|
| **安全性** | 4 | -4 | 认证后 SQL 注入 + pprof 公开 |
| **可靠性** | 4 | -2 | 数据竞争 + 静默错误吞没 |
| **可扩展性** | 7 | -1 | 微服务边界合理，但膨胀风险 |
| **可维护性** | 6 | -2 | 跨服务重复代码增多 |
| **性能** | 6 | -1 | N+1 查询 + 每请求新 Chrome 进程 |
| **可观测性** | 4 | — | 仍无链路追踪 |
| **部署自动化** | 6 | — | 容器化完成，CI/CD 不完整 |

**综合评分: 6.2/10** (上次: 8.2/10)

评分下降原因：审计深度增加暴露了此前未检出的安全和可靠性问题，而非架构退化。

---

## 🎯 行动计划（更新）

### Phase 1: 安全止血 (1-2 周)

- [ ] **SQL 注入修复**: 将所有 Page() 方法的 sort 参数改为白名单校验（loan.go 模式）
- [ ] **pprof 下线**: 3 个服务加 env guard 或内部端口隔离
- [ ] **邀请路由修正**: `accept_invitation`/`check_invitation` 移出 ApiAuthMiddleware

### Phase 2: 可靠性加固 (2-4 周)

- [ ] **error-handling 审计**: 替换所有 `_ =` 错误吞没为显式检查 + 日志
- [ ] **pdf_generator 并发修复**: allocCtx Mutex 保护 + Chrome tab 复用 + 忙等超时
- [ ] **RPC 层认证补全**: 5 个缺失 token 校验的端点补 `GetUserClaimFromCtx`

### Phase 3: 架构演进 (1-3 个月)

- [ ] **N+1 查询消除**: Page/Info 热点路径改为批量预加载
- [ ] **跨服务重复消除**: register/reset 逻辑提取到 domains/zeebee/logic
- [ ] **gormx.UpdateOne 修正**: `Save()` → `Updates()`（避免零值字段覆写）
- [ ] **可观测性**: 引入 OpenTelemetry 链路追踪

### Phase 4: 长期规划 (3-6 个月)

- [ ] 服务发现（Etcd/Consul）
- [ ] API 网关统一入口
- [ ] 读写分离 + Redis 多层缓存
- [ ] Kubernetes 编排

---

## 📝 结论

Zeebee 后端仍是设计良好的微服务项目，但安全审计暴露了不容忽视的风险敞口。**SQL 注入和 pprof 公开必须立即修复**（Phase 1），随后系统性地消除静默错误吞没和并发缺陷（Phase 2）。代码重复和 N+1 查询可安排在正常迭代中逐步收敛。

上次分析（2025-11）的 Phase 1/2/3 行动计划均未落地——本次建议缩窄范围、提高可执行性，每阶段不超过 3 项。

---

*分析完成时间: 2026年6月24日（全量深度审计，121 agent，1254 工具调用，19 分钟）*
*上次分析: 2025年11月6日*
*下次评估建议: Phase 1 完成后（2 周内）*
