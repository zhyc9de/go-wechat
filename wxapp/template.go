package wxapp

import (
	"bytes"
	"fmt"
	"net/http"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

const tplMsgURL = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=%s"

func NewTemplateMsg(openId, templateId, page, formId string, data weUtil.TemplateMsgData) weUtil.WxaTemplate {
	return weUtil.WxaTemplate{
		TemplateMsg: weUtil.TemplateMsg{
			ToUser:     openId,
			TemplateId: templateId,
			Data:       data,
		},
		Page:   page,
		FormId: formId,
	}
}

// https://developers.weixin.qq.com/miniprogram/dev/api/open-api/template-message/sendTemplateMessage.html
func (client *Client) SendTemplateMsg(t weUtil.WxaTemplate) (err error) {
	u := fmt.Sprintf(tplMsgURL, client.GetToken())

	data, _ := json.Marshal(t)
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))

	err = weUtil.DoRequestJson(req, nil)
	return
}
