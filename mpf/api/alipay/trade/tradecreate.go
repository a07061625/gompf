/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 12:03
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

// 统一收单交易创建接口
type tradeCreate struct {
    alipay.BaseAliPay
    outTradeNo         string  // 商户订单号
    totalAmount        float32 // 订单总金额,单位为元
    discountableAmount float32 // 可打折金额,单位为元
    subject            string  // 订单标题
}

func (tc *tradeCreate) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        tc.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (tc *tradeCreate) SetTotalAmount(totalAmount float32) {
    if (totalAmount > 0) && (totalAmount <= 100000000) {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalAmount), 64)
        tc.totalAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不合法", nil))
    }
}

func (tc *tradeCreate) SetDiscountableAmount(discountableAmount float32) {
    if (discountableAmount > 0) && (discountableAmount <= 100000000) {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", discountableAmount), 64)
        tc.discountableAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "可打折金额不合法", nil))
    }
}

func (tc *tradeCreate) SetSubject(subject string) {
    if len(subject) > 0 {
        trueName := []rune(subject)
        tc.subject = string(trueName[:128])
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不合法", nil))
    }
}

func (tc *tradeCreate) SetBody(body string) {
    if (len(body) > 0) && (len(body) <= 128) {
        tc.BizContent["body"] = body
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商品描述不合法", nil))
    }
}

func (tc *tradeCreate) SetBuyerId(buyerId string) {
    match, _ := regexp.MatchString(project.RegexDigit, buyerId)
    if match {
        tc.BizContent["buyer_id"] = buyerId
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "买家支付宝用户ID不合法", nil))
    }
}

func (tc *tradeCreate) SetGoodsDetail(goodsDetail map[string]interface{}) {
    if len(goodsDetail) > 0 {
        tc.BizContent["goods_detail"] = goodsDetail
    } else {
        delete(tc.BizContent, "goods_detail")
    }
}

func (tc *tradeCreate) SetOperatorId(operatorId string) {
    if len(operatorId) > 0 {
        tc.BizContent["operator_id"] = operatorId
    } else {
        delete(tc.BizContent, "operator_id")
    }
}

func (tc *tradeCreate) SetStoreId(storeId string) {
    if len(storeId) > 0 {
        tc.BizContent["store_id"] = storeId
    } else {
        delete(tc.BizContent, "store_id")
    }
}

func (tc *tradeCreate) SetTerminalId(terminalId string) {
    if len(terminalId) > 0 {
        tc.BizContent["terminal_id"] = terminalId
    } else {
        delete(tc.BizContent, "terminal_id")
    }
}

func (tc *tradeCreate) SetExtendParams(extendParams map[string]interface{}) {
    if len(extendParams) > 0 {
        tc.BizContent["extend_params"] = extendParams
    } else {
        delete(tc.BizContent, "extend_params")
    }
}

func (tc *tradeCreate) SetTimeoutExpress(timeoutExpress string) {
    if len(timeoutExpress) > 0 {
        tc.BizContent["timeout_express"] = timeoutExpress
    } else {
        delete(tc.BizContent, "timeout_express")
    }
}

func (tc *tradeCreate) SetSettleInfo(settleInfo map[string]interface{}) {
    if len(settleInfo) > 0 {
        tc.BizContent["settle_info"] = settleInfo
    } else {
        delete(tc.BizContent, "settle_info")
    }
}

func (tc *tradeCreate) SetBusinessParams(businessParams map[string]interface{}) {
    if len(businessParams) > 0 {
        tc.BizContent["business_params"] = mpf.JsonMarshal(businessParams)
    } else {
        delete(tc.BizContent, "business_params")
    }
}

func (tc *tradeCreate) SetReceiverAddressInfo(receiverAddressInfo map[string]interface{}) {
    if len(receiverAddressInfo) > 0 {
        tc.BizContent["receiver_address_info"] = receiverAddressInfo
    } else {
        delete(tc.BizContent, "receiver_address_info")
    }
}

func (tc *tradeCreate) SetLogisticsDetail(logisticsDetail map[string]interface{}) {
    if len(logisticsDetail) > 0 {
        tc.BizContent["logistics_detail"] = logisticsDetail
    } else {
        delete(tc.BizContent, "logistics_detail")
    }
}

func (tc *tradeCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tc.outTradeNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不能为空", nil))
    }
    if tc.totalAmount <= 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不能为空", nil))
    }
    if tc.discountableAmount >= tc.totalAmount {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "可打折金额必须小于订单总金额", nil))
    }
    if len(tc.subject) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不能为空", nil))
    }
    tc.BizContent["out_trade_no"] = tc.outTradeNo
    tc.BizContent["total_amount"] = tc.totalAmount
    tc.BizContent["subject"] = tc.subject
    if tc.discountableAmount > 0 {
        tc.BizContent["discountable_amount"] = tc.discountableAmount
    }

    return tc.GetRequest()
}

func NewTradeCreate(appId string) *tradeCreate {
    conf := alipay.NewConfig().GetAccount(appId)
    tc := &tradeCreate{alipay.NewBase(appId), "", 0.00, 0.00, ""}
    tc.totalAmount = 0.00
    tc.discountableAmount = 0.00
    tc.BizContent["seller_id"] = conf.GetSellerId()
    tc.SetMethod("alipay.trade.create")
    tc.SetUrlNotify(true)
    return tc
}
