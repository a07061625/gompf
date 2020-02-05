/**
 * Createg by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package endpoint

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取endpoint列表
type endpointList struct {
    mpiot.BaseBaiDu
}

func (el *endpointList) SetOrder(order string) {
    if (order == "asc") || (order == "desc") {
        el.ReqData["order"] = order
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (el *endpointList) SetOrderBy(orderBy string) {
    if (orderBy == "createTime") || (orderBy == "name") {
        el.ReqData["orderBy"] = orderBy
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序字段不合法", nil))
    }
}

func (el *endpointList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        el.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (el *endpointList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 200) {
        el.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (el *endpointList) SetQuery(query string) {
    if len(query) > 0 {
        el.ReqData["q"] = query
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "模糊查询内容不合法", nil))
    }
}

func (el *endpointList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    el.ReqUrl = el.GetServiceUrl() + "?" + mpf.HttpCreateParams(el.ReqData, "none", 1)

    return el.GetRequest()
}

func NewEndpointList() *endpointList {
    el := &endpointList{mpiot.NewBaseBaiDu()}
    el.ServiceUri = "/v1/endpoint"
    el.ReqData["order"] = "desc"
    el.ReqData["orderBy"] = "createTime"
    el.ReqData["pageNo"] = "1"
    el.ReqData["pageSize"] = "50"
    return el
}
