/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 11:45
 */
package trade

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 统一收单交易结算接口
type orderSettle struct {
    alipay.BaseAliPay
    outRequestNo  string                 // 结算请求流水号
    tradeNo       string                 // 订单号
    royaltyParams map[string]interface{} // 分账明细信息
}

func (os *orderSettle) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        os.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "结算请求流水号不合法", nil))
    }
}

func (os *orderSettle) SetTradeNo(tradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, tradeNo)
    if match {
        os.tradeNo = tradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单号不合法", nil))
    }
}

func (os *orderSettle) SetRoyaltyParams(royaltyParams map[string]interface{}) {
    if len(royaltyParams) > 0 {
        os.royaltyParams = royaltyParams
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "分账明细信息不合法", nil))
    }
}

func (os *orderSettle) SetOperatorId(operatorId string) {
    if len(operatorId) > 0 {
        os.BizContent["operator_id"] = operatorId
    } else {
        delete(os.BizContent, "operator_id")
    }
}

func (os *orderSettle) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(os.outRequestNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "结算请求流水号不能为空", nil))
    }
    if len(os.tradeNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单号不能为空", nil))
    }
    if len(os.royaltyParams) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "分账明细信息不能为空", nil))
    }
    os.BizContent["out_request_no"] = os.outRequestNo
    os.BizContent["trade_no"] = os.tradeNo
    os.BizContent["royalty_parameters"] = os.royaltyParams

    return os.GetRequest()
}

func NewOrderSettle(appId string) *orderSettle {
    os := &orderSettle{alipay.NewBase(appId), "", "", make(map[string]interface{})}
    os.SetMethod("alipay.trade.order.settle")
    return os
}
