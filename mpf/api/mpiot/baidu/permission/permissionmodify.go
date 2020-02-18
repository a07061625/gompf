/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package permission

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新已有的topic设置
type permissionModify struct {
    mpiot.BaseBaiDu
    endpointName   string   // endpoint名称
    permissionUuid string   // permissionID
    operationList  []string // 操作列表
    topic          string   // 主题名称
}

func (pm *permissionModify) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pm.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pm *permissionModify) SetPermissionUuid(permissionUuid string) {
    if len(permissionUuid) > 0 {
        pm.permissionUuid = permissionUuid
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "permissionID不合法", nil))
    }
}

func (pm *permissionModify) SetTopic(topic string) {
    if len(topic) > 0 {
        pm.topic = topic
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名称不合法", nil))
    }
}

func (pm *permissionModify) SetOperationList(operationList []string) {
    if len(operationList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "操作列表不能为空", nil))
    }
    pm.operationList = make([]string, 0)
    for _, v := range operationList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            pm.operationList = append(pm.operationList, v)
        }
    }
}

func (pm *permissionModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pm.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pm.permissionUuid) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "permissionID不能为空", nil))
    }
    if len(pm.topic) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名称不能为空", nil))
    }
    if len(pm.operationList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "操作列表不能为空", nil))
    }
    pm.ServiceUri = "/v1/endpoint/" + pm.endpointName + "/permission/" + pm.permissionUuid
    pm.ExtendData["topic"] = pm.topic
    pm.ExtendData["operations"] = pm.operationList

    pm.ReqUrl = pm.GetServiceUrl()

    reqBody := mpf.JSONMarshal(pm.ExtendData)
    client, req := pm.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPermissionModify() *permissionModify {
    pm := &permissionModify{mpiot.NewBaseBaiDu(), "", "", make([]string, 0), ""}
    pm.ReqContentType = project.HTTPContentTypeJSON
    pm.ReqMethod = fasthttp.MethodPut
    return pm
}
