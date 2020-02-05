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

// 在policy下创建topic
type permissionCreate struct {
    mpiot.BaseBaiDu
    endpointName  string   // endpoint名称
    policyName    string   // policy名称
    operationList []string // 操作列表
    topic         string   // 主题名称
}

func (pc *permissionCreate) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pc.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pc *permissionCreate) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        pc.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (pc *permissionCreate) SetTopic(topic string) {
    if len(topic) > 0 {
        pc.topic = topic
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名称不合法", nil))
    }
}

func (pc *permissionCreate) SetOperationList(operationList []string) {
    if len(operationList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "操作列表不能为空", nil))
    }
    pc.operationList = make([]string, 0)
    for _, v := range operationList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            pc.operationList = append(pc.operationList, v)
        }
    }
}

func (pc *permissionCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pc.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    if len(pc.topic) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "主题名称不能为空", nil))
    }
    if len(pc.operationList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "操作列表不能为空", nil))
    }
    pc.ServiceUri = "/v1/endpoint/" + pc.endpointName + "/permission"
    pc.ExtendData["policyName"] = pc.policyName
    pc.ExtendData["topic"] = pc.topic
    pc.ExtendData["operations"] = pc.operationList

    pc.ReqUrl = pc.GetServiceUrl()

    reqBody := mpf.JsonMarshal(pc.ExtendData)
    client, req := pc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPermissionCreate() *permissionCreate {
    pc := &permissionCreate{mpiot.NewBaseBaiDu(), "", "", make([]string, 0), ""}
    pc.ReqContentType = project.HttpContentTypeJson
    pc.ReqMethod = fasthttp.MethodPost
    return pc
}
