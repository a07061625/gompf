/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/16 0016
 * Time: 13:17
 */
package amali

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mplogistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 单号查快递公司名
type companyFetch struct {
    mplogistics.BaseAMAli
    nu  string // 快递单号
}

func (cf *companyFetch) SetNu(nu string) {
    if len(nu) > 0 {
        cf.nu = nu
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "快递单号不合法", nil))
    }
}

func (cf *companyFetch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cf.nu) == 0 {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "快递单号不能为空", nil))
    }
    cf.ReqData["nu"] = cf.nu
    cf.ServiceUri = "/fetchCom?" + mpf.HTTPCreateParams(cf.ReqData, "none", 1)

    return cf.GetRequest()
}

func NewCompanyFetch() *companyFetch {
    cf := &companyFetch{mplogistics.NewBaseAMAli(), ""}
    return cf
}
