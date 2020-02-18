/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package permission

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除已有的topic
type permissionDelete struct {
    mpiot.BaseBaiDu
    endpointName   string // endpoint名称
    permissionUuid string // policy名称
}

func (pd *permissionDelete) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pd.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pd *permissionDelete) SetPermissionUuid(permissionUuid string) {
    if len(permissionUuid) > 0 {
        pd.permissionUuid = permissionUuid
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "permissionID不合法", nil))
    }
}

func (pd *permissionDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pd.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pd.permissionUuid) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "permissionID不能为空", nil))
    }
    pd.ServiceUri = "/v1/endpoint/" + pd.endpointName + "/permission/" + pd.permissionUuid

    pd.ReqURI = pd.GetServiceUrl()

    return pd.GetRequest()
}

func NewPermissionDelete() *permissionDelete {
    pd := &permissionDelete{mpiot.NewBaseBaiDu(), "", ""}
    pd.ReqMethod = fasthttp.MethodDelete
    return pd
}
