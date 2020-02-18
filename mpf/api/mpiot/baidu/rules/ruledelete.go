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

// 删除规则
type ruleDelete struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (rd *ruleDelete) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        rd.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (rd *ruleDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rd.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    rd.ServiceUri = "/v3/iot/rules/device/" + rd.deviceName

    rd.ReqURI = rd.GetServiceUrl()

    return rd.GetRequest()
}

func NewRuleDelete() *ruleDelete {
    rd := &ruleDelete{mpiot.NewBaseBaiDu(), ""}
    rd.ReqMethod = fasthttp.MethodDelete
    return rd
}
