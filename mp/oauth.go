package mp

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/zhyc9de/go-wechat"
)

type (
	OauthToken struct {
		weUtil.ErrResp
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		Scope        string `json:"scope"`
	}

	OauthUserInfo struct {
		weUtil.ErrResp
		OpenId     string   `json:"openid"`
		Nickname   string   `json:"nickname"`
		Sex        int64    `json:"sex"`
		Province   string   `json:"province"`
		Country    string   `json:"country"`
		HeadImgUrl string   `json:"headimgurl"`
		Privilege  []string `json:"privilege"`
		UnionId    string   `json:"unionid,omitempty"`
	}
)

// 获取重定向地址
func (client *Client) RedirectOauth(callbackURL, state string) string {
	if state != "" {
		state = fmt.Sprintf("&state=%s", state)
	}
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo%s#wechat_redirect",
		client.GetAppId(), url.QueryEscape(callbackURL), state, // 这里就在里边做urlEncode了
	)
}

// oauth的token
func (client *Client) GetOauthToken(code string) (token OauthToken, err error) {
	u := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		client.GetAppId(), client.GetAppSecret(), code)
	req, _ := http.NewRequest("GET", u, nil)
	err = weUtil.DoRequestJson(req, &token)
	return
}

// oauth获取用户具体信息
func (client *Client) GetUserInfoByOauth(accessToken, openId, lang string) (userInfo OauthUserInfo, err error) {
	u := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=%s", accessToken, openId, lang)
	req, _ := http.NewRequest("GET", u, nil)
	err = weUtil.DoRequestJson(req, &userInfo)
	return
}
