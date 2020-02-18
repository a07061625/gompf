/**
 * Createg by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package thing

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

// 获取thing列表
type thingList struct {
    mpiot.BaseBaiDu
    endpointName string // endpoint名称
}

func (tl *thingList) SetEndpointName(endpointName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, endpointName)
    if match {
        tl.endpointName = endpointName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不合法", nil))
    }
}

func (tl *thingList) SetOrder(order string) {
    if (order == "asc") || (order == "desc") {
        tl.ReqData["order"] = order
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (tl *thingList) SetOrderBy(orderBy string) {
    if (orderBy == "createTime") || (orderBy == "name") {
        tl.ReqData["orderBy"] = orderBy
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序字段不合法", nil))
    }
}

func (tl *thingList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        tl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (tl *thingList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 200) {
        tl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (tl *thingList) SetQuery(query string) {
    if len(query) > 0 {
        tl.ReqData["q"] = query
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "模糊查询内容不合法", nil))
    }
}

func (tl *thingList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tl.endpointName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "endpoint名称不能为空", nil))
    }
    tl.ServiceUri = "/v1/endpoint/" + tl.endpointName + "/thing"

    tl.ReqURI = tl.GetServiceUrl() + "?" + mpf.HTTPCreateParams(tl.ReqData, "none", 1)

    return tl.GetRequest()
}

func NewThingList() *thingList {
    tl := &thingList{mpiot.NewBaseBaiDu(), ""}
    tl.ReqData["order"] = "desc"
    tl.ReqData["orderBy"] = "createTime"
    tl.ReqData["pageNo"] = "1"
    tl.ReqData["pageSize"] = "50"
    return tl
}
