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

// 更改密钥
type secretKeyUpdate struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (sku *secretKeyUpdate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        sku.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (sku *secretKeyUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sku.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    sku.ServiceUri = "/v3/iot/management/device/" + sku.deviceName

    sku.ReqURI = sku.GetServiceUrl() + "?updateSecretKey"

    return sku.GetRequest()
}

func NewSecretKeyUpdate() *secretKeyUpdate {
    sku := &secretKeyUpdate{mpiot.NewBaseBaiDu(), ""}
    sku.ReqData["updateSecretKey"] = ""
    sku.ReqMethod = fasthttp.MethodPut
    return sku
}
