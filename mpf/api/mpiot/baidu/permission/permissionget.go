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

// 获取指定topic的信息
type permissionGet struct {
    mpiot.BaseBaiDu
    endpointName   string // endpoint名称
    permissionUuid string // policy名称
}

func (pg *permissionGet) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pg.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pg *permissionGet) SetPermissionUuid(permissionUuid string) {
    if len(permissionUuid) > 0 {
        pg.permissionUuid = permissionUuid
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "permissionID不合法", nil))
    }
}

func (pg *permissionGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pg.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pg.permissionUuid) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "permissionID不能为空", nil))
    }
    pg.ServiceUri = "/v1/endpoint/" + pg.endpointName + "/permission/" + pg.permissionUuid

    pg.ReqURI = pg.GetServiceUrl()

    return pg.GetRequest()
}

func NewPermissionGet() *permissionGet {
    pg := &permissionGet{mpiot.NewBaseBaiDu(), "", ""}
    return pg
}
