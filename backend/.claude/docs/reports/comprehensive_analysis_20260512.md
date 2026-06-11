# Zeebee Backend 综合分析报告

**分析日期**: 2026-05-12 | **分支**: feat/backend/borrower/property-transaction
**代码规模**: 978 Go 文件 | 75 API 定义 | 15 测试文件

---

## 总评

| 维度 | 评分 | 状态 |
|------|------|------|
| 🏗️ 架构 | 8.5/10 | 优秀 — DDD 分层清晰，共享领域层设计合理 |
| 🔒 安全 | 5.5/10 | 需改进 — 硬编码密钥、CORS 通配符、TLS 跳过 |
| ⚡ 性能 | 7.0/10 | 良好 — N+1 查询风险、缺少索引 |
| 🧪 测试 | 2.0/10 | 严重不足 — 15 测试文件，~1.5% 覆盖率 |
| 📐 代码质量 | 7.0/10 | 良好 — 错误处理不一致、低测试覆盖 |

---

## 🔒 安全问题

### 🔴 严重

| # | 问题 | 文件 | 风险 |
|---|------|------|------|
| S1 | **生产配置硬编码密钥** — JWT secret、DB 密码、AWS S3 AK/SK 明文 | `app/admin/etc/config-production.yaml` | 凭证泄露 |
| S2 | **CORS 通配符** — `rest.WithCors("*")` 允许任意来源 | `app/admin/api/admin.go:23`, `app/borrower/api/borrower.go:22` | CSRF 攻击 |
| S3 | **Redis TLS 跳过** — `InsecureSkipVerify: true` | `common/limiter/redis_limiter.go:26`, `common/redisx/client.go:40`, `common/rbac/enforcer.go:64` | 中间人攻击 |

### 🟡 中等

| # | 问题 | 文件 | 风险 |
|---|------|------|------|
| S4 | **无 JWT 黑名单机制** — token 无法撤销 | `common/utils/jwt.go` | 失效 token 仍可用 |
| S5 | **错误栈暴露** — `errorx.GetErrorStack()` 包含内部细节 | `common/errorx/error.go` | 信息泄露 |

### ✅ 安全优势

- GORM 参数化查询，无 SQL 注入风险
- JWT HMAC-SHA256 签名，bcrypt 密码加密
- 全端点输入验证（`common/validator/`）
- 认证中间件正确实现（Bearer token + gRPC 验证）

---

## ⚡ 性能问题

### 🔴 严重

| # | 问题 | 文件 | 影响 |
|---|------|------|------|
| P1 | **N+1 查询** — 分页中每条 loan 触发 2 次额外查询（property + applicants） | `domains/zeebee/repo/loan.go:235-242, 386-393` | 20 条/页 = 41 次查询 |
| P2 | **缺少数据库索引** — `borrower_organization_id`, `lender_organization_id`, `stage`, `pending_status` 无索引 | `domains/zeebee/model/loan.go` | 全表扫描 |

### 🟡 中等

| # | 问题 | 文件 | 影响 |
|---|------|------|------|
| P3 | **FindAll 无分页** — `base_repo.go` FindAll() 可加载全表 | `common/gormx/base_repo.go:134-141` | 内存溢出风险 |
| P4 | **Redis 无连接池配置** | `common/redisx/client.go:12-24` | 连接开销 |
| P5 | **S3 操作无超时** — `context.Background()` 无 timeout | `app/admin/rpc/internal/global/s3.go:11-18` | 挂起风险 |
| P6 | **循环中字符串拼接** — 未使用 `strings.Builder` | `app/admin/rpc/internal/logic/common/sync_data_logic.go:86-107` | 内存浪费 |

---

## 🏗️ 架构质量

### ✅ 优势

- DDD 分层清晰：API → RPC → Domain → Infrastructure
- 共享领域层（27 models, 27 repos, 26 DTOs）实现代码复用
- Repository 模式统一数据访问
- go-zero 代码生成保持一致性
- copierx 泛型转换类型安全

### 🟡 问题

| # | 问题 | 详情 |
|---|------|------|
| A1 | **循环依赖风险** — `common/types` 引用 `domains/zeebee/constants`，反向也有引用 | 打破分层方向 |
| A2 | **错误处理不一致** — 混用 `errorx.New()`, `errors.New()`, 未包装错误 | 多处 RPC 逻辑 |
| A3 | **AuthorityMiddleware 未实现** — 权限控制空实现（admin 已在 routes.go 中注释） | 权限缺失 |

---

## 🧪 测试覆盖

| 模块 | 测试文件数 | 覆盖率估计 |
|------|-----------|-----------|
| common/ | ~8 | ~15% |
| domains/ | ~5 | ~3% |
| app/ 服务层 | ~2 | <1% |
| **总计** | **15** | **~1.5%** |

**目标**: 80%+ 覆盖率（宪法 §V 要求）

---

## 优先修复建议

### 立即（1-3 天）

1. **S1**: 移除生产配置中的硬编码密钥，改用环境变量
2. **S2**: CORS 限制为具体域名白名单
3. **S3**: 生产环境启用 Redis TLS 证书验证
4. **P1**: 用 Preload 或 JOIN 替代 N+1 查询

### 短期（1-2 周）

5. **P2**: 为高频查询字段添加数据库索引
6. **P3**: FindAll 添加默认分页限制
7. **A2**: 统一错误处理，全部使用 errorx
8. **S4**: 实现 JWT token 黑名单

### 中期（1-2 月）

9. 测试覆盖率提升至 50%+
10. A1: 解决 common/domains 循环依赖
11. P5/P6: 外部调用超时配置

---

*本报告由 /sc:analyze 自动生成 | 2026-05-12*
