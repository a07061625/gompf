/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package thing

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建thing
type thingCreate struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    thingName    string // thing名称
}

func (tc *thingCreate) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        tc.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (tc *thingCreate) SetThingName(thingName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, thingName)
    if match {
        tc.thingName = thingName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不合法", nil))
    }
}

func (tc *thingCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tc.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(tc.thingName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不能为空", nil))
    }
    tc.ServiceUri = "/v1/endpoint/" + tc.endpointName + "/thing"
    tc.ExtendData["thingName"] = tc.thingName

    tc.ReqUrl = tc.GetServiceUrl()

    reqBody := mpf.JsonMarshal(tc.ExtendData)
    client, req := tc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewThingCreate() *thingCreate {
    tc := &thingCreate{mpiot.NewBaseBaiDu(), "", ""}
    tc.ReqContentType = project.HttpContentTypeJson
    tc.ReqMethod = fasthttp.MethodPost
    return tc
}
