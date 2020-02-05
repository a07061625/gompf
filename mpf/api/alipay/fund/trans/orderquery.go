/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 11:54
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

// 查询转账订单接口
type orderQuery struct {
    alipay.BaseAliPay
    outBizNo string // 商户转账单号
    orderId  string // 支付宝转账单号
}

func (oq *orderQuery) SetOutBizNo(outBizNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outBizNo)
    if match {
        oq.outBizNo = outBizNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户转账单号不合法", nil))
    }
}

func (oq *orderQuery) SetOrderId(orderId string) {
    match, _ := regexp.MatchString(project.RegexDigit, orderId)
    if match {
        oq.orderId = orderId
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "支付宝转账单号不合法", nil))
    }
}

func (oq *orderQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(oq.orderId) > 0 {
        oq.BizContent["order_id"] = oq.orderId
    } else if len(oq.outBizNo) > 0 {
        oq.BizContent["out_biz_no"] = oq.outBizNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户转账单号和支付宝转账单号不能同时为空", nil))
    }

    return oq.GetRequest()
}

func NewOrderQuery(appId string) *orderQuery {
    oq := &orderQuery{alipay.NewBase(appId), "", ""}
    oq.SetMethod("alipay.fund.trans.order.query")
    return oq
}
