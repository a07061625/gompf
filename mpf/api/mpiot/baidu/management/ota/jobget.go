/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 15:13
 */
package ota

import (
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询升级任务
type jobGet struct {
    mpiot.BaseBaiDu
    jobId string // 任务ID
}

func (jg *jobGet) SetJobId(jobId string) {
    if len(jobId) > 0 {
        jg.jobId = jobId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务ID不合法", nil))
    }
}

func (jg *jobGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(jg.jobId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务ID不能为空", nil))
    }
    jg.ServiceUri = "/v3/iot/management/ota/job/" + jg.jobId

    jg.ReqUrl = jg.GetServiceUrl()

    return jg.GetRequest()
}

func NewJobGet() *jobGet {
    jg := &jobGet{mpiot.NewBaseBaiDu(), ""}
    return jg
}
