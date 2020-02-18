/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package thing

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取指定的thing信息
type thingGet struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    thingName    string // thing名称
}

func (tg *thingGet) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        tg.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (tg *thingGet) SetThingName(thingName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, thingName)
    if match {
        tg.thingName = thingName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不合法", nil))
    }
}

func (tg *thingGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tg.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(tg.thingName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不能为空", nil))
    }
    tg.ServiceUri = "/v1/endpoint/" + tg.endpointName + "/thing/" + tg.thingName

    tg.ReqURI = tg.GetServiceUrl()

    return tg.GetRequest()
}

func NewThingGet() *thingGet {
    tg := &thingGet{mpiot.NewBaseBaiDu(), "", ""}
    return tg
}
