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

// 启用一条规则
type ruleEnable struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (re *ruleEnable) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        re.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (re *ruleEnable) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(re.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    re.ServiceUri = "/v3/iot/rules/device/" + re.deviceName

    re.ReqURI = re.GetServiceUrl() + "?enable"

    return re.GetRequest()
}

func NewRuleEnable() *ruleEnable {
    re := &ruleEnable{mpiot.NewBaseBaiDu(), ""}
    re.ReqData["enable"] = ""
    re.ReqMethod = fasthttp.MethodPut
    return re
}
