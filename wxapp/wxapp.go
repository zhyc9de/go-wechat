package wxapp

import (
	"gitee.com/hzsuoyi/go-wechat.git/util"
	"github.com/json-iterator/go"
)

// 禁止转义html
var json = jsoniter.ConfigFastest

//------------------------------------------------------------------------------

type Client struct {
	*weUtil.Client
}

// 初始化sdk
func NewClient(client *weUtil.Client) *Client {
	return &Client{
		Client: client,
	}
}
