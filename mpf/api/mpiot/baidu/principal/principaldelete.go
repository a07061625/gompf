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

// 删除principal
type principalDelete struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    principalName string // principal名称
}

func (pd *principalDelete) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pd.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pd *principalDelete) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        pd.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (pd *principalDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pd.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pd.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    pd.ServiceUri = "/v1/endpoint/" + pd.endpointName + "/principal/" + pd.principalName

    pd.ReqUrl = pd.GetServiceUrl()

    return pd.GetRequest()
}

func NewPrincipalDelete() *principalDelete {
    pd := &principalDelete{mpiot.NewBaseBaiDu(), "", ""}
    pd.ReqMethod = fasthttp.MethodPut
    return pd
}
