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

// 更新单个设备注册表信息
type registryUpdate struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
    deviceDesc string // 描述
    schemaId   string // 物模型ID
    favourite  bool   // 收藏标识
}

func (ru *registryUpdate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        ru.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (ru *registryUpdate) SetDeviceDesc(deviceDesc string) {
    if len(deviceDesc) > 0 {
        ru.ExtendData["description"] = deviceDesc
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "描述不合法", nil))
    }
}

func (ru *registryUpdate) SetSchemaId(schemaId string) {
    if len(schemaId) > 0 {
        ru.ExtendData["schemaId"] = schemaId
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "物模型ID不合法", nil))
    }
}

func (ru *registryUpdate) SetFavourite(favourite bool) {
    ru.ExtendData["favourite"] = favourite
}

func (ru *registryUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ru.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    ru.ServiceUri = "/v3/iot/management/device/" + ru.deviceName

    ru.ReqURI = ru.GetServiceUrl() + "?updateRegistry"

    reqBody := mpf.JSONMarshal(ru.ExtendData)
    client, req := ru.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRegistryUpdate() *registryUpdate {
    ru := &registryUpdate{mpiot.NewBaseBaiDu(), "", "", "", true}
    ru.ReqData["updateRegistry"] = ""
    ru.ExtendData["favourite"] = true
    ru.ReqContentType = project.HTTPContentTypeJSON
    ru.ReqMethod = fasthttp.MethodPut
    return ru
}
