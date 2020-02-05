/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 18:37
 */
package amyiyuan

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpcurrency"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 十大银行实时汇率表
type rateListBank struct {
    mpcurrency.BaseAMYiYuan
    bankCode string // 银行编码
}

func (rlb *rateListBank) SetBankCode(bankCode string) {
    _, ok := bankTypes[bankCode]
    if ok {
        rlb.bankCode = bankCode
    } else {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "银行编码不合法", nil))
    }
}

func (rlb *rateListBank) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rlb.bankCode) == 0 {
        panic(mperr.NewCurrencyAMYiYuan(errorcode.CurrencyAMYiYuanParam, "银行编码不能为空", nil))
    }
    rlb.ReqData["bankCode"] = rlb.bankCode
    rlb.ServiceUri = "/bank10?" + mpf.HttpCreateParams(rlb.ReqData, "none", 1)

    return rlb.GetRequest()
}

func NewRateListBank() *rateListBank {
    rlb := &rateListBank{mpcurrency.NewBaseAMYiYuan(), ""}
    return rlb
}
