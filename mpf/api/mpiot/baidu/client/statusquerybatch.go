/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:31
 */
package client

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取所有MQTT客户端在线状态
type statusQueryBatch struct {
    mpiot.BaseBaiDu
    endpointName string   // endpoint名称
    clientList   []string // 客户端列表
}

func (sqb *statusQueryBatch) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        sqb.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (sqb *statusQueryBatch) SetClientList(clientList []string) {
    if len(clientList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "客户端列表不合法", nil))
    }

    sqb.clientList = make([]string, 0)
    for _, v := range clientList {
        if len(v) > 0 {
            sqb.clientList = append(sqb.clientList, v)
        }
    }
}

func (sqb *statusQueryBatch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sqb.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(sqb.clientList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "客户端列表不能为空", nil))
    }
    sqb.ServiceUri = "/v2/endpoint/" + sqb.endpointName + "/batch-client-query/status"
    sqb.ExtendData["mqttID"] = sqb.clientList

    sqb.ReqUrl = sqb.GetServiceUrl()

    reqBody := mpf.JSONMarshal(sqb.ExtendData)
    client, req := sqb.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewStatusQueryBatch() *statusQueryBatch {
    sqb := &statusQueryBatch{mpiot.NewBaseBaiDu(), "", make([]string, 0)}
    sqb.ReqContentType = project.HTTPContentTypeJSON
    sqb.ReqMethod = fasthttp.MethodPost
    return sqb
}
