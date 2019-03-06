package weUtil

type (
	TemplateMsg struct {
		ToUser     string          `json:"touser"`
		TemplateId string          `json:"template_id"`
		Data       TemplateMsgData `json:"data"`
	}

	TemplateMsgData map[string]TemplateMsgField

	TemplateMsgField struct {
		Value string `json:"value"`
		Color string `json:"color"`
	}

	// 公众号发送模板消息
	MpTemplate struct {
		TemplateMsg
		Url         string     `json:"url"`
		MiniProgram *WXAppPath `json:"miniprogram,omitempty"`
	}

	// 小程序端发送模板消息
	WxaTemplate struct {
		TemplateMsg
		Page            string `json:"page"`
		FormId          string `json:"form_id"`
		EmphasisKeyword string `json:"emphasis_keyword"` // 放大关键词
	}
)

// 给模板消息加内容
func (t TemplateMsgData) Set(k, v, color string) {
	t[k] = TemplateMsgField{
		v, color,
	}
}
