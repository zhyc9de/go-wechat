package weUtil

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

const contactURL = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token="

const (
	ContactTextMaxLimit = 600

	// 小程序消息类型推送
	MsgTypText            = "text"
	MsgTypImage           = "image"
	MsgTypVoice           = "voice"
	MsgTypVideo           = "video"
	MsgTypShortvideo      = "shortvideo"
	MsgTypLocation        = "location"
	MsgTypLink            = "link"
	MsgTypEvent           = "event"
	MsyTypMiniprogramPage = "miniprogrampage"

	// 如果 MsgTypEvent = event
	EventTypSubscribe   = "subscribe"
	EventTypUnsubscribe = "unsubscribe"
	EventTypScan        = "SCAN"
	EventTypLocation    = "LOCATION"
	// 小程序用户进入客服消息
	EventTypUserEnterTempsession = "user_enter_tempsession"
	EventTypTemplateJobFinish    = "TEMPLATESENDJOBFINISH" // 模板消息发送成功
	// 菜单事件
	EventTypClick           = "CLICK"
	EventTypView            = "VIEW"
	EventTypViewMiniprogram = "view_miniprogram"
	EventTypScancodePush    = "scancode_push"
	EventTypScancodeWaitmsg = "scancode_waitmsg"
	EventTypPicSysphoto     = "pic_sysphoto"
	EventTypPicPhotoOrAlbum = "pic_photo_or_album"
	EventTypPicWeixin       = "pic_weixin"
	EventTypLocationSelect  = "location_select"

	// 返回内容类型
	ReturnMsgText            = "text"
	ReturnMsgImage           = "image"
	ReturnMsgVoice           = "voice"
	ReturnMsgVideo           = "video"
	ReturnMsgMusic           = "music"
	ReturnMsgNews            = "news"
	ReturnMsgMpnews          = "mpnews"
	ReturnMsgWxcard          = "wxcard"
	ReturnMsgMiniprogramPage = "miniprogrampage"
	ReturnMsgLink            = "link"

	ReturnMsgOk = "success"
)

type (
	CDATA struct {
		Text string `xml:",cdata"`
	}

	// 接受微信发来的消息 xml
	// 关于重试的消息排重，有msgid的消息推荐使用msgid排重。事件类型消息推荐使用FromUserName + CreateTime 排重。
	MessageXml struct {
		BaseMsg
		CommonMsg
		EventMsg
	}

	BaseMsg struct {
		XMLName      xml.Name `xml:"xml"`
		ToUserName   CDATA    `xml:"ToUserName"`
		FromUserName CDATA    `xml:"FromUserName"`
		CreateTime   int64    `xml:"CreateTime"`
		MsgType      CDATA    `xml:"MsgType"`
	}

	// 接收普通消息
	CommonMsg struct {
		// 文本消息
		Content CDATA `xml:"Content"`
		// 图片消息
		PicUrl  CDATA `xml:"PicUrl"`
		MediaId CDATA `xml:"MediaId"`
		// 语音消息
		Format      CDATA `xml:"Format"`
		Recognition CDATA `xml:"Recognition"`
		// TODO 其他普通消息
	}

	// 接收事件推送
	EventMsg struct {
		Event        CDATA `xml:"Event"`
		EventKey     CDATA `xml:"EventKey"`
		Ticket       CDATA `xml:"Ticket"`
		MenuID       int64 `xml:"MenuID"`
		ScanCodeInfo struct {
			ScanType   CDATA `xml:"ScanType"`
			ScanResult CDATA `xml:"ScanResult"`
		} `xml:"ScanCodeInfo"`
		SendPicsInfo struct {
			Count   int64 `xml:"Count"`
			PicList []struct {
				Item struct {
					PicMd5Sum CDATA `xml:"PicMd5Sum"`
				} `xml:"item"`
			} `xml:"PicList"`
		}
	}

	// 接受微信发来的消息
	Message struct {
		ToUserName   string `json:"ToUserName"`
		FromUserName string `json:"FromUserName"`
		CreateTime   int64  `json:"CreateTime"`
		MsgId        int64  `json:"MsgId"`
		MsgType      string `json:"MsgType"`
		// 图片
		PicUrl  string `json:"PicUrl"`
		MediaId string `json:"MediaId"`
		// 文本
		Content string `json:"Content"`
		// 卡片消息
		Title        string `json:"Title"`
		AppId        string `json:"AppId"`
		PagePath     string `json:"PagePath"`
		ThumbUrl     string `json:"ThumbUrl"`
		ThumbMediaId string `json:"ThumbMediaId"`
		//	进入回话
		Event       string `json:"Event"`
		SessionFrom string `json:"SessionFrom"`
	}

	HrefData struct {
		AppId string `json:"miniprogram_appid"`
		Path  string `json:"miniprogram-path"`
	}

	TextHref struct {
		Text string
		Path string
		Data HrefData
	}

	// 客服消息结构体
	SendMessage struct {
		ToUserName      string          `json:"touser"`
		MsgType         string          `json:"msgtype"`
		Text            *ReturnTextMsg  `json:"text,omitempty"`
		Image           *ReturnImageMsg `json:"image,omitempty"`
		Link            *ReturnLinkMsg  `json:"link,omitempty"`
		MiniProgramPage *ReturnWXAppMsg `json:"miniprogrampage,omitempty"`
	}

	// 文本消息
	ReturnTextMsg struct {
		Content string `json:"content"`
	}

	// 图片消息
	ReturnImageMsg struct {
		MediaId string `json:"media_id"`
	}

	// 发送链接(小程序特有?
	ReturnLinkMsg struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		ThumbUrl    string `json:"thumb_url"`
	}

	// 小程序卡片
	ReturnWXAppMsg struct {
		Title        string `json:"title"`
		AppId        string `json:"appid,omitempty"`
		PagePath     string `json:"pagepath"`
		ThumbMediaId string `json:"thumb_media_id"`
	}

	// 返回客服消息用的
	RespMessage struct {
		ToUserName   string `json:"ToUserName"`
		FromUserName string `json:"FromUserName"`
		CreateTime   int64  `json:"CreateTime"`
		MsgType      string `json:"MsgType"`
	}
)

//------------------------------------------------------------------------------

// 给消息体增加一些通用函数
type ReceiveMessage interface {
	GetFromUserName() string
	GetToUserName() string
	GetCreateTime() int64
	ContainsKeywords([]string) bool // 判断是否包含关键词
	EqualKeywords([]string) bool    // 完全等于keywords
}

//------------------------------------------------------------------------------

func (m *MessageXml) GetFromUserName() string {
	return m.FromUserName.Text
}

func (m *MessageXml) GetToUserName() string {
	return m.ToUserName.Text
}

func (m *MessageXml) GetCreateTime() int64 {
	return m.CreateTime
}

// 扫二维码的场景值
func (m *MessageXml) SceneValue() string {
	if m.Event.Text == EventTypSubscribe && strings.Index(m.EventKey.Text, "qrscene_") >= 0 {
		return m.EventKey.Text[8:] // 即去掉qrscene_
	} else if m.Event.Text == EventTypScan {
		return m.EventKey.Text
	} else {
		return ""
	}
}

// 判断是否来自扫二维码
func (m *MessageXml) IsFromScan() bool {
	return m.MsgType.Text == MsgTypEvent && (
		(m.Event.Text == EventTypSubscribe && strings.Index(m.EventKey.Text, "qrscene_") >= 0) || // 扫描关注
			(m.Event.Text == EventTypScan)) // 已关注的扫描
}

func (m *MessageXml) ContainsKeywords(keywords []string) bool {
	for _, word := range keywords {
		if strings.Contains(m.Content.Text, word) {
			return true
		}
	}

	return false
}

func (m *MessageXml) EqualKeywords(keywords []string) bool {
	for _, word := range keywords {
		if m.Content.Text == word {
			return true
		}
	}

	return false
}

//------------------------------------------------------------------------------

func (m *Message) GetFromUserName() string {
	return m.FromUserName
}

func (m *Message) GetToUserName() string {
	return m.ToUserName
}

func (m *Message) GetCreateTime() int64 {
	return m.CreateTime
}

func (m *Message) ContainsKeywords(keywords []string) bool {
	for _, word := range keywords {
		if strings.Contains(m.Content, word) {
			return true
		}
	}

	return false
}

func (m *Message) EqualKeywords(keywords []string) bool {
	for _, word := range keywords {
		if m.Content == word {
			return true
		}
	}

	return false
}

//------------------------------------------------------------------------------

// 格式化CDATA
func (c CDATA) String() string {
	return c.Text
}

// 格式化超链接
func (a TextHref) String() string {
	var dataMiniProgram string
	if a.Data.AppId != "" {
		dataMiniProgram = fmt.Sprintf(`data-miniprogram-appid="%s" data-miniprogram-path="%s"`, a.Data.AppId, a.Data.Path)
	}
	return fmt.Sprintf(`<a href="%s" %s>%s</a>`, a.Path, dataMiniProgram, a.Text)
}

//------------------------------------------------------------------------------

type ContactMgr interface {
	TransferCustomerService(msg ReceiveMessage) RespMessage
	NewLink(msg ReceiveMessage, title, desc, url, thumbUrl string) SendMessage
	NewText(msg ReceiveMessage, text string) SendMessage
	NewImage(msg ReceiveMessage, mediaId string) SendMessage
	// Deprecated
	NewCard(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) SendMessage
	NewMiniProgramPage(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) SendMessage
}

type ContactHook interface {
	AfterSend(message SendMessage, extra interface{})
}

type Contact struct {
}

//------------------------------------------------------------------------------

// 直接转发客服消息
func (*Contact) TransferCustomerService(msg ReceiveMessage) RespMessage {
	return RespMessage{
		ToUserName:   msg.GetFromUserName(),
		FromUserName: msg.GetToUserName(),
		CreateTime:   msg.GetCreateTime(),
		MsgType:      "transfer_customer_service",
	}
}

// 发送链接
func (c *Contact) NewLink(msg ReceiveMessage, title, desc, url, thumbUrl string) SendMessage {
	return SendMessage{
		ToUserName: msg.GetFromUserName(),
		MsgType:    ReturnMsgLink,
		Link: &ReturnLinkMsg{
			title, desc, url, thumbUrl,
		},
	}
}

// 发送文字
func (c *Contact) NewText(msg ReceiveMessage, text string) SendMessage {
	return SendMessage{
		ToUserName: msg.GetFromUserName(),
		MsgType:    ReturnMsgText,
		Text: &ReturnTextMsg{
			Content: text,
		},
	}
}

// 发送图片
func (c *Contact) NewImage(msg ReceiveMessage, mediaId string) SendMessage {
	return SendMessage{
		ToUserName: msg.GetFromUserName(),
		MsgType:    ReturnMsgImage,
		Image: &ReturnImageMsg{
			MediaId: mediaId,
		},
	}
}

// Deprecated
// 发送小程序卡片
func (c *Contact) NewCard(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) SendMessage {
	return c.NewMiniProgramPage(msg, title, appId, pagePath, thumbMediaId)
}

// 发送小程序卡片
func (c *Contact) NewMiniProgramPage(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) SendMessage {
	return SendMessage{
		ToUserName: msg.GetFromUserName(),
		MsgType:    ReturnMsgMiniprogramPage,
		MiniProgramPage: &ReturnWXAppMsg{
			Title:        title,
			AppId:        appId,
			PagePath:     pagePath,
			ThumbMediaId: thumbMediaId,
		},
	}
}

//------------------------------------------------------------------------------

type SendContactMgr interface {
	// Deprecated
	SetContactHook(hook ContactHook)
	PostText(msg ReceiveMessage, text string) (SendMessage, error)
	PostImage(msg ReceiveMessage, mediaId string) (SendMessage, error)
	// Deprecated
	PostCard(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) (SendMessage, error)
	PostMiniProgramPage(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) (SendMessage, error)
	// Deprecated
	Send(body SendMessage, extra interface{}) (err error)
	SendMessage(body SendMessage, extra interface{}) (err error)
}

func (c *Client) SetContactHook(hook ContactHook) {
	c.contactHook = hook
}

func (c *Client) PostText(msg ReceiveMessage, text string) (SendMessage, error) {
	message := c.NewText(msg, text)
	return message, c.Send(message, nil)
}

func (c *Client) PostImage(msg ReceiveMessage, mediaId string) (SendMessage, error) {
	message := c.NewImage(msg, mediaId)
	return message, c.Send(message, nil)
}

func (c *Client) PostCard(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) (SendMessage, error) {
	return c.PostMiniProgramPage(msg, title, appId, pagePath, thumbMediaId)
}

func (c *Client) PostMiniProgramPage(msg ReceiveMessage, title, appId, pagePath, thumbMediaId string) (SendMessage, error) {
	message := c.NewMiniProgramPage(msg, title, appId, pagePath, thumbMediaId)
	return message, c.Send(message, nil)
}

// Deprecated
// 发送客服消息
func (c *Client) Send(body SendMessage, extra interface{}) (err error) {
	return c.SendMessage(body, extra)
}

// 发送客服消息
func (c *Client) SendMessage(body SendMessage, extra interface{}) (err error) {
	defer func() { // 发送完后回调
		if err == nil && c.contactHook != nil {
			c.contactHook.AfterSend(body, extra)
		}
	}()

	postBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", contactURL+c.GetToken(), bytes.NewBuffer(postBody))
	if err = DoRequestJson(req, nil); err != nil {
		return
	}
	return nil
}
