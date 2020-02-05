/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package rules

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取带TSDB格式的规则详情
type formatGet struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (fg *formatGet) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        fg.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (fg *formatGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(fg.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    fg.ServiceUri = "/v3/iot/rules/device/" + fg.deviceName + "/format"

    fg.ReqUrl = fg.GetServiceUrl()

    return fg.GetRequest()
}

func NewFormatGet() *formatGet {
    fg := &formatGet{mpiot.NewBaseBaiDu(), ""}
    return fg
}
