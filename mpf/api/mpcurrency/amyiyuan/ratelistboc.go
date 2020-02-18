/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:16
 */
package amyiyuan

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpcurrency"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 中国银行的实时汇率表
type rateListBOC struct {
    mpcurrency.BaseAMYiYuan
    currencyCode string // 货币类型
}

func (boc *rateListBOC) SetCurrencyCode(currencyCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, currencyCode)
    if match {
        boc.currencyCode = currencyCode
    } else {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "货币类型不合法", nil))
    }
}

func (boc *rateListBOC) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    boc.ServiceUri = "/waihui-list"
    if len(boc.currencyCode) > 0 {
        boc.ReqData["code"] = boc.currencyCode
        boc.ServiceUri += "?" + mpf.HTTPCreateParams(boc.ReqData, "none", 1)
    }

    return boc.GetRequest()
}

func NewRateListBOC() *rateListBOC {
    boc := &rateListBOC{mpcurrency.NewBaseAMYiYuan(), ""}
    return boc
}
