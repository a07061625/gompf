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

// 删除权限组
type domainDelete struct {
    mpiot.BaseBaiDu
    domainName string // 权限组名称
}

func (dd *domainDelete) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        dd.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不合法", nil))
    }
}

func (dd *domainDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dd.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不能为空", nil))
    }
    dd.ServiceUri = "/v3/iot/management/domain/" + dd.domainName

    dd.ReqUrl = dd.GetServiceUrl()

    return dd.GetRequest()
}

func NewDomainDelete() *domainDelete {
    dd := &domainDelete{mpiot.NewBaseBaiDu(), ""}
    dd.ReqMethod = fasthttp.MethodDelete
    return dd
}
