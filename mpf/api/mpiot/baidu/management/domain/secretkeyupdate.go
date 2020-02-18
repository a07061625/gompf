/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package domain

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新权限组密钥
type secretKeyUpdate struct {
    mpiot.BaseBaiDu
    domainName string // 权限组名称
}

func (sku *secretKeyUpdate) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        sku.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不合法", nil))
    }
}

func (sku *secretKeyUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sku.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不能为空", nil))
    }
    sku.ServiceUri = "/v3/iot/management/domain/" + sku.domainName

    sku.ReqURI = sku.GetServiceUrl() + "?updateSecretKey"

    return sku.GetRequest()
}

func NewSecretKeyUpdate() *secretKeyUpdate {
    sku := &secretKeyUpdate{mpiot.NewBaseBaiDu(), ""}
    sku.ReqData["updateSecretKey"] = ""
    sku.ReqMethod = fasthttp.MethodPut
    return sku
}
