/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 21:57
 */
package action

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 从一个Thing移除一个Principal
type thingPrincipalRemove struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    thingName     string // thing名称
    principalName string // principal名称
}

func (tpa *thingPrincipalRemove) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        tpa.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (tpa *thingPrincipalRemove) SetThingName(thingName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, thingName)
    if match {
        tpa.thingName = thingName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不合法", nil))
    }
}

func (tpa *thingPrincipalRemove) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        tpa.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (tpa *thingPrincipalRemove) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tpa.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(tpa.thingName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不能为空", nil))
    }
    if len(tpa.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    tpa.ExtendData["endpointName"] = tpa.endpointName
    tpa.ExtendData["thingName"] = tpa.thingName
    tpa.ExtendData["principalName"] = tpa.principalName

    tpa.ReqUrl = tpa.GetServiceUrl()

    reqBody := mpf.JsonMarshal(tpa.ExtendData)
    client, req := tpa.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewThingPrincipalRemove() *thingPrincipalRemove {
    tpa := &thingPrincipalRemove{mpiot.NewBaseBaiDu(), "", "", ""}
    tpa.ServiceUri = "/v1/action/remove-thing-principal"
    tpa.ReqContentType = project.HTTPContentTypeJSON
    tpa.ReqMethod = fasthttp.MethodPost
    return tpa
}
