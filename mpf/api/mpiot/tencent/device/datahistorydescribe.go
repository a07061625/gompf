/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 19:40
 */
package device

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取设备历史数据
type dataHistoryDescribe struct {
    mpiot.BaseTencent
    startTime  int64  // 开始时间
    endTime    int64  // 结束时间
    productId  string // 产品ID
    deviceName string // 设备名称
    fieldName  string // 属性字段名称
}

func (dhd *dataHistoryDescribe) SetStartAndEndTime(startTime, endTime int) {
    if startTime <= 1000000000 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "开始时间不合法", nil))
    } else if endTime <= startTime {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "开始时间必须小于结束时间", nil))
    }

    dhd.startTime = int64(1000 * startTime)
    dhd.endTime = int64(1000 * endTime)
}

func (dhd *dataHistoryDescribe) SetProductId(productId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, productId)
    if match {
        dhd.productId = productId
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不合法", nil))
    }
}

func (dhd *dataHistoryDescribe) SetDeviceName(deviceName string) {
    if len(deviceName) > 0 {
        dhd.deviceName = deviceName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "设备名称不合法", nil))
    }
}

func (dhd *dataHistoryDescribe) SetFieldName(fieldName string) {
    if len(fieldName) > 0 {
        dhd.fieldName = fieldName
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "属性字段名称不合法", nil))
    }
}

func (dhd *dataHistoryDescribe) SetLimit(limit int) {
    if limit > 0 {
        dhd.ReqData["Limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "返回记录数不合法", nil))
    }
}

func (dhd *dataHistoryDescribe) SetContext(context string) {
    dhd.ReqData["Context"] = strings.TrimSpace(context)
}

func (dhd *dataHistoryDescribe) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if dhd.startTime <= 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "开始时间不能为空", nil))
    }
    if len(dhd.productId) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "产品ID不能为空", nil))
    }
    if len(dhd.deviceName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "设备名称不能为空", nil))
    }
    if len(dhd.fieldName) == 0 {
        panic(mperr.NewIotTencent(errorcode.IotTencentParam, "属性字段名称不能为空", nil))
    }
    dhd.ReqData["MinTime"] = strconv.FormatInt(dhd.startTime, 10)
    dhd.ReqData["MaxTime"] = strconv.FormatInt(dhd.endTime, 10)
    dhd.ReqData["ProductId"] = dhd.productId
    dhd.ReqData["DeviceName"] = dhd.deviceName
    dhd.ReqData["FieldName"] = dhd.fieldName

    return dhd.GetRequest()
}

func NewDataHistoryDescribe() *dataHistoryDescribe {
    dhd := &dataHistoryDescribe{mpiot.NewBaseTencent(), 0, 0, "", "", ""}
    dhd.ReqData["Limit"] = "10"
    dhd.ReqHeader["X-TC-Action"] = "DescribeDeviceDataHistory"
    return dhd
}
