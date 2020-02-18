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

// 创建带TSDB格式的规则
type formatCreate struct {
    mpiot.BaseBaiDu
    deviceName      string                   // 设备名称
    ruleName        string                   // 规则名称
    sourceList      []map[string]interface{} // 存储数据配置列表
    destinationList []map[string]interface{} // 写入数据配置列表
    format          map[string]interface{}   // TSDB数据格式
}

func (fc *formatCreate) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        fc.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (fc *formatCreate) SetRuleName(ruleName string) {
    if len(ruleName) > 0 {
        fc.ruleName = ruleName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称不合法", nil))
    }
}

func (fc *formatCreate) SetSourceList(sourceList []map[string]interface{}) {
    if len(sourceList) > 0 {
        fc.sourceList = sourceList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "存储数据配置列表不合法", nil))
    }
}

func (fc *formatCreate) SetDestinationList(destinationList []map[string]interface{}) {
    if len(destinationList) > 0 {
        fc.destinationList = destinationList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "写入数据配置列表不合法", nil))
    }
}

func (fc *formatCreate) SetFormat(format map[string]interface{}) {
    fc.format = format
}

func (fc *formatCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(fc.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if len(fc.ruleName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称不能为空", nil))
    }
    if len(fc.sourceList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "存储数据配置不能为空", nil))
    }
    if len(fc.destinationList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "写入数据配置不能为空", nil))
    }
    if len(fc.format) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "TSDB数据格式不能为空", nil))
    }
    fc.ServiceUri = "/v3/iot/rules/device/" + fc.deviceName + "/format"
    fc.ExtendData["name"] = fc.ruleName
    fc.ExtendData["sources"] = fc.sourceList
    fc.ExtendData["destinations"] = fc.destinationList
    fc.ExtendData["format"] = fc.format

    fc.ReqURI = fc.GetServiceUrl()

    reqBody := mpf.JSONMarshal(fc.ExtendData)
    client, req := fc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewFormatCreate() *formatCreate {
    fc := &formatCreate{mpiot.NewBaseBaiDu(), "", "", make([]map[string]interface{}, 0), make([]map[string]interface{}, 0), make(map[string]interface{})}
    fc.ReqContentType = project.HTTPContentTypeJSON
    fc.ReqMethod = fasthttp.MethodPost
    return fc
}
