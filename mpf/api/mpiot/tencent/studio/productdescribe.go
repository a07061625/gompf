/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package studio

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取产品详情
type productDescribe struct {
    mpiot.BaseTencent
    productId string // 产品ID
}

func (pd *productDescribe) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        pd.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (pd *productDescribe) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pd.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    pd.ReqData["ProductId"] = pd.productId

    return pd.GetRequest()
}

func NewProductDescribe() *productDescribe {
    pd := &productDescribe{mpiot.NewBaseTencent(), ""}
    pd.ReqHeader["X-TC-Action"] = "DescribeStudioProduct"
    return pd
}
