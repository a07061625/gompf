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

// 线上资金授权冻结接口
type orderAppFreeze struct {
    alipay.BaseAliPay
    outOrderNo   string                   // 商户授权资金单号
    outRequestNo string                   // 商户授权资金操作流水号
    orderTitle   string                   // 标题
    amount       float32                  // 冻结金额
    payChannels  []map[string]interface{} // 支付渠道列表
}

func (oaf *orderAppFreeze) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outOrderNo)
    if match {
        oaf.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        oaf.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetOrderTitle(orderTitle string) {
    if len(orderTitle) > 0 {
        trueTitle := []rune(orderTitle)
        oaf.orderTitle = string(trueTitle[:50])
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标题不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetAmount(amount float32) {
    if amount >= 0.01 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
        oaf.amount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "冻结金额不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetPayeeLogonId(payeeLogonId string) {
    if len(payeeLogonId) > 0 {
        oaf.BizContent["payee_logon_id"] = payeeLogonId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "收款方支付宝账号不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetPayeeUserId(payeeUserId string) {
    match, _ := regexp.MatchString(project.RegexDigit, payeeUserId)
    if match {
        oaf.BizContent["payee_user_id"] = payeeUserId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "收款方支付宝用户号不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetPayTimeout(payTimeout string) {
    if len(payTimeout) > 0 {
        oaf.BizContent["pay_timeout"] = payTimeout
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "最晚付款时间不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetExtraParam(extraParam map[string]interface{}) {
    if len(extraParam) > 0 {
        oaf.BizContent["extra_param"] = mpf.JsonMarshal(extraParam)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "扩展信息不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetSceneCode(sceneCode string) {
    switch sceneCode {
    case "O2O_AUTH_COMMON_SCENE": // 当面预授权通用场景
        oaf.BizContent["scene_code"] = sceneCode
    case "ONLINE_AUTH_COMMON_SCENE": // 支付宝预授权通用场景
        oaf.BizContent["scene_code"] = sceneCode
    case "OVERSEAS_O2O_AUTH_COMMON_SCENE": // 境外当面预授权通用场景
        oaf.BizContent["scene_code"] = sceneCode
    case "OVERSEAS_ONLINE_AUTH_COMMON_SCENE": // 境外支付宝预授权通用场景
        oaf.BizContent["scene_code"] = sceneCode
    default:
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "场景码不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetCurrencyTrans(currencyTrans string) {
    _, ok := alipay.FundCurrencyList[currencyTrans]
    if ok {
        oaf.BizContent["trans_currency"] = currencyTrans
    } else if currencyTrans == alipay.FundCurrencyCNY {
        oaf.BizContent["trans_currency"] = currencyTrans
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标价币种不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetCurrencySettle(currencySettle string) {
    _, ok := alipay.FundCurrencyList[currencySettle]
    if ok {
        oaf.BizContent["settle_currency"] = currencySettle
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "结算币种不合法", nil))
    }
}

func (oaf *orderAppFreeze) SetPayChannels(payChannels []map[string]interface{}) {
    if len(payChannels) > 0 {
        oaf.payChannels = payChannels
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付渠道不合法", nil))
    }
}

func (oaf *orderAppFreeze) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(oaf.outOrderNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不能为空", nil))
    }
    if len(oaf.outRequestNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不能为空", nil))
    }
    if len(oaf.orderTitle) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "标题不能为空", nil))
    }
    if oaf.amount <= 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "冻结金额不能为空", nil))
    }
    if len(oaf.payChannels) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付渠道不能为空", nil))
    }
    oaf.BizContent["out_order_no"] = oaf.outOrderNo
    oaf.BizContent["out_request_no"] = oaf.outRequestNo
    oaf.BizContent["order_title"] = oaf.orderTitle
    oaf.BizContent["amount"] = oaf.amount
    oaf.BizContent["enable_pay_channels"] = mpf.JsonMarshal(oaf.payChannels)

    return oaf.GetRequest()
}

func NewOrderAppFreeze(appId string) *orderAppFreeze {
    oaf := &orderAppFreeze{alipay.NewBase(appId), "", "", "", 0.00, make([]map[string]interface{}, 0)}
    oaf.amount = 0.00
    oaf.SetMethod("alipay.fund.auth.order.app.freeze")
    oaf.BizContent["product_code"] = "PRE_AUTH_ONLINE"
    return oaf
}
