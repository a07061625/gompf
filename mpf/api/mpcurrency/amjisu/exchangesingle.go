/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:08
 */
package amjisu

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpcurrency"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 单个货币查询
type exchangeSingle struct {
    mpcurrency.BaseAMJiSu
    currency string // 货币类型
}

func (es *exchangeSingle) SetCurrency(currency string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, currency)
    if match {
        es.currency = currency
    } else {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "货币类型不合法", nil))
    }
}

func (es *exchangeSingle) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(es.currency) == 0 {
        panic(mperr.NewCurrencyAMJiSu(errorcode.CurrencyAMJiSuParam, "货币类型不能为空", nil))
    }
    es.ReqData["currency"] = es.currency
    es.ServiceUri = "/exchange/single?" + mpf.HTTPCreateParams(es.ReqData, "none", 1)

    return es.GetRequest()
}

func NewExchangeSingle() *exchangeSingle {
    es := &exchangeSingle{mpcurrency.NewBaseAMJiSu(), ""}
    return es
}
