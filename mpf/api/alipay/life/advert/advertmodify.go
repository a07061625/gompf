/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:10
 */
package advert

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 生活号广告位修改接口
type advertModify struct {
    alipay.BaseAliPay
    advertId string                   // 广告ID
    items    []map[string]interface{} // 广告内容列表
}

func (am *advertModify) SetAdvertId(advertId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,20}$`, advertId)
    if match {
        am.advertId = advertId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告ID不合法", nil))
    }
}

func (am *advertModify) SetItems(items []map[string]interface{}) {
    if len(items) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告内容列表不合法", nil))
    } else if len(items) > 3 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告内容列表超过限制", nil))
    }
    am.items = items
}

func (am *advertModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(am.advertId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告ID不能为空", nil))
    }
    if len(am.items) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告内容列表不能为空", nil))
    }
    am.BizContent["advert_id"] = am.advertId
    am.BizContent["advert_items"] = am.items

    return am.GetRequest()
}

func NewAdvertModify(appId string) *advertModify {
    am := &advertModify{alipay.NewBase(appId), "", make([]map[string]interface{}, 0)}
    am.SetMethod("alipay.open.public.advert.modify")
    return am
}
