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

// 重置设备影子
type deviceReset struct {
    mpiot.BaseBaiDu
    deviceList []string // 设备名称列表
}

func (dr *deviceReset) SetDeviceList(deviceList []string) {
    if len(deviceList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称列表不合法", nil))
    }
    dr.deviceList = make([]string, 0)
    for _, v := range deviceList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            dr.deviceList = append(dr.deviceList, v)
        }
    }
}

func (dr *deviceReset) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dr.deviceList) == 0 {
        panic(mperr.NewIotBaiDu(errorcode.IotBaiDuParam, "设备名称列表不能为空", nil))
    }
    dr.ExtendData["devices"] = dr.deviceList

    dr.ReqUrl = dr.GetServiceUrl() + "?reset"

    reqBody := mpf.JSONMarshal(dr.ExtendData)
    client, req := dr.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceReset() *deviceReset {
    dr := &deviceReset{mpiot.NewBaseBaiDu(), make([]string, 0)}
    dr.ServiceUri = "/v3/iot/management/device"
    dr.ReqData["reset"] = ""
    dr.ReqContentType = project.HTTPContentTypeJSON
    dr.ReqMethod = fasthttp.MethodPut
    return dr
}
