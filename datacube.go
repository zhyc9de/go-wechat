package weUtil

import (
	"bytes"
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

var AnalysisActions = []string{AnalysisDailyRetain, AnalysisMonthlyRetain, AnalysisWeeklyRetain, AnalysisDailySummary,
	AnalysisDailyVisitTrend, AnalysisMonthlyVisitTrend, AnalysisWeeklyVisitTrend, AnalysisUserPortrait,
	AnalysisVisitDistribution, AnalysisVisitPage}

// 获取小程序或者公众号数据
type AnalysisArgs struct {
	BeginDate string `json:"begin_date"`
	EndDate   string `json:"end_date"`
}

func (c *Client) GetAnalysisData(action, begin, end string) (rb []byte, err error) {
	u := "https://api.weixin.qq.com/datacube/" + action + "?access_token=" + c.GetToken()
	body, _ := json.Marshal(AnalysisArgs{
		BeginDate: begin,
		EndDate:   end,
	})
	req, _ := http.NewRequest("POST", u, bytes.NewReader(body))
	return DoRequest(req)
}

// 判断data cube action是否是获取小程序数据
func IsDataCudeForWxa(action string) bool {
	for i := range AnalysisActions {
		if action == AnalysisActions[i] {
			return true
		}
	}
	return false
}
