/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 11:37
 */
package trade

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 统一收单交易退款接口
type tradeRefund struct {
    alipay.BaseAliPay
    outTradeNo   string  // 商户订单号
    tradeNo      string  // 支付宝交易号
    outRequestNo string  // 退款单号
    refundAmount float32 // 退款的金额,该金额不能大于订单金额,单位为元
}

func (tr *tradeRefund) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        tr.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (tr *tradeRefund) SetTradeNo(tradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, tradeNo)
    if match {
        tr.tradeNo = tradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付宝交易号不合法", nil))
    }
}

func (tr *tradeRefund) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        tr.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "退款单号不合法", nil))
    }
}

func (tr *tradeRefund) SetRefundAmount(refundAmount float32) {
    if refundAmount > 0 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", refundAmount), 64)
        tr.refundAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "退款金额不合法", nil))
    }
}

func (tr *tradeRefund) SetRefundReason(refundReason string) {
    if len(refundReason) > 0 {
        trueReason := []rune(refundReason)
        tr.BizContent["refund_reason"] = string(trueReason[:80])
    } else {
        delete(tr.BizContent, "refund_reason")
    }
}

func (tr *tradeRefund) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tr.tradeNo) > 0 {
        tr.BizContent["trade_no"] = tr.tradeNo
    } else if len(tr.outTradeNo) > 0 {
        tr.BizContent["out_trade_no"] = tr.outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号和支付宝交易号不能都为空", nil))
    }
    if len(tr.outRequestNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "退款单号不能为空", nil))
    }
    if tr.refundAmount <= 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "退款金额不能为空", nil))
    }
    tr.BizContent["out_request_no"] = tr.outRequestNo
    tr.BizContent["refund_amount"] = tr.refundAmount

    return tr.GetRequest()
}

func NewTradeRefund(appId string) *tradeRefund {
    tr := &tradeRefund{alipay.NewBase(appId), "", "", "", 0.00}
    tr.refundAmount = 0.00
    tr.SetMethod("alipay.trade.refund")
    return tr
}
