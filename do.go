package weUtil

import "time"

type Token struct {
	Token   string
	Expires int64
}

func (token Token) IsNotExpire() bool {
	return token.Expires > time.Now().Unix()
}

//--------------------------------------------------------------------

// 微信收货地址
type Address struct {
	UserName     string `json:"userName"`     // 收货人姓名
	PostalCode   string `json:"postalCode"`   // 邮编
	ProvinceName string `json:"provinceName"` // 国标收货地址第一级地址
	CityName     string `json:"cityName"`     // 国标收货地址第一级地址
	CountyName   string `json:"countyName"`   // 国标收货地址第一级地址
	DetailInfo   string `json:"detailInfo"`   // 详细收货地址信息
	NationalCode string `json:"nationalCode"` // 收货地址国家码
	TelNumber    string `json:"telNumber"`    // 收货人手机号码
}
