package mch

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

type Client struct {
	MchId      string       // 商户Id
	MchKey     string       // 商户key
	NotifyURL  string       // 支付回调地址
	RefundURL  string       // 退款回调地址
	http       *http.Client // 微信支付http客户端
	certConfig *tls.Config  // 证书设置
	debug      bool         // 是否开启沙盒模式
}

func InitClient(mchId, mchKey, notifyURL, refundURL string) *Client {
	return &Client{
		MchId:     mchId,
		MchKey:    mchKey,
		NotifyURL: notifyURL,
		RefundURL: refundURL,
		http:      weUtil.GetHttpClient(),
	}
}

// 设置开发环境
func (client *Client) Debug() {
	client.debug = true
}

// 设置支付证书
func (client *Client) SetCert(certFile, keyFile string) error {
	if certFile == "" || keyFile == "" {
		return fmt.Errorf("check file path, certFile=%s, keyFile=%s", certFile, keyFile)
	}
	// 读取证书
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}
	// parse cert
	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}
	// 设置client
	client.certConfig = &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
	}
	client.http = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       client.certConfig,
		},
	}
	return nil
}

func (client *Client) DoRequest(req *http.Request) (rb []byte, err error) {
	resp, err := client.http.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
