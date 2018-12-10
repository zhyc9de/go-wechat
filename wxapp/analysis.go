package wxapp

import (
	"bytes"
	"github.com/zhyc9de/go-wechat"
	"net/http"
)

const (
	AnalysisDailyRetain       = "getweanalysisappiddailyretaininfo"   // 获取用户访问小程序日留存
	AnalysisMonthlyRetain     = "getweanalysisappidmonthlyretaininfo" // 获取用户访问小程序月留存
	AnalysisWeeklyRetain      = "getweanalysisappidweeklyretaininfo"  // 获取用户访问小程序周留存
	AnalysisDailySummary      = "getweanalysisappiddailysummarytrend" // 获取用户访问小程序数据概况
	AnalysisDailyVisitTrend   = "getweanalysisappiddailyvisittrend"   // 获取用户访问小程序数据日趋势
	AnalysisMonthlyVisitTrend = "getweanalysisappidmonthlyvisittrend" // 获取用户访问小程序数据月趋势
	AnalysisWeeklyVisitTrend  = "getweanalysisappidweeklyvisittrend"  // 获取用户访问小程序数据周趋势
	AnalysisUserPortrait      = "getweanalysisappiduserportrait"      // 获取小程序新增或活跃用户的画像分布数据
	AnalysisVisitDistribution = "getweanalysisappidvisitdistribution" // 获取用户小程序访问分布数据
	AnalysisVisitPage         = "getweanalysisappidvisitpage"         // 访问页面
)

type AnalysisDateRange struct {
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
}

func (client *Client) GetAnalysisData(action, begin, end string) (rb []byte, err error) {
	u := "https://api.weixin.qq.com/datacube/" + action + "?access_token=" + client.GetToken()
	body, _ := json.Marshal(AnalysisDateRange{
		BeginDate: begin,
		EndDate:   end,
	})
	req, _ := http.NewRequest("POST", u, bytes.NewReader(body))
	return weUtil.DoRequest(req)
}
