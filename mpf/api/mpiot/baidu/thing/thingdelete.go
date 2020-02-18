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

// 删除thing
type thingDelete struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    thingName    string // thing名称
}

func (td *thingDelete) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        td.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (td *thingDelete) SetThingName(thingName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, thingName)
    if match {
        td.thingName = thingName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不合法", nil))
    }
}

func (td *thingDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(td.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(td.thingName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "thing名称不能为空", nil))
    }
    td.ServiceUri = "/v1/endpoint/" + td.endpointName + "/thing/" + td.thingName

    td.ReqURI = td.GetServiceUrl()

    return td.GetRequest()
}

func NewThingDelete() *thingDelete {
    td := &thingDelete{mpiot.NewBaseBaiDu(), "", ""}
    td.ReqMethod = fasthttp.MethodDelete
    return td
}
