/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 12:45
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

// 资金授权操作查询接口
type operationDetailQuery struct {
    alipay.BaseAliPay
    authNo       string // 支付宝授权资金单号
    outOrderNo   string // 商户授权资金单号
    operationId  string // 支付宝授权资金操作流水号
    outRequestNo string // 商户授权资金操作流水号
}

func (odq *operationDetailQuery) SetAuthNo(authNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, authNo)
    if match {
        odq.authNo = authNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金单号不合法", nil))
    }
}

func (odq *operationDetailQuery) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outOrderNo)
    if match {
        odq.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金单号不合法", nil))
    }
}

func (odq *operationDetailQuery) SetOperationId(operationId string) {
    match, _ := regexp.MatchString(project.RegexDigit, operationId)
    if match {
        odq.operationId = operationId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金操作流水号不合法", nil))
    }
}

func (odq *operationDetailQuery) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        odq.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户授权资金操作流水号不合法", nil))
    }
}

func (odq *operationDetailQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(odq.authNo) > 0 {
        odq.BizContent["auth_no"] = odq.authNo
    } else if len(odq.outOrderNo) > 0 {
        odq.BizContent["out_order_no"] = odq.outOrderNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金单号和商户授权资金单号不能同时为空", nil))
    }
    if len(odq.operationId) > 0 {
        odq.BizContent["operation_id"] = odq.operationId
    } else if len(odq.outRequestNo) > 0 {
        odq.BizContent["out_request_no"] = odq.outRequestNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝授权资金操作流水号和商户授权资金操作流水号不能同时为空", nil))
    }

    return odq.GetRequest()
}

func NewOperationDetailQuery(appId string) *operationDetailQuery {
    odq := &operationDetailQuery{alipay.NewBase(appId), "", "", "", ""}
    odq.SetMethod("alipay.fund.auth.operation.detail.query")
    return odq
}
