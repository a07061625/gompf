/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/1 0001
 * Time: 23:49
 */
package ota

import (
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询升级任务结果
type jobStatistics struct {
    mpiot.BaseBaiDu
    jobId string // 任务ID
}

func (js *jobStatistics) SetJobId(jobId string) {
    if len(jobId) > 0 {
        js.jobId = jobId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务ID不合法", nil))
    }
}

func (js *jobStatistics) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(js.jobId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务ID不能为空", nil))
    }
    js.ServiceUri = "/v3/iot/management/ota/job/" + js.jobId + "/statistics"

    js.ReqUrl = js.GetServiceUrl()

    return js.GetRequest()
}

func NewJobStatistics() *jobStatistics {
    js := &jobStatistics{mpiot.NewBaseBaiDu(), ""}
    return js
}
