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

// 创建权限组
type domainCreate struct {
    mpiot.BaseBaiDu
    domainName string // 权限组名称
    domainDesc string // 权限组描述
    domainType string // 权限组类型
}

func (dc *domainCreate) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        dc.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "名称不合法", nil))
    }
}

func (dc *domainCreate) SetDomainDesc(domainDesc string) {
    if len(domainDesc) > 0 {
        dc.domainDesc = domainDesc
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不合法", nil))
    }
}

func (dc *domainCreate) SetDomainType(domainType string) {
    if (domainType == "ROOT") || (domainType == "NORMAL") {
        dc.domainType = domainType
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "类型不合法", nil))
    }
}

func (dc *domainCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dc.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "名称不能为空", nil))
    }
    if len(dc.domainDesc) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不能为空", nil))
    }
    if len(dc.domainType) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "类型不能为空", nil))
    }
    dc.ExtendData["name"] = dc.domainName
    dc.ExtendData["description"] = dc.domainDesc
    dc.ExtendData["type"] = dc.domainType

    dc.ReqURI = dc.GetServiceUrl()

    reqBody := mpf.JSONMarshal(dc.ExtendData)
    client, req := dc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDomainCreate() *domainCreate {
    dc := &domainCreate{mpiot.NewBaseBaiDu(), "", "", ""}
    dc.ServiceUri = "/v3/iot/management/domain"
    dc.ReqContentType = project.HTTPContentTypeJSON
    dc.ReqMethod = fasthttp.MethodPost
    return dc
}
