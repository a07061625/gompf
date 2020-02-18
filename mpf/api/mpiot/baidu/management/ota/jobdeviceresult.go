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

// 查询升级任务各设备升级结果
type jobDeviceResult struct {
    mpiot.BaseBaiDu
    jobId string // 任务ID
}

func (jdr *jobDeviceResult) SetJobId(jobId string) {
    if len(jobId) > 0 {
        jdr.jobId = jobId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务ID不合法", nil))
    }
}

func (jdr *jobDeviceResult) SetPageNo(pageNo int) {
    if pageNo > 0 {
        jdr.ReqData["pageNo"] = strconv.Itoa(pageNo)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "页码不合法", nil))
    }
}

func (jdr *jobDeviceResult) SetPageSize(pageSize int) {
    if (pageSize > 0) && (pageSize <= 100) {
        jdr.ReqData["pageSize"] = strconv.Itoa(pageSize)
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "每页个数不合法", nil))
    }
}

func (jdr *jobDeviceResult) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(jdr.jobId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务ID不能为空", nil))
    }
    jdr.ServiceUri = "/v3/iot/management/ota/job/" + jdr.jobId + "/device-result"

    jdr.ReqURI = jdr.GetServiceUrl() + "?" + mpf.HTTPCreateParams(jdr.ReqData, "none", 1)

    return jdr.GetRequest()
}

func NewJobDeviceResult() *jobDeviceResult {
    jdr := &jobDeviceResult{mpiot.NewBaseBaiDu(), ""}
    jdr.ReqData["pageNo"] = "1"
    jdr.ReqData["pageSize"] = "10"
    return jdr
}
