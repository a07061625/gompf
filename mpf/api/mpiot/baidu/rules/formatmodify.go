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

// 修改带TSDB格式的规则
type formatModify struct {
    mpiot.BaseBaiDu
    deviceName string // 设备名称
}

func (fm *formatModify) SetDeviceName(deviceName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, deviceName)
    if match {
        fm.deviceName = deviceName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不合法", nil))
    }
}

func (fm *formatModify) SetRuleName(ruleName string) {
    if len(ruleName) > 0 {
        fm.ExtendData["name"] = ruleName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称不合法", nil))
    }
}

func (fm *formatModify) SetSourceList(sourceList []map[string]interface{}) {
    if len(sourceList) > 0 {
        fm.ExtendData["sources"] = sourceList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "存储数据配置列表不合法", nil))
    }
}

func (fm *formatModify) SetDestinationList(destinationList []map[string]interface{}) {
    if len(destinationList) > 0 {
        fm.ExtendData["destinations"] = destinationList
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "写入数据配置列表不合法", nil))
    }
}

func (fm *formatModify) SetFormat(format map[string]interface{}) {
    fm.ExtendData["format"] = format
}

func (fm *formatModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(fm.deviceName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称不能为空", nil))
    }
    if len(fm.ExtendData) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "规则名称,存储数据配置,写入数据配置和TSDB数据格式不能都为空", nil))
    }
    fm.ServiceUri = "/v3/iot/rules/device/" + fm.deviceName + "/format"

    fm.ReqUrl = fm.GetServiceUrl()

    reqBody := mpf.JsonMarshal(fm.ExtendData)
    client, req := fm.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewFormatModify() *formatModify {
    fm := &formatModify{mpiot.NewBaseBaiDu(), ""}
    fm.ReqContentType = project.HttpContentTypeJson
    fm.ReqMethod = fasthttp.MethodPut
    return fm
}
