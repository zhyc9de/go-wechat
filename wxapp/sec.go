package wxapp

import (
	"bytes"
	"fmt"
	"image"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

// https://developers.weixin.qq.com/miniprogram/dev/api/open-api/sec-check/imgSecCheck.html
// 传入的图片
func (client *Client) ImgSecCheck(media []byte) (err error) {
	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)

	var filename string
	var contentTyp string
	_, ext, _ := image.Decode(bytes.NewBuffer(media))
	filename = fmt.Sprintf("%s.%s", weUtil.Md5(string(media)), ext)
	if ext == "jpeg" {
		contentTyp = "image/jpeg"
	} else if ext == "png" {
		contentTyp = "image/png"
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media"; filename="%s"`, filename))
	h.Set("Content-Type", contentTyp)
	part, _ := writer.CreatePart(h)
	part.Write(media)

	writer.Close()

	u := fmt.Sprintf("https://api.weixin.qq.com/wxa/img_sec_check?access_token=%s", client.GetToken())
	req, _ := http.NewRequest("POST", u, b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = weUtil.DoRequestJson(req, nil)
	return
}

// https://developers.weixin.qq.com/miniprogram/dev/api/open-api/sec-check/msgSecCheck.html
func (client *Client) MsgSecCheck(msg string) (err error) {
	data, _ := json.Marshal(struct {
		Content string `json:"content"`
	}{
		Content: msg,
	})
	u := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s", client.GetToken())
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))

	err = weUtil.DoRequestJson(req, nil)
	return
}
