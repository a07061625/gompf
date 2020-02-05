/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 14:39
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

// 资金授权解冻接口
type orderUnfreeze struct {
    alipay.BaseAliPay
    authNo       string  // 支付宝授权资金单号
    outRequestNo string  // 商户授权资金操作流水号
    amount       float32 // 解冻金额
    remark       string  // 描述
}

func (ou *orderUnfreeze) SetAuthNo(authNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, authNo)
    if match {
        ou.authNo = authNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金单号不合法", nil))
    }
}

func (ou *orderUnfreeze) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        ou.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不合法", nil))
    }
}

func (ou *orderUnfreeze) SetAmount(amount float32) {
    if amount >= 0.01 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
        ou.amount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "解冻金额不合法", nil))
    }
}

func (ou *orderUnfreeze) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        ou.remark = string(trueRemark[:50])
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "描述不合法", nil))
    }
}

func (ou *orderUnfreeze) SetExtraParam(extraParam map[string]interface{}) {
    if len(extraParam) > 0 {
        ou.BizContent["extra_param"] = mpf.JsonMarshal(extraParam)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "扩展信息不合法", nil))
    }
}

func (ou *orderUnfreeze) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ou.authNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金单号不能为空", nil))
    }
    if len(ou.outRequestNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不能为空", nil))
    }
    if ou.amount <= 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "解冻金额不能为空", nil))
    }
    if len(ou.remark) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "描述不能为空", nil))
    }
    ou.BizContent["auth_no"] = ou.authNo
    ou.BizContent["out_request_no"] = ou.outRequestNo
    ou.BizContent["amount"] = ou.amount
    ou.BizContent["remark"] = ou.remark

    return ou.GetRequest()
}

func NewOrderUnfreeze(appId string) *orderUnfreeze {
    ou := &orderUnfreeze{alipay.NewBase(appId), "", "", 0.00, ""}
    ou.amount = 0.00
    ou.SetMethod("alipay.fund.auth.order.unfreeze")
    return ou
}
