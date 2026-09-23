# Zeebee Backend Agent Guide

本文件是 `backend/` 的项目工作入口。根目录 `AGENTS.md`、`CLAUDE.md` 只负责把 Agent 引到这里；项目规则不要在入口文件重复维护。

## 1. 必读文档与优先级

开始任何后端任务前，按以下顺序读取：

1. 本文件：项目流程与技术约束。
2. [`constitution.md`](constitution.md)：代码风格与工程行为宪章，**每个任务都必须读取并遵守**。
3. [`guides/architecture.md`](guides/architecture.md)：代码变更必读，说明分层、数据流、目录和架构决策。
4. [`INDEX.md`](INDEX.md)：只在需要定位模块、组件或服务事实时读取。

若上下文被压缩、任务切换或继续旧任务，重新读取 `constitution.md` 和本文件中与当前任务相关的章节；不要依赖之前的记忆。若规则冲突，用户明确要求优先；未明确覆盖时，`constitution.md` 优先于本文件，本文件优先于索引性文档和现有习惯。

## 2. 开始任务

先把需求写成可验证目标，再开始修改：

1. 确定服务、模块和层次：HTTP API、RPC、domain、repository、model、配置或文档。
2. 判断目标文件是否由 `.api`、`.proto`、模板或脚本生成。
3. 用 codebase-memory 的 `search_graph` 查找定义和相邻实现；需要调用关系时用 `trace_path`，需要源代码时先搜索再用 `get_code_snippet`。字符串、配置、脚本和 Markdown 才使用文本搜索。
4. 先找可复用的类型、转换器、错误处理和实现范例，再决定是否新增文件或抽象。
5. 明确验证方式：受影响包的测试、构建、代码生成检查或文档检查。

只改完成目标所需的文件。不要为了“以后可能需要”增加接口、配置、抽象或兼容层。

## 3. 分层边界

标准数据流：

```text
HTTP Handler → API Logic → RPC Client → RPC Logic → Domain Logic → Repository → Model → Database
```

- API 层负责 HTTP 参数绑定、边界校验、调用 RPC 和响应转换。
- RPC 层负责鉴权、业务编排和跨服务调用；RPC logic 不直接访问数据库。
- 领域层负责可复用业务规则，输入和返回均使用 `data/dto`，不暴露 `model` 数据；repository 负责 GORM 查询和持久化；model 只描述数据库模型。
- HTTP `types`、Proto message、domain DTO、BO 和 model 不能跨层混用。
- 数据库访问必须经过 `baseLogic → domains/zeebee/logic → repo`，不要在 RPC logic 中绕过领域层。
- domain logic 不直接操作数据库；可通过 repo 调用其嵌入的 `BaseRepo` 数据操作方法，但不获取 `GetDB`、`Client` 等 ORM 连接来拼接 SQL。
- repo 优先复用 `BaseRepo` 已封装的方法；仅在基础方法无法表达业务查询时新增专用仓储查询。

## 4. API、Proto 与代码生成

生成文件的定义源只有 schema：HTTP 使用 `.api`，gRPC 使用 `.proto`。先改定义，再生成，再填写 logic；生成器输出不能手工仿写或补齐。

### HTTP

- 模块化 API 放在 `app/<service>/api/desc/<module>.api`，由 `app/<service>/api/<service>.api` 导入。
- 请求、响应和 entity 属于 `.api`；不要手写生成的 `internal/types`、handler、routes 或 Swagger。
- 修改 `.api` 后运行 `make create-api-service <service>`。
- 需要同时生成 HTTP 服务与 Swagger 时运行 `make create-server SERVICE=<service>`；该命令不生成 RPC。

### RPC

- 模块方法和不可复用 message 放在 `app/<service>/rpc/pb/*.proto`。
- 可复用的公共 entity 放在 `domains/zeebee/proto/entity.proto`；服务 proto 只保留不可复用的 message 和 service 定义。
- Proto 响应按项目约定统一包裹 `data`；列表是否直接为 repeated field 以现有模块约定为准，不自行发明另一层 `lists`。
- 修改 `.proto` 后运行 `make create-rpc-service <service>`。
- 不要手写或修改生成的 pb、grpc client、server 文件；生成的 `rpc/internal/logic/*_logic.go` 是实现模板，只在生成后填充业务逻辑。

### 两层同时变更

需要 API、RPC 和 Swagger 全部更新时运行：

```bash
make create-service <service>
```

没有对应 schema 变更时不要为了完成清单而运行生成命令。`create-api-service`、`create-rpc-service`、`create-service` 和 `create-server` 的既有行为必须保持不变。

### 分页契约

- API 分页请求固定保留 `page`、`limit`、`filter`、`sorter` 四个层次字段；领域筛选新增模块专用 `*PageFilter`，不要把筛选条件平铺到请求根部。
- Proto 的 `*PageRequestData` 镜像该层次：`page=1`、`limit=2`、模块 `filter=3`、`sorter=4`。排序字段与 `proto.CommonPageSorter` 一致时直接复用它，不重复定义本模块 sorter message。
- Domain DTO 嵌入 `BasePageRequestDto`，并保留模块 `Filter` 和 `BasePageSorter`；API → Proto → DTO 的字段名和嵌套结构保持一致，才能直接使用 `copierx.Trans`。
- Repository 接收完整分页 DTO，按白名单将 `Filter` 和 `Sorter` 组装为 GORM 查询；排序字段不能直接作为客户端 SQL 输入。若查询只支持固定排序，契约中不要暴露 `sorter`。
- 分页响应继续按现有模块约定返回 `total`、`page`、`limit`、`total_page`、`list`，不要为列表额外包装一层。

## 5. 类型转换与错误处理

- 跨层结构体和切片转换统一使用 `common/copierx` 的 `copierx.Trans[T, U]`；`T`、`U` 使用非指针类型，源值按函数签名传入。
- 切片（包括嵌套 `Stages`、`Transitions` 等）优先一次性转换，例如 `copierx.Trans[[]*model.T, []*proto.T](&list, ...)`；不要逐项循环或为单一模块保留手写转换 helper。
- `copierx.Trans`/`Assign` 的非 nil 切片或 map 是完整替换：非 nil 空集合会清空目标，nil 表示未提供并保留目标。涉及集合的修改应有转换回归测试。
- 时间字段转字符串时使用 `copierx.TimeToRFC3339Str()`。
- 只有在 copierx 无法表达特殊业务转换时才手工处理，并把原因留在代码结构或短注释中。
- 输入在 HTTP/RPC 边界校验；错误使用 `common/errorx`，保留堆栈和上下文，不向外暴露敏感实现细节。
- 日志使用 go-zero `logx`，不打印 token、密钥或个人敏感信息。

## 6. 配置与服务入口

检查配置前先读 `app/<service>/cmd/main.go` 和实际启动命令，确认 `-f` 指向的运行时 YAML；配置结构看 `app/<service>/config/config.go`，运行值看入口实际加载的 `app/<service>/etc/config-<env>.yaml`。不要仅凭目录名猜配置来源，也不要把生成目录下的 `etc/` 当成运行时配置。

## 7. 完成任务

按变更风险选择最小验证集：

- Go 代码：`gofmt`、受影响包的 `go test`；跨包、生成或接口变更再补 `go build`。
- API/Proto：确认 schema 已生成，检查生成文件和 logic 是否匹配，必要时运行对应服务构建。
- 领域逻辑、权限、持久化或跨服务行为：补充回归/集成测试，并运行相关测试。
- 纯文档：检查链接、代码块和前后规则一致性，不强制运行代码生成。
- 所有任务：检查 `git diff --check`，确认没有越界文件、临时文件或手工生成代码。

完成标准不是“文件已修改”，而是目标、生成顺序、测试和差异检查都已满足；无法验证的项要明确说明原因。

## 8. Git

提交前确认当前目录是 `backend/` 仓库根目录。提交遵守 `constitution.md` 的 Conventional Commits 和多 commit 规则；本项目 scope 使用：

```text
<type>(<scope>-<service-name>-backend): <description>
```

示例：

```text
feat(applicant-borrower-backend): add loan application endpoint
fix(document-lender-backend): fix document upload validation error
```

不要把 Claude、Codex、Copilot 或其他 LLM 写成 committer/co-author。

## 9. 快速定位

| 需求 | 位置 |
|------|------|
| HTTP 定义 | `app/<service>/api/desc/`、`app/<service>/api/<service>.api` |
| API logic | `app/<service>/api/internal/logic/` |
| RPC 定义 | `app/<service>/rpc/pb/`、`domains/zeebee/proto/` |
| RPC logic | `app/<service>/rpc/internal/logic/` |
| 领域逻辑 | `domains/zeebee/logic/` |
| DTO | `domains/zeebee/data/dto/` |
| Repository | `domains/zeebee/repo/` |
| Model | `domains/zeebee/model/` |
| 公共转换 | `common/copierx/` |
| 统一错误 | `common/errorx/` |
| 服务入口 | `app/<service>/cmd/main.go` |

项目服务矩阵、目录细节和完整范例见 `guides/architecture.md`；模块事实和组件索引见 `INDEX.md`，不要在本文件复制维护易变清单。
