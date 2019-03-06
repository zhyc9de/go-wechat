package weUtil

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
)

type Config struct {
	AppId     string `yaml:"appId" json:"appId"`
	AppSecret string `yaml:"appSecret" json:"appSecret"`
	Username  string `yaml:"username" json:"username"`
	Name      string `yaml:"name" json:"name"`
	IsMp      bool   `yaml:"isMp" json:"isMp"`
}

type TokenMgr interface {
	// Deprecated
	GetAppId() string // 获取appId
	// Deprecated
	GetAppSecret() string // 获取appSecret

	GetConfig() Config                                    // 获取配置
	Set(k, v string, expire int64) error                  // 手动设置session
	Get(k string) string                                  // 获取session
	GetOrNewToken() (token string, isNew bool, err error) // 获取accessToken
}

const KeyToken = "token"
const KeyJSApi = "jsApi"
const tokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

// 传回access_token，有效期提前10分钟
func fetchToken(appId, appSecret string) (token Token, err error) {
	u := fmt.Sprintf(tokenURL, appId, appSecret)

	req, _ := http.NewRequest("GET", u, nil)
	r := new(struct {
		ErrResp
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	})
	if err = DoRequestJson(req, r); err != nil {
		return
	}

	return Token{
		Token:   r.AccessToken,
		Expires: r.ExpiresIn - 600, // 提前10分钟过期
	}, nil
}

//------------------------------------------------------------------------------

// access_token 维护在内存中的基础版
type WxTokenMgr struct {
	Config
	session map[string]atomic.Value
}

func NewWxTokenMgr(appId, appSecret string) *WxTokenMgr {
	return &WxTokenMgr{
		Config: Config{
			AppId:     appId,
			AppSecret: appSecret,
		},
		session: make(map[string]atomic.Value),
	}
}

func NewWxTokenMgrByConfig(c Config) *WxTokenMgr {
	return &WxTokenMgr{
		Config:  c,
		session: make(map[string]atomic.Value),
	}
}

func (mgr *WxTokenMgr) GetAppId() string {
	return mgr.AppId
}

func (mgr *WxTokenMgr) GetAppSecret() string {
	return mgr.AppSecret
}

func (mgr *WxTokenMgr) GetConfig() Config {
	return mgr.Config
}

func (mgr *WxTokenMgr) GetOrNewToken() (token string, isNew bool, err error) {
	token = mgr.Get(KeyToken)
	if token != "" {
		return
	}

	newToken, err := fetchToken(mgr.AppId, mgr.AppSecret)
	if err != nil {
		return
	}
	isNew = true
	token = newToken.Token
	err = mgr.Set(KeyToken, newToken.Token, newToken.Expires)
	return
}

func (mgr *WxTokenMgr) Set(k, v string, expire int64) error {
	if expire == 0 {
		expire = 99999999999
	}
	if s, ok := mgr.session[k]; ok {
		s.Store(&Token{
			Token:   v,
			Expires: time.Now().Unix() + expire,
		})
	} else {
		var s atomic.Value
		s.Store(&Token{
			Token:   v,
			Expires: time.Now().Unix() + expire,
		})
		mgr.session[k] = s
	}
	return nil
}

func (mgr *WxTokenMgr) Get(k string) string {
	if token, ok := mgr.session[KeyToken]; ok {
		token := token.Load()
		if token != nil {
			token := token.(*Token)
			if token.IsNotExpire() {
				return token.Token
			}
		}
	}

	return ""
}

//------------------------------------------------------------------------------

// 基于 go-redis 维护access_token
type RTokenMgr struct {
	Config
	session *redis.Client
}

func NewRTokenMgr(appId, appSecret string, session *redis.Client) *RTokenMgr {
	return &RTokenMgr{
		Config: Config{
			AppId:     appId,
			AppSecret: appSecret,
		},
		session: session,
	}
}

func NewRTokenMgrByConfig(c Config, session *redis.Client) *RTokenMgr {
	return &RTokenMgr{
		Config:  c,
		session: session,
	}
}

func (mgr *RTokenMgr) GetAppId() string {
	return mgr.AppId
}

func (mgr *RTokenMgr) GetAppSecret() string {
	return mgr.AppSecret
}

func (mgr *RTokenMgr) GetConfig() Config {
	return mgr.Config
}

func (mgr *RTokenMgr) GetOrNewToken() (token string, isNew bool, err error) {
	tokenKey := fmt.Sprintf("%s:%s", KeyToken, mgr.AppId)
	token = mgr.session.Get(tokenKey).Val()
	if token != "" {
		return
	}

	newToken, err := fetchToken(mgr.AppId, mgr.AppSecret)
	if err != nil {
		return
	}
	// redis 就依赖expires
	isNew = true
	token = newToken.Token
	err = mgr.session.Set(tokenKey, newToken.Token, time.Duration(newToken.Expires)*time.Second).Err()
	return
}

func (mgr *RTokenMgr) Set(k, v string, expire int64) error {
	tokenKey := fmt.Sprintf("%s:%s", k, mgr.AppId)
	return mgr.session.Set(tokenKey, v, time.Duration(expire)*time.Second).Err()
}

func (mgr *RTokenMgr) Get(k string) string {
	tokenKey := fmt.Sprintf("%s:%s", k, mgr.AppId)
	return mgr.session.Get(tokenKey).Val()
}
