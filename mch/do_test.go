package mch

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestCouponMap_MarshalXML(t *testing.T) {
	data := make(CouponMap)
	data[0] = Coupon{
		CouponFee: 1,
	}
	data[1] = Coupon{
		CouponFee: 1,
	}
	if s, err := xml.Marshal(data); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(s))
	}
}
