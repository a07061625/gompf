/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package rules

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改规则
type ruleModify struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (rm *ruleModify) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        rm.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (rm *ruleModify) SetRuleName(ruleName string) {
    if len(ruleName) > 0 {
        rm.ExtendData["name"] = ruleName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称不合法", nil))
    }
}

func (rm *ruleModify) SetSourceList(sourceList []map[string]interface{}) {
    if len(sourceList) > 0 {
        rm.ExtendData["sources"] = sourceList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "存储数据配置列表不合法", nil))
    }
}

func (rm *ruleModify) SetDestinationList(destinationList []map[string]interface{}) {
    if len(destinationList) > 0 {
        rm.ExtendData["destinations"] = destinationList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "写入数据配置列表不合法", nil))
    }
}

func (rm *ruleModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rm.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if len(rm.ExtendData) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称,存储数据配置和写入数据配置不能都为空", nil))
    }
    rm.ServiceUri = "/v3/iot/rules/device/" + rm.deviceName

    rm.ReqUrl = rm.GetServiceUrl()

    reqBody := mpf.JsonMarshal(rm.ExtendData)
    client, req := rm.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRuleModify() *ruleModify {
    rm := &ruleModify{mpiot.NewBaseBaiDu(), ""}
    rm.ReqContentType = project.HTTPContentTypeJSON
    rm.ReqMethod = fasthttp.MethodPut
    return rm
}
