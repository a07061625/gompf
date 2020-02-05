/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:12
 */
package amyiyuan

import (
    "github.com/a07061625/gompf/mpf/api/mpcurrency"
    "github.com/valyala/fasthttp"
)

// 支持的外汇币种列表
type currencyList struct {
    mpcurrency.BaseAMYiYuan
}

func (cl *currencyList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    cl.ServiceUri = "/list"

    return cl.GetRequest()
}

func NewCurrencyList() *currencyList {
    cl := &currencyList{mpcurrency.NewBaseAMYiYuan()}
    return cl
}
