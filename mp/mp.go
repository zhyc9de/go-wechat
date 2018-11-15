package mp

import (
	"fmt"
	"net/http"

	"gitee.com/hzsuoyi/go-wechat.git/util"
	"github.com/json-iterator/go"
)

// 禁止转义html
var json = jsoniter.ConfigFastest

//------------------------------------------------------------------------------

type Client struct {
	*weUtil.Client
}

// 初始化sdk
func NewClient(client *weUtil.Client) *Client {
	return &Client{
		Client: client,
	}
}

// 获取jsApi ticket，用于公众号网页jssdk的
func (client *Client) GetJsApiTicket() string {
	token := client.Get(weUtil.KeyJSApi)
	if token != "" {
		return token
	}

	ticketURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", client.GetToken())
	req, _ := http.NewRequest("GET", ticketURL, nil)

	r := new(struct {
		weUtil.ErrResp
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
	})
	if err := weUtil.DoRequestJson(req, r); err != nil {
		return ""
	}

	client.Set(weUtil.KeyJSApi, r.Ticket, r.ExpiresIn-600)
	client.Logger.Infof("wx-auth mp jsApiTicket, appId=%s, token=%s", client.GetAppId(), r.Ticket)
	return r.Ticket
}
