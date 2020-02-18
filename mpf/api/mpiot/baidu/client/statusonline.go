/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:31
 */
package client

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取指定MQTT客户端在线状态
type statusOnline struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    clientId     string // 客户端ID
}

func (so *statusOnline) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        so.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (so *statusOnline) SetClientId(clientId string) {
    if len(clientId) > 0 {
        so.clientId = clientId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "客户端ID不合法", nil))
    }
}

func (so *statusOnline) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(so.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(so.clientId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "客户端ID不能为空", nil))
    }
    so.ServiceUri = "/v2/endpoint/" + so.endpointName + "/client/" + so.clientId + "/status/online"
    so.ReqURI = so.GetServiceUrl()

    return so.GetRequest()
}

func NewStatusOnline() *statusOnline {
    so := &statusOnline{mpiot.NewBaseBaiDu(), "", ""}
    return so
}
