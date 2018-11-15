package mp

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

type JsSdkConfig struct {
	AppId     string      `json:"appId"`
	Timestamp int64       `json:"timestamp"`
	NonceStr  string      `json:"nonceStr"`
	Signature string      `json:"signature"`
	JsApiList []JsApiName `json:"jsApiList"`
}

type JsApiName string

const ( // api 权限列表
	ApiOnMenuShareTimeline     JsApiName = "onMenuShareTimeline"
	ApiOnMenuShareAppMessage             = "onMenuShareAppMessage"
	ApiOnMenuShareQQ                     = "onMenuShareQQ"
	ApiOnMenuShareWeibo                  = "onMenuShareWeibo"
	ApiOnMenuShareQZone                  = "onMenuShareQZone"
	ApiStartRecord                       = "startRecord"
	ApiStopRecord                        = "stopRecord"
	ApiOnVoiceRecordEnd                  = "onVoiceRecordEnd"
	ApiPlayVoice                         = "playVoice"
	ApiPauseVoice                        = "pauseVoice"
	ApiStopVoice                         = "stopVoice"
	ApiOnVoicePlayEnd                    = "onVoicePlayEnd"
	ApiUploadVoice                       = "uploadVoice"
	ApiDownloadVoice                     = "downloadVoice"
	ApiChooseImage                       = "chooseImage"
	ApiPreviewImage                      = "previewImage"
	ApiUploadImage                       = "uploadImage"
	ApiDownloadImage                     = "downloadImage"
	ApiTranslateVoice                    = "translateVoice"
	ApiGetNetworkType                    = "getNetworkType"
	ApiOpenLocation                      = "openLocation"
	ApiGetLocation                       = "getLocation"
	ApiHideOptionMenu                    = "hideOptionMenu"
	ApiShowOptionMenu                    = "showOptionMenu"
	ApiHideMenuItems                     = "hideMenuItems"
	ApiShowMenuItems                     = "showMenuItems"
	ApiHideAllNonBaseMenuItem            = "hideAllNonBaseMenuItem"
	ApiShowAllNonBaseMenuItem            = "showAllNonBaseMenuItem"
	ApiCloseWindow                       = "closeWindow"
	ApiScanQRCode                        = "scanQRCode"
	ApiChooseWXPay                       = "chooseWXPay"
	ApiOpenProductSpecificView           = "openProductSpecificView"
	ApiAddCard                           = "addCard"
	ApiChooseCard                        = "chooseCard"
	ApiOpenCard                          = "openCard"
)

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421141115
func (client *Client) NewJsSdkConfig(uri string, apiList []JsApiName) JsSdkConfig {
	// 生成一个随机数
	noncestr := weUtil.GetRndString(16)
	timestamp := time.Now().Unix()

	signature := signTicket(client.GetJsApiTicket(), noncestr, uri, timestamp)
	return JsSdkConfig{
		AppId:     client.GetAppId(),
		Timestamp: timestamp,
		NonceStr:  noncestr,
		Signature: signature,
		JsApiList: apiList,
	}
}

func signTicket(ticket, noncestr, uri string, timestamp int64) string {
	s := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, noncestr, timestamp, uri)
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
