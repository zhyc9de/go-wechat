package wxapp

import (
	"bytes"
	"github.com/zhyc9de/go-wechat"
	"net/http"
)

type AnalysisDateRange struct {
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
}

func (client *Client) GetAnalysisData(action, begin, end string) (rb []byte, err error) {
	u := "https://api.weixin.qq.com/datacube/" + action + "?access_token=" + client.GetToken()
	body, _ := json.Marshal(AnalysisDateRange{
		BeginDate: begin,
		EndDate:   end,
	})
	req, _ := http.NewRequest("POST", u, bytes.NewReader(body))
	return weUtil.DoRequest(req)
}
