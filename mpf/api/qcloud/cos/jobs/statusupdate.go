/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 15:36
 */
package jobs

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新任务状态
type statusUpdate struct {
    qcloud.BaseCos
    jobId     string // 任务ID
    jobStatus string // 任务状态
}

func (su *statusUpdate) SetJobId(jobId string) {
    if len(jobId) > 0 {
        su.ReqUri = "/jobs/" + jobId + "/status"
        su.jobId = jobId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (su *statusUpdate) SetJobStatus(jobStatus string) {
    if (jobStatus == "Ready") || (jobStatus == "Cancelled") {
        su.SetParamData("requestedJobStatus", jobStatus)
        su.jobStatus = jobStatus
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务状态不合法", nil))
    }
}

func (su *statusUpdate) SetUpdateReason(updateReason string) {
    if (len(updateReason) > 0) && (len(updateReason) <= 256) {
        su.SetParamData("statusUpdateReason", updateReason)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "更新原因不合法", nil))
    }
}

func (su *statusUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(su.jobId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不能为空", nil))
    }
    if len(su.jobStatus) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务状态不能为空", nil))
    }
    su.ReqUrl = "http://" + su.ReqHeader["Host"] + su.ReqUri + "?" + mpf.HttpCreateParams(su.ReqData, "none", 1)
    return su.GetRequest()
}

func NewStatusUpdate() *statusUpdate {
    conf := qcloud.NewConfig().GetCos()
    su := &statusUpdate{qcloud.NewCos(), "", ""}
    su.ReqMethod = fasthttp.MethodPost
    su.SetHeaderData("Host", conf.GetControlHost())
    su.SetHeaderData("x-cos-appid", conf.GetAppId())
    return su
}
