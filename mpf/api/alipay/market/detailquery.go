/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 18:27
 */
package market

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询单个门店信息接口
type detailQuery struct {
    alipay.BaseAliPay
    shopId string // 门店ID
}

func (dq *detailQuery) SetShopId(shopId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, shopId)
    if match {
        dq.shopId = shopId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店ID不合法", nil))
    }
}

func (dq *detailQuery) SetOpRole(opRole string) {
    if (opRole == "MERCHANT") || (opRole == "PROVIDER") {
        dq.BizContent["op_role"] = opRole
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "操作人角色不合法", nil))
    }
}

func (dq *detailQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dq.shopId) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店ID不能为空", nil))
    }
    dq.BizContent["shop_id"] = dq.shopId

    return dq.GetRequest()
}

func NewDetailQuery(appId string) *detailQuery {
    dq := &detailQuery{alipay.NewBase(appId), ""}
    dq.SetMethod("alipay.offline.market.shop.querydetail")
    return dq
}
