/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 15:17
 */
package jobs

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取任务信息
type jobDescribe struct {
    qcloud.BaseCos
    jobId string // 任务ID
}

func (jd *jobDescribe) SetJobId(jobId string) {
    if len(jobId) > 0 {
        jd.ReqUri = "/jobs/" + jobId
        jd.jobId = jobId
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不合法", nil))
    }
}

func (jd *jobDescribe) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(jd.jobId) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务ID不能为空", nil))
    }

    jd.ReqUrl = "http://" + jd.ReqHeader["Host"] + jd.ReqUri
    return jd.GetRequest()
}

func NewJobDescribe() *jobDescribe {
    conf := qcloud.NewConfig().GetCos()
    jd := &jobDescribe{qcloud.NewCos(), ""}
    jd.SetHeaderData("Host", conf.GetControlHost())
    jd.SetHeaderData("x-cos-appid", conf.GetAppId())
    return jd
}
