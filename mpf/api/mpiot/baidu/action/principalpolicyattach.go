/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 21:57
 */
package action

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 给一个Principal添加一个Policy
type principalPolicyAttach struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    principalName string // principal名称
    policyName    string // policy名称
}

func (ppa *principalPolicyAttach) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        ppa.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (ppa *principalPolicyAttach) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        ppa.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (ppa *principalPolicyAttach) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        ppa.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (ppa *principalPolicyAttach) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ppa.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(ppa.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    if len(ppa.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    ppa.ExtendData["endpointName"] = ppa.endpointName
    ppa.ExtendData["principalName"] = ppa.principalName
    ppa.ExtendData["policyName"] = ppa.policyName

    ppa.ReqUrl = ppa.GetServiceUrl()

    reqBody := mpf.JSONMarshal(ppa.ExtendData)
    client, req := ppa.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPrincipalPolicyAttach() *principalPolicyAttach {
    ppa := &principalPolicyAttach{mpiot.NewBaseBaiDu(), "", "", ""}
    ppa.ServiceUri = "/v1/action/attach-principal-policy"
    ppa.ReqContentType = project.HTTPContentTypeJSON
    ppa.ReqMethod = fasthttp.MethodPost
    return ppa
}
