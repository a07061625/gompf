/**
 * Created by GoLand.
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

// 获取权限组列表
type domainList struct {
    mpiot.BaseBaiDu
}

func (dl *domainList) SetOrder(order string) {
    if (order == "asc") || (order == "desc") {
        dl.ReqData["order"] = order
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (dl *domainList) SetOrderBy(orderBy string) {
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

func (dl *domainList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        dl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (dl *domainList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 200) {
        dl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (dl *domainList) SetKey(key string) {
    if len(key) > 0 {
        dl.ReqData["key"] = key
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "查询关键字不合法", nil))
    }
}

func (dl *domainList) SetDomainType(domainType string) {
    switch domainType {
    case "ROOT":
        dl.ReqData["type"] = domainType
    case "NORMAL":
        dl.ReqData["type"] = domainType
    case "ALL":
        dl.ReqData["type"] = domainType
    default:
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组类型不合法", nil))
    }
}

func (dl *domainList) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        dl.ReqData["deviceName"] = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (dl *domainList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    dl.ReqUrl = dl.GetServiceUrl() + "?" + mpf.HttpCreateParams(dl.ReqData, "none", 1)

    return dl.GetRequest()
}

func NewDomainList() *domainList {
    dl := &domainList{mpiot.NewBaseBaiDu()}
    dl.ServiceUri = "/v3/iot/management/domain"
    dl.ReqData["order"] = "desc"
    dl.ReqData["orderBy"] = "createTime"
    dl.ReqData["pageNo"] = "1"
    dl.ReqData["pageSize"] = "10"
    dl.ReqData["type"] = "ALL"
    return dl
}
