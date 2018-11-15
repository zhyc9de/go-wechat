package wxapp

import (
	"github.com/json-iterator/go"
	"github.com/zhyc9de/go-wechat"
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
