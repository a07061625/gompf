/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 22:01
 */
package advert

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 生活号广告位添加接口
type advertCreate struct {
    alipay.BaseAliPay
    items []map[string]interface{} // 广告内容列表
}

func (ac *advertCreate) SetItems(items []map[string]interface{}) {
    if len(items) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告内容列表不合法", nil))
    } else if len(items) > 3 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告内容列表超过限制", nil))
    }
    ac.items = items
}

func (ac *advertCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ac.items) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "广告内容列表不能为空", nil))
    }
    ac.BizContent["advert_items"] = ac.items

    return ac.GetRequest()
}

func NewAdvertCreate(appId string) *advertCreate {
    ac := &advertCreate{alipay.NewBase(appId), make([]map[string]interface{}, 0)}
    ac.SetMethod("alipay.open.public.advert.create")
    return ac
}
