/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 8:02
 */
package permission

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

// 获取policy下所有topic信息
type permissionList struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
    policyName   string // policy名称
}

func (pl *permissionList) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        pl.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (pl *permissionList) SetPolicyName(policyName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, policyName)
    if match {
        pl.policyName = policyName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不合法", nil))
    }
}

func (pl *permissionList) SetOrder(order string) {
    if (order == "asc") || (order == "desc") {
        pl.ReqData["order"] = order
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (pl *permissionList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        pl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (pl *permissionList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 200) {
        pl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (pl *permissionList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pl.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    if len(pl.policyName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "policy名称不能为空", nil))
    }
    pl.ServiceUri = "/v1/endpoint/" + pl.endpointName + "/permission"
    pl.ReqData["policyName"] = pl.policyName

    pl.ReqUrl = pl.GetServiceUrl() + "?" + mpf.HttpCreateParams(pl.ReqData, "none", 1)

    return pl.GetRequest()
}

func NewPermissionList() *permissionList {
    pl := &permissionList{mpiot.NewBaseBaiDu(), "", ""}
    pl.ReqData["order"] = "desc"
    pl.ReqData["orderBy"] = "createTime"
    pl.ReqData["pageNo"] = "1"
    pl.ReqData["pageSize"] = "50"
    return pl
}
