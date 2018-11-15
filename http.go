package weUtil

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/json-iterator/go"
)

type ErrResp struct {
	ErrCode int64  `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (e ErrResp) Error() string {
	return fmt.Sprintf("errCode=%d, errMsg=%s", e.ErrCode, e.ErrMsg)
}

//------------------------------------------------------------------------------

// 只是单纯的去掉了proxy
var netTransport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

func GetHttpClient() *http.Client {
	return netClient
}

func DoRequest(req *http.Request) (rb []byte, err error) {
	resp, err := netClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func DoRequestJson(req *http.Request, v interface{}) (err error) {
	rb, err := DoRequest(req)
	if err != nil {
		return
	}
	// 先反序列化到resp
	var resp ErrResp
	resp.ErrCode = jsoniter.Get(rb, "errcode").ToInt64()
	resp.ErrMsg = jsoniter.Get(rb, "errmsg").ToString()
	if resp.ErrCode != 0 {
		return resp
	}

	if v == nil {
		return
	}
	err = json.Unmarshal(rb, v)
	if err != nil {
		err = fmt.Errorf("parse json err=%s, body=%s", err.Error(), string(rb))
	}
	return
}
