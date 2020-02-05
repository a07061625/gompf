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

// 发布产品
type productRelease struct {
    mpiot.BaseTencent
    productId string // 产品ID
}

func (pr *productRelease) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        pr.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (pr *productRelease) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pr.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    pr.ReqData["ProductId"] = pr.productId

    return pr.GetRequest()
}

func NewProductRelease() *productRelease {
    pr := &productRelease{mpiot.NewBaseTencent(), ""}
    pr.ReqData["DevStatus"] = "released"
    pr.ReqHeader["X-TC-Action"] = "ReleaseStudioProduct"
    return pr
}
