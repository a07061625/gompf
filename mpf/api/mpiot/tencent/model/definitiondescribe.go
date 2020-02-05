/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package model

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询产品数据模板
type definitionDescribe struct {
    mpiot.BaseTencent
    productId string // 产品ID
}

func (dd *definitionDescribe) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        dd.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (dd *definitionDescribe) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dd.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    dd.ReqData["ProductId"] = dd.productId

    return dd.GetRequest()
}

func NewDefinitionDescribe() *definitionDescribe {
    dd := &definitionDescribe{mpiot.NewBaseTencent(), ""}
    dd.ReqHeader["X-TC-Action"] = "DescribeModelDefinition"
    return dd
}
