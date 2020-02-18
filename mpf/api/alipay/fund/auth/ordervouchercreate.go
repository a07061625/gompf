/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 13:29
 */
package fund

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

// 资金授权发码接口
type orderVoucherCreate struct {
    alipay.BaseAliPay
    outOrderNo   string                   // 商户授权资金单号
    outRequestNo string                   // 商户授权资金操作流水号
    orderTitle   string                   // 标题
    amount       float32                  // 冻结金额
    payChannels  []map[string]interface{} // 支付渠道列表
}

func (ovc *orderVoucherCreate) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outOrderNo)
    if match {
        ovc.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        ovc.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetOrderTitle(orderTitle string) {
    if len(orderTitle) > 0 {
        trueTitle := []rune(orderTitle)
        ovc.orderTitle = string(trueTitle[:50])
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标题不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetAmount(amount float32) {
    if amount >= 0.01 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
        ovc.amount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "冻结金额不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetPayeeLogonId(payeeLogonId string) {
    if len(payeeLogonId) > 0 {
        ovc.BizContent["payee_logon_id"] = payeeLogonId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "收款方支付宝账号不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetPayeeUserId(payeeUserId string) {
    match, _ := regexp.MatchString(project.RegexDigit, payeeUserId)
    if match {
        ovc.BizContent["payee_user_id"] = payeeUserId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "收款方支付宝用户号不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetPayTimeout(payTimeout string) {
    if len(payTimeout) > 0 {
        ovc.BizContent["pay_timeout"] = payTimeout
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "最晚付款时间不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetExtraParam(extraParam map[string]interface{}) {
    if len(extraParam) > 0 {
        ovc.BizContent["extra_param"] = mpf.JSONMarshal(extraParam)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "扩展信息不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetCurrencyTrans(currencyTrans string) {
    _, ok := alipay.FundCurrencyList[currencyTrans]
    if ok {
        ovc.BizContent["trans_currency"] = currencyTrans
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标价币种不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetCurrencySettle(currencySettle string) {
    _, ok := alipay.FundCurrencyList[currencySettle]
    if ok {
        ovc.BizContent["settle_currency"] = currencySettle
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "结算币种不合法", nil))
    }
}

func (ovc *orderVoucherCreate) SetPayChannels(payChannels []map[string]interface{}) {
    if len(payChannels) > 0 {
        ovc.payChannels = payChannels
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付渠道不合法", nil))
    }
}

func (ovc *orderVoucherCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ovc.outOrderNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不能为空", nil))
    }
    if len(ovc.outRequestNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不能为空", nil))
    }
    if len(ovc.orderTitle) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标题不能为空", nil))
    }
    if ovc.amount <= 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "冻结金额不能为空", nil))
    }
    if len(ovc.payChannels) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付渠道不能为空", nil))
    }
    ovc.BizContent["out_order_no"] = ovc.outOrderNo
    ovc.BizContent["out_request_no"] = ovc.outRequestNo
    ovc.BizContent["order_title"] = ovc.orderTitle
    ovc.BizContent["amount"] = ovc.amount
    ovc.BizContent["enable_pay_channels"] = mpf.JSONMarshal(ovc.payChannels)

    return ovc.GetRequest()
}

func NewOrderVoucherCreate(appId string) *orderVoucherCreate {
    ovc := &orderVoucherCreate{alipay.NewBase(appId), "", "", "", 0.00, make([]map[string]interface{}, 0)}
    ovc.amount = 0.00
    ovc.SetMethod("alipay.fund.auth.order.voucher.create")
    ovc.BizContent["product_code"] = "PRE_AUTH"
    return ovc
}
