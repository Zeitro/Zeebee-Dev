# Tenant 部署检查清单

- tenant `10401` 只允许 borrower、lender、admin 的 RPC 工作负载访问；禁止公网与 API 网关入口。
- 三端 `RpcServiceConf.Tenant` 指向 tenant 内网地址，且健康检查通过。
- tenant 的 `JwtConf` 与三端 API 中间件使用同一组签名配置。
- 上线后分别用三端 token 验证：本组织订阅只能读取自身；admin 可跨组织读取。
- tenant 的三个订阅状态任务均已注册，三端不再注册同名任务。
