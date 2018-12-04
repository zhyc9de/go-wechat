package wxapp

import (
	"bytes"
	"github.com/zhyc9de/go-wechat"
	"net/http"
)

type CreateActivityIdResp struct {
	weUtil.ErrResp
	ActivityId     string `json:"activity_id"`
	ExpirationTime int64  `json:"expiration_time"`
}

const (
	ParameterNameMember = "member_count"
	ParameterNameRoom   = "room_limit"
	ParameterNamePath   = "path"
	ParameterNameVer    = "version_type"
)

type Parameter struct {
	Name  string `json:"name"`  // 要修改的参数名
	Value string `json:"value"` // 修改后的参数值
}

type updateMsgParams struct {
	ActivityId   string      `json:"activity_id"`
	TargetState  int         `json:"target_state"`
	TemplateInfo []Parameter `json:"template_info"`
}

func (client *Client) CreateActivityId() (string, error) {
	var resp CreateActivityIdResp

	u := "https://api.weixin.qq.com/cgi-bin/message/wxopen/activityid/create?access_token=" + client.GetToken()
	req, _ := http.NewRequest("GET", u, nil)
	err := weUtil.DoRequestJson(req, &resp)
	if err != nil {
		return "", err
	}
	if resp.ErrCode != 0 {
		return "", resp
	}
	return resp.ActivityId, nil
}

func (client *Client) UpdateMsg(activityId string, targetState int, info []Parameter) error {
	u := "https://api.weixin.qq.com/cgi-bin/message/wxopen/updatablemsg/send?access_token=" + client.GetToken()

	body, err := json.Marshal(updateMsgParams{
		ActivityId:   activityId,
		TargetState:  targetState,
		TemplateInfo: info,
	})
	if err != nil {
		return err
	}

	var resp weUtil.ErrResp

	req, _ := http.NewRequest("POST", u, bytes.NewReader(body))
	err = weUtil.DoRequestJson(req, &resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return resp
	}

	return nil
}
