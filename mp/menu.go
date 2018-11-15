package mp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zhyc9de/go-wechat"
)

const (
	MenuTypClick           = "click"
	MenuTypView            = "view"
	MenuTypScancodePush    = "scancode_push"
	MenuTypScancodeWaitmsg = "scancode_waitmsg"
	MenuTypPicSysphoto     = "pic_sysphoto"
	MenuTypPicPhotoOrAlbum = "pic_photo_or_album"
	MenuTypPicWeixin       = "pic_weixin"
	MenuTypLocationSelect  = "location_select"
	MenuTypMediaId         = "media_id"
	MenuTypViewLimited     = "view_limited"
)

type (
	Button struct {
		Type string `json:"type"`
		Name string `json:"name"`
		// 点击事件
		Key string `json:"key,omitempty"`
		// 网页
		Url string `json:"url,omitempty"`
		// 小程序
		AppId    string `json:"appid,omitempty"`
		PagePath string `json:"pagepath,omitempty"`
		// 子菜单
		SubButton *[]Button `json:"sub_button,omitempty"`
	}

	Matchrule struct {
		TagId              int64  `json:"tag_id,omitempty"`
		Sex                int64  `json:"sex,omitempty"`
		Country            string `json:"country,omitempty"`
		Province           string `json:"province,omitempty"`
		City               string `json:"city,omitempty"`
		ClientPlatformType int64  `json:"client_platform_type,omitempty"`
		Language           string `json:"language,omitempty"`
	}

	Menu struct {
		weUtil.ErrResp
		Button    []Button   `json:"button"`
		Matchrule *Matchrule `json:"matchrule,omitempty"`
	}
)

//------------------------------------------------------------------------------

const createMenuURL = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
const addConditionalMenuURL = "https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token=%s"
const tryMenuURL = "https://api.weixin.qq.com/cgi-bin/menu/trymatch?access_token=%s"
const deleteMenuURL = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s"

// 读取文件更改菜单
func (client *Client) CreateMenuWithFile(filePath string) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	u := fmt.Sprintf(createMenuURL, client.GetToken())
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(b))
	return weUtil.DoRequestJson(req, nil)
}

// 更新菜单
func (client *Client) CreateMenu(menu Menu) error {
	menuBytes, _ := json.Marshal(menu)

	u := fmt.Sprintf(createMenuURL, client.GetToken())
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(menuBytes))
	return weUtil.DoRequestJson(req, nil)
}

// 创建个性化菜单
func (client *Client) AddConditionalMenu(menu Menu) error {
	menuBytes, _ := json.Marshal(menu)

	u := fmt.Sprintf(addConditionalMenuURL, client.GetToken())
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(menuBytes))
	return weUtil.DoRequestJson(req, nil)
}

// 测试是否匹配
func (client *Client) TryMenu(openId string) (menu Menu, err error) {
	data, _ := json.Marshal(struct {
		UserId string `json:"user_id"`
	}{
		UserId: openId,
	})

	u := fmt.Sprintf(tryMenuURL, client.GetToken())
	req, _ := http.NewRequest("POST", u, bytes.NewReader(data))
	err = weUtil.DoRequestJson(req, &menu)
	return
}

// 删除个性化菜单
func (client *Client) DeleteMenu() error {
	u := fmt.Sprintf(deleteMenuURL, client.GetToken())
	req, _ := http.NewRequest("GET", u, nil)
	return weUtil.DoRequestJson(req, nil)
}
