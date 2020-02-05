/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 12:49
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

// 资金授权冻结接口
type orderFreeze struct {
    alipay.BaseAliPay
    authCode     string  // 支付授权码
    outOrderNo   string  // 商户授权资金单号
    outRequestNo string  // 商户授权资金操作流水号
    orderTitle   string  // 标题
    amount       float32 // 冻结金额
}

func (of *orderFreeze) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(project.RegexDigit, authCode)
    if match {
        of.authCode = authCode
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付授权码不合法", nil))
    }
}

func (of *orderFreeze) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outOrderNo)
    if match {
        of.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不合法", nil))
    }
}

func (of *orderFreeze) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        of.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不合法", nil))
    }
}

func (of *orderFreeze) SetOrderTitle(orderTitle string) {
    if len(orderTitle) > 0 {
        trueTitle := []rune(orderTitle)
        of.orderTitle = string(trueTitle[:50])
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标题不合法", nil))
    }
}

func (of *orderFreeze) SetAmount(amount float32) {
    if amount >= 0.01 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
        of.amount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "冻结金额不合法", nil))
    }
}

func (of *orderFreeze) SetPayeeLogonId(payeeLogonId string) {
    if len(payeeLogonId) > 0 {
        of.BizContent["payee_logon_id"] = payeeLogonId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "收款方支付宝账号不合法", nil))
    }
}

func (of *orderFreeze) SetPayeeUserId(payeeUserId string) {
    match, _ := regexp.MatchString(project.RegexDigit, payeeUserId)
    if match {
        of.BizContent["payee_user_id"] = payeeUserId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "收款方支付宝用户号不合法", nil))
    }
}

func (of *orderFreeze) SetPayTimeout(payTimeout string) {
    if len(payTimeout) > 0 {
        of.BizContent["pay_timeout"] = payTimeout
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "最晚付款时间不合法", nil))
    }
}

func (of *orderFreeze) SetExtraParam(extraParam map[string]interface{}) {
    if len(extraParam) > 0 {
        of.BizContent["extra_param"] = mpf.JsonMarshal(extraParam)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "扩展信息不合法", nil))
    }
}

func (of *orderFreeze) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(of.authCode) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付授权码不能为空", nil))
    }
    if len(of.outOrderNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不能为空", nil))
    }
    if len(of.outRequestNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不能为空", nil))
    }
    if len(of.orderTitle) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标题不能为空", nil))
    }
    if of.amount <= 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "冻结金额不能为空", nil))
    }
    of.BizContent["auth_code"] = of.authCode
    of.BizContent["out_order_no"] = of.outOrderNo
    of.BizContent["out_request_no"] = of.outRequestNo
    of.BizContent["order_title"] = of.orderTitle
    of.BizContent["amount"] = of.amount

    return of.GetRequest()
}

func NewOrderFreeze(appId string) *orderFreeze {
    of := &orderFreeze{alipay.NewBase(appId), "", "", "", "", 0.00}
    of.amount = 0.00
    of.SetMethod("alipay.fund.auth.order.freeze")
    of.BizContent["auth_code_type"] = "bar_code"
    of.BizContent["product_code"] = "PRE_AUTH"
    return of
}
