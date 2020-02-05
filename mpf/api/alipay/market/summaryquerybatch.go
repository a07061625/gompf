/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 18:02
 */
package market

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 门店摘要信息批量查询接口
type summaryQueryBatch struct {
    alipay.BaseAliPay
    opRole     string   // 调用方身份
    queryType  string   // 查询类型
    statusList []string // 门店状态
}

func (sqb *summaryQueryBatch) SetOpRole(opRole string) {
    if (opRole == "ISV") || (opRole == "PROVIDER") {
        sqb.opRole = opRole
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "调用方身份不合法", nil))
    }
}

func (sqb *summaryQueryBatch) SetQueryType(queryType string) {
    switch queryType {
    case "BRAND_RELATION":
        sqb.queryType = queryType
    case "MALL_SELF":
        sqb.queryType = queryType
    case "MALL_RELATION":
        sqb.queryType = queryType
    case "MERCHANT_SELF":
        sqb.queryType = queryType
    case "KB_PROMOTER":
        sqb.queryType = queryType
    default:
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "查询类型不合法", nil))
    }
}

func (sqb *summaryQueryBatch) SetRelatedPartnerId(relatedPartnerId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,16}$`, relatedPartnerId)
    if match {
        sqb.BizContent["related_partner_id"] = relatedPartnerId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "关联商户PID不合法", nil))
    }
}

func (sqb *summaryQueryBatch) SetShopId(shopId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, shopId)
    if match {
        sqb.BizContent["shop_id"] = shopId
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店ID不合法", nil))
    }
}

func (sqb *summaryQueryBatch) SetStatusList(statusList []string) {
    if len(statusList) > 0 {
        sqb.statusList = statusList
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "门店状态不合法", nil))
    }
}

func (sqb *summaryQueryBatch) SetPageNo(pageNo int) {
    if pageNo > 0 {
        sqb.BizContent["page_no"] = pageNo
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "页码不合法", nil))
    }
}

func (sqb *summaryQueryBatch) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 100) {
        sqb.BizContent["page_size"] = pageSize
    } else {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "每页记录数不合法", nil))
    }
}

func (sqb *summaryQueryBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sqb.opRole) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "调用方身份不能为空", nil))
    }
    if len(sqb.queryType) == 0 {
        panic(mperr.NewAliPayMarket(errorcode.AliPayMarketParam, "查询类型不能为空", nil))
    }
    sqb.BizContent["op_role"] = sqb.opRole
    sqb.BizContent["query_type"] = sqb.queryType
    if len(sqb.statusList) > 0 {
        sqb.BizContent["shop_status"] = strings.Join(sqb.statusList, ",")
    }

    return sqb.GetRequest()
}

func NewSummaryQueryBatch(appId string) *summaryQueryBatch {
    sqb := &summaryQueryBatch{alipay.NewBase(appId), "", "", make([]string, 0)}
    sqb.BizContent["page_no"] = 1
    sqb.BizContent["page_size"] = 20
    sqb.SetMethod("alipay.offline.market.shop.summary.batchquery")
    return sqb
}
