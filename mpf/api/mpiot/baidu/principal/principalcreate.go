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

// 创建principal
type principalCreate struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    principalName string // principal名称
}

func (pc *principalCreate) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pc.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pc *principalCreate) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        pc.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (pc *principalCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pc.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    pc.ServiceUri = "/v1/endpoint/" + pc.endpointName + "/principal"
    pc.ExtendData["principalName"] = pc.principalName

    pc.ReqUrl = pc.GetServiceUrl()

    reqBody := mpf.JsonMarshal(pc.ExtendData)
    client, req := pc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPrincipalCreate() *principalCreate {
    pc := &principalCreate{mpiot.NewBaseBaiDu(), "", ""}
    pc.ReqContentType = project.HTTPContentTypeJSON
    pc.ReqMethod = fasthttp.MethodPost
    return pc
}
