/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 12:24
 */
package trade

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 统一收单交易支付接口
type tradePay struct {
    alipay.BaseAliPay
    outTradeNo         string  // 商户订单号
    authCode           string  // 支付授权码
    totalAmount        float32 // 订单总金额,单位为元
    discountableAmount float32 // 可打折金额,单位为元
    subject            string  // 订单标题
}

func (tp *tradePay) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        tp.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (tp *tradePay) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(project.RegexDigit, authCode)
    if match {
        tp.authCode = authCode
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付授权码不合法", nil))
    }
}

func (tp *tradePay) SetTotalAmount(totalAmount float32) {
    if (totalAmount > 0) && (totalAmount <= 100000000) {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalAmount), 64)
        tp.totalAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不合法", nil))
    }
}

func (tp *tradePay) SetDiscountableAmount(discountableAmount float32) {
    if (discountableAmount > 0) && (discountableAmount <= 100000000) {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", discountableAmount), 64)
        tp.discountableAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "可打折金额不合法", nil))
    }
}

func (tp *tradePay) SetSubject(subject string) {
    if len(subject) > 0 {
        trueName := []rune(subject)
        tp.subject = string(trueName[:128])
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不合法", nil))
    }
}

func (tp *tradePay) SetScene(scene string) {
    if (scene == "bar_code") || (scene == "wave_code") {
        tp.BizContent["scene"] = scene
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付场景不合法", nil))
    }
}

func (tp *tradePay) SetProductCode(productCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productCode)
    if match {
        tp.BizContent["product_code"] = productCode
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "销售产品码不合法", nil))
    }
}

func (tp *tradePay) SetBuyerId(buyerId string) {
    match, _ := regexp.MatchString(project.RegexDigit, buyerId)
    if match {
        tp.BizContent["buyer_id"] = buyerId
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "买家支付宝用户ID不合法", nil))
    }
}

func (tp *tradePay) SetCurrencyTrans(currencyTrans string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, currencyTrans)
    if match {
        tp.BizContent["trans_currency"] = currencyTrans
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "标价币种不合法", nil))
    }
}

func (tp *tradePay) SetCurrencySettle(currencySettle string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, currencySettle)
    if match {
        tp.BizContent["settle_currency"] = currencySettle
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "结算币种不合法", nil))
    }
}

func (tp *tradePay) SetBody(body string) {
    if (len(body) > 0) && (len(body) <= 128) {
        tp.BizContent["body"] = body
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商品描述不合法", nil))
    }
}

func (tp *tradePay) SetGoodsDetail(goodsDetail map[string]interface{}) {
    if len(goodsDetail) > 0 {
        tp.BizContent["goods_detail"] = goodsDetail
    } else {
        delete(tp.BizContent, "goods_detail")
    }
}

func (tp *tradePay) SetOperatorId(operatorId string) {
    if len(operatorId) > 0 {
        tp.BizContent["operator_id"] = operatorId
    } else {
        delete(tp.BizContent, "operator_id")
    }
}

func (tp *tradePay) SetStoreId(storeId string) {
    if len(storeId) > 0 {
        tp.BizContent["store_id"] = storeId
    } else {
        delete(tp.BizContent, "store_id")
    }
}

func (tp *tradePay) SetTerminalId(terminalId string) {
    if len(terminalId) > 0 {
        tp.BizContent["terminal_id"] = terminalId
    } else {
        delete(tp.BizContent, "terminal_id")
    }
}

func (tp *tradePay) SetTerminalParams(terminalParams map[string]interface{}) {
    if len(terminalParams) > 0 {
        tp.BizContent["terminal_params"] = mpf.JSONMarshal(terminalParams)
    } else {
        delete(tp.BizContent, "terminal_params")
    }
}

func (tp *tradePay) SetExtendParams(extendParams map[string]interface{}) {
    if len(extendParams) > 0 {
        tp.BizContent["extend_params"] = extendParams
    } else {
        delete(tp.BizContent, "extend_params")
    }
}

func (tp *tradePay) SetTimeoutExpress(timeoutExpress string) {
    if len(timeoutExpress) > 0 {
        tp.BizContent["timeout_express"] = timeoutExpress
    } else {
        delete(tp.BizContent, "timeout_express")
    }
}

func (tp *tradePay) SetAuthConfirmMode(authConfirmMode string) {
    if (authConfirmMode == "COMPLETE") || (authConfirmMode == "NOT_COMPLETE") {
        tp.BizContent["auth_confirm_mode"] = authConfirmMode
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "预授权确认模式不合法", nil))
    }
}

func (tp *tradePay) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tp.outTradeNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不能为空", nil))
    }
    if len(tp.authCode) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付授权码不能为空", nil))
    }
    if tp.totalAmount <= 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不能为空", nil))
    }
    if tp.discountableAmount >= tp.totalAmount {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "可打折金额必须小于订单总金额", nil))
    }
    if len(tp.subject) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不能为空", nil))
    }
    tp.BizContent["out_trade_no"] = tp.outTradeNo
    tp.BizContent["auth_code"] = tp.authCode
    tp.BizContent["total_amount"] = tp.totalAmount
    tp.BizContent["subject"] = tp.subject
    if tp.discountableAmount > 0 {
        tp.BizContent["discountable_amount"] = tp.discountableAmount
    }

    return tp.GetRequest()
}

func NewTradePay(appId string) *tradePay {
    conf := alipay.NewConfig().GetAccount(appId)
    tp := &tradePay{alipay.NewBase(appId), "", "", 0.00, 0.00, ""}
    tp.totalAmount = 0.00
    tp.discountableAmount = 0.00
    tp.BizContent["seller_id"] = conf.GetSellerId()
    tp.BizContent["scene"] = "bar_code"
    tp.BizContent["trans_currency"] = "CNY"
    tp.BizContent["settle_currency"] = "CNY"
    tp.SetMethod("alipay.trade.pay")
    tp.SetUrlNotify(true)
    return tp
}
