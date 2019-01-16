package weUtil

import (
	"bytes"
	"github.com/zhyc9de/go-wechat/wxapp"
	"net/http"
)

type CommonClient interface {
	Logger
	TokenMgr
	MediaMgr
	ContactMgr
	SendContactMgr
}

type Client struct {
	Logger
	TokenMgr
	MediaMgr
	ContactMgr
	SendContactMgr

	contactHook ContactHook
}

func NewClient(tokenMgr TokenMgr, logger Logger) *Client {
	return &Client{
		Logger:     logger,
		TokenMgr:   tokenMgr,
		ContactMgr: &Contact{},
	}
}

// 获取access token
func (c *Client) GetToken() string {
	token, isNew, err := c.TokenMgr.GetOrNewToken()
	if err != nil {
		c.Logger.Errorf("wx-auth accessToken appId=%s, err=%s", err.Error(), c.GetAppId())
		return ""
	}
	if isNew {
		c.Logger.Infof("wx-auth accessToken, appId=%s, token=%s", c.GetAppId(), token)
	}
	return token
}

// 获取小程序或者公众号数据
type AnalysisArgs struct {
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
}

func (c *Client) GetAnalysisData(action, begin, end string) (rb []byte, err error) {
	u := "https://api.weixin.qq.com/datacube/" + action + "?access_token=" + c.GetToken()
	body, _ := json.Marshal(AnalysisArgs{
		BeginDate: begin,
		EndDate:   end,
	})
	req, _ := http.NewRequest("POST", u, bytes.NewReader(body))
	return DoRequest(req)
}

// 判断data cube action是否是获取小程序数据
func IsDataCudeForWxa(action string) bool {
	for i := range wxapp.AnalysisActions {
		if action == wxapp.AnalysisActions[i] {
			return true
		}
	}
	return false
}
