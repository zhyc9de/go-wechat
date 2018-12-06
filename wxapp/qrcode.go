package wxapp

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/zhyc9de/go-wechat"
)

const codeURL = "https://api.weixin.qq.com/wxa/getwxacode?access_token=%s"
const codeUnlimitURL = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s"

type LineColor struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

type WxaCode struct {
	Path      string     `json:"path"`
	Scene     string     `json:"scene,omitempty"` // 永久二维码，path不能带参数
	Width     int        `json:"width,omitempty"`
	AutoColor bool       `json:"auto_color"`
	LineColor *LineColor `json:"line_color,omitempty"`
	IsHyaline bool       `json:"is_hyaline,omitempty"` // 是否需要透明底色， is_hyaline 为true时，生成透明底色的小程序码
}

// 生成小程序二维码
// https://developers.weixin.qq.com/miniprogram/dev/api/open-api/qr-code/createWXAQRCode.html
func (client *Client) WxaCode(param WxaCode, tmp bool) (wxacode []byte, err error) {
	var u string
	var data []byte
	if tmp {
		u = fmt.Sprintf(codeURL, client.GetToken())
	} else {
		u = fmt.Sprintf(codeUnlimitURL, client.GetToken())
	}

	if data, err = json.Marshal(param); err != nil {
		return
	}

	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))
	wxacode, err = weUtil.DoRequest(req)
	// 判断是否是json
	if bytes.Index(wxacode, []byte("{")) == 0 {
		resp := new(weUtil.ErrResp)
		if err = json.Unmarshal(wxacode, resp); err == nil && resp.ErrCode != 0 {
			err = errors.New(resp.Error())
		}
	}
	return
}
