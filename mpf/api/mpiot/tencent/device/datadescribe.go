/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
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

// 获取设备属性数据
type dataDescribe struct {
    mpiot.BaseTencent
    productId  string // 产品ID
    deviceName string // 设备名称
}

func (dd *dataDescribe) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        dd.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (dd *dataDescribe) SetDeviceName(deviceName string) {
    if len(deviceName) > 0 {
        dd.deviceName = deviceName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "设备名称不合法", nil))
    }
}

func (dd *dataDescribe) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dd.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    if len(dd.deviceName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "设备名称不能为空", nil))
    }
    dd.ReqData["ProductId"] = dd.productId
    dd.ReqData["DeviceName"] = dd.deviceName

    return dd.GetRequest()
}

func NewDataDescribe() *dataDescribe {
    dd := &dataDescribe{mpiot.NewBaseTencent(), "", ""}
    dd.ReqHeader["X-TC-Action"] = "DescribeDeviceData"
    return dd
}
