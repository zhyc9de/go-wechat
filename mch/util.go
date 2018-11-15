package mch

import (
	"fmt"
	"github.com/zhyc9de/go-wechat"
	"net"
	"time"
)

// LocalIP 获取机器的IP
func LocalIP() string {
	info, _ := net.InterfaceAddrs()
	for _, addr := range info {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}

// 要求32个字符内，只能是数字、大小写字母_-|* 且在同一个商户号下唯一
func NewTradeNo(suffix string) string {
	return fmt.Sprintf("%s%s%s", time.Now().Format("20060102150405"), weUtil.GetRndString(6), suffix)
}
