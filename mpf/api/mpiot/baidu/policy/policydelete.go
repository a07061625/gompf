/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 0:18
 */
package policy

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除policy
type policyDelete struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    policyName   string // policy名称
}

func (pd *policyDelete) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pd.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pd *policyDelete) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        pd.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (pd *policyDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pd.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pd.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    pd.ServiceUri = "/v1/endpoint/" + pd.endpointName + "/policy/" + pd.policyName

    pd.ReqUrl = pd.GetServiceUrl()

    return pd.GetRequest()
}

func NewPolicyDelete() *policyDelete {
    pd := &policyDelete{mpiot.NewBaseBaiDu(), "", ""}
    pd.ReqMethod = fasthttp.MethodDelete
    return pd
}
