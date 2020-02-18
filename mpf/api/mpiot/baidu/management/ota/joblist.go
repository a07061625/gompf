/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package ota

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询升级任务列表
type jobList struct {
    mpiot.BaseBaiDu
}

func (jl *jobList) SetOrderBy(orderBy string) {
    if (orderBy == "asc") || (orderBy == "desc") {
        jl.ReqData["orderBy"] = orderBy
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "排序方式不合法", nil))
    }
}

func (jl *jobList) SetPageNo(pageNo int) {
    if pageNo > 0 {
        jl.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (jl *jobList) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 100) {
        jl.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (jl *jobList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    jl.ReqURI = jl.GetServiceUrl() + "?" + mpf.HTTPCreateParams(jl.ReqData, "none", 1)

    return jl.GetRequest()
}

func NewJobList() *jobList {
    jl := &jobList{mpiot.NewBaseBaiDu()}
    jl.ServiceUri = "/v3/iot/management/ota/job"
    jl.ReqData["order"] = "createdAt"
    jl.ReqData["orderBy"] = "desc"
    jl.ReqData["pageNo"] = "1"
    jl.ReqData["pageSize"] = "10"
    return jl
}
