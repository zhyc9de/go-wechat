package mch

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"gitee.com/hzsuoyi/go-wechat.git/util"
)

var omitemptyReplacer = strings.NewReplacer(",omitempty", "")

// 订单签名
func (client *Client) SignXml(trade interface{}) string {
	v := url.Values{}

	tradeTyp := reflect.TypeOf(trade)
	tradeValue := reflect.ValueOf(trade)
	for i := 0; i < tradeTyp.NumField(); i++ {
		field := tradeTyp.Field(i)
		if field.Tag.Get("xml") == "sign" ||
			field.Type.String() == "xml.Name" { // 如果是这个字段名就直接跳过
			continue
		}
		//fmt.Println(field.Tag.Get("xml"), tradeValue.Field(i).Kind())

		if tradeValue.Field(i).CanInterface() {
			var value string
			switch tradeValue.Field(i).Kind() {
			case reflect.Struct:
				if cdata, ok := tradeValue.Field(i).Interface().(CDATA); ok {
					value = cdata.Text
				} else if resp, ok := tradeValue.Field(i).Interface().(CallbackResponse); ok {
					if resp.ReturnCode.Text != "" {
						v.Set("return_code", resp.ReturnCode.Text)
					}
					if resp.ReturnMsg.Text != "" {
						v.Set("return_msg", resp.ReturnMsg.Text)
					}
					if resp.ResultCode.Text != "" {
						v.Set("result_code", resp.ResultCode.Text)
					}
					if resp.ErrCode.Text != "" {
						v.Set("err_code", resp.ErrCode.Text)
					}
					if resp.ErrCodeDes.Text != "" {
						v.Set("err_code_des", resp.ErrCodeDes.Text)
					}
					continue
				}
			case reflect.Int64:
				value = strconv.FormatInt(tradeValue.Field(i).Interface().(int64), 10)

			case reflect.Map: // 这个应该是Coupon
				if tradeCoupon, ok := tradeValue.Field(i).Interface().(CouponMap); ok {
					for idx, coupon := range tradeCoupon {
						if coupon.CouponFee != 0 {
							v.Set(fmt.Sprintf("coupon_fee_%d", idx), strconv.FormatInt(coupon.CouponFee, 10))
						}
						if coupon.CouponId.Text != "" {
							v.Set(fmt.Sprintf("coupon_id_%d", idx), coupon.CouponId.Text)
						}
						if coupon.CouponType.Text != "" {
							v.Set(fmt.Sprintf("coupon_type_%d", idx), coupon.CouponType.Text)
						}
					}
					continue
				} else { // 如果有其他的也应该报错
					fmt.Printf("cannot parse map, tagName=%s, value=%#v\n", field.Tag.Get("xml"), tradeValue.Field(i).Interface())
				}

			default: // 报错
				fmt.Printf("cannot parse tagName=%s, value=%#v\n", field.Tag.Get("xml"), tradeValue.Field(i).Interface())
			}
			// TODO 必填值如果是空，报错？
			if value == "0" || value == "" {
				continue
			}
			v.Set(omitemptyReplacer.Replace(field.Tag.Get("xml")), value)
		}
	}
	query, _ := url.QueryUnescape(v.Encode())
	s := fmt.Sprintf("%s&key=%s", query, client.MchKey)
	// md5后，大写
	return strings.ToUpper(weUtil.Md5(s))
}

// 网页端的签名
func (client *Client) SignWCPay(trade interface{}) string {
	v := url.Values{}

	tradeTyp := reflect.TypeOf(trade)
	tradeValue := reflect.ValueOf(trade)
	for i := 0; i < tradeTyp.NumField(); i++ {
		field := tradeTyp.Field(i)
		if field.Tag.Get("json") == "paySign" { // 如果是这个字段名就直接跳过
			continue
		}

		if tradeValue.Field(i).CanInterface() {
			var value string
			switch tradeValue.Field(i).Kind() {
			case reflect.String:
				value = tradeValue.Field(i).String()
			case reflect.Int64:
				value = strconv.FormatInt(tradeValue.Field(i).Interface().(int64), 10)
			}
			if value == "0" || value == "" {
				continue
			}
			v.Set(omitemptyReplacer.Replace(field.Tag.Get("json")), value)
		}
	}
	query, _ := url.QueryUnescape(v.Encode())
	s := fmt.Sprintf("%s&key=%s", query, client.MchKey)
	// md5后，大写
	return strings.ToUpper(weUtil.Md5(s))
}
