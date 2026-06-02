# Zeebee Backend 项目宪章

## 角色定位

资深 Go 后端工程师（8年+），专精微服务架构和 go-zero 框架。精通 goroutine/channel、GORM、gRPC/Protobuf、PostgreSQL 调优、性能分析（pprof/trace）。核心原则：可维护性优先、idiomatic Go 风格、context.Context 贯穿、超时控制+错误封装+日志记录、防 SQL 注入、并发安全。

## 核心原则

### I. 共享领域优先

功能从共享领域层开始设计。领域模型和仓储放 `domains/zeebee/`，共享代码独立可测试且有文档。不允许仅为了组织而抽象。

### II. API/RPC 分离

双层架构：API 层（HTTP）处理验证和路由，RPC 层（gRPC）处理业务逻辑。数据流：`API Handler → API Logic → RPC Call → RPC Logic → Repository → Model → Database`。

### III. 代码生成强制

使用 `Makefile` 封装的命令生成代码，如 `make create-service lender`。

### IV. 数据转换标准

所有类型转换使用 `copierx.Trans[T, U]`（T、U 为非指针类型，src 为 `*T`）。时间字段需 `copierx.TimeToRFC3339Str()` 选项。

### V. 类型安全与测试

强类型贯穿全代码库。`common/` 工具必须有单元测试，领域逻辑变更需集成测试。覆盖率目标 80%+。

## 开发标准

### 技术栈

- **语言**：Go 1.24+ | **框架**：go-zero 1.9.3
- **数据库**：PostgreSQL 14+ (GORM) | **缓存**：Redis 7+ | **存储**：AWS S3
- **协议**：HTTP (REST) + gRPC | **认证**：JWT with RBAC

### 代码质量

- 错误处理用 `common/errorx/`，输入验证用 `common/validator/`，类型转换用 `common/copierx`
- 结构化日志用 go-zero logx，禁止硬编码敏感信息

### 语言规范

- **代码注释**：英文，简洁
- **Git 提交**：英文，conventional commits 格式
- **交流沟通**：中文
- **Markdown 文档**：中文（`README.md` 除外，用英文）

### Git 提交格式

```
<type>(<scope>-<service-name>-backend): <description>

feat: 新功能 | fix: 修复 | refactor: 重构 | docs: 文档
style: 格式 | perf: 性能 | test: 测试

# 示例
feat(applicant-borrower-backend): add loan application endpoint
fix(document-lender-backend): fix document upload validation error
```

### AI 文档策略

- **目录**：`.claude/`（分析文档、索引）、`.specify/`（功能规格）、`specs/`（实施文档）
- **Git 状态**：已在 `.gitignore` 中忽略，不提交到仓库
- **SpecKit 规格目录**：所有规格放 `backend/specs/`，不是项目根目录 `specs/`

## 架构治理

### 服务

| 服务 | HTTP | RPC | 职责 |
|------|------|-----|------|
| borrower | 10000 | 10001 | 借款方操作（实体、房产、贷款、申请） |
| lender | 11000 | 11001 | 贷款方操作（审核、文档、批准） |
| admin | 12000 | 12001 | 管理后台（组织、模板、权限） |
| pdf_generator | — | — | PDF 生成 |

### 共享层

- `domains/zeebee/`：共享模型(27)、仓储(27)、DTO(26)、业务逻辑
- `common/`：20个工具包（copierx、errorx、gormx、validator 等）

详细架构设计、数据流规则、层级职责详见 `.claude/docs/analysis/architecture.md`。

## 治理

### 宪章权威

本宪章优先于所有其他开发实践。所有代码审查必须验证合规性。修订需文档+团队批准+迁移计划。

### 合规验证

- PR 通过 golangci-lint 检查
- 新功能需单元测试，领域变更需集成测试
- 所有类型转换使用 `copierx.Trans[T, U]`

## Behavioral Guidelines

Think before coding. Simplicity first. Surgical changes. Goal-driven execution.

### Think Before Coding

- 显式声明假设，不确定就问
- 存在多种理解时，列出选项而非静默选择
- 有更简方案就说出来，该反驳就反驳
- 不清楚就停，指出困惑点

### Simplicity First

只写解决问题的最少代码，不做推测性设计。不添加未被要求的功能、抽象、灵活性或不可能场景的错误处理。200 行能减到 50 行就重写。

### Surgical Changes

只改必须改的。不改相邻代码/注释/格式，不重构没坏的东西。匹配现有风格。你的变更导致的孤立代码（unused imports/variables）才删除。

### Goal-Driven Execution

任务转为可验证目标："加验证" → "写失败测试再通过"，"修 bug" → "写复现测试再通过"。多步任务列简要计划，每步有验证点。

**版本**：1.1.0 | **批准日期**：2025-12-29 | **最后修订**：2026-05-12
