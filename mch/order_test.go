package mch_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/zhyc9de/go-wechat"
	"github.com/zhyc9de/go-wechat/mch"
	"github.com/zhyc9de/go-wechat/wxapp"
)

func TestClient_SignTrade(t *testing.T) {
	client := mch.InitClient("test", "test2", "http://1270.0.0.1", "http://1270.0.01")
	appClient := wxapp.NewClient(weUtil.NewClient(weUtil.NewWxTokenMgr("", ""), nil))
	trade := client.NewOrder(appClient.GetAppId(), "body", mch.NewTradeNo("fntt"), mch.TradeTypApp, "", 100)
	sign := client.SignXml(trade)
	fmt.Printf("%#v\n%s", trade, sign)
}

func TestTradeCoupon_MarshalXML(t *testing.T) {
	data := []byte (`<xml><appid><![CDATA[wxd21a26f18f315d68]]></appid>
	<bank_type><![CDATA[CMB_DEBIT]]></bank_type>
	<cash_fee><![CDATA[981]]></cash_fee>
	<coupon_count><![CDATA[1]]></coupon_count>
	<coupon_fee>9</coupon_fee>
	<coupon_fee_0><![CDATA[9]]></coupon_fee_0>
	<coupon_id_0><![CDATA[2000000045944299135]]></coupon_id_0>
	<fee_type><![CDATA[CNY]]></fee_type>
	<is_subscribe><![CDATA[Y]]></is_subscribe>
	<mch_id><![CDATA[1486223512]]></mch_id>
	<nonce_str><![CDATA[ogI5uvCHBoMny188DKpvhzZwOgpo7Ir2]]></nonce_str>
	<openid><![CDATA[oMSIv1oziYFbEGDQujq3i7-1e0gc]]></openid>
	<out_trade_no><![CDATA[1537030841dpstKL5e3vbasi]]></out_trade_no>
	<result_code><![CDATA[SUCCESS]]></result_code>
	<return_code><![CDATA[SUCCESS]]></return_code>
	<sign><![CDATA[52F86E5F8526BAFCF2128FC4BE079617]]></sign>
	<time_end><![CDATA[20180916010054]]></time_end>
	<total_fee>990</total_fee>
	<trade_type><![CDATA[JSAPI]]></trade_type>
	<transaction_id><![CDATA[4200000194201809165945099528]]></transaction_id>
	</xml>`)
	var callback mch.TradeCallback
	err := xml.Unmarshal(data, &callback)
	if err != nil {
		fmt.Println(err.Error())
	}

	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// 测试完一定记得删除key，不要提交到commit里
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	client := mch.InitClient("", "", "http://127.0.0.0.1", "http://1270.0.01")
	fmt.Println(client.SignXml(callback) == callback.Sign.Text)
}
