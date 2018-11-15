package mp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zhyc9de/go-wechat"
)

const qrCodeURL = "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s"

type (
	QrCodeParam struct {
		ExpireSeconds int64  `json:"expire_seconds,omitempty"`
		ActionName    string `json:"action_name"`
		ActionInfo    struct {
			Scene QrCodeScene `json:"scene"`
		} `json:"action_info"`
	}

	QrCodeScene struct {
		SceneId  int64  `json:"scene_id,omitempty"`
		SceneStr string `json:"scene_str,omitempty"`
	}

	QrCodeTicket struct {
		weUtil.ErrResp
		Ticket        string `json:"ticket"`
		ExpireSeconds int64  `json:"expire_seconds"`
		Url           string `json:"url,omitempty"`
	}
)

const (
	QrCodeMaxExpire = 2592000 // 临时最久

	ActionTmp        = "QR_SCENE"           // 临时的整型参数值
	ActionTmpStr     = "QR_STR_SCENE"       // 临时的字符串参数值
	ActionUnlimit    = "QR_LIMIT_SCENE"     // 永久的整型参数值
	ActionStrUnlimit = "QR_LIMIT_STR_SCENE" // 永久的字符串参数值
)

func NewTmpScendId(sceneId int64) QrCodeParam {
	return QrCodeParam{
		ExpireSeconds: QrCodeMaxExpire,
		ActionName:    ActionTmp,
		ActionInfo:    struct{ Scene QrCodeScene `json:"scene"` }{Scene: QrCodeScene{SceneId: sceneId}},
	}
}

func NewTmpScendStr(sceneStr string) QrCodeParam {
	return QrCodeParam{
		ExpireSeconds: QrCodeMaxExpire,
		ActionName:    ActionTmpStr,
		ActionInfo:    struct{ Scene QrCodeScene `json:"scene"` }{Scene: QrCodeScene{SceneStr: sceneStr}},
	}
}

func (client *Client) QrCode(param QrCodeParam) (ticket string, err error) {
	data, err := json.Marshal(param)
	if err != nil {
		return
	}
	u := fmt.Sprintf(qrCodeURL, client.GetToken())
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))
	resp := new(QrCodeTicket)
	err = weUtil.DoRequestJson(req, resp)
	if err == nil {
		ticket = resp.Ticket
	}
	return
}

// 直接请求腾讯
func (client *Client) GetQrCode(ticket string) (rb []byte, err error) {
	u := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", url.QueryEscape(ticket))
	req, _ := http.NewRequest("GET", u, nil)
	return weUtil.DoRequest(req)
}
