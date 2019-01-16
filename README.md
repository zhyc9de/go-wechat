# go-wechat

golang微信sdk，支持微信公众号，微信小程序，微信支付

`go get github.com/zhyc9de/go-wechat`

## 文档

### pakcage name
核心模块 **weUtil**
公众号 **mp**
小程序 **wxapp**
微信支付 **mch**

### 公众号和小程序初始化

公众号和小程序有一定的通用请求，比如获取access token，主动发送客服消息，上传临时媒体文件等。

所以sdk在专门为小程序和公众号服务的client上，有一个`weUtil.CommonClient`接口，这么设计是为了某些在小程序和公众号有相同处理逻辑的部分，分装成一个函数，传入这个接口即可。`weUtil.Client`实现了`weUtil.CommonClient`接口

```go
// 支持两种方式存储access token

// 使用atomic.Value，直接存储在内存中
tokenMgr := weUtil.NewWxTokenMgr(appId, appSecret)

// 存储在redis中, redisClient使用`github.com/go-redis/redis`
tokenMgr := weUtil.NewRTokenMgr(appId, appSecret, redisClient)

// 新建一个weUtil.CommonClient
// 最后参数用于记录sdk内部的日志，weUtil.Logger接口
weClient := weUtil.NewClient(tokenMgr, log)
// 新建一个公众号专用的client
mpClient := mp.NewClient(weClient)
// 新建一个小程序专用的client
wxaClient := wxa.NewClient(weClient)

// debug时，可以强行指定access token
weClient.Set(weUtil.KeyToken, token, 0)
```

### 客服消息

#### 接受客服消息

TODO

#### 被动回复消息

TODO

#### 主动回复消息

TODO

### 公众号

TODO

### 小程序

TODO

### 微信支付

新建一个微信支付client
```go
// notifyURL 回调地址
// refundURL 退款地址
mchClient = mch.NewClient(mchId, mchKey, notifyURL, refundURL）

// 如果需要设置证书
mchClient.SetCert(certFile, keyFile)

// 统一下单
var trade mch.Trade // 可以通过mchClient.NewOrder新建
order, _ := mchClient.UnifiedOrder(trade)
// 进行签名，然后就可以传入可以JSAPI或者小程序支付了
wcpay := mchClient.WCPayRequest(order)

// TODO 退款、转账

```

微信支付回调通过自定义`UnmarshalXML`实现解析struct，由于微信支付返回的xml的值写在`cdata`里，所以string的值被定义为`mch.CDATA`

```go
// cdata
CDATA struct {
    Text string `xml:",cdata"`
}
// 获取string类型的值，例如微信支付订单号
var callback mch.TradeCallback
transactionId := callback.TransactionId.Text
```

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

### TODO
- [ ] test
- [ ] benchmark