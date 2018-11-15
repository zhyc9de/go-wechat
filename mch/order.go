package mch

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/zhyc9de/go-wechat"
)

// 新建订单
// 只填写基础数据，new出来之后直接修改字段吧
func (client *Client) NewOrder(appId, body, outTradeNo, tradeTyp, openId string, fee int64) Trade {
	return Trade{
		AppId:          CDATA{appId},
		MchId:          CDATA{client.MchId},
		NonceStr:       CDATA{weUtil.GetRndString(32)},
		Body:           CDATA{body},
		OutTradeNo:     CDATA{outTradeNo},
		TotalFee:       fee,
		SpbillCreateIp: CDATA{LocalIP()}, // 默认的普通支付，如果是native
		NotifyUrl:      CDATA{client.NotifyURL},
		TradeType:      CDATA{tradeTyp},
		OpenId:         CDATA{openId},
	}
}

// 统一下单
func (client *Client) UnifiedOrder(trade Trade) (tradeResult TradeResult, err error) {
	trade.Sign = CDATA{client.SignXml(trade)}

	data, err := xml.Marshal(trade)
	if err != nil {
		return
	}

	var u string
	if client.debug {
		u = "https://api.mch.weixin.qq.com/sandbox/pay/unifiedorder"
	} else {
		u = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	}

	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))
	rb, err := client.DoRequest(req)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(rb, &tradeResult); err != nil {
		err = fmt.Errorf("parse xml err=%s, body=%s", err.Error(), string(rb))
		return
	} else if !tradeResult.Valid() {
		err = tradeResult
		return
	}
	// TODO 校验签名
	return
}

// JSAPI/小程序 支付参数
func (client *Client) WCPayRequest(trade TradeResult) WCPay {
	wcPay := WCPay{
		AppId:     trade.AppId.Text,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  weUtil.GetRndString(32),
		Package:   fmt.Sprintf("prepay_id=%s", trade.PrepayId.Text),
		SignType:  "MD5",
	}
	wcPay.PaySign = client.SignWCPay(wcPay)
	return wcPay
}
