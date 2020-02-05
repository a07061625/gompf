/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package model

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改产品数据模板
type definitionModify struct {
    mpiot.BaseTencent
    productId   string                 // 产品ID
    modelSchema map[string]interface{} // 数据模板定义
}

func (dm *definitionModify) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        dm.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (dm *definitionModify) SetModelSchema(modelSchema map[string]interface{}) {
    if len(modelSchema) > 0 {
        dm.modelSchema = modelSchema
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "数据模板定义不合法", nil))
    }
}

func (dm *definitionModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dm.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    if len(dm.modelSchema) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "数据模板定义不能为空", nil))
    }
    dm.ReqData["ProductId"] = dm.productId
    dm.ReqData["ModelSchema"] = mpf.JsonMarshal(dm.modelSchema)

    return dm.GetRequest()
}

func NewDefinitionModify() *definitionModify {
    dm := &definitionModify{mpiot.NewBaseTencent(), "", make(map[string]interface{})}
    dm.ReqHeader["X-TC-Action"] = "ModifyModelDefinition"
    return dm
}
