/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 17:36
 */
package ota

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建升级任务
type jobCreate struct {
    mpiot.BaseBaiDu
    jobName    string // 任务名称
    firmwareId string // 固件包ID
}

func (jc *jobCreate) SetJobName(jobName string) {
    if len(jobName) > 0 {
        jc.jobName = jobName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务名称不合法", nil))
    }
}

func (jc *jobCreate) SetFirmwareId(firmwareId string) {
    if len(firmwareId) > 0 {
        jc.firmwareId = firmwareId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包ID不合法", nil))
    }
}

func (jc *jobCreate) SetDescription(description string) {
    if len(description) > 0 {
        jc.ExtendData["description"] = description
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不合法", nil))
    }
}

func (jc *jobCreate) SetDeviceList(deviceList []string) {
    devices := make([]string, 0)
    for _, v := range deviceList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            devices = append(devices, v)
        }
    }

    if len(devices) > 0 {
        jc.ExtendData["devices"] = devices
    }
}

func (jc *jobCreate) SetVersionList(versionList []string) {
    versions := make([]string, 0)
    for _, v := range versionList {
        if len(v) > 0 {
            versions = append(versions, v)
        }
    }

    if len(versions) > 0 {
        jc.ExtendData["versionFilter"] = versions
    }
}

func (jc *jobCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(jc.jobName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "任务名称不能为空", nil))
    }
    if len(jc.firmwareId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "固件包ID不能为空", nil))
    }
    jc.ExtendData["jobName"] = jc.jobName
    jc.ExtendData["firmwareId"] = jc.firmwareId

    jc.ReqURI = jc.GetServiceUrl()

    reqBody := mpf.JSONMarshal(jc.ExtendData)
    client, req := jc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewJobCreate() *jobCreate {
    jc := &jobCreate{mpiot.NewBaseBaiDu(), "", ""}
    jc.ServiceUri = "/v3/iot/management/ota/job"
    jc.ReqContentType = project.HTTPContentTypeJSON
    jc.ReqMethod = fasthttp.MethodPost
    return jc
}
