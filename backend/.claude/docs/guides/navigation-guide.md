# 文档导航指南

> 知识库文档导航中心，帮助快速定位所需信息。
> 最后更新：2026-06-24

---

## 快速开始路径

### 新开发者入门

```
INDEX.md → architecture.md → CLAUDE.md → project_structure_guide.md
```

1. **INDEX.md** — 知识库总索引，了解文档全貌
2. **architecture.md** — 系统架构总览（必读）
3. **CLAUDE.md** — 开发规范、数据转换规则、Git 提交规范
4. **project_structure_guide.md** — 项目目录结构和代码组织

### API 开发

```
api_documentation_index.md → CLAUDE.md（数据转换规范） → domain_architecture_guide.md
```

### 领域建模

```
domain_architecture_guide.md → domains/zeebee/model/ → domains/zeebee/dto/
```

### 安全/性能审查

- 安全: `reports/security/security_deep_analysis_20260226.md`
- 性能: `reports/performance/performance_deep_analysis_20260226.md`

---

## 文档目录

### analysis/ — 系统分析

| 文件 | 说明 |
|------|------|
| `architecture.md` | **必读**。系统架构总览 |
| `analysis.md` | 全量深度审计报告（2026-06-24） |
| `architecture_analysis_20250124.md` | 2025-01 架构分析快照 |
| `project-research/research_project_analysis_20251106.md` | 项目现状调研 |

### guides/ — 开发指导

| 文件 | 说明 |
|------|------|
| `project_structure_guide.md` | 项目目录结构、代码组织 |
| `api_documentation_index.md` | API 接口文档索引 |
| `technical_architecture_overview.md` | 技术架构详解 |
| `readme_index.md` | 各模块 README 索引 |
| `readme_index.md` | 各模块 README 索引 |

### reports/ — 分析报告

| 文件 | 说明 |
|------|------|
| `architecture/architecture_deep_analysis_20260226.md` | 架构深度分析 (8.5/10) |
| `code-quality/multi_domain_analysis_20260226.md` | 多领域综合分析 |
| `code-quality/backend_quality_analysis_20260226.md` | 后端质量分析 |
| `security/security_deep_analysis_20260226.md` | 安全深度分析 (6.5/10) |
| `performance/performance_deep_analysis_20260226.md` | 性能深度分析 (7.5/10) |

### planning/ — 项目规划

| 文件 | 说明 |
|------|------|
| `progress_tracker.md` | 项目进度跟踪 |
| `project_roadmap.md` | 项目路线图 |

---

## 项目关键参数速查

**四大服务**：

| 服务 | HTTP | RPC |
|------|------|-----|
| borrower | 10000 | 10001 |
| lender | 11000 | 11001 |
| admin | 12000 | 12001 |
| pdf_generator | 13000 | — |

**共享数据层**：`domains/zeebee/` — 38 models, 38 repos, 26 DTOs
**通用工具**：`common/` — 22 个工具包
**架构模式**：DDD + 微服务 + 共享数据层（PostgreSQL）
**分层依赖**：API → RPC → Domain → Infrastructure（单向向下）

---

## 搜索技巧

| 场景 | 搜索关键词 |
|------|-----------|
| 架构设计 | architecture, 分层, 微服务, DDD |
| API 开发 | api, endpoint, handler, logic |
| 数据模型 | model, dto, repo, domain |
| 数据转换 | copierx, Trans, DTO, BO |
| 错误处理 | errorx, error |
| 安全 | security, auth, JWT |
| 性能 | performance, 缓存 |

### 按目录定位

- **接口定义**：`app/{service}/api/desc/`
- **业务逻辑**：`app/{service}/rpc/internal/logic/`
- **数据模型**：`domains/zeebee/model/`
- **数据访问**：`domains/zeebee/repo/`
- **配置文件**：`app/{service}/etc/`
- **工具函数**：`common/`

---

## 注意事项

- `architecture.md` 是所有开发者必读文档
- 分析报告按日期后缀区分版本，优先阅读最新
- `domains/zeebee/` 是跨服务共享层，修改需评估全服务影响
- 文档中文编写，代码注释英文
