# 微信支付

## 使用文档

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