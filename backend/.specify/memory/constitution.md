# 代码风格宪章 (Code Style Constitution)

> 项目无关的代码风格与工程行为规范。适用于任何项目，可原样迁移。
> 项目特定规范（技术栈、目录、命名、scope 格式）见各项目 `CLAUDE.md`。

## 语言规范

- **代码注释**：英文，简洁，不使用中文
- **Git 提交**：英文，conventional commits 格式
- **交流沟通**：中文
- **Markdown 文档**：中文（`README.md` 除外，用英文）
- **LLM 中间文档**（plan、analysis、design）：中文

## Git 提交

### 格式

```
<type>(<scope>): <description>

feat | fix | refactor | docs | style | perf | test | chore
```

### 约束

- 提交前确认在项目根目录，不在则先 `cd` 到根目录
- **禁止将 LLM（Claude、Copilot 等）列为 committer 或 co-author**
- **多功能 → 多 commit**：按功能分批提交。同一类功能（纯文档整理、纯重构）合并为一个 commit，body 逐项列出文件变化，不拆得过细
- 不同 type 必须分开：一个 `feat` 代码改动 + 一个 `docs` 文档整理 = 2 个 commit

## 整洁代码

> 开发团队的代码质量是团队的代码质量，所以要写整洁的代码。
>
> 简单 > 抽象，组合 > 继承，可读 > 优雅，显式 > 魔法

### 命名: 名字就是最好的注释
> 好的命名就是最好的文档。days 比 d 好一万倍。如果发现一个变量名或函数名需要额外注释来解释，那通常是名字取得不对。在 Go 的哲学里，省下敲键盘的时间，远不如省下读代码的人理解的时间重要。
>
> 命名要短但准确
**坏味道**
```go
func get(d []string) map[string]string { // 返回用户ID到名字的映射
    m := make(map[string]string)
    for _, id := range d {
        // 从数据库查名字...
        m[id] = name
    }
    return m
}

func CreateUser(u User){
	
}
// tmp1 tmp2 aaa bbb 这样的变量/参数缩写，命名不好，代码可读性差
```
**好味道**
```go
// 不需要注释，名字本身说明一切
func getUserNameMap(userIDs []string) map[string]string {
nameMap := make(map[string]string)
for _, id := range userIDs {
// 从数据库查名字...
nameMap[id] = name
}
return nameMap
}

func CreateUser(user User) error {
	
}

// cfg ctx tx db req resp 这些是受社区接受的变量/参数缩写
```

### 保持一致
> 函数，参数，变量，常量，包等命名应该保持一致，一致性比最优命名更重要。
> 不要一会一用大写字母，一会用小写字母，一会用下划线
> 整个项目用CreateUser DeleteUser UpdateUser，就不要突然用 InsertUser

### 函数：函数应该有单一职责
> 一个函数应该只做一件事情。
>
> 小而精的函数是代码可读性和可测试性的基石。

**坏味道**
```go
func processOrder(order *Order) error {
    if !order.IsValid() {
        return errors.New("invalid order")
    }
    // 计算总额
    total := 0.0
    for _, item := range order.Items {
        total += item.Price * float64(item.Quantity)
    }
    order.Total = total

    // 保存到数据库
    if err := db.Save(order); err != nil {
        return err
    }

    // 发送邮件
    return email.SendConfirmation(order)
}
```
**好味道**
```go
func processOrder(order *Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    if err := calculateTotal(order); err != nil {
        return err
    }
    if err := saveOrder(order); err != nil {
        return err
    }
    return sendConfirmation(order)
}

func calculateTotal(order *Order) error {
    total := 0.0
    for _, item := range order.Items {
        total += item.Price * float64(item.Quantity)
    }
    order.Total = total
    return nil
}
```
### 提前返回
> 提前返回可以提高代码可读性，并减少嵌套。
> 主线逻辑是清晰的、平铺的。每个 if 都是一个"卫语句"（Guard Clause），处理完特殊情况就立即返回，让你专注于核心的、"快乐路径"（Happy Path）的逻辑。代码的缩进层次代表了逻辑的嵌套深度，缩进越少，心智负担越小。

**坏味道**
```go
func getDiscount(user *User) string {
    if user != nil {
        if user.IsPremium() {
            if user.Age() > 60 {
                return "20%"
            } else {
                return "10%"
            }
        } else {
            return "0%"
        }
    }
    return "0%"
}
```

**好味道**
```go
func getDiscount(user *User) string {
    if user == nil {
        return "0%"
    }
    if !user.IsPremium() {
        return "0%"
    }
    if user.Age() > 60 {
        return "20%"
    }
    return "10%"
}
```

### 函数不要太长
> 20 行以内最好
>
> 40 行还能接受
>
> 超过 80 行就应该考虑拆

**坏味道**
```go
func HandleOrder() {
    ...
    200 行
}
```
**好味道**
```go
func HandleOrder() {
    order := loadOrder()

    check(order)

    pay(order)

    notify(order)
}
```

### 注释：解释为什么，而不是是什么
> 注释应该解释为什么这样做，而不是是什么这样做。
>
> 注释用来解释"为什么这样设计"（比如这个特殊的业务规则为什么存在，或者这个 workaround 是为了绕过哪个第三方库的 bug），而不是解释"这段代码做了什么"。 公共 API（如包、函数、类型的声明）当然需要注释来生成文档，但那是另一回事。对于实现细节，优先考虑用代码本身来表达意图。


**坏味道**
```go
// 检查用户是否有资格获得高级折扣
if user.Points > 1000 && user.Age > 18 && !user.IsBlocked {
    // 应用 15% 折扣
    price = price * 0.85
}
```

**好味道**
```go
if user.IsEligibleForPremiumDiscount() {
    price = price.ApplyDiscount(0.85)
}
```

## 函数约束

- 签名：≤4 参数，≤3 返回值
- 单函数长度：20 行最佳，40 行可接受，80 行必拆

## 错误处理原则

- 禁止裸 `errors.New` / 裸返回底层错误，统一封装为带上下文的错误类型
- 错误需携带堆栈与来源信息，便于排查
- 错误发生处用结构化日志记录，包含足够定位信息
- 对外错误消息对调用方友好，不泄露实现细节与敏感信息

## 测试原则

- **TDD**：先写失败测试，再写实现使其通过
- 工具/公共函数须有单元测试
- 领域逻辑变更须集成测试
- 覆盖率目标 80%+

## 工程行为准则

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

## 治理

本宪章为代码风格的权威来源，优先于个人偏好。修订需文档说明。

---

**版本**：2.0.0 | **最后修订**：2026-07-13
