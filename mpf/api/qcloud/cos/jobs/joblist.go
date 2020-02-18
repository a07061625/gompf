/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 23:30
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

// 批量处理任务列表
type jobList struct {
    qcloud.BaseCos
}

func (jl *jobList) SetJobStatus(jobStatus string) {
    if len(jobStatus) > 0 {
        jl.SetParamData("jobStatuses", jobStatus)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务状态不合法", nil))
    }
}

func (jl *jobList) SetMaxResult(maxResult int) {
    if (maxResult > 0) && (maxResult <= 1000) {
        jl.SetParamData("maxResults", strconv.Itoa(maxResult))
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务最大数量不合法", nil))
    }
}

func (jl *jobList) SetNextToken(nextToken string) {
    if (len(nextToken) > 0) && (len(nextToken) <= 64) {
        jl.SetParamData("nextToken", nextToken)
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "翻页符不合法", nil))
    }
}

func (jl *jobList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    jl.ReqURI = "http://" + jl.ReqHeader["Host"] + jl.ReqUri + "?" + mpf.HTTPCreateParams(jl.ReqData, "none", 1)
    return jl.GetRequest()
}

func NewJobList() *jobList {
    conf := qcloud.NewConfig().GetCos()
    jl := &jobList{qcloud.NewCos()}
    jl.ReqUri = "/jobs"
    jl.SetHeaderData("Host", conf.GetControlHost())
    jl.SetHeaderData("x-cos-appid", conf.GetAppId())
    jl.SetParamData("maxResults", "100")
    return jl
}
