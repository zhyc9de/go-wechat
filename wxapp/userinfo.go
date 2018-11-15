package wxapp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/zhyc9de/go-wechat"
)

type UserInfo struct {
	OpenId    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int64  `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
	UnionId   string `json:"unionId"`
	Watermark struct {
		AppId     string `json:"appId"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

// 前端传入的用户信息
type EncData struct {
	//RawData       string `json:"rawData"`
	//Signature     string `json:"signature"`
	EncryptedData string `json:"encryptedData"`
	Iv            string `json:"iv"`
}

// GetSession 需要使用到的结构体
type AuthResp struct {
	weUtil.ErrResp
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

const jsCode2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

// 小程序拿code换openid
func (client *Client) GetSession(code string) (authResp AuthResp, err error) {
	u := fmt.Sprintf(jsCode2sessionURL, client.GetAppId(), client.GetAppSecret(), code)

	req, _ := http.NewRequest("GET", u, nil)
	err = weUtil.DoRequestJson(req, &authResp)
	return
}

// 解密前端传入的小程序用户信息
func DecData(sessionKey string, info *EncData, decData interface{}) (err error) {
	if info == nil {
		err = fmt.Errorf("encData is nil")
		return
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return
	}
	cipherText, err := base64.StdEncoding.DecodeString(info.EncryptedData)
	if err != nil {
		return
	} else if len(cipherText)%8 != 0 {
		err = errors.New("cipherText cannot exact division by 8")
		return
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	iv, err := base64.StdEncoding.DecodeString(info.Iv)
	if err != nil {
		return
	}
	if len(iv) == 0 || len(iv)%8 != 0 {
		err = fmt.Errorf("invalid IV length, iv=%s", info.Iv)
		return
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	trailingPadding := uint(cipherText[len(cipherText)-1])
	if trailingPadding > 16 {
		err = errors.New("invalid trailing padding")
		return
	}
	data := cipherText[:len(cipherText)-int(trailingPadding)]

	if err = json.Unmarshal(data, decData); err != nil {
		return
	}
	return
}
