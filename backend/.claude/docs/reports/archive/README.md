# 已归档文档

> **本目录文档为历史快照，非权威规范。**
>
> 当前权威文档：
> - 架构、目录结构、共享层、ADR → `.claude/docs/analysis/architecture.md`
> - 开发规范、代码示例 → `backend/CLAUDE.md`
> - 项目规则 → `.specify/memory/constitution.md`
> - 知识库导航 → `.claude/docs/INDEX.md`

## 归档原因

### guides/（步骤 3 合并）

| 文档 | 原位置 | 归档原因 |
|------|--------|----------|
| `guides/technical_architecture_overview.md` | `guides/` | ≈85% 与 `architecture.md` 重复，精华（ADR）已并入架构文档 |
| `guides/domain_architecture_guide.md` | `guides/` | BaseRepo 模式、扩展流程已并入 `architecture.md` |
| `guides/project_structure_guide.md` | `guides/` | 目录结构、依赖方向已并入 `architecture.md` |

### reports/ + analysis/ 历史快照（步骤 5 归档）

| 文档 | 原位置 | 归档原因 |
|------|--------|----------|
| `architecture/architecture_analysis_report_20251107.md` | `reports/` | 2025-11 快照，已被 2026-02-26 深度分析取代 |
| `code-quality/code_quality_analysis_20251107.md` | `reports/` | 2025-11 快照，已被 2026-02-26 取代 |
| `code-quality/analysis_report_2025-12-22.md` | `reports/` | 2025-12 快照 |
| `performance/performance_analysis_report_20251107.md` | `reports/` | 2025-11 快照，已被 2026-02-26 取代 |
| `comprehensive-system-analysis/*` | `reports/` | 2025-11 全量分析，已被 `comprehensive_analysis_20260512.md` 取代 |
| `documentation_quality_validation_2025-12-22.md` | `reports/` | 2025-12 文档质量校验快照 |
| `analysis/architecture_analysis_20250124.md` | `analysis/` | 2025-01 架构快照 |
| `analysis/project-research/research_project_analysis_20251106.md` | `analysis/` | 2025-11 项目调研 |

> 保留在活跃层的最新报告：`analysis/analysis.md`（2026-06-24 全量审计）、`reports/*_20260226.md` 系列、`reports/comprehensive_analysis_20260512.md`。

## 注意事项

- 本目录文档**仅供历史溯源**，其中的数字、模块清单可能已过时
- 引用架构信息请以 `.claude/docs/analysis/architecture.md` 为准
- 后续全量分析报告（时间点快照）亦归档至 `reports/archive/` 对应子目录
