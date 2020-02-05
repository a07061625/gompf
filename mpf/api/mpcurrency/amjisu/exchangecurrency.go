/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:05
 */
package amjisu

import (
    "github.com/a07061625/gompf/mpf/api/mpcurrency"
    "github.com/valyala/fasthttp"
)

// 所有货币查询
type exchangeCurrency struct {
    mpcurrency.BaseAMJiSu
}

func (ec *exchangeCurrency) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    ec.ServiceUri = "/exchange/currency"

    return ec.GetRequest()
}

func NewExchangeCurrency() *exchangeCurrency {
    ec := &exchangeCurrency{mpcurrency.NewBaseAMJiSu()}
    return ec
}
