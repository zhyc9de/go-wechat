package weUtil

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
