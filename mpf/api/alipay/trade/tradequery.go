/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 10:32
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

// 统一收单线下交易查询
type tradeQuery struct {
    alipay.BaseAliPay
    outTradeNo string // 商户订单号
    tradeNo    string // 支付宝交易号
}

func (tq *tradeQuery) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        tq.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (tq *tradeQuery) SetTradeNo(tradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, tradeNo)
    if match {
        tq.tradeNo = tradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付宝交易号不合法", nil))
    }
}

func (tq *tradeQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tq.tradeNo) > 0 {
        tq.BizContent["trade_no"] = tq.tradeNo
    } else if len(tq.outTradeNo) > 0 {
        tq.BizContent["out_trade_no"] = tq.outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号和支付宝交易号不能都为空", nil))
    }

    return tq.GetRequest()
}

func NewTradeQuery(appId string) *tradeQuery {
    tq := &tradeQuery{alipay.NewBase(appId), "", ""}
    tq.SetMethod("alipay.trade.query")
    return tq
}
