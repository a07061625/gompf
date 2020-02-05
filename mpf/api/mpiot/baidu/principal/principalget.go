/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package principal

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取指定的principal信息
type principalGet struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    principalName string // principal名称
}

func (pg *principalGet) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pg.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pg *principalGet) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        pg.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (pg *principalGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pg.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pg.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    pg.ServiceUri = "/v1/endpoint/" + pg.endpointName + "/principal/" + pg.principalName

    pg.ReqUrl = pg.GetServiceUrl()

    return pg.GetRequest()
}

func NewPrincipalGet() *principalGet {
    pg := &principalGet{mpiot.NewBaseBaiDu(), "", ""}
    return pg
}
