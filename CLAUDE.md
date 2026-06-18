# CLAUDE.md - Monorepo 工作路径指南

本文件为 Claude Code 在此 monorepo 中工作时提供路径路由。

## 项目结构

```
Zeebee-Dev/
├── backend/      # Go + go-zero 微服务后端
├── frontend/     # 前端应用
└── env/          # 环境配置
```

## 路由规则

**后端任务** → `cd backend` → 读 `backend/CLAUDE.md` → 在 `backend/` 下工作

对于后端任务，请保持根目录，并进入 `backend/` 目录，将`backend/`目录作为工作目录，除非显示要求，请勿进入其他目录。

触发词：后端、backend、API、服务、微服务、Go、数据库、gRPC

**前端任务** → `cd frontend` → 读 `frontend/CLAUDE.md` → 在 `frontend/` 下工作

对于前端任务，请保持根目录，并进入 `frontend/` 目录，将`frontend/`目录作为工作目录，除非显示要求，请勿进入其他目录。

触发词：前端、frontend、UI、页面、组件、React、Vue、样式

**全局任务** → 保持根目录 → 读 `README.md`

触发词：CI/CD、monorepo 配置、跨项目协调

## 子项目指导

- **后端**：`backend/CLAUDE.md` — 开发规范、宪法引用、代码示例
- **前端**：`frontend/CLAUDE.md` — 前端规范（如存在）

## Monorepo Rules

本文件仅做路径路由，具体技术规范见各子项目 CLAUDE.md。
