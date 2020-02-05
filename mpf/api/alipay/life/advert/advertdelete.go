/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:06
 */
package advert

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 生活号广告位删除接口
type advertDelete struct {
    alipay.BaseAliPay
    advertId string // 广告ID
}

func (ad *advertDelete) SetAdvertId(advertId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,20}$`, advertId)
    if match {
        ad.advertId = advertId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告ID不合法", nil))
    }
}

func (ad *advertDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ad.advertId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告ID不能为空", nil))
    }
    ad.BizContent["advert_id"] = ad.advertId

    return ad.GetRequest()
}

func NewAdvertDelete(appId string) *advertDelete {
    ad := &advertDelete{alipay.NewBase(appId), ""}
    ad.SetMethod("alipay.open.public.advert.delete")
    return ad
}
