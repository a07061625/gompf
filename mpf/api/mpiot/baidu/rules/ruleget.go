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

// 获取规则详情
type ruleGet struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (rg *ruleGet) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        rg.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (rg *ruleGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rg.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    rg.ServiceUri = "/v3/iot/rules/device/" + rg.deviceName

    rg.ReqURI = rg.GetServiceUrl()

    return rg.GetRequest()
}

func NewRuleGet() *ruleGet {
    rg := &ruleGet{mpiot.NewBaseBaiDu(), ""}
    return rg
}
