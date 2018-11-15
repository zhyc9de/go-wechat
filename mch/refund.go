package mch

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zhyc9de/go-wechat"
)

// 默认退全款
func (client *Client) NewRefund(appId, txnId, refundNo, desc string, fee int64) Refund {
	return Refund{
		AppId:         CDATA{appId},
		MchId:         CDATA{client.MchId},
		NonceStr:      CDATA{weUtil.GetRndString(32)},
		TransactionId: CDATA{txnId},
		OutRefundNo:   CDATA{refundNo},
		TotalFee:      fee,
		RefundFee:     fee,
		RefundDesc:    CDATA{desc},
		NotifyUrl:     CDATA{client.RefundURL},
	}
}

// 退款
func (client *Client) Refund(refund Refund) (refundResult RefundResult, err error) {
	if client.certConfig == nil {
		err = errors.New("no cert")
		return
	}

	refund.Sign = CDATA{client.SignXml(refund)}

	data, err := xml.Marshal(refund)
	if err != nil {
		return
	}

	var u string
	if client.debug {
		u = "https://api.mch.weixin.qq.com/sandbox/secapi/pay/refund"
	} else {
		u = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	}

	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))
	rb, err := client.DoRequest(req)
	if err != nil {
		return
	}

	if err = xml.Unmarshal(rb, &refundResult); err != nil {
		err = fmt.Errorf("parse xml err=%s, body=%s", err.Error(), string(rb))
		return
	} else if !refundResult.Valid() {
		err = refundResult
		return
	}
	// TODO 校验签名
	return
}

// 解码退款用户信息
// 解密步骤如下：
//（1）对加密串A做base64解码，得到加密串B
//（2）对商户key做md5，得到32位小写key* ( key设置路径：微信商户平台(pay.weixin.qq.com)-->账户设置-->API安全-->密钥设置 )
//（3）用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
func (client *Client) DecRefundReqInfo(reqInfo string) (info RefundInfo, err error) {
	cipherText, err := base64.StdEncoding.DecodeString(reqInfo)
	if err != nil {
		return
	}
	key := []byte(strings.ToLower(weUtil.Md5(client.MchKey)))

	plainText := DecryptAes256Ecb(cipherText, key)
	// 反序列化
	if err = xml.Unmarshal([]byte(plainText), &info); err != nil {
		return
	}
	return
}

func DecryptAes256Ecb(data []byte, key []byte) string {
	cipher, _ := aes.NewCipher([]byte(key))

	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	paddingSize := int(decrypted[len(decrypted)-1])
	return string(decrypted[0 : len(decrypted)-paddingSize])
}
