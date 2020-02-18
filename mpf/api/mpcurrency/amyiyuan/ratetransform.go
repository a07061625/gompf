/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:22
 */
package amyiyuan

import (
    "fmt"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpcurrency"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 汇率转换
type rateTransform struct {
    mpcurrency.BaseAMYiYuan
    fromCode string  // 源货币类型
    toCode   string  // 目标货币类型
    money    float32 // 转换金额,单位为元
}

func (rt *rateTransform) SetFromCode(fromCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, fromCode)
    if match {
        rt.fromCode = fromCode
    } else {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "源货币类型不合法", nil))
    }
}

func (rt *rateTransform) SetToCode(toCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, toCode)
    if match {
        rt.toCode = toCode
    } else {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "目标货币类型不合法", nil))
    }
}

func (rt *rateTransform) SetMoney(money float32) {
    if money > 0 {
        rt.money = money
    } else {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "转换金额不合法", nil))
    }
}

func (rt *rateTransform) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rt.fromCode) == 0 {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "源货币类型不能为空", nil))
    }
    if len(rt.toCode) == 0 {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "目标货币类型不能为空", nil))
    }
    if rt.money <= 0 {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "转换金额不能为空", nil))
    }
    rt.ReqData["fromCode"] = rt.fromCode
    rt.ReqData["toCode"] = rt.toCode
    rt.ReqData["money"] = fmt.Sprintf("%.2f", rt.money)
    rt.ServiceUri = "/waihui-transform?" + mpf.HTTPCreateParams(rt.ReqData, "none", 1)

    return rt.GetRequest()
}

func NewRateTransform() *rateTransform {
    rt := &rateTransform{mpcurrency.NewBaseAMYiYuan(), "", "", 0.00}
    rt.money = 0.00
    return rt
}
