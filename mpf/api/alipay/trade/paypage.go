/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 9:36
 */
package trade

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 统一收单下单并支付页面接口
type payPage struct {
    alipay.BaseAliPay
    outTradeNo  string  // 商户订单号
    totalAmount float32 // 订单总金额,单位为元
    subject     string  // 订单标题
}

func (pp *payPage) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outTradeNo)
    if match {
        pp.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (pp *payPage) SetTotalAmount(totalAmount float32) {
    if (totalAmount > 0) && (totalAmount <= 100000000) {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalAmount), 64)
        pp.totalAmount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不合法", nil))
    }
}

func (pp *payPage) SetSubject(subject string) {
    if len(subject) > 0 {
        trueName := []rune(subject)
        pp.subject = string(trueName[:128])
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不合法", nil))
    }
}

func (pp *payPage) SetTimeExpire(timeExpire int) {
    if timeExpire > time.Now().Second() {
        et := time.Unix(int64(timeExpire), 0)
        pp.BizContent["time_expire"] = et.Format("2006-01-02 03-04")
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "绝对超时时间不合法", nil))
    }
}

func (pp *payPage) SetGoodsDetail(goodsDetail map[string]interface{}) {
    if len(goodsDetail) > 0 {
        pp.BizContent["goods_detail"] = goodsDetail
    } else {
        delete(pp.BizContent, "goods_detail")
    }
}

func (pp *payPage) SetPassBackParams(passBackParams map[string]string) {
    if len(passBackParams) > 0 {
        pp.BizContent["passback_params"] = mpf.HttpCreateParams(passBackParams, "none", 1)
    } else {
        delete(pp.BizContent, "passback_params")
    }
}

func (pp *payPage) SetExtendParams(extendParams map[string]interface{}) {
    if len(extendParams) > 0 {
        pp.BizContent["extend_params"] = extendParams
    } else {
        delete(pp.BizContent, "extend_params")
    }
}

func (pp *payPage) SetGoodsType(goodsType string) {
    if (goodsType == "0") || (goodsType == "1") {
        pp.BizContent["goods_type"] = goodsType
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商品主类型不合法", nil))
    }
}

func (pp *payPage) SetTimeoutExpress(timeoutExpress string) {
    if len(timeoutExpress) > 0 {
        pp.BizContent["timeout_express"] = timeoutExpress
    } else {
        delete(pp.BizContent, "timeout_express")
    }
}

func (pp *payPage) SetPromoParams(promoParams map[string]interface{}) {
    if len(promoParams) > 0 {
        pp.BizContent["promo_params"] = mpf.JsonMarshal(promoParams)
    } else {
        delete(pp.BizContent, "promo_params")
    }
}

func (pp *payPage) SetRoyaltyInfo(royaltyInfo map[string]interface{}) {
    if len(royaltyInfo) > 0 {
        pp.BizContent["royalty_info"] = royaltyInfo
    } else {
        delete(pp.BizContent, "royalty_info")
    }
}

func (pp *payPage) SetSubMerchant(subMerchant map[string]interface{}) {
    if len(subMerchant) > 0 {
        pp.BizContent["sub_merchant"] = subMerchant
    } else {
        delete(pp.BizContent, "sub_merchant")
    }
}

func (pp *payPage) SetDisablePayChannels(payChannels []string) {
    if len(payChannels) > 0 {
        pp.BizContent["disable_pay_channels"] = strings.Join(payChannels, ",")
    } else {
        delete(pp.BizContent, "disable_pay_channels")
    }
}

func (pp *payPage) SetEnablePayChannels(payChannels []string) {
    if len(payChannels) > 0 {
        pp.BizContent["enable_pay_channels"] = strings.Join(payChannels, ",")
    } else {
        delete(pp.BizContent, "enable_pay_channels")
    }
}

func (pp *payPage) SetStoreId(storeId string) {
    if len(storeId) > 0 {
        pp.BizContent["store_id"] = storeId
    } else {
        delete(pp.BizContent, "store_id")
    }
}

func (pp *payPage) SetQrPayMode(qrPayMode int) {
    if (qrPayMode >= 0) && (qrPayMode <= 4) {
        pp.BizContent["qr_pay_mode"] = strconv.Itoa(qrPayMode)
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "扫码支付方式不合法", nil))
    }
}

func (pp *payPage) SetQrCodeWidth(qrCodeWidth int) {
    if qrCodeWidth > 0 {
        pp.BizContent["qrcode_width"] = qrCodeWidth
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "二维码宽度不合法", nil))
    }
}

func (pp *payPage) SetSettleInfo(settleInfo map[string]interface{}) {
    if len(settleInfo) > 0 {
        pp.BizContent["settle_info"] = settleInfo
    } else {
        delete(pp.BizContent, "settle_info")
    }
}

func (pp *payPage) SetInvoiceInfo(invoiceInfo map[string]interface{}) {
    if len(invoiceInfo) > 0 {
        pp.BizContent["invoice_info"] = invoiceInfo
    } else {
        delete(pp.BizContent, "invoice_info")
    }
}

func (pp *payPage) SetAgreementSignParams(agreementSignParams map[string]interface{}) {
    if len(agreementSignParams) > 0 {
        pp.BizContent["agreement_sign_params"] = agreementSignParams
    } else {
        delete(pp.BizContent, "agreement_sign_params")
    }
}

func (pp *payPage) SetIntegrationType(integrationType string) {
    if (integrationType == "ALIAPP") || (integrationType == "PCWEB") {
        pp.BizContent["integration_type"] = integrationType
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "页面集成方式不合法", nil))
    }
}

func (pp *payPage) SetRequestFromUrl(requestFromUrl string) {
    if len(requestFromUrl) > 0 {
        pp.BizContent["request_from_url"] = requestFromUrl
    } else {
        delete(pp.BizContent, "request_from_url")
    }
}

func (pp *payPage) SetBusinessParams(businessParams map[string]interface{}) {
    if len(businessParams) > 0 {
        pp.BizContent["business_params"] = mpf.JsonMarshal(businessParams)
    } else {
        delete(pp.BizContent, "business_params")
    }
}

func (pp *payPage) SetExtUserInfo(extUserInfo map[string]interface{}) {
    if len(extUserInfo) > 0 {
        pp.BizContent["ext_user_info"] = extUserInfo
    } else {
        delete(pp.BizContent, "ext_user_info")
    }
}

func (pp *payPage) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := pp.ReqData["return_url"]
    if !ok {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "同步通知地址不能为空", nil))
    }
    if len(pp.outTradeNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不能为空", nil))
    }
    if pp.totalAmount <= 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单总金额不能为空", nil))
    }
    if len(pp.subject) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "订单标题不能为空", nil))
    }
    pp.BizContent["out_trade_no"] = pp.outTradeNo
    pp.BizContent["total_amount"] = pp.totalAmount
    pp.BizContent["subject"] = pp.subject

    return pp.GetRequest()
}

func NewPayPage(appId string) *payPage {
    pp := &payPage{alipay.NewBase(appId), "", 0.00, ""}
    pp.totalAmount = 0.00
    pp.BizContent["product_code"] = "FAST_INSTANT_TRADE_PAY"
    pp.BizContent["goods_type"] = "1"
    pp.BizContent["qr_pay_mode"] = "2"
    pp.BizContent["integration_type"] = "PCWEB"
    pp.SetMethod("alipay.trade.page.pay")
    pp.SetUrlNotify(true)
    return pp
}
