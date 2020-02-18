/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 11:54
 */
package trade

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 支付宝订单信息同步接口
type orderSync struct {
    alipay.BaseAliPay
    outRequestNo string                 // 商户订单号
    bizType      string                 // 业务类型
    bizInfo      map[string]interface{} // 同步信息
}

func (os *orderSync) SetOutRequestNo(outRequestNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outRequestNo)
    if match {
        os.outRequestNo = outRequestNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不合法", nil))
    }
}

func (os *orderSync) SetBizType(bizType string) {
    if len(bizType) > 0 {
        os.bizType = bizType
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "业务类型不合法", nil))
    }
}

func (os *orderSync) SetBizInfo(bizInfo map[string]interface{}) {
    if len(bizInfo) > 0 {
        os.bizInfo = bizInfo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "同步信息不合法", nil))
    }
}

func (os *orderSync) SetTradeNo(tradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, tradeNo)
    if match {
        os.BizContent["trade_no"] = tradeNo
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "支付宝交易号不合法", nil))
    }
}

func (os *orderSync) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(os.outRequestNo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "商户订单号不能为空", nil))
    }
    if len(os.bizType) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "业务类型不能为空", nil))
    }
    if len(os.bizInfo) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "同步信息不能为空", nil))
    }
    os.BizContent["out_request_no"] = os.outRequestNo
    os.BizContent["biz_type"] = os.bizType
    os.BizContent["order_biz_info"] = mpf.JSONMarshal(os.bizInfo)

    return os.GetRequest()
}

func NewOrderSync(appId string) *orderSync {
    os := &orderSync{alipay.NewBase(appId), "", "", make(map[string]interface{})}
    os.SetMethod("alipay.trade.orderinfo.sync")
    return os
}
