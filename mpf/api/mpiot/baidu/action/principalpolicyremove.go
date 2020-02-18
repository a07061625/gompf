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

// 从一个Principal移除一个Policy
type principalPolicyRemove struct {
    mpiot.BaseBaiDu
    endpointName  string // endpoint名称
    principalName string // principal名称
    policyName    string // policy名称
}

func (ppr *principalPolicyRemove) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        ppr.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (ppr *principalPolicyRemove) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        ppr.principalName = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (ppr *principalPolicyRemove) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        ppr.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (ppr *principalPolicyRemove) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ppr.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(ppr.principalName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不能为空", nil))
    }
    if len(ppr.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    ppr.ExtendData["endpointName"] = ppr.endpointName
    ppr.ExtendData["principalName"] = ppr.principalName
    ppr.ExtendData["policyName"] = ppr.policyName

    ppr.ReqURI = ppr.GetServiceUrl()

    reqBody := mpf.JSONMarshal(ppr.ExtendData)
    client, req := ppr.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewPrincipalPolicyRemove() *principalPolicyRemove {
    ppr := &principalPolicyRemove{mpiot.NewBaseBaiDu(), "", "", ""}
    ppr.ServiceUri = "/v1/action/remove-principal-policy"
    ppr.ReqContentType = project.HTTPContentTypeJSON
    ppr.ReqMethod = fasthttp.MethodPost
    return ppr
}
