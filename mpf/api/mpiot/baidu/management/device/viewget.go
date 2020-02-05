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

// 获取设备Profile和模型合并后的View
type viewGet struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (vg *viewGet) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        vg.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (vg *viewGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(vg.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    vg.ServiceUri = "/v3/iot/management/deviceView/" + vg.deviceName

    vg.ReqUrl = vg.GetServiceUrl()

    return vg.GetRequest()
}

func NewViewGet() *viewGet {
    vg := &viewGet{mpiot.NewBaseBaiDu(), ""}
    return vg
}
