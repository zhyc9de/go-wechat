package mp

import (
	"bytes"
	"fmt"
	"net/http"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

const userInfoURL = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=%s"

// 用户信息
type UserInfo struct {
	weUtil.ErrResp
	Subscribe      int64   `json:"subscribe"`
	OpenId         string  `json:"openid"`
	Nickname       string  `json:"nickname"`
	Sex            int64   `json:"sex"`
	Language       string  `json:"language"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	Country        string  `json:"country"`
	HeadImgUrl     string  `json:"headimgurl"`
	SubscribeTime  int64   `json:"subscribe_time"`
	UnionId        string  `json:"unionid,omitempty"`
	Remark         string  `json:"remark"`
	GroupId        int64   `json:"groupid"`
	TagIdList      []int64 `json:"tagid_list"`
	SubscribeScene string  `json:"subscribe_scene"`
	QrScene        int64   `json:"qr_scene"`
	QrSceneStr     string  `json:"qr_scene_str"`
}

// 获取用户信息
func (client *Client) GetUserInfo(openId, lang string) (userInfo UserInfo, err error) {
	u := fmt.Sprintf(userInfoURL, client.GetToken(), openId, lang)

	req, _ := http.NewRequest("GET", u, nil)
	err = weUtil.DoRequestJson(req, &userInfo)
	return
}

//------------------------------------------------------------------------------

const createTagURL = "https://api.weixin.qq.com/cgi-bin/tags/create?access_token=%s"
const getTagURL = "https://api.weixin.qq.com/cgi-bin/tags/get?access_token=%s"
const batchTaggingURL = "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=%s"
const batchUnTaggingURL = "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=%s"

type UserTag struct {
	Id    int64  `json:"id,omitempty"`
	Name  string `json:"name"`
	Count int64  `json:"count,omitempty"`
}

func (client *Client) GetOrCreateTagId(name string) (int64, error) {
	tags, _ := client.GetTag()
	for _, t := range tags {
		if t.Name == name {
			return t.Id, nil
		}
	}
	if tag, err := client.CreateTag(name); err != nil {
		return 0, err
	} else {
		return tag.Id, nil
	}
}

// 创建标签
func (client *Client) CreateTag(tag string) (newTag UserTag, err error) {
	u := fmt.Sprintf(createTagURL, client.GetToken())

	data, _ := json.Marshal(struct {
		Tag UserTag `json:"tag"`
	}{
		Tag: UserTag{
			Name: tag,
		},
	})
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))

	err = weUtil.DoRequestJson(req, &newTag)
	return
}

// 获取标签列表
func (client *Client) GetTag() (tags []UserTag, err error) {
	u := fmt.Sprintf(getTagURL, client.GetToken())
	req, _ := http.NewRequest("GET", u, nil)

	var tagList struct {
		Tags []UserTag `json:"tags"`
	}
	if err = weUtil.DoRequestJson(req, &tagList); err != nil {
		return
	}
	tags = tagList.Tags
	return
}

// 批量打标签
func (client *Client) BatchTagging(openIdList []string, tagId int64) (err error) {
	u := fmt.Sprintf(batchTaggingURL, client.GetToken())

	data, _ := json.Marshal(struct {
		OpenIdList []string `json:"openid_list"`
		TagId      int64    `json:"tagid"`
	}{
		OpenIdList: openIdList,
		TagId:      tagId,
	})
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))

	return weUtil.DoRequestJson(req, nil)
}

// 批量取消标签
func (client *Client) BatchUnTagging(openIdList []string, tagId int64) (err error) {
	u := fmt.Sprintf(batchUnTaggingURL, client.GetToken())

	data, _ := json.Marshal(struct {
		OpenIdList []string `json:"openid_list"`
		TagId      int64    `json:"tagid"`
	}{
		OpenIdList: openIdList,
		TagId:      tagId,
	})
	req, _ := http.NewRequest("POST", u, bytes.NewBuffer(data))

	return weUtil.DoRequestJson(req, nil)
}
