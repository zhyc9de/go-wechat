package mch

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/zhyc9de/go-wechat"
)

func (client *Client) NewTransfers(appId, openId, tradeNo, desc, ip string, amount int64) Transfers {
	return Transfers{
		MchAppId:       CDATA{appId},
		MchId:          CDATA{client.MchId},
		NonceStr:       CDATA{weUtil.GetRndString(32)},
		PartnerTradeNo: CDATA{tradeNo},
		OpenId:         CDATA{openId},
		CheckName:      CDATA{CheckNameNo}, // 默认不填写姓名
		Amount:         amount,
		Desc:           CDATA{desc},
		SpbillCreateIp: CDATA{ip},
	}
}

// 付款
func (client *Client) Transfers(t Transfers) (res TransfersResult, err error) {
	t.Sign = CDATA{client.SignXml(t)}

	data, err := xml.Marshal(t)
	if err != nil {
		return
	}

	u := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))
	rb, err := client.DoRequest(req)
	if err != nil {
		return
	}
	if err = xml.Unmarshal(rb, &res); err != nil {
		err = fmt.Errorf("parse xml err=%s, body=%s", err.Error(), string(rb))
		return
	} else if !res.Valid() {
		err = res
		return
	}
	// TODO 校验签名
	return
}
