/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package domain

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新权限组注册信息
type domainModify struct {
    mpiot.BaseBaiDu
    domainName string // 权限组名称
    domainDesc string // 权限组描述
}

func (dm *domainModify) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        dm.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "名称不合法", nil))
    }
}

func (dm *domainModify) SetDomainDesc(domainDesc string) {
    if len(domainDesc) > 0 {
        dm.domainDesc = domainDesc
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不合法", nil))
    }
}

func (dm *domainModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dm.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不能为空", nil))
    }
    if len(dm.domainDesc) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不能为空", nil))
    }
    dm.ServiceUri = "/v3/iot/management/domain/" + dm.domainName
    dm.ExtendData["description"] = dm.domainDesc

    dm.ReqUrl = dm.GetServiceUrl()

    reqBody := mpf.JsonMarshal(dm.ExtendData)
    client, req := dm.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDomainModify() *domainModify {
    dm := &domainModify{mpiot.NewBaseBaiDu(), "", ""}
    dm.ReqContentType = project.HttpContentTypeJson
    dm.ReqMethod = fasthttp.MethodPut
    return dm
}
