/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package domain

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 权限组中更改设备
type domainDevice struct {
    mpiot.BaseBaiDu
    domainName    string   // 权限组名称
    addDevices    []string // 添加设备列表
    removeDevices []string // 移除设备列表
}

func (dd *domainDevice) SetDomainName(domainName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, domainName)
    if match {
        dd.domainName = domainName
    } else {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "名称不合法", nil))
    }
}

func (dd *domainDevice) SetAddDevices(addDevices []string) {
    if len(addDevices) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "添加设备列表不合法", nil))
    }
    dd.addDevices = make([]string, 0)
    for _, v := range addDevices {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            dd.addDevices = append(dd.addDevices, v)
        }
    }
}

func (dd *domainDevice) SetRemoveDevices(removeDevices []string) {
    if len(removeDevices) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "移除设备列表不合法", nil))
    }
    dd.removeDevices = make([]string, 0)
    for _, v := range removeDevices {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            dd.removeDevices = append(dd.removeDevices, v)
        }
    }
}

func (dd *domainDevice) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dd.domainName) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "权限组名称不能为空", nil))
    }
    if (len(dd.addDevices) == 0) && (len(dd.removeDevices) == 0) {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "添加和移除的设备列表不能都为空", nil))
    }
    dd.ServiceUri = "/v3/iot/management/domain/" + dd.domainName
    if len(dd.addDevices) > 0 {
        dd.ExtendData["addedDevices"] = dd.addDevices
    }
    if len(dd.removeDevices) > 0 {
        dd.ExtendData["removedDevices"] = dd.removeDevices
    }

    dd.ReqUrl = dd.GetServiceUrl() + "?modify"

    reqBody := mpf.JSONMarshal(dd.ExtendData)
    client, req := dd.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDomainDevice() *domainDevice {
    dd := &domainDevice{mpiot.NewBaseBaiDu(), "", make([]string, 0), make([]string, 0)}
    dd.ReqData["modify"] = ""
    dd.ReqContentType = project.HTTPContentTypeJSON
    dd.ReqMethod = fasthttp.MethodPut
    return dd
}
