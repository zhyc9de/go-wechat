package mch

import "errors"

// 支付回调
var ErrCallbackFee = errors.New("valid callback fee not equal")
var ErrCallbackSign = errors.New("valid callback sign not equal")

// 校验支付回调签名和金额
func (client *Client) ValidateOrder(res TradeCallback, fee int64) error {
	// 校验金额
	if res.TotalFee != fee {
		return ErrCallbackFee
	}
	// 校验签名
	if client.SignXml(res) != res.Sign.Text {
		return ErrCallbackSign
	}

	return nil
}

func CallbackSuccess() CallbackResponse {
	return CallbackResponse{
		ReturnCode: CDATA{ReturnSuccess},
		ReturnMsg:  CDATA{ReturnMsgOk},
	}
}

func CallbackFail(err string) CallbackResponse {
	return CallbackResponse{
		ReturnCode: CDATA{ReturnFail},
		ReturnMsg:  CDATA{err},
	}
}
