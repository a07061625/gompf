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

// 获取设备接入详情
type accessDetail struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (ad *accessDetail) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        ad.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (ad *accessDetail) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ad.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    ad.ServiceUri = "/v3/iot/management/device/" + ad.deviceName + "/accessDetail"

    ad.ReqURI = ad.GetServiceUrl()

    return ad.GetRequest()
}

func NewAccessDetail() *accessDetail {
    ad := &accessDetail{mpiot.NewBaseBaiDu(), ""}
    return ad
}
