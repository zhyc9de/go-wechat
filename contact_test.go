package weUtil_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"gitee.com/hzsuoyi/go-wechat.git/util"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

func TestContactMgr_PostText(t *testing.T) {
	s, _ := json.Marshal(weUtil.SendMessage{
		ToUserName: "",
		MsgType:    "text",
		Text: &weUtil.ReturnTextMsg{
			Content: fmt.Sprintf(`您还未对小程序进行授权，<a href="%s">点击授权</a>。授权后即可进入公众号领取金币`, "test"),
		},
	})
	fmt.Println(string(s))

}

func TestContactMgr_PostText2(t *testing.T) {
	rawText := []byte(`<xml><ToUserName><![CDATA[gh_e136c6e50636]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408090816</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[pic_photo_or_album]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<SendPicsInfo><Count>1</Count>
<PicList><item><PicMd5Sum><![CDATA[5a75aaca956d97be686719218f275c6b]]></PicMd5Sum>
</item>
</PicList>
</SendPicsInfo>
</xml>`)
	var message weUtil.MessageXml
	err := xml.Unmarshal(rawText, &message)
	if err != nil {
		fmt.Println(err.Error())
	}
	s, _ := xml.Marshal(message)
	fmt.Println(string(s))
}
