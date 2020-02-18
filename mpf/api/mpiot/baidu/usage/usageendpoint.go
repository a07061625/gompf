/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:05
 */
package usage

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取当前账单月内特定实例的使用量
type usageEndpoint struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
}

func (ue *usageEndpoint) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        ue.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (ue *usageEndpoint) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ue.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    ue.ServiceUri = "/v1/endpoint/" + ue.endpointName + "/usage"

    ue.ReqURI = ue.GetServiceUrl()

    return ue.GetRequest()
}

func NewUsageEndpoint() *usageEndpoint {
    ue := &usageEndpoint{mpiot.NewBaseBaiDu(), ""}
    return ue
}
