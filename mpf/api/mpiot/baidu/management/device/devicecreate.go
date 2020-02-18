/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建单个设备
type deviceCreate struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
    deviceDesc string // 描述
    schemaId   string // 物模型ID
}

func (dc *deviceCreate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        dc.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (dc *deviceCreate) SetDeviceDesc(deviceDesc string) {
    if len(deviceDesc) > 0 {
        dc.deviceDesc = deviceDesc
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不合法", nil))
    }
}

func (dc *deviceCreate) SetSchemaId(schemaId string) {
    if len(schemaId) > 0 {
        dc.schemaId = schemaId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不合法", nil))
    }
}

func (dc *deviceCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dc.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if len(dc.deviceDesc) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不能为空", nil))
    }
    if len(dc.schemaId) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不能为空", nil))
    }
    dc.ExtendData["deviceName"] = dc.deviceName
    dc.ExtendData["description"] = dc.deviceDesc
    dc.ExtendData["schemaId"] = dc.schemaId

    dc.ReqURI = dc.GetServiceUrl()

    reqBody := mpf.JSONMarshal(dc.ExtendData)
    client, req := dc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceCreate() *deviceCreate {
    dc := &deviceCreate{mpiot.NewBaseBaiDu(), "", "", ""}
    dc.ServiceUri = "/v3/iot/management/device"
    dc.ReqContentType = project.HTTPContentTypeJSON
    dc.ReqMethod = fasthttp.MethodPost
    return dc
}
