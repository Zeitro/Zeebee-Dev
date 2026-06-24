# Claude AI 生成的项目文档

本目录包含所有 Claude AI 工具生成的项目文档，为 Zeebee Backend 项目提供完整的知识库支持。

## 📁 目录结构

### 🏗️ analysis/ — 架构分析

- `architecture.md` — **必读**。系统架构总览：分层职责、数据流、代码组织、开发流程
- `analysis.md` — 全量深度审计报告（2026-06-24），90 个已验证发现
- `architecture_analysis_20250124.md` — 2025-01 架构分析快照
- `project-research/` — 项目深度研究分析

### 📚 guides/ — 开发指导

- `project_structure_guide.md` — 项目目录结构和代码组织
- `api_documentation_index.md` — API 接口文档索引
- `domain_architecture_guide.md` — 共享数据层设计
- `technical_architecture_overview.md` — 技术架构详解
- `navigation-guide.md` — 文档导航指南
- `readme_index.md` — 各模块 README 索引

### 📊 reports/ — 分析报告

- `code-quality/` — 代码质量分析报告
- `architecture/` — 架构深度分析
- `security/` — 安全深度分析
- `performance/` — 性能深度分析

### 📋 planning/ — 项目规划

- `progress_tracker.md` — 项目进度跟踪
- `project_roadmap.md` — 项目路线图

### 📍 INDEX.md — 项目文档索引

整个项目知识库的导航中心，提供完整的文档交叉引用和快速导航。

---

## 🔄 文档维护

### 更新原则

- 保持文档与实际代码同步
- 确保交叉引用链接有效
- 维护统一的文档格式和风格
- 分析报告按日期后缀区分版本，优先阅读最新

### 使用指南

- **新成员**：从 `INDEX.md` 开始 → `analysis/architecture.md` → `CLAUDE.md`
- **开发人员**：`CLAUDE.md` + 对应模块的 `guides/`
- **审查人员**：最新审计报告 → `analysis/analysis.md`

---

*最后更新：2026-06-24*
