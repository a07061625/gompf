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

// 创建规则
type ruleCreate struct {
    mpiot.BaseBaiDu
    deviceName      string                   // 设备名称
    ruleName        string                   // 规则名称
    sourceList      []map[string]interface{} // 存储数据配置列表
    destinationList []map[string]interface{} // 写入数据配置列表
}

func (rc *ruleCreate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        rc.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (rc *ruleCreate) SetRuleName(ruleName string) {
    if len(ruleName) > 0 {
        rc.ruleName = ruleName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称不合法", nil))
    }
}

func (rc *ruleCreate) SetSourceList(sourceList []map[string]interface{}) {
    if len(sourceList) > 0 {
        rc.sourceList = sourceList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "存储数据配置列表不合法", nil))
    }
}

func (rc *ruleCreate) SetDestinationList(destinationList []map[string]interface{}) {
    if len(destinationList) > 0 {
        rc.destinationList = destinationList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "写入数据配置列表不合法", nil))
    }
}

func (rc *ruleCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rc.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if len(rc.ruleName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称不能为空", nil))
    }
    if len(rc.sourceList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "存储数据配置不能为空", nil))
    }
    if len(rc.destinationList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "写入数据配置不能为空", nil))
    }
    rc.ServiceUri = "/v3/iot/rules/device/" + rc.deviceName
    rc.ExtendData["name"] = rc.ruleName
    rc.ExtendData["sources"] = rc.sourceList
    rc.ExtendData["destinations"] = rc.destinationList

    rc.ReqUrl = rc.GetServiceUrl()

    reqBody := mpf.JSONMarshal(rc.ExtendData)
    client, req := rc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRuleCreate() *ruleCreate {
    rc := &ruleCreate{mpiot.NewBaseBaiDu(), "", "", make([]map[string]interface{}, 0), make([]map[string]interface{}, 0)}
    rc.ReqContentType = project.HTTPContentTypeJSON
    rc.ReqMethod = fasthttp.MethodPost
    return rc
}
