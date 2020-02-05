/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package endpoint

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除endpoint
type endpointDelete struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
}

func (ed *endpointDelete) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        ed.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (ed *endpointDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ed.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    ed.ServiceUri = "/v1/endpoint/" + ed.endpointName

    ed.ReqUrl = ed.GetServiceUrl()

    return ed.GetRequest()
}

func NewEndpointDelete() *endpointDelete {
    ed := &endpointDelete{mpiot.NewBaseBaiDu(), ""}
    ed.ReqMethod = fasthttp.MethodDelete
    return ed
}
