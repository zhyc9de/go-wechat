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

func (c *Client) GetToken() string {
	token := c.TokenMgr.GetToken()
	c.Logger.Infof("wx-auth accessToken, appId=%s, token=%s", c.GetAppId(), token)
	return token
}
