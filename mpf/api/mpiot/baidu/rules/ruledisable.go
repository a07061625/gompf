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

// 禁用一条规则
type ruleDisable struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (rd *ruleDisable) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        rd.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (rd *ruleDisable) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rd.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    rd.ServiceUri = "/v3/iot/rules/device/" + rd.deviceName

    rd.ReqUrl = rd.GetServiceUrl() + "?disable"

    return rd.GetRequest()
}

func NewRuleDisable() *ruleDisable {
    rd := &ruleDisable{mpiot.NewBaseBaiDu(), ""}
    rd.ReqData["disable"] = ""
    rd.ReqMethod = fasthttp.MethodPut
    return rd
}
