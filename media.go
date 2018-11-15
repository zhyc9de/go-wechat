package weUtil

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

type MediaMgr interface {
	UploadMedia(mediaTyp string, media []byte) (mediaId string, err error)
	UploadMediaFile(mediaTyp, filePath string) (mediaId string, err error)
	GetMedia(mediaId string) (media []byte, err error)
}

// 上传媒体文件返回的结果
type UploadMediaResponse struct {
	ErrResp
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}

const MediaTypImage = "image"
const MediaTypVoice = "voice"
const MediaTypVideo = "video"
const MediaTypThumb = "thumb"

// 上传临时素材
func (c *Client) UploadMedia(mediaTyp string, media []byte) (mediaId string, err error) {
	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)

	var filename string
	var contentTyp string
	if mediaTyp == MediaTypImage {
		_, ext, _ := image.Decode(bytes.NewBuffer(media))
		filename = fmt.Sprintf("%s.%s", Md5(string(media)), ext)
		if ext == "jpeg" {
			contentTyp = "image/jpeg"
		} else if ext == "png" {
			contentTyp = "image/png"
		}
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="media"; filename="%s"`, filename))
	h.Set("Content-Type", contentTyp)
	part, _ := writer.CreatePart(h)
	part.Write(media)

	writer.Close()

	u := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s", c.GetToken(), mediaTyp)
	req, _ := http.NewRequest("POST", u, b)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	respJson := new(UploadMediaResponse)
	if err = DoRequestJson(req, respJson); err != nil {
		return
	}
	mediaId = respJson.MediaId
	return
}

// 上传本地文件，同时进行缓存下
func (c *Client) UploadMediaFile(mediaTyp, filePath string) (mediaId string, err error) {
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	cacheKey := fmt.Sprintf("wx:media:%s", Md5(string(f)))
	mediaId = c.Get(cacheKey)
	if mediaId != "" {
		return
	}
	// 如果没有的话，那么上传一下
	mediaId, err = c.UploadMedia(mediaTyp, f)
	c.Set(cacheKey, mediaId, 86400)
	return
}

// 下载临时素材
func (c *Client) GetMedia(mediaId string) (media []byte, err error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s", c.GetToken(), mediaId), nil)
	media, err = DoRequest(req)
	if bytes.Index(media, []byte("{")) == 0 {
		resp := new(ErrResp)
		if err = json.Unmarshal(media, resp); err == nil && resp.ErrCode != 0 {
			err = errors.New(resp.Error())
		}
	}
	return
}
