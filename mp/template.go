package mp

import (
	"bytes"
	"fmt"
	"net/http"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

const tplMsgURL = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"

// https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1433751277
func (client *Client) SendTemplateMsg(t weUtil.MpTemplate) error {
	u := fmt.Sprintf(tplMsgURL, client.GetToken())

	data, _ := json.Marshal(t)
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))

	return weUtil.DoRequestJson(req, nil)
}
