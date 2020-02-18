/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 15:24
 */
package jobs

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新任务优先级
type priorityUpdate struct {
    qcloud.BaseCos
    jobId    string // 任务ID
    priority int    // 任务优先级
}

func (pu *priorityUpdate) SetJobId(jobId string) {
    if len(jobId) > 0 {
        pu.ReqUri = "/jobs/" + jobId + "/priority"
        pu.jobId = jobId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (pu *priorityUpdate) SetPriority(priority int) {
    if priority >= 0 {
        pu.SetParamData("priority", strconv.Itoa(priority))
        pu.priority = priority
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (pu *priorityUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pu.jobId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不能为空", nil))
    }
    if pu.priority < 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务优先级不能为空", nil))
    }
    pu.ReqURI = "http://" + pu.ReqHeader["Host"] + pu.ReqUri + "?" + mpf.HTTPCreateParams(pu.ReqData, "none", 1)
    return pu.GetRequest()
}

func NewPriorityUpdate() *priorityUpdate {
    conf := qcloud.NewConfig().GetCos()
    pu := &priorityUpdate{qcloud.NewCos(), "", 0}
    pu.priority = -1
    pu.ReqMethod = fasthttp.MethodPost
    pu.SetHeaderData("Host", conf.GetControlHost())
    pu.SetHeaderData("x-cos-appid", conf.GetAppId())
    return pu
}
