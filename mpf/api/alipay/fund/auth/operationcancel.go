/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 12:29
 */
package fund

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 资金授权撤销接口
type operationCancel struct {
    alipay.BaseAliPay
    authNo       string // 支付宝授权资金单号
    outOrderNo   string // 商户授权资金单号
    operationId  string // 支付宝授权资金操作流水号
    outRequestNo string // 商户授权资金操作流水号
    remark       string // 描述
}

func (oc *operationCancel) SetAuthNo(authNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, authNo)
    if match {
        oc.authNo = authNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金单号不合法", nil))
    }
}

func (oc *operationCancel) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outOrderNo)
    if match {
        oc.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不合法", nil))
    }
}

func (oc *operationCancel) SetOperationId(operationId string) {
    match, _ := regexp.MatchString(project.RegexDigit, operationId)
    if match {
        oc.operationId = operationId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金操作流水号不合法", nil))
    }
}

func (oc *operationCancel) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        oc.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不合法", nil))
    }
}

func (oc *operationCancel) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        oc.remark = string(trueRemark[:50])
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "描述不合法", nil))
    }
}

func (oc *operationCancel) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(oc.authNo) > 0 {
        oc.BizContent["auth_no"] = oc.authNo
    } else if len(oc.outOrderNo) > 0 {
        oc.BizContent["out_order_no"] = oc.outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金单号和商户授权资金单号不能同时为空", nil))
    }
    if len(oc.operationId) > 0 {
        oc.BizContent["operation_id"] = oc.operationId
    } else if len(oc.outRequestNo) > 0 {
        oc.BizContent["out_request_no"] = oc.outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金操作流水号和商户授权资金操作流水号不能同时为空", nil))
    }
    if len(oc.remark) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "描述不能为空", nil))
    }
    oc.BizContent["remark"] = oc.remark

    return oc.GetRequest()
}

func NewOperationCancel(appId string) *operationCancel {
    oc := &operationCancel{alipay.NewBase(appId), "", "", "", "", ""}
    oc.SetMethod("alipay.fund.auth.operation.cancel")
    return oc
}
