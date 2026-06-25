# Claude Code 规范文档整合分析报告

> 范围：`backend/` 下与 Claude Code / 架构 / 开发规范相关的全部 Markdown 文档。
> 方法：全量读取 `CLAUDE.md`、`PROJECT_INDEX.md`、`README.md`、`constitution.md`、`architecture.md`、`INDEX.md`、`docs/README.md`，并抽样读取 `guides/` 全部 6 个文档 + 2 个 API 索引 + 导航指南。
> 约束：本次仅分析，不修改任何代码或文档。

---

## 0. 执行摘要

当前后端规范文档存在 **三核心问题**：

1. **结构性冗余**：同一信息（架构图、数据流、服务端口、`common/` 包列表、`domains/zeebee/` 结构、`copierx.Trans`）在 5–7 个文件中重复定义。
2. **权威源缺失**：版本号、模型数、模块数等"事实型数据"在多个文件中**互相矛盾**，无单一权威源。
3. **导航层过载**：`INDEX.md`、`docs/README.md`、`PROJECT_INDEX.md`、`navigation-guide.md` 四份文档均承担"导航"职责，互相覆盖。

**17,000 行** Markdown 文档中，对 LLM 高价值的"规则 + 架构 + 示例"类内容不足 **20%**；80% 为历史快照、规划、重复描述。建议整合后规范文档总量压缩约 **60%**。

---

## 1. 文档清单与现状

### 1.1 规范 / 架构层文档（LLM 高价值）

| 文档 | 行数 | 职责定位 | 当前角色 |
|------|------|----------|----------|
| `CLAUDE.md` | ~120 | 开发标准 + 代码示例 | 主入口（路由 + 示例） |
| `.specify/memory/constitution.md` | 120 | 项目宪章（规则） | 规则权威 |
| `.claude/docs/analysis/architecture.md` | 282 | 架构（结构） | 必读架构 |
| `PROJECT_INDEX.md` | ~145 | Repo 索引（事实快照） | 事实清单 |
| `README.md` | ~166 | 项目概览 + Quick Start | 对外入口 |

### 1.2 导航 / 索引层（互相重叠）

| 文档 | 行数 | 职责 | 问题 |
|------|------|------|------|
| `.claude/docs/INDEX.md` | 194 | 知识库导航 | 与 README、PROJECT_INDEX、navigation-guide 重叠 |
| `.claude/docs/README.md` | 58 | 知识库目录说明 | 与 INDEX.md 导航部分重叠 |
| `.claude/docs/guides/navigation-guide.md` | 126 | 文档导航 + 速查 | 与 INDEX.md + PROJECT_INDEX 重叠 |
| `.claude/docs/guides/readme_index.md` | 547 | 各模块 README 索引 | 低价值，多数模块无 README |

### 1.3 guides 层（重复严重）

| 文档 | 行数 | 与 architecture.md 重叠度 |
|------|------|---------------------------|
| `technical_architecture_overview.md` | **1057** | **高（≈85%）** — 系统图/分层/数据流/技术栈全覆盖 |
| `project_structure_guide.md` | 190 | 中（≈60%）— 目录结构 + 依赖方向 |
| `domain_architecture_guide.md` | 271 | 中（≈40%）— 共享层 + BaseRepo |
| `api_documentation_index.md` | 138 | 低（专项，保留） |
| `api_documentation_index_2025.md` | 371 | **冗余**（旧版，已被上一文件取代） |

### 1.4 历史 / 报告层（LLM 低价值，存档）

`.claude/docs/reports/`（10+ 文件，含 2025-11-07、2025-12-22、2026-02-26、2026-05-12 多个时间戳的架构/安全/性能/代码质量分析）、`.claude/docs/planning/`、`.claude/docs/tasks/`、`.claude/docs/analysis/analysis.md`、`architecture_analysis_20250124.md`。

> 结论：这些是**时间点快照**，对 LLM 日常工作无持续参考价值，应归档而非与规范并列。

---

## 2. 内容重复分析

### 2.1 高频重复内容清单

| # | 重复内容 | 出现文件（来源） | 重复次数 | 建议单一权威源 |
|---|----------|------------------|----------|----------------|
| 1 | **服务端口表**（borrower 10000/lender 11000/admin 12000/pdf 13000） | constitution.md、architecture.md、project_structure_guide.md、README.md、PROJECT_INDEX.md、api_documentation_index.md、technical_architecture_overview.md | **7** | `architecture.md` §服务矩阵 |
| 2 | **数据流链**（`Handler→API Logic→RPC→Repo→Model→DB`） | constitution.md §II、architecture.md、project_structure_guide.md、README.md、technical_architecture_overview.md | **5** | `architecture.md` §数据流 |
| 3 | **架构图**（ASCII/Mermaid） | architecture.md、technical_architecture_overview.md、analysis.md、project_structure_guide.md、README.md | **5** | `architecture.md` §架构概览（仅保留 1 张） |
| 4 | **`common/` 包列表**（22 包） | architecture.md、INDEX.md、project_structure_guide.md、PROJECT_INDEX.md、README.md、domain_architecture_guide.md | **6** | `architecture.md` §代码组织 |
| 5 | **`domains/zeebee/` 目录结构** | architecture.md、INDEX.md、project_structure_guide.md、domain_architecture_guide.md、PROJECT_INDEX.md、README.md | **6** | `architecture.md` §代码组织 |
| 6 | **`copierx.Trans` 转换规范** | CLAUDE.md、constitution.md §IV、architecture.md、project_structure_guide.md、domain_architecture_guide.md | **5** | `CLAUDE.md` §数据转换（API + 示例） |
| 7 | **依赖方向规则**（API→RPC→Domain→Infra） | architecture.md §关键约束、project_structure_guide.md §6、constitution.md §II | **3** | `architecture.md` §关键约束 |
| 8 | **BaseRepo 泛型仓储模式** | architecture.md §Repository 模式、domain_architecture_guide.md | **2** | `architecture.md` |
| 9 | **快速定位表 / 速查表** | project_structure_guide.md §7、navigation-guide.md §搜索技巧 | **2** | 合并入 `INDEX.md` |

### 2.2 完全冗余文件

| 文件 | 冗余原因 | 处置 |
|------|----------|------|
| `guides/api_documentation_index_2025.md` | `api_documentation_index.md` 已是其更新版 | **删除**（或移入 `reports/archive/`） |
| `guides/navigation-guide.md` | 与 `INDEX.md` 职责完全重合 | **合并入 `INDEX.md` 后删除** |
| `docs/README.md` | `INDEX.md` 已是知识库导航 | **合并入 `INDEX.md` 后删除** |
| `guides/readme_index.md` | 索引大量不存在的模块 README | **删除** |

---

## 3. 架构一致性检查

### 3.1 事实型数据矛盾（最高优先级）

| 数据项 | constitution.md | README.md | PROJECT_INDEX.md | INDEX.md | architecture.md | 实际 |
|--------|-----------------|-----------|------------------|----------|-----------------|------|
| **Go 版本** | 1.24+ | **1.26** | **1.23** | — | — | 需统一（建议以 `go.mod` 为准） |
| **go-zero 版本** | 1.9.3 | 1.10.1 | 1.10.1 | — | — | go.mod: 1.10.1 |
| **model 数** | — | 36 | 38 | **38** | — | 需核对 |
| **repo 数** | — | 36 | 38 | — | — | 需核对 |
| **DTO 数** | — | 33 | — | **26** | — | 需核对 |
| **common 包数** | 21 | 21 | 22 | 22 | 21（标题）/ 列出 22 | **22**（实际枚举） |
| **Admin API 模块数** | — | 13 | — | **12** | — | 10（structure guide）/ 12（INDEX） |
| **Borrower API 模块数** | — | 19 | — | **18** | — | 17（structure guide）/ 18（INDEX） |

> **影响**：LLM 每次读取不同文档得到不同"事实"，无法形成稳定心智模型。`go.mod` 是唯一可信源，应成为版本数据的**唯一引用对象**。

### 3.2 分层结构无冲突

handler / logic / repo / domain 四层职责定义在多处一致（API 层禁业务逻辑、RPC 层核心业务、Domain 共享、Common 工具）。**无设计风格混用**，但描述详略不一。

### 3.3 命名一致性

- 服务命名统一（borrower/lender/admin/pdf_generator）。✅
- 目录职责清晰，无命名歧义。✅
- 唯一问题：`repo` 层职责在 `constitution.md`（"数据库操作"）与 `CLAUDE.md` §12（"logic 不直接操作数据库，通过 repo"）措辞略有差异，但语义一致。

### 3.4 模块计数不一致根因

`project_structure_guide.md`（手写、较旧）与 `INDEX.md`（自动生成、较新）模块数不同。**根因**：guide 未随代码生成同步更新。建议模块清单**只保留自动生成版**（`INDEX.md`），guide 不手写模块数。

---

## 4. LLM 可读性评估

| 维度 | 评分 | 问题 |
|------|------|------|
| 结构清晰度 | 6/10 | `CLAUDE.md` 路由清晰，但下游文档职责边界模糊 |
| 明确分层 | 5/10 | rules（constitution）/ structure（architecture）/ examples（CLAUDE）三层存在，但 guides 层与 architecture 层大面积重叠 |
| 可执行规则 | 7/10 | constitution 规则明确，CLAUDE 示例可执行 |
| 冗余描述 | 3/10 | **重度冗余**，LLM 需在 7 个文件中比对同一事实 |
| 入口收敛 | 4/10 | 4 个"导航"文档竞争，LLM 不知从何读起 |
| 事实稳定 | 2/10 | 版本/计数矛盾，无权威源 |

**核心痛点**：LLM 每次会话被注入大量重复且矛盾的上下文，消耗 token 同时降低决策置信度。

---

## 5. 优化建议

### 5.1 建议的最终结构（清晰分层）

```
backend/
├── CLAUDE.md                          # 唯一入口：路由 + 代码示例 + LLM 工作规范
├── README.md                          # 对外概览（英文，Quick Start）
├── .specify/memory/
│   └── constitution.md                # 唯一规则源：原则 + 标准约束 + Git 格式
└── .claude/docs/
    ├── INDEX.md                       # 唯一导航：文档地图 + 事实速查（自动生成部分）
    └── analysis/
        └── architecture.md            # 唯一架构源：分层 + 数据流 + 目录 + 约束
                                        #         + 服务矩阵（含端口，权威）
```

**分层原则**：
- **入口层**（`CLAUDE.md`）：路由 + 可执行示例，不含事实型数据
- **规则层**（`constitution.md`）：what must be done
- **结构层**（`architecture.md`）：how code is organized + 权威事实（端口、规模）
- **事实层**（`INDEX.md`）：自动生成的清单（模块、模型、DTO），随代码同步
- **历史层**（`reports/archive/`）：时间点快照，不进 LLM 常规上下文

### 5.2 合并 / 拆分 / 删除方案

#### ✅ 合并

| 源文件 | 目标 | 合并内容 |
|--------|------|----------|
| `guides/navigation-guide.md` | `INDEX.md` | 速查表 + 搜索技巧 + 快速开始路径 |
| `docs/README.md` | `INDEX.md` | 知识库目录说明（作为 INDEX 顶部介绍） |
| `guides/project_structure_guide.md` | `architecture.md` | 依赖方向规则、快速定位表（§6/§7 合入架构约束） |
| `guides/domain_architecture_guide.md` | `architecture.md` | BaseRepo 模式细节（避免与 §Repository 重复，仅补差异） |
| `technical_architecture_overview.md` 精华 | `architecture.md` | ADR（架构决策记录）章节 — **仅此部分有增量价值**，其余删除 |
| `PROJECT_INDEX.md` 事实部分 | `INDEX.md` + `architecture.md` | 端口/规模入 architecture，导航入 INDEX |

#### 🔀 拆分

| 源文件 | 拆分结果 |
|--------|----------|
| `technical_architecture_overview.md`（1057 行） | ADR 章节 → `architecture.md`；安全/性能/部署章节 → `reports/` 归档；删除系统图/分层/数据流（已在 architecture） |

#### 🗑️ 删除 / 归档

| 文件 | 处置 | 理由 |
|------|------|------|
| `guides/api_documentation_index_2025.md` | 删除 | 已有更新版 |
| `guides/readme_index.md` | 删除 | 索引大量不存在的 README |
| `guides/navigation-guide.md`（合并后） | 删除 | 已并入 INDEX |
| `docs/README.md`（合并后） | 删除 | 已并入 INDEX |
| `technical_architecture_overview.md`（拆分后） | 删除 | 精华已并入 architecture |
| `architecture_analysis_20250124.md` | 移入 `reports/archive/` | 历史快照 |
| `reports/` 下 2025 年度报告 | 移入 `reports/archive/` | 历史 |
| `planning/progress_tracker.md`、`project_roadmap.md` | 移入 `reports/archive/` 或 `.specify/` | 非规范文档 |

**预期结果**：`guides/` 目录从 7 文件 → 仅保留 `api_documentation_index.md`（自动生成）。

### 5.3 精简后 CLAUDE.md 结构草案

```markdown
# CLAUDE.md

> 本文件为 LLM 工作入口。规则见 constitution，结构见 architecture，事实见 INDEX。
> 修改任何 .api/.proto 后须 `make create-service <name>` 重生成。

## Constitution & Architecture（指针，不复制内容）
- 规则：`.specify/memory/constitution.md`
- 架构：`.claude/docs/analysis/architecture.md`
- 事实清单：`.claude/docs/INDEX.md`

## 工作目录
Monorepo，工作范围限 `backend/`。Specs 入 `backend/specs/`。

## 数据转换规范（copierx.Trans）—— 唯一示例源
[保留当前表格 + 4 个示例代码块]
[删去 constitution/architecture 中的简版，全部指向此处]

## 错误处理规范 —— 唯一示例源
[保留当前 defer + errorx 示例]

## 数据库操作约束
logic 不直接操作数据库，必须经 `domains/zeebee/repo`。

## 代码质量约束
- 注释：英文 | Git：英文 conventional | 交流：中文 | 文档：中文（README 英文）
- 函数：≤4 参 ≤3 返回值
- Slice 转换：用 `copierx.Trans`

## Git 提交格式
<type>(<scope>-<service-name>-backend): <description>
feat | fix | refactor | docs | style | perf | test
[示例 2 条]
⚠️ 禁止 LLM 作为 committer/co-author

## 文档阅读顺序（首次必读）
1. 本文件
2. constitution.md
3. architecture.md

## Related Docs（表格，保留）
```

**变更要点**：
- 删除与 constitution/architecture 重复的"角色定位""核心原则""服务矩阵""分层职责"——改为**指针引用**。
- 保留唯一增量价值：**代码示例**（copierx / errorx）+ **Git 格式** + **阅读顺序**。
- `CLAUDE.md` 从 ~120 行精简至 ~80 行，且**零事实型数据**（版本/端口全部外置）。

---

## 6. 落地步骤建议

> 以下为推荐执行顺序，**本次未执行任何修改**，待确认后分步实施。

1. **统一事实源**（最高优先级）
   - 以 `go.mod` 为准修正 Go/go-zero 版本（constitution.md、README.md、PROJECT_INDEX.md）
   - 以实际代码统计修正 model/repo/DTO/common 包数（统一写入 `INDEX.md`，其他文档改为引用）
   - 统一模块计数口径（仅 `INDEX.md` 保留，自动生成）

2. **合并导航层**
   - `navigation-guide.md` + `docs/README.md` → `INDEX.md`
   - 删除两个源文件

3. **合并架构层**
   - `project_structure_guide.md` §6/§7 + `domain_architecture_guide.md` BaseRepo → `architecture.md`
   - `technical_architecture_overview.md` ADR → `architecture.md` 新增章节
   - 删除 3 个 guide 源文件

4. **精简 CLAUDE.md**
   - 按草案重构，删除重复段，改为指针引用

5. **归档历史**
   - `reports/archive/` 收纳所有 2025 年度 + 时间戳报告
   - `planning/` 评估是否保留

6. **删除冗余**
   - `api_documentation_index_2025.md`、`readme_index.md`

**预期收益**：规范文档总量从 ~17,000 行降至 ~7,000 行（**-60%**），LLM 上下文无矛盾、无重复，事实型数据单一权威源。

---

**生成时间**：2026-06-25 | **状态**：分析完成，待评审实施
