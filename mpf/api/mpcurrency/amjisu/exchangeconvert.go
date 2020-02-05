/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/16 0016
 * Time: 15:04
 */
package amjisu

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
type exchangeConvert struct {
    mpcurrency.BaseAMJiSu
    from   string  // 源货币类型
    to     string  // 目标货币类型
    amount float32 // 转换金额,单位为元
}

func (ec *exchangeConvert) SetFrom(from string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, from)
    if match {
        ec.from = from
    } else {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "源货币类型不合法", nil))
    }
}

func (ec *exchangeConvert) SetTo(to string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, to)
    if match {
        ec.to = to
    } else {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "目标货币类型不合法", nil))
    }
}

func (ec *exchangeConvert) SetAmount(amount float32) {
    if amount > 0 {
        ec.amount = amount
    } else {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "转换金额不合法", nil))
    }
}

func (ec *exchangeConvert) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ec.from) == 0 {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "源货币类型不能为空", nil))
    }
    if len(ec.to) == 0 {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "目标货币类型不能为空", nil))
    }
    if ec.amount <= 0 {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "转换金额不能为空", nil))
    }
    ec.ReqData["from"] = ec.from
    ec.ReqData["to"] = ec.to
    ec.ReqData["amount"] = fmt.Sprintf("%.2f", ec.amount)
    ec.ServiceUri = "/exchange/convert?" + mpf.HttpCreateParams(ec.ReqData, "none", 1)

    return ec.GetRequest()
}

func NewExchangeConvert() *exchangeConvert {
    ec := &exchangeConvert{mpcurrency.NewBaseAMJiSu(), "", "", 0.00}
    ec.amount = 0.00
    return ec
}
