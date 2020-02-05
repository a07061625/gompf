/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 10:12
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

// 统一收单线下交易预创建
type preCreate struct {
    alipay.BaseAliPay
    outTradeNo  string  // 商户订单号
    totalAmount float32 // 订单总金额,单位为元
    subject     string  // 订单标题
}

func (pc *preCreate) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        pc.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (pc *preCreate) SetTotalAmount(totalAmount float32) {
    if totalAmount > 0 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalAmount), 64)
        pc.totalAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不合法", nil))
    }
}

func (pc *preCreate) SetSubject(subject string) {
    if len(subject) > 0 {
        trueName := []rune(subject)
        pc.subject = string(trueName[:128])
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不合法", nil))
    }
}

func (pc *preCreate) SetTimeoutExpress(timeoutExpress string) {
    if len(timeoutExpress) > 0 {
        pc.BizContent["timeout_express"] = timeoutExpress
    } else {
        delete(pc.BizContent, "timeout_express")
    }
}

func (pc *preCreate) SetBody(body string) {
    if (len(body) > 0) && (len(body) <= 128) {
        pc.BizContent["body"] = body
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商品描述不合法", nil))
    }
}

func (pc *preCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.outTradeNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不能为空", nil))
    }
    if pc.totalAmount <= 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不能为空", nil))
    }
    if len(pc.subject) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不能为空", nil))
    }
    pc.BizContent["out_trade_no"] = pc.outTradeNo
    pc.BizContent["total_amount"] = pc.totalAmount
    pc.BizContent["subject"] = pc.subject

    return pc.GetRequest()
}

func NewPreCreate(appId string) *preCreate {
    conf := alipay.NewConfig().GetAccount(appId)
    pc := &preCreate{alipay.NewBase(appId), "", 0.00, ""}
    pc.totalAmount = 0.00
    pc.BizContent["seller_id"] = conf.GetSellerId()
    pc.SetMethod("alipay.trade.precreate")
    pc.SetUrlNotify(true)
    return pc
}
