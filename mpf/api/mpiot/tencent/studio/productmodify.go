/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package studio

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改产品
type productModify struct {
    mpiot.BaseTencent
    productId   string // 产品ID
    productName string // 产品名称
    productDesc string // 产品描述
    moduleId    int    // 模型ID
}

func (pm *productModify) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        pm.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (pm *productModify) SetProductName(productName string) {
    if len(productName) > 0 {
        pm.productName = productName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品名称不合法", nil))
    }
}

func (pm *productModify) SetProductDesc(productDesc string) {
    if len(productDesc) > 0 {
        pm.productDesc = productDesc
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品描述不合法", nil))
    }
}
func (pm *productModify) SetModuleId(moduleId int) {
    if moduleId > 0 {
        pm.moduleId = moduleId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "模型ID不合法", nil))
    }
}

func (pm *productModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pm.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    if len(pm.productName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品名称不能为空", nil))
    }
    if len(pm.productDesc) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品描述不能为空", nil))
    }
    if pm.moduleId <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "模型ID不能为空", nil))
    }
    pm.ReqData["ProductId"] = pm.productId
    pm.ReqData["ProjectName"] = pm.productName
    pm.ReqData["ProjectDesc"] = pm.productDesc
    pm.ReqData["ModuleId"] = strconv.Itoa(pm.moduleId)

    return pm.GetRequest()
}

func NewProductModify() *productModify {
    pm := &productModify{mpiot.NewBaseTencent(), "", "", "", 0}
    pm.ReqHeader["X-TC-Action"] = "ModifyStudioProduct"
    return pm
}
