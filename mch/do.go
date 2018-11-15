package mch

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	// 交易类型
	TradeTypJsApi  = "JSAPI"
	TradeTypNative = "NATIVE"
	TradeTypApp    = "APP"
	TradeTypWeb    = "MWEB"

	ReturnSuccess = "SUCCESS"
	ReturnFail    = "FAIL"

	ReturnMsgOk = "OK"

	IsSubscribe = "Y"
	NoSubscribe = "N"

	CheckNameForce = "FORCE_CHECK"
	CheckNameNo    = "NO_CHECK"

	SystemError = "SYSTEMERROR" // 企业付款时的错误，遇到要拿旧订单进行重试
)

type (
	CDATA struct {
		Text string `xml:",cdata"`
	}

	Coupon struct {
		CouponType CDATA // 代金券类型
		CouponId   CDATA // 代金券ID
		CouponFee  int64 // 单个代金券支付金额
	}

	CouponMap map[int]Coupon
)

func (m CouponMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	for i := range m {
		start.Name = xml.Name{
			Local: fmt.Sprintf("coupon_type_%d", i),
		}
		if err := e.EncodeElement(m[i].CouponType, start); err != nil {
			return err
		}
		start.Name = xml.Name{
			Local: fmt.Sprintf("coupon_id_%d", i),
		}
		if err := e.EncodeElement(m[i].CouponId, start); err != nil {
			return err
		}
		start.Name = xml.Name{
			Local: fmt.Sprintf("coupon_fee_%d", i),
		}
		if err := e.EncodeElement(m[i].CouponFee, start); err != nil {
			return err
		}
	}
	return nil
}

//------------------------------------------------------------------------------

// 回调完成返回结构体
type CallbackResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode CDATA    `xml:"return_code"`
	ReturnMsg  CDATA    `xml:"return_msg"`
	ResultCode CDATA    `xml:"result_code,omitempty"`  // 业务结果
	ErrCode    CDATA    `xml:"err_code,omitempty"`     // 错误代码
	ErrCodeDes CDATA    `xml:"err_code_des,omitempty"` // 错误代码描述
}

func (res CallbackResponse) Valid() bool {
	return res.ReturnCode.Text == ReturnSuccess && res.ResultCode.Text == ReturnSuccess
}

func (res CallbackResponse) Error() string {
	return fmt.Sprintf("return_code=%s, return_msg=%s, err_code=%s, err_code_des=%s",
		res.ReturnCode.Text, res.ReturnMsg.Text, res.ErrCode.Text, res.ErrCodeDes.Text)
}

func (res CallbackResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type callback struct {
		XMLName    xml.Name `xml:"xml"`
		ReturnCode CDATA    `xml:"return_code"`
		ReturnMsg  CDATA    `xml:"return_msg"`
	}
	return e.EncodeElement(callback{
		ReturnCode: res.ReturnCode,
		ReturnMsg:  res.ReturnMsg,
	}, start)
}

//------------------------------------------------------------------------------

// 网页端支付参数
type WCPay struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

//------------------------------------------------------------------------------

// 统一下单
type (
	Trade struct {
		XMLName        xml.Name `xml:"xml"`
		AppId          CDATA    `xml:"appid"`                 // 公众账号ID
		MchId          CDATA    `xml:"mch_id"`                // 商户号
		DeviceInfo     CDATA    `xml:"device_info,omitempty"` // 设备号
		NonceStr       CDATA    `xml:"nonce_str"`             // 随机字符串
		Sign           CDATA    `xml:"sign"`                  // 签名
		SignType       CDATA    `xml:"sign_type,omitempty"`   // 签名类型
		Body           CDATA    `xml:"body"`                  // 商品描述
		Detail         CDATA    `xml:"detail,omitempty"`      // 商品详情
		Attach         CDATA    `xml:"attach,omitempty"`      // 附加数据
		OutTradeNo     CDATA    `xml:"out_trade_no"`          // 商户订单号
		FeeType        CDATA    `xml:"fee_type,omitempty"`    // 标价币种
		TotalFee       int64    `xml:"total_fee"`             // 标价金额
		SpbillCreateIp CDATA    `xml:"spbill_create_ip"`      // 终端IP
		TimeStart      int64    `xml:"time_start,omitempty"`  // 交易起始时间
		TimeExpire     int64    `xml:"time_expire,omitempty"` // 交易结束时间
		GoodsTag       CDATA    `xml:"goods_tag,omitempty"`   // 订单优惠标记
		NotifyUrl      CDATA    `xml:"notify_url"`            // 通知地址
		TradeType      CDATA    `xml:"trade_type"`            // 交易类型
		ProductId      CDATA    `xml:"product_id,omitempty"`  // 商品ID trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
		LimitPay       CDATA    `xml:"limit_pay,omitempty"`   // 指定支付方式 上传此参数no_credit--可限制用户不能使用信用卡支付
		OpenId         CDATA    `xml:"openid,omitempty"`      // 用户标识 trade_type=JSAPI时（即公众号支付），此参数必传
		// 还有一个场景值，就暂时不填了 https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
	}

	// 统一下单返回结构体
	TradeResult struct {
		CallbackResponse
		// 以下字段在return_code为SUCCESS的时候有返回
		AppId      CDATA `xml:"appid"`                 // 公众账号ID
		MchId      CDATA `xml:"mch_id"`                // 商户号
		DeviceInfo CDATA `xml:"device_info,omitempty"` // 设备号
		NonceStr   CDATA `xml:"nonce_str"`             // 随机字符串
		Sign       CDATA `xml:"sign"`                  // 签名
		// 以下字段在return_code 和result_code都为SUCCESS的时候有返回
		TradeType CDATA `xml:"trade_type"`         // 交易类型
		PrepayId  CDATA `xml:"prepay_id"`          // 预支付交易会话标识
		CodeUrl   CDATA `xml:"code_url,omitempty"` // 二维码链接
	}
)

//------------------------------------------------------------------------------

// 支付回调返回结构体
type (
	TradeCallback struct {
		CallbackResponse
		AppId              CDATA     `xml:"appid"`                          // 公众账号ID
		MchId              CDATA     `xml:"mch_id"`                         // 商户号
		DeviceInfo         CDATA     `xml:"device_info,omitempty"`          // 设备号
		NonceStr           CDATA     `xml:"nonce_str"`                      // 随机字符串
		Sign               CDATA     `xml:"sign"`                           // 签名
		SignType           CDATA     `xml:"sign_type,omitempty"`            // 签名类型
		OpenId             CDATA     `xml:"openid"`                         // 用户标识
		IsSubscribe        CDATA     `xml:"is_subscribe,omitempty"`         // 是否关注公众账号
		TradeType          CDATA     `xml:"trade_type"`                     // 交易类型
		BankType           CDATA     `xml:"bank_type"`                      // 付款银行
		TotalFee           int64     `xml:"total_fee"`                      // 订单金额
		SettlementTotalFee int64     `xml:"settlement_total_fee,omitempty"` // 应结订单金额
		FeeType            CDATA     `xml:"fee_type,omitempty"`             // 货币种类
		CashFee            int64     `xml:"cash_fee"`                       // 现金支付金额
		CashFeeType        CDATA     `xml:"cash_fee_type,omitempty"`        // 现金支付货币类型
		CouponFee          int64     `xml:"coupon_fee,omitempty"`           // 总代金券金额
		CouponCount        int64     `xml:"coupon_count,omitempty"`         // 代金券使用数量
		TransactionId      CDATA     `xml:"transaction_id"`                 // 微信支付订单号
		OutTradeNo         CDATA     `xml:"out_trade_no"`                   // 商户订单号
		Attach             CDATA     `xml:"attach,omitempty"`               // 商家数据包
		TimeEnd            CDATA     `xml:"time_end"`                       // 支付完成时间
		TradeCoupon        CouponMap `xml:"coupon_map"`                     // 代金券
	}
)

func (trade *TradeCallback) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tradeCoupon := make(CouponMap)
	tagToStructName := map[string]string{
		"return_code":  "ReturnCode",
		"return_msg":   "ReturnMsg",
		"result_code":  "ResultCode",
		"err_code":     "ErrCode",
		"err_code_des": "ErrCodeDes",
	}

	typ := reflect.TypeOf(*trade)
	for i := 0; i < typ.NumField(); i++ {
		tagName := omitemptyReplacer.Replace(typ.Field(i).Tag.Get("xml"))
		tagToStructName[tagName] = typ.Field(i).Name
	}
	//fmt.Println(tagToStructName)
	ps := reflect.ValueOf(trade)
	s := ps.Elem()

	var startOffset *xml.StartElement
	for {
		t, _ := d.Token()
		switch tt := t.(type) {

		case xml.StartElement:
			//fmt.Println(">", tt.Name)
			startOffset = &tt

		case xml.EndElement:
			//fmt.Println("<", tt.Name)
			if tt.Name == start.Name {
				trade.TradeCoupon = tradeCoupon
				return nil
			}
			startOffset = nil

		case xml.CharData:
			if startOffset != nil {
				tagName := startOffset.Name.Local
				//fmt.Printf("struct name %s, offset name %s, value %s\n", tagToStructName[tagName], tagName, string(tt))
				field := s.FieldByName(tagToStructName[startOffset.Name.Local])
				if field.IsValid() && field.CanSet() {
					switch field.Kind() {
					case reflect.Int64:
						x, _ := strconv.ParseInt(string(tt), 10, 64)
						field.SetInt(x)
					case reflect.Struct:
						v := reflect.ValueOf(CDATA{string(tt)})
						field.Set(v)
					}
				} else if strings.Index(tagName, "coupon_type") >= 0 {
					k := parseCoupon(tagName)
					if _, ok := tradeCoupon[k]; ok {
						v := tradeCoupon[k]
						v.CouponType = CDATA{string(tt)}
						tradeCoupon[k] = v
					} else {
						tradeCoupon[k] = Coupon{
							CouponType: CDATA{string(tt)},
						}
					}
				} else if strings.Index(tagName, "coupon_id") >= 0 {
					k := parseCoupon(tagName)
					if _, ok := tradeCoupon[k]; ok {
						v := tradeCoupon[k]
						v.CouponId = CDATA{string(tt)}
						tradeCoupon[k] = v
					} else {
						tradeCoupon[k] = Coupon{
							CouponId: CDATA{string(tt)},
						}
					}
				} else if strings.Index(tagName, "coupon_fee") >= 0 {
					k := parseCoupon(tagName)
					x, _ := strconv.ParseInt(string(tt), 10, 64)
					if _, ok := tradeCoupon[k]; ok {
						v := tradeCoupon[k]
						v.CouponFee = x
						tradeCoupon[k] = v
					} else {
						tradeCoupon[k] = Coupon{
							CouponFee: x,
						}
					}
				} else {
					return fmt.Errorf("cannot parse tag: %s", tagName)
				}
			}

		}
	}

}

func parseCoupon(tag string) int {
	s := strings.Split(tag, "_")
	k, _ := strconv.Atoi(s[len(s)-1])
	return k
}

//------------------------------------------------------------------------------

// 退款结构体
type (
	Refund struct {
		XMLName       xml.Name `xml:"xml"`
		AppId         CDATA    `xml:"appid"`                     // 公众账号ID
		MchId         CDATA    `xml:"mch_id"`                    // 商户号
		NonceStr      CDATA    `xml:"nonce_str"`                 // 随机字符串
		Sign          CDATA    `xml:"sign"`                      // 签名
		SignType      CDATA    `xml:"sign_type,omitempty"`       // 签名类型
		TransactionId CDATA    `xml:"transaction_id,omitempty"`  // 交易单号
		OutTradeNo    CDATA    `xml:"out_trade_no,omitempty"`    // 商户订单号 和交易单号二选一
		OutRefundNo   CDATA    `xml:"out_refund_no"`             // 退款单号
		TotalFee      int64    `xml:"total_fee"`                 // 订单金额
		RefundFee     int64    `xml:"refund_fee"`                // 退款金额
		RefundFeeType CDATA    `xml:"refund_fee_type,omitempty"` // 货币种类
		RefundDesc    CDATA    `xml:"refund_desc,omitempty"`     // 退款原因
		RefundAccount CDATA    `xml:"refund_account,omitempty"`  // 退款资金来源
		NotifyUrl     CDATA    `xml:"notify_url,omitempty"`      // 退款结果通知url
	}

	RefundResult struct {
		CallbackResponse
		AppId               CDATA     `xml:"appid"`                           // 公众账号ID
		MchId               CDATA     `xml:"mch_id"`                          // 商户号
		NonceStr            CDATA     `xml:"nonce_str"`                       // 随机字符串
		Sign                CDATA     `xml:"sign"`                            // 签名
		TransactionId       CDATA     `xml:"transaction_id"`                  // 交易单号
		OutTradeNo          CDATA     `xml:"out_trade_no"`                    // 商户订单号
		OutRefundNo         CDATA     `xml:"out_refund_no"`                   // 退款单号
		RefundId            CDATA     `xml:"refund_id"`                       // 微信退款单号
		RefundFee           int64     `xml:"refund_fee"`                      // 退款金额
		SettlementRefundFee int64     `xml:"settlement_refund_fee,omitempty"` // 应结退款金额
		TotalFee            int64     `xml:"total_fee"`                       // 标价金额
		SettlementTotalFee  int64     `xml:"settlement_total_fee,omitempty"`  // 应结订单金额
		FeeType             CDATA     `xml:"fee_type,omitempty"`              // 标价币种
		CashFee             int64     `xml:"cash_fee"`                        // 现金支付金额
		CashFeeType         CDATA     `xml:"cash_fee_type,omitempty"`         // 现金支付币种
		CashRefundFee       int64     `xml:"cash_refund_fee,omitempty"`       // 现金退款金额
		CouponRefundFee     int64     `xml:"coupon_refund_fee,omitempty"`     // 代金券退款总金额
		CouponRefundCount   int64     `xml:"coupon_refund_count,omitempty"`   // 退款代金券使用数量
		RefundChannel       CDATA     `xml:"refund_channel"`                  // 退款渠道
		CouponRefund        CouponMap `xml:"coupon_map"`
	}
)

func (trade *RefundResult) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tradeCoupon := make(CouponMap)
	tagToStructName := map[string]string{
		"return_code":  "ReturnCode",
		"return_msg":   "ReturnMsg",
		"result_code":  "ResultCode",
		"err_code":     "ErrCode",
		"err_code_des": "ErrCodeDes",
	}

	typ := reflect.TypeOf(*trade)
	for i := 0; i < typ.NumField(); i++ {
		tagName := omitemptyReplacer.Replace(typ.Field(i).Tag.Get("xml"))
		tagToStructName[tagName] = typ.Field(i).Name
	}
	//fmt.Println(tagToStructName)
	ps := reflect.ValueOf(trade)
	s := ps.Elem()

	var startOffset *xml.StartElement
	for {
		t, _ := d.Token()
		switch tt := t.(type) {

		case xml.StartElement:
			//fmt.Println(">", tt.Name)
			startOffset = &tt

		case xml.EndElement:
			//fmt.Println("<", tt.Name)
			if tt.Name == start.Name {
				trade.CouponRefund = tradeCoupon
				return nil
			}
			startOffset = nil

		case xml.CharData:
			if startOffset != nil {
				tagName := startOffset.Name.Local
				//fmt.Printf("struct name %s, offset name %s, value %s\n", tagToStructName[tagName], tagName, string(tt))
				field := s.FieldByName(tagToStructName[startOffset.Name.Local])
				if field.IsValid() && field.CanSet() {
					switch field.Kind() {
					case reflect.Int64:
						x, _ := strconv.ParseInt(string(tt), 10, 64)
						field.SetInt(x)
					case reflect.Struct:
						v := reflect.ValueOf(CDATA{string(tt)})
						field.Set(v)
					}
				} else if strings.Index(tagName, "coupon_type") >= 0 {
					k := parseCoupon(tagName)
					if _, ok := tradeCoupon[k]; ok {
						v := tradeCoupon[k]
						v.CouponType = CDATA{string(tt)}
						tradeCoupon[k] = v
					} else {
						tradeCoupon[k] = Coupon{
							CouponType: CDATA{string(tt)},
						}
					}
				} else if strings.Index(tagName, "coupon_refund_id") >= 0 {
					k := parseCoupon(tagName)
					if _, ok := tradeCoupon[k]; ok {
						v := tradeCoupon[k]
						v.CouponId = CDATA{string(tt)}
						tradeCoupon[k] = v
					} else {
						tradeCoupon[k] = Coupon{
							CouponId: CDATA{string(tt)},
						}
					}
				} else if strings.Index(tagName, "coupon_refund_fee") >= 0 {
					k := parseCoupon(tagName)
					x, _ := strconv.ParseInt(string(tt), 10, 64)
					if _, ok := tradeCoupon[k]; ok {
						v := tradeCoupon[k]
						v.CouponFee = x
						tradeCoupon[k] = v
					} else {
						tradeCoupon[k] = Coupon{
							CouponFee: x,
						}
					}
				} else {
					return fmt.Errorf("cannot parse tag: %s", tagName)
				}
			}

		}
	}
}

//------------------------------------------------------------------------------

// 退款回调
type (
	RefundCallback struct {
		XMLName    xml.Name `xml:"xml"`
		ReturnCode string   `xml:"return_code"`
		ReturnMsg  CDATA    `xml:"return_msg"`
		// 以下字段在return_code为SUCCESS的时候有返回
		AppId    CDATA `xml:"appid"`     // 公众账号ID
		MchId    CDATA `xml:"mch_id"`    // 商户号
		NonceStr CDATA `xml:"nonce_str"` // 随机字符串
		ReqInfo  CDATA `xml:"req_info"`  // 加密信息
	}

	RefundInfo struct {
		TransactionId       CDATA `xml:"transaction_id"`
		OutTradeNo          CDATA `xml:"out_trade_no"`
		RefundId            CDATA `xml:"refund_id"`
		OutRefundNo         CDATA `xml:"out_refund_no"`
		TotalFee            int64 `xml:"total_fee"`
		SettlementTotalFee  int64 `xml:"settlement_total_fee,omitempty"`
		RefundFee           int64 `xml:"refund_fee"`
		SettlementRefundFee int64 `xml:"settlement_refund_fee,omitempty"`
		RefundStatus        CDATA `xml:"refund_status"`
		SuccessTime         CDATA `xml:"success_time,omitempty"`
		RefundRecvAccout    CDATA `xml:"refund_recv_accout"`
		RefundAccount       CDATA `xml:"refund_account"`
		RefundRequestSource CDATA `xml:"refund_request_source"`
	}
)

//------------------------------------------------------------------------------

// 企业付款
type (
	Transfers struct {
		XMLName        xml.Name `xml:"xml"`
		MchAppId       CDATA    `xml:"mch_appid"`
		MchId          CDATA    `xml:"mchid"`
		DeviceInfo     CDATA    `xml:"device_info,omitempty"`
		NonceStr       CDATA    `xml:"nonce_str"`
		Sign           CDATA    `xml:"sign"`
		PartnerTradeNo CDATA    `xml:"partner_trade_no"`       // 商户订单号
		OpenId         CDATA    `xml:"openid"`                 // 用户openid
		CheckName      CDATA    `xml:"check_name"`             // 校验用户姓名选项
		ReUserName     CDATA    `xml:"re_user_name,omitempty"` // 收款用户姓名
		Amount         int64    `xml:"amount"`                 // 金额
		Desc           CDATA    `xml:"desc"`                   // 企业付款备注
		SpbillCreateIp CDATA    `xml:"spbill_create_ip"`       // Ip地址
	}

	TransfersResult struct {
		CallbackResponse
		MchAppId       CDATA `xml:"mch_appid"`
		MchId          CDATA `xml:"mch_id"`
		DeviceInfo     CDATA `xml:"device_info,omitempty"`
		NonceStr       CDATA `xml:"nonce_str"`
		PartnerTradeNo CDATA `xml:"partner_trade_no"` // 商户订单号，需保持历史全局唯一性(只能是字母或者数字，不能包含有其他字符)
		PaymentNo      CDATA `xml:"payment_no"`       // 企业付款成功，返回的微信付款单号
		PaymentTime    CDATA `xml:"payment_time"`     // 企业付款成功时间
	}
)
