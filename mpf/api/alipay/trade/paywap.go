/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 10:39
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

// 手机网站支付接口2.0
type payWap struct {
    alipay.BaseAliPay
    outTradeNo  string  // 商户订单号
    totalAmount float32 // 订单总金额,单位为元
    subject     string  // 订单标题
}

func (pw *payWap) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        pw.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (pw *payWap) SetTotalAmount(totalAmount float32) {
    if totalAmount > 0 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalAmount), 64)
        pw.totalAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不合法", nil))
    }
}

func (pw *payWap) SetSubject(subject string) {
    if len(subject) > 0 {
        trueName := []rune(subject)
        pw.subject = string(trueName[:128])
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不合法", nil))
    }
}

func (pw *payWap) SetTimeoutExpress(timeoutExpress string) {
    if len(timeoutExpress) > 0 {
        pw.BizContent["timeout_express"] = timeoutExpress
    } else {
        delete(pw.BizContent, "timeout_express")
    }
}

func (pw *payWap) SetBody(body string) {
    if (len(body) > 0) && (len(body) <= 128) {
        pw.BizContent["body"] = body
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "交易描述不合法", nil))
    }
}

func (pw *payWap) SetGoodsType(goodsType string) {
    if (goodsType == "0") || (goodsType == "1") {
        pw.BizContent["goods_type"] = goodsType
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商品主类型不合法", nil))
    }
}

func (pw *payWap) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := pw.ReqData["return_url"]
    if !ok {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "同步通知地址不能为空", nil))
    }
    if len(pw.outTradeNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不能为空", nil))
    }
    if pw.totalAmount <= 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不能为空", nil))
    }
    if len(pw.subject) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不能为空", nil))
    }
    pw.BizContent["out_trade_no"] = pw.outTradeNo
    pw.BizContent["total_amount"] = pw.totalAmount
    pw.BizContent["subject"] = pw.subject

    return pw.GetRequest()
}

func NewPayWap(appId string) *payWap {
    conf := alipay.NewConfig().GetAccount(appId)
    pw := &payWap{alipay.NewBase(appId), "", 0.00, ""}
    pw.totalAmount = 0.00
    pw.BizContent["seller_id"] = conf.GetSellerId()
    pw.BizContent["product_code"] = "QUICK_WAP_PAY"
    pw.BizContent["goods_type"] = "1"
    pw.SetMethod("alipay.trade.wap.pay")
    pw.SetUrlNotify(true)
    return pw
}
