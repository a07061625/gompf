/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package principal

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 重新生成密钥
type keyRefresh struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    principalName string // principal名称
}

func (kr *keyRefresh) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        kr.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (kr *keyRefresh) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        kr.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (kr *keyRefresh) SetTarget(target string) {
    switch target {
    case "all":
        kr.ExtendData["target"] = target
    case "password":
        kr.ExtendData["target"] = target
    case "cert":
        kr.ExtendData["target"] = target
    default:
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "生成类型不合法", nil))
    }
}

func (kr *keyRefresh) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(kr.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(kr.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    kr.ServiceUri = "/v1/endpoint/" + kr.endpointName + "/principal/" + kr.principalName

    kr.ReqUrl = kr.GetServiceUrl()

    reqBody := mpf.JsonMarshal(kr.ExtendData)
    client, req := kr.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewKeyRefresh() *keyRefresh {
    kr := &keyRefresh{mpiot.NewBaseBaiDu(), "", ""}
    kr.ExtendData["target"] = "all"
    kr.ReqContentType = project.HttpContentTypeJson
    kr.ReqMethod = fasthttp.MethodPost
    return kr
}
