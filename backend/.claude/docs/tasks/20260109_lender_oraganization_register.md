# lender端组织修改需求开发
## 需求描述
在后端实现lender的组织修改功能开发,包含api与rpc服务
请勿从主分支开发，从当前分支`feat-backend-lender_organization_register`切出新分支，具体可以参考user 注册流程
## 需求
### 需求1：获取当前用户所在组织信息
1. 获取当前用户所在组织信息，用户信息通过ctx来获取
### 需求2：更新组织信息
1. 更新当前所在组织的信息，只有组织owner才能更新组织信息
### 需求3：创建新组织
1. 创建新组织，组织创建成功后，当前用户为新组织owner，组织的状态为 `Inactive`,`owner_type`为`ThirdParty`，具体可以参考user 注册流程