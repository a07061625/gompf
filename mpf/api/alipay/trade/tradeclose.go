/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 11:30
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

// 统一收单交易关闭接口
type tradeClose struct {
    alipay.BaseAliPay
    outTradeNo string // 商户订单号
    tradeNo    string // 支付宝交易号
}

func (tc *tradeClose) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        tc.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (tc *tradeClose) SetTradeNo(tradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, tradeNo)
    if match {
        tc.tradeNo = tradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付宝交易号不合法", nil))
    }
}

func (tc *tradeClose) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tc.tradeNo) > 0 {
        tc.BizContent["trade_no"] = tc.tradeNo
    } else if len(tc.outTradeNo) > 0 {
        tc.BizContent["out_trade_no"] = tc.outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号和支付宝交易号不能都为空", nil))
    }

    return tc.GetRequest()
}

func NewTradeClose(appId string) *tradeClose {
    tc := &tradeClose{alipay.NewBase(appId), "", ""}
    tc.SetMethod("alipay.trade.close")
    tc.SetUrlNotify(true)
    return tc
}
