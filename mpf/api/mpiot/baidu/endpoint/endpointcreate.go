/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package endpoint

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建endpoint
type endpointCreate struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
}

func (ec *endpointCreate) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        ec.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (ec *endpointCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ec.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    ec.ExtendData["endpointName"] = ec.endpointName

    ec.ReqUrl = ec.GetServiceUrl()

    reqBody := mpf.JSONMarshal(ec.ExtendData)
    client, req := ec.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewEndpointCreate() *endpointCreate {
    ec := &endpointCreate{mpiot.NewBaseBaiDu(), ""}
    ec.ServiceUri = "/v1/endpoint"
    ec.ReqContentType = project.HTTPContentTypeJSON
    ec.ReqMethod = fasthttp.MethodPost
    return ec
}
