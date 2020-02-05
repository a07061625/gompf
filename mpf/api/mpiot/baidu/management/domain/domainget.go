/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package domain

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取权限组详情
type domainGet struct {
    mpiot.BaseBaiDu
    domainName string // 权限组名称
}

func (dg *domainGet) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        dg.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不合法", nil))
    }
}

func (dg *domainGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dg.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不能为空", nil))
    }
    dg.ServiceUri = "/v3/iot/management/domain/" + dg.domainName

    dg.ReqUrl = dg.GetServiceUrl()

    return dg.GetRequest()
}

func NewDomainGet() *domainGet {
    dg := &domainGet{mpiot.NewBaseBaiDu(), ""}
    return dg
}
