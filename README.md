# go-wechat

golang微信sdk，支持微信公众号，微信小程序，微信支付

`go get github.com/zhyc9de/go-wechat`

## 功能模块

### 公众号
- 模板消息
    - [x] 发送
- 消息管理
    - [x] 接受普通消息
    - [x] 接受事件推送
- 客服消息
    - [x] 发送文本
    - [x] 发送图片
    - [x] 发送小程序卡片
- 网页开发
    - [x] OAuth 网页授权
    - [x] JS-SDK
- 自定义菜单
    - [x] 创建自定义菜单
    - [x] 新增个性化菜单
    - [x] 测试用户是否匹配
    - [x] 删除个性化菜单
- 素材管理
    - [x] 临时素材上传
    - [x] 获取临时素材
- 用户管理
    - [x] 获取用户信息
    - [x] 获取标签
    - [x] 新增标签
    - [x] 批量标记标签
    - [x] 批量取消标记标签
- 二维码
    - [x] 临时二维码
    - [x] 永久二维码
- 数据分析

### 小程序
- 模板消息
    - [x] 发送
- 消息管理
    - [x] 接受普通消息
    - [x] 接受事件推送
- 客服消息
    - [x] 发送文本
    - [x] 发送图片
- 二维码
    - [x] 临时小程序码
    - [x] 永久小程序码
- 素材管理
    - [x] 临时素材上传
    - [x] 获取临时素材
- 解码用户信息
    - [x] 用户信息
    - [x] 分享信息
    - [ ] 手机号
- 内容安全
    - [x] 图片
    - [x] 文本
- 动态消息
    - [x] 新建
    - [x] 更新
- 数据分析

### 商户平台
- [ ] 统一下单
    - [x] 二维码支付
    - [x] 小程序支付
    - [x] 公众号网页支付
- [x] 支付回调
- [x] 退款
- [x] 退款回调
- [x] 付款到零钱

## 文档

### pakcage name

- 核心模块 **weUtil**
    - 公众号[文档](docs/mp.md) **mp**
    - 小程序[文档](docs/mxa.md) **wxapp**
    - 微信支付[文档](docs/mch.md) **mch**

### 公众号和小程序初始化

公众号和小程序有部分一样的请求，比如获取access token，主动发送客服消息，上传临时媒体文件等。

所以sdk在专门为小程序和公众号服务的client上，有一个`weUtil.CommonClient`接口，这么设计是为了某些在小程序和公众号有相同处理逻辑的部分，封装成一个函数，传入这个接口即可。`weUtil.Client`实现了`weUtil.CommonClient`接口

```go
// 支持两种方式存储access token

// 1. 使用atomic.Value，直接存储在内存中
tokenMgr := weUtil.NewWxTokenMgr(appId, appSecret)

// 2. 存储在redis中, redisClient使用`github.com/go-redis/redis`
tokenMgr := weUtil.NewRTokenMgr(appId, appSecret, redisClient)

// 新建一个weUtil.CommonClient
// 最后参数用于记录sdk内部的日志，weUtil.Logger接口
weClient := weUtil.NewClient(tokenMgr, log)
// 新建一个公众号专用的client
mpClient := mp.NewClient(weClient)
// 新建一个小程序专用的client
wxaClient := wxa.NewClient(weClient)

// debug时，可以手动设置access token
weClient.Set(weUtil.KeyToken, token, expires)
```

### 客服消息

TODO 接收客服消息 && 被动回复消息 && 主动回复消息

### TODO
- [ ] test
- [ ] benchmark