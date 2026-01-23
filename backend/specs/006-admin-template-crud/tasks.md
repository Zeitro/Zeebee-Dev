# Admin 模板 CRUD - 实施任务清单

**功能分支**: `6-admin-template-crud`
**生成日期**: 2026-01-23
**状态**: 待实施

## 任务概述

本文档包含 Admin 模板 CRUD 功能的完整实施任务清单。任务按照 go-zero 框架的正确工作流程编排：
1. **先定义契约**（.api 和 .proto 文件）
2. **通过 Makefile 生成模板代码**
3. **填充业务逻辑**

**用户故事**:
- **US1 (P1)**: 管理基础模板 - 核心功能
- **US2 (P2)**: 级联删除行为 - 数据完整性
- **US3 (P3)**: 模板变量架构管理 - 高级功能

**总任务数**: 70
**并行机会**: 多个任务可并行执行（标记 [P]）
**MVP 范围**: Phase 0-4 (契约 + 生成 + 基础层 + RPC + API)

---

## Phase 0: 契约定义与代码生成

**目标**: 定义 API 和 RPC 契约，使用 goctl 生成代码骨架

### 0.1 API 契约定义

- [ ] T001 复制 specs/006-admin-template-crud/contracts/template.api 到 app/admin/api/desc/base_template.api

### 0.2 RPC 契约定义

- [ ] T002 在 app/admin/rpc/pb/admin.proto 中添加模板相关的 protobuf 定义
  - 定义 CreateBaseTemplateReq/Resp
  - 定义 UpdateBaseTemplateReq/Resp
  - 定义 GetBaseTemplateReq/Resp
  - 定义 BaseTemplatePageReq/Resp
  - 定义 SwitchBaseTemplateStatusReq/Resp
  - 定义 DeleteBaseTemplateReq/Resp
  - 定义 EmptyResp

### 0.3 生成代码骨架

- [ ] T003 运行 `make service admin` 生成 API与Rpc 层代码骨架
  ```bash
  make create-service admin
  ```
- [ ] T005 验证生成的代码结构完整性
  - 检查 app/admin/api/internal/handler/base_template/ 目录
  - 检查 app/admin/api/internal/logic/base_template/ 目录
  - 检查 app/admin/api/internal/types/types.go 更新
  - 检查 app/admin/rpc/internal/logic/template/ 目录
  - 检查 app/admin/rpc/internal/server/ 目录
  - 检查 app/admin/rpc/client/ 目录

---

## Phase 1: 常量准备

**目标**: 更新和添加必要的常量定义

- [ ] T006 更新 domains/zeebee/constants/template.go，将 TemplateCategoryContact 改为 TemplateCategoryContract
- [ ] T007 在 domains/zeebee/constants/ 中添加 UserType 常量（如不存在）
- [ ] T008 在 domains/zeebee/constants/ 中添加 CommonStatus 常量（如不存在）

---

## Phase 2: 领域层实现

**目标**: 实现共享领域层的模型、仓储、业务逻辑和 BO，为 RPC 和 API 层提供基础

### Model 层

- [ ] T009 [P] 创建 domains/zeebee/model/base_template.go，定义 BaseTemplate 结构体
- [ ] T010 [P] 创建 domains/zeebee/model/base_template.go，定义 VariableGroupSchema 和 TemplateVariableGroup 结构体
- [ ] T011 [P] 创建 domains/zeebee/model/base_template.go，定义 TemplateVariable 结构体和 Validate 方法
- [ ] T012 [P] 创建 domains/zeebee/model/organization_template.go，定义 OrganizationTemplate 结构体

### Repo 层

- [ ] T013 [P] 创建 domains/zeebee/repo/base_template.go，定义 BaseTemplateRepo 接口
- [ ] T014 [P] 创建 domains/zeebee/repo/base_template.go，实现 baseTemplateRepo 结构体
- [ ] T015 [P] 创建 domains/zeebee/repo/base_template.go，实现 Create 方法
- [ ] T016 [P] 创建 domains/zeebee/repo/base_template.go，实现 GetByID 方法
- [ ] T017 [P] 创建 domains/zeebee/repo/base_template.go，实现 Update 方法
- [ ] T018 [P] 创建 domains/zeebee/repo/base_template.go，实现 UpdateStatus 方法
- [ ] T019 [P] 创建 domains/zeebee/repo/base_template.go，实现 Delete 方法（硬删除）
- [ ] T020 [P] 创建 domains/zeebee/repo/base_template.go，实现 Page 方法
- [ ] T021 [P] 创建 domains/zeebee/repo/base_template.go，实现 ExistsByTag 方法
- [ ] T022 [P] 创建 domains/zeebee/repo/base_template.go，实现 Transaction 方法
- [ ] T023 [P] 创建 domains/zeebee/repo/organization_template.go，定义 OrganizationTemplateRepo 接口
- [ ] T024 [P] 创建 domains/zeebee/repo/organization_template.go，实现 organizationTemplateRepo 结构体
- [ ] T025 [P] 创建 domains/zeebee/repo/organization_template.go，实现 CascadeStatus 方法
- [ ] T026 [P] 创建 domains/zeebee/repo/organization_template.go，实现 CascadeDelete 方法

### Logic 层

- [ ] T027 [P] 创建 domains/zeebee/logic/template.go，定义 TemplateService 接口
- [ ] T028 [P] 创建 domains/zeebee/logic/template.go，实现 templateService 结构体
- [ ] T029 [P] 创建 domains/zeebee/logic/template.go，实现 CascadeStatus 方法
- [ ] T030 [P] 创建 domains/zeebee/logic/template.go，实现 CascadeHardDelete 方法
- [ ] T031 [P] 创建 domains/zeebee/logic/template.go，实现 ValidateVariableGroupSchema 函数

### BO 层

- [ ] T032 [P] 创建 domains/zeebee/data/bo/base_template.go，定义 BaseTemplateBo 结构体
- [ ] T033 [P] 创建 domains/zeebee/data/bo/base_template.go，定义 BaseTemplatePageBo 结构体

---

## Phase 3: RPC 业务逻辑填充

**目标**: 在生成的 RPC logic 文件中填充业务逻辑

### RPC Logic 实现

- [ ] T034 [P] 填充 app/admin/rpc/internal/logic/template/create_logic.go，实现创建逻辑
  - 验证标签唯一性
  - 调用 Repo.Create
  - 数据转换（Model → Proto）
- [ ] T035 [P] 填充 app/admin/rpc/internal/logic/template/update_logic.go，实现更新逻辑
  - 验证标签不可变性（FR-016）
  - 调用 Repo.Update
  - 数据转换
- [ ] T036 [P] 填充 app/admin/rpc/internal/logic/template/detail_logic.go，实现详情查询逻辑
  - 调用 Repo.GetByID
  - 数据转换
- [ ] T037 [P] 填充 app/admin/rpc/internal/logic/template/page_logic.go，实现分页查询逻辑
  - 构建查询条件
  - 调用 Repo.Page
  - 数据转换（切片）
- [ ] T038 [P] 填充 app/admin/rpc/internal/logic/template/switch_status_logic.go，实现状态切换逻辑
  - 在事务中执行更新和级联
  - 调用 Repo.UpdateStatus
  - 调用 TemplateService.CascadeStatus
- [ ] T039 [P] 填充 app/admin/rpc/internal/logic/template/delete_logic.go，实现删除逻辑
  - 在事务中执行级联删除
  - 调用 TemplateService.CascadeHardDelete
  - 调用 Repo.Delete

### RPC Server 注册

- [ ] T040 在 app/admin/rpc/internal/server/template_server.go 中注册 RPC 服务（如需要手动调整）
- [ ] T041 [P] 创建 app/admin/rpc/client/template/template.go，实现 RPC 客户端（如 goctl 未生成）

---

## Phase 4: API 业务逻辑填充

**目标**: 在生成的 API logic 文件中填充业务逻辑

### API Logic 实现

- [ ] T042 [P] 填充 app/admin/api/internal/logic/base_template/create_logic.go
  - 数据验证（使用 validator）
  - 调用 TemplateRPC.Create
  - 数据转换（Proto → API Entity）
- [ ] T043 [P] 填充 app/admin/api/internal/logic/base_template/update_logic.go
  - 数据验证
  - 调用 TemplateRPC.Update
  - 数据转换
- [ ] T044 [P] 填充 app/admin/api/internal/logic/base_template/detail_logic.go
  - 调用 TemplateRPC.Detail
  - 数据转换
- [ ] T045 [P] 填充 app/admin/api/internal/logic/base_template/page_logic.go
  - 调用 TemplateRPC.Page
  - 数据转换（切片）
- [ ] T046 [P] 填充 app/admin/api/internal/logic/base_template/switch_status_logic.go
  - 调用 TemplateRPC.SwitchStatus
- [ ] T047 [P] 填充 app/admin/api/internal/logic/base_template/delete_logic.go
  - 调用 TemplateRPC.Delete

### 服务依赖注册

- [ ] T048 在 app/admin/api/internal/svc/service_context.go 中注册 TemplateRPC 客户端
  - 添加 TemplateRPC 字段
  - 在 InitializeServiceContext 中初始化

### API Handler 验证

- [ ] T049 验证所有生成的 handler 文件正确引用 logic
- [ ] T050 验证 routes.go 包含新增的路由定义

---

## Phase 5: 用户故事 2 - 级联删除行为 [US2]

**目标**: 确保删除基础模板时正确级联删除组织模板

**独立测试标准**: 删除有组织模板的基础模板 → 组织模板也被删除

**注意**: 级联删除已在 Phase 3 的 delete_logic.go 中实现，此阶段主要是验证和测试

- [ ] T051 编写级联删除的单元测试
- [ ] T052 编写级联删除的集成测试

---

## Phase 6: 用户故事 3 - 模板变量架构管理 [US3]

**目标**: 支持复杂的 JSONB 变量架构验证

**独立测试标准**: 创建包含嵌套对象和数组的变量架构 → 正确保存和检索

**注意**: JSONB 验证已在 Phase 2 的 TemplateVariable.Validate 方法中实现，此阶段主要是验证和测试

- [ ] T053 编写变量架构验证的单元测试
- [ ] T054 编写复杂嵌套结构的集成测试

---

## Phase 7: 收尾与验证

**目标**: 验证所有功能正常工作，代码质量达标

### 代码质量

- [ ] T055 运行 golangci-lint 检查代码质量
- [ ] T056 修复所有 lint 问题
- [ ] T057 确保所有 exported 函数有 godoc 注释
- [ ] T058 确保错误处理使用 common/errorx 包

### 功能验证

- [ ] T059 手动测试创建基础模板功能
- [ ] T060 手动测试分页查询功能
- [ ] T061 手动测试状态切换功能
- [ ] T062 手动测试删除功能
- [ ] T063 验证级联状态切换是否正常工作
- [ ] T064 验证级联删除是否正常工作
- [ ] T065 验证 JSONB 变量架构验证是否正常工作

### 文档更新

- [ ] T066 更新 .claude/docs/INDEX.md，添加模板功能相关文档链接
- [ ] T067 更新 .claude/docs/analysis/architecture.md，添加模板相关架构说明

---

## 依赖关系图

```
Phase 0 (契约 + 代码生成)
    ↓
Phase 1 (常量) - T006-T008
    ↓
Phase 2 (领域层) - T009-T033
    ↓
Phase 3 (RPC Logic) - T034-T041
    ↓
Phase 4 (API Logic) - T042-T050
    ↓
Phase 5 (US2 测试) - T051-T052 ──→ Phase 6 (US3 测试) - T053-T054
                                    ↓
                            Phase 7 (收尾) - T055-T067
```

---

## 并行执行机会

### Phase 2 - 领域层并行组

```bash
# 并行组 1: Model 层
T009, T010, T011, T012

# 并行组 2: Repo 层
T013-T022 (base_template.go)
T023-T026 (organization_template.go)

# 并行组 3: Logic/BO 层
T027-T033 (template.go, base_template.go)
```

### Phase 3 - RPC Logic 并行组

```bash
# RPC Logic 层可并行
T034, T035, T036, T037, T038, T039
```

### Phase 4 - API Logic 并行组

```bash
# API Logic 层可并行
T042, T043, T044, T045, T046, T047
```

---

## 实施策略建议

### MVP (最小可行产品) - Phase 0-4

**范围**: Phase 0 + Phase 1 + Phase 2 + Phase 3 + Phase 4
**功能**: 基础模板的完整 CRUD 功能和级联状态切换
**时间**: 核心功能优先

### 增量交付顺序

1. **Sprint 1**: Phase 0-1 (契约生成 + 常量) - 代码骨架就绪
2. **Sprint 2**: Phase 2 (领域层) - 基础层就绪
3. **Sprint 3**: Phase 3 (RPC Logic) - 业务逻辑实现
4. **Sprint 4**: Phase 4 (API Logic) - HTTP 接口实现 → MVP 可交付
5. **Sprint 5**: Phase 5-6 (US2-US3) - 完整功能验证
6. **Sprint 6**: Phase 7 (收尾) - 生产就绪

---

## 注意事项

1. **go-zero 工作流**: 先定义契约 → goctl 生成代码 → 填充逻辑
2. **标签不可变性**: 更新模板时必须验证标签未更改 (FR-016)
3. **应用层级联**: 所有级联操作在应用层代码实现，不使用数据库触发器
4. **事务一致性**: 级联操作必须在数据库事务中执行
5. **代码风格**: 遵循 admin 端现有的代码风格，使用 BasePageRequest/BasePageResponse
6. **数据转换**: 所有数据类型转换使用 `copierx.Trans[T, U]`

---

## Makefile 命令参考

```bash
# API 代码生成
make create-api-service

# RPC 代码生成
make create-rpc-service

# 同时生成 API 和 RPC
make create-service
```

---

**下一步**: 开始执行 Phase 0 任务 - 定义契约并生成代码骨架
