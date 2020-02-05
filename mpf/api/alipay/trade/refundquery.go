/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 11:23
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

// 统一收单交易退款查询
type refundQuery struct {
    alipay.BaseAliPay
    outTradeNo   string // 商户订单号
    tradeNo      string // 支付宝交易号
    outRequestNo string // 退款单号
}

func (rq *refundQuery) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        rq.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (rq *refundQuery) SetTradeNo(tradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, tradeNo)
    if match {
        rq.tradeNo = tradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付宝交易号不合法", nil))
    }
}

func (rq *refundQuery) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        rq.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "退款单号不合法", nil))
    }
}

func (rq *refundQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rq.tradeNo) > 0 {
        rq.BizContent["trade_no"] = rq.tradeNo
    } else if len(rq.outTradeNo) > 0 {
        rq.BizContent["out_trade_no"] = rq.outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号和支付宝交易号不能都为空", nil))
    }
    if len(rq.outRequestNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "退款单号不能为空", nil))
    }
    rq.BizContent["out_request_no"] = rq.outRequestNo

    return rq.GetRequest()
}

func NewRefundQuery(appId string) *refundQuery {
    rq := &refundQuery{alipay.NewBase(appId), "", "", ""}
    rq.SetMethod("alipay.trade.fastpay.refund.query")
    return rq
}
