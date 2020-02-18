/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取设备Profile
type deviceGet struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (dg *deviceGet) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        dg.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (dg *deviceGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dg.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    dg.ServiceUri = "/v3/iot/management/device/" + dg.deviceName

    dg.ReqURI = dg.GetServiceUrl()

    return dg.GetRequest()
}

func NewDeviceGet() *deviceGet {
    dg := &deviceGet{mpiot.NewBaseBaiDu(), ""}
    return dg
}
