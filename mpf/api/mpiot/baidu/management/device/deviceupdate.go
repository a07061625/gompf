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

// 修改设备属性
type deviceUpdate struct {
    mpiot.BaseBaiDu
    deviceName      string                 // 设备名称
    deviceAttribute map[string]interface{} // 设备属性
    attributeList   map[string]interface{} // 属性列表
}

func (du *deviceUpdate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        du.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (du *registryUpdate) SetDeviceAttribute(deviceAttribute map[string]interface{}) {
    du.ExtendData["device"] = deviceAttribute
}

func (du *registryUpdate) SetAttributeList(attributeList map[string]interface{}) {
    du.ExtendData["attributes"] = attributeList
}

func (du *deviceUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(du.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if len(du.ExtendData) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备属性和属性列表不能都为空", nil))
    }
    du.ServiceUri = "/v3/iot/management/device/" + du.deviceName

    du.ReqUrl = du.GetServiceUrl() + "?updateProfile"

    reqBody := mpf.JSONMarshal(du.ExtendData)
    client, req := du.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceUpdate() *deviceUpdate {
    du := &deviceUpdate{mpiot.NewBaseBaiDu(), "", make(map[string]interface{}), make(map[string]interface{})}
    du.ReqData["updateProfile"] = ""
    du.ReqContentType = project.HTTPContentTypeJSON
    du.ReqMethod = fasthttp.MethodPut
    return du
}
