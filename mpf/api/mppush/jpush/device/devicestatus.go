/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 12:19
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户在线状态
type deviceStatus struct {
    mppush.BaseJPush
    registrationList []string // 设备列表
}

func (ds *deviceStatus) SetRegistrationList(registrationList []string) {
    if len(registrationList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能为空", nil))
    } else if len(registrationList) > 1000 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能超过1000个", nil))
    }
    ds.registrationList = make([]string, 0)
    for _, v := range registrationList {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            ds.registrationList = append(ds.registrationList, v)
        }
    }
}

func (ds *deviceStatus) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ds.registrationList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能为空", nil))
    }
    ds.ExtendData["registration_ids"] = ds.registrationList

    ds.ReqURI = ds.GetServiceUrl()

    reqBody := mpf.JSONMarshal(ds.ExtendData)
    client, req := ds.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceStatus(key string) *deviceStatus {
    ds := &deviceStatus{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), make([]string, 0)}
    ds.ServiceUri = "/v3/devices/status"
    ds.ReqContentType = project.HTTPContentTypeJSON
    ds.ReqMethod = fasthttp.MethodPost
    return ds
}
