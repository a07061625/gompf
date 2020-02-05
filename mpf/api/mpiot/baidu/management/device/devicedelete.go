/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/2 0002
 * Time: 9:29
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mpiot"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除设备
type deviceDelete struct {
    mpiot.BaseBaiDu
    deviceList []string // 设备名称列表
}

func (dd *deviceDelete) SetDeviceList(deviceList []string) {
    if len(deviceList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称列表不合法", nil))
    }
    dd.deviceList = make([]string, 0)
    for _, v := range deviceList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            dd.deviceList = append(dd.deviceList, v)
        }
    }
}

func (dd *deviceDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dd.deviceList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称列表不能为空", nil))
    }
    dd.ExtendData["devices"] = dd.deviceList

    dd.ReqUrl = dd.GetServiceUrl() + "?remove"

    reqBody := mpf.JsonMarshal(dd.ExtendData)
    client, req := dd.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceDelete() *deviceDelete {
    dd := &deviceDelete{mpiot.NewBaseBaiDu(), make([]string, 0)}
    dd.ServiceUri = "/v3/iot/management/device"
    dd.ReqData["remove"] = ""
    dd.ReqContentType = project.HttpContentTypeJson
    dd.ReqMethod = fasthttp.MethodPut
    return dd
}
