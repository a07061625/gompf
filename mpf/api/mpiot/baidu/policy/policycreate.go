/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package policy

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建policy
type policyCreate struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    policyName   string // policy名称
}

func (pc *policyCreate) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pc.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pc *policyCreate) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        pc.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (pc *policyCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pc.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pc.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    pc.ServiceUri = "/v1/endpoint/" + pc.endpointName + "/policy"
    pc.ExtendData["policyName"] = pc.policyName

    pc.ReqUrl = pc.GetServiceUrl()

    reqBody := mpf.JsonMarshal(pc.ExtendData)
    client, req := pc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPolicyCreate() *policyCreate {
    pc := &policyCreate{mpiot.NewBaseBaiDu(), "", ""}
    pc.ReqContentType = project.HTTPContentTypeJSON
    pc.ReqMethod = fasthttp.MethodPost
    return pc
}
