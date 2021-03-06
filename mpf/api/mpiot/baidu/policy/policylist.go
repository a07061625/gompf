/**
 * Createg by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package policy

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取policy列表
type policyList struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
}

func (pl *policyList) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pl.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pl *policyList) SetPrincipalName(principalName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, principalName)
    if match {
        pl.ReqData["principalName"] = principalName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "principal名称不合法", nil))
    }
}

func (pl *policyList) SetOrder(order string) {
    if (order == "asc") || (order == "desc") {
        pl.ReqData["order"] = order
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (pl *policyList) SetOrderBy(orderBy string) {
    if (orderBy == "createTime") || (orderBy == "name") {
        pl.ReqData["orderBy"] = orderBy
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序字段不合法", nil))
    }
}

func (pl *policyList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        pl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (pl *policyList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 200) {
        pl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (pl *policyList) SetQuery(query string) {
    if len(query) > 0 {
        pl.ReqData["q"] = query
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "模糊查询内容不合法", nil))
    }
}

func (pl *policyList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pl.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    pl.ServiceUri = "/v1/endpoint/" + pl.endpointName + "/policy"

    pl.ReqURI = pl.GetServiceUrl() + "?" + mpf.HTTPCreateParams(pl.ReqData, "none", 1)

    return pl.GetRequest()
}

func NewPolicyList() *policyList {
    pl := &policyList{mpiot.NewBaseBaiDu(), ""}
    pl.ReqData["order"] = "desc"
    pl.ReqData["orderBy"] = "createTime"
    pl.ReqData["pageNo"] = "1"
    pl.ReqData["pageSize"] = "50"
    return pl
}
