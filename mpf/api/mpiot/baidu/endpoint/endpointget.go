/**
 * Createg by GoLand.
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

// 获取指定的endpoint信息
type endpointGet struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
}

func (eg *endpointGet) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        eg.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (eg *endpointGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(eg.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    eg.ServiceUri = "/v1/endpoint/" + eg.endpointName

    eg.ReqURI = eg.GetServiceUrl()

    return eg.GetRequest()
}

func NewEndpointGet() *endpointGet {
    eg := &endpointGet{mpiot.NewBaseBaiDu(), ""}
    return eg
}
