/**
 * Createg by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package domain

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

// 获取权限组设备列表
type deviceList struct {
    mpiot.BaseBaiDu
    domainName string // 权限组名称
}

func (dl *deviceList) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        dl.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不合法", nil))
    }
}

func (dl *deviceList) SetOrder(order string) {
    if (order == "asc") || (order == "desc") {
        dl.ReqData["order"] = order
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (dl *deviceList) SetOrderBy(orderBy string) {
    switch orderBy {
    case "name":
        dl.ReqData["orderBy"] = orderBy
    case "createTime":
        dl.ReqData["orderBy"] = orderBy
    case "lastUpdatedTime":
        dl.ReqData["orderBy"] = orderBy
    default:
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序字段不合法", nil))
    }
}

func (dl *deviceList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        dl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (dl *deviceList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 200) {
        dl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (dl *deviceList) SetName(name string) {
    if len(name) > 0 {
        dl.ReqData["name"] = name
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "属性名不合法", nil))
    }
}

func (dl *deviceList) SetValue(value string) {
    if len(value) > 0 {
        dl.ReqData["value"] = value
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "属性名对应值不合法", nil))
    }
}

func (dl *deviceList) SetFavourite(favourite string) {
    switch favourite {
    case "true":
        dl.ReqData["favourite"] = favourite
    case "false":
        dl.ReqData["favourite"] = favourite
    case "all":
        dl.ReqData["favourite"] = favourite
    default:
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "收藏标识不合法", nil))
    }
}

func (dl *deviceList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dl.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不能为空", nil))
    }
    dl.ServiceUri = "/v3/iot/management/domain/" + dl.domainName + "/devices"

    dl.ReqUrl = dl.GetServiceUrl() + "?" + mpf.HttpCreateParams(dl.ReqData, "none", 1)

    return dl.GetRequest()
}

func NewDeviceList() *deviceList {
    dl := &deviceList{mpiot.NewBaseBaiDu(), ""}
    dl.ReqData["order"] = "asc"
    dl.ReqData["orderBy"] = "name"
    dl.ReqData["pageNo"] = "1"
    dl.ReqData["pageSize"] = "10"
    dl.ReqData["favourite"] = "all"
    return dl
}
