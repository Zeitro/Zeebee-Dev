# 后端开发admin中的组织管理CRUD

## 需求描述

在后端实现admin中的组织管理CRUD
请勿从主分支开发，从当前分支`feat-backend-lender-register`切出新分支，代码风格参考以实现的功能代码
`@app/admin/api/internal/logic/organization/organization_page_logic.go` ,
`@app/admin/rpc/internal/logic/organization/page_logic.go`
现在已经生成了模板，但是未填充逻辑，路径为`@app/admin/api/internal/logic/organization/`，
`@app/admin/rpc/internal/logic/organization/`

## 需求1： 实现创建lender组织

- admin创建的lender组织不需要激活，organization，user，organization_user的状态都是Active状态
- 创建逻辑参考`@app/domain/zeebee/logic/user.go`中注册逻辑`RegisterWithOrganization`

## 需求2：编辑lender组织

- 只允许编辑数据表organizations中的字段，允许的字段已经通过参数传入，不需要再新增

## 需求3：获取组织详情

## 需求4：删除组织

- 删除组织前需要删除组织下的用户，角色以及中间表的关联关系

## 需求5：切换状态

- 切换组织状态，如果切换为`Active`状态，则需要将organization_user中的`Disabled`状态改为`Active`
- 如果切换状态为`Disabled`,则需要就对应的organization_user中的`Active`状态改为`Disabled`