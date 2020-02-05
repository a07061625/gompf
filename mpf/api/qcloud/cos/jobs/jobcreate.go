/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/21 0021
 * Time: 13:05
 */
package jobs

import (
    "github.com/a07061625/gompf/mpf/api/qcloud"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/clbanning/mxj"
    "github.com/valyala/fasthttp"
)

// 创建批量处理任务
type jobCreate struct {
    qcloud.BaseCos
    jobInfo map[string]interface{}
}

func (jc *jobCreate) SetJobInfo(jobInfo map[string]interface{}) {
    if len(jobInfo) > 0 {
        jc.jobInfo = jobInfo
    } else {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务信息不合法", nil))
    }
}

func (jc *jobCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(jc.jobInfo) == 0 {
        panic(mperr.NewQCloudCos(errorcode.QCloudCosParam, "任务信息不能为空", nil))
    }
    xmlData := mxj.Map(jc.jobInfo)
    xmlInfo, _ := xmlData.Xml("CreateJobRequest")
    reqBody := project.DataPrefixXml + string(xmlInfo) + ""
    jc.ReqUrl = "http://" + jc.ReqHeader["Host"] + jc.ReqUri
    client, req := jc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewJobCreate() *jobCreate {
    conf := qcloud.NewConfig().GetCos()
    jc := &jobCreate{qcloud.NewCos(), make(map[string]interface{})}
    jc.ReqUri = "/jobs"
    jc.ReqMethod = fasthttp.MethodPost
    jc.SetHeaderData("Host", conf.GetControlHost())
    jc.SetHeaderData("x-cos-appid", conf.GetAppId())
    return jc
}
