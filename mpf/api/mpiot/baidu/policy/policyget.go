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

// 获取指定的policy信息
type policyGet struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    policyName   string // policy名称
}

func (pg *policyGet) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pg.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pg *policyGet) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        pg.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (pg *policyGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pg.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pg.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    pg.ServiceUri = "/v1/endpoint/" + pg.endpointName + "/policy/" + pg.policyName

    pg.ReqUrl = pg.GetServiceUrl()

    return pg.GetRequest()
}

func NewPolicyGet() *policyGet {
    pg := &policyGet{mpiot.NewBaseBaiDu(), "", ""}
    return pg
}
