# zeitro list 接口转发

## 需求描述

这是一个后端 borrower转发接口的请求，
我现在需要转发zeitro提供的接口，并将结果转化为本项目的数据结构。
现在主要是borrower端需要该数据，为了避免重复开发，我将在domains/zeebee/zeitro中对zeitro的接口封装转换。
请勿从主分支开发，从当前分支`feat-backend-lender-register`切出新分支

## 需求

### 需求1：配置获取

在config中配置zeitro相关配置，配置参数如下

```yaml
Zeitro:
  Url: http://localhost:8080
  Username: admin
  Password: <PASSWORD>
```

### 需求2：接口封装

在`domains/zeebee/zeitro`中对接口进行封装，由于token会过期，所以需要添加一个token缓存，当token过期时，重新获取token并缓存。

### 需求3：接口转换

请求方式还是通过 web->api-rpc->zeitro这样的方式进行，但是接口返回的数据结构不同，所以需要将数据结构进行转换，转换成borrower返回给前端的数据结构。
代码风格需要与borrower保持一致。

## 接口文档

> 文档中给出的数据结构不一定准确，需要根据实际情况修改。

```
Zeitro Listing 接口说明
接口：POST  https://lo.zeitro.com/auth/access-token
  - Authorization: Basic Auth
    - Username - “cli_hu4tKiltzCfQfheOa3”
    - Password - “sec_8SyshMSf1ZORJrH6mT38FMtDAItSLAox8ysd”
  - Headers
    - Content-Type - “application/x-www-form-urlencoded”
  - Body： x-www-form-urlencoded
    - “grant_type”  - “client_credentials”
  - 返回数据参考:
    - {
    -    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2Fub2ZmaWNlcmlkIjoiIiwiY3VzdG9tZXJpZCI6InJhdGVkbmEiLCJyb2xlcyI6bnVsbCwiaXNoYXJkbW9uZXlsb2FuY3VzdG9tZXIiOmZhbHNlLCJleHAiOjE3Njc5MzEwMTgsImlhdCI6MTc2NzkyNzQxOH0.BIF86xnxxVFBW1Z1xfwVpVqzk9OdBAc6HmP6jwbqIuc",
    -    "token_type": "Bearer",
    -    "expires_in": %!d(float64=3600)
    - }
接口：POST  https://lo.zeitro.com/los/openhousesearch
  - Authorization: Bear Token
    - 使用/auth/access-token接口返回的“access_token”值
  - Body: raw
    - 请求数据参考：
    - {
    -    "Lat":37.33874,
    -    "Lon":-121.8852525,
    -    "Dist":5,
    -    "Addr":"San Jose, CA, USA"
    - }
  - 返回数据参考：
    - [
    -    {
    -        "FullAddress": "523 Union Ave, Campbell, CA 95008",
    -        "State": "CA",
    -        "County": "Santa Clara",
    -        "City": "Campbell",
    -        "ZipCode": "95008",
    -        "PropertyType": "Townhouse",
    -        "Price": "1699000                             ",
    -        "IncomeIndicator": "4",
    -        "MinorityPercent": "93.91",
    -        "TotalPopulation": "3662",
    -        "Hispanic": "3302",
    -        "Black": "9",
    -        "Asian": "75",
    -        "Pacific": "10",
    -        "MsaCode": "31084",
    -        "DaysOnMarket": "31",
    -        "AgentName": "Jina Lee",
    -        "AgentPhone": "3109855462",
    -        "AgentEmail": "jina@metisre.com",
    -        "Distance": "4"
    -    },
    - ]
           const response: Response = await fetch("/los/openhousesearch", {
               method: "POST",
               headers: {
                   "Content-Type": "application/json",
                   Authorization: "Bearer " + token,
               },
               body: JSON.stringify({
                   Lat: location.lat,
                   Lon: location.lng,
                   Dist: parseFloat(form.getFieldValue("radius")),
                   Addr: location.address,
               })
           });

```