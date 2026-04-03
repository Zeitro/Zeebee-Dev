### 前端项目简介

---

##### 公共组件/功能/方法

> 存放位置 packages/shared-methods
>
> 引入使用@repo/shared-methods
> - components --存放公共封装小组件
> - componentsUi --存放的是只有admin、lender、borrower端才能使用的功能大组件，方便三端复用，例如：loan details
> - 如果能用Ts封装组件是最好的（不强制要求）

##### 基础开发规范

> 1. 封装组件文件名大写开头，其余都小写

> 2. 使用tailwindcss，如需创建样式文件，可创建.scss文件

> 3. 单个文件代码行数尽量少建议不超过 300 行，可封装子组件、hooks，方便维护。

> 4. 使用Prettier格式化代码，需下载插件。

> 5. 使用Eslint检查代码，需下载插件。

> 6. 代码需使用>=ES6的语法

> 7. 复杂逻辑需添加代码注释

> 8. 公共组件需编写使用注释

> 9. **布尔值变量**：以 `is`/`has`/`can`/`should` 开头，如 `isLoading`、`hasPermission`

> 10. **Ref 对象**：以 `Ref` 结尾


##### Git提交规范
> 1. 提交commit需要填写以下类型
> - feature: 新功能
> - fix: bug 修复
> - docs: 文档更新
> - style: 样式信息
> - refactor: 已有的的方法重构（既不是 bug 也不是新增）
> - chore: 构建过程或辅助工具的变动
> - revert: 撤销，版本回退
> - perf: 性能优化
> - test：测试
> - improvement: 改进
> - build: 打包
> - ci: 持续集成

> 1. 从 master 分支创建新开发分支 `feature-xxx`
> 2. 紧急修复使用 `hotfix-xxx` 分支