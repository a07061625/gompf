/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
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

// 设备远程控制
type dataControl struct {
    mpiot.BaseTencent
    productId  string                 // 产品ID
    deviceName string                 // 设备名称
    data       map[string]interface{} // 属性数据
}

func (dc *dataControl) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        dc.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (dc *dataControl) SetDeviceName(deviceName string) {
    if len(deviceName) > 0 {
        dc.deviceName = deviceName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "设备名称不合法", nil))
    }
}

func (dc *dataControl) SetData(data map[string]interface{}) {
    if len(data) > 0 {
        dc.data = data
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "属性数据不合法", nil))
    }
}

func (dc *dataControl) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dc.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    if len(dc.deviceName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "设备名称不能为空", nil))
    }
    if len(dc.data) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "属性数据不能为空", nil))
    }
    dc.ReqData["ProductId"] = dc.productId
    dc.ReqData["DeviceName"] = dc.deviceName
    dc.ReqData["Data"] = mpf.JsonMarshal(dc.data)

    return dc.GetRequest()
}

func NewDataControl() *dataControl {
    dc := &dataControl{mpiot.NewBaseTencent(), "", "", make(map[string]interface{})}
    dc.ReqHeader["X-TC-Action"] = "ControlDeviceData"
    return dc
}
