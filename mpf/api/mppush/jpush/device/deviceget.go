/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 12:19
 */
package device

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询设备的别名与标签
type deviceGet struct {
    mppush.BaseJPush
    registrationId string // 设备ID
}

func (dg *deviceGet) SetRegistrationId(registrationId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, registrationId)
    if match {
        dg.registrationId = registrationId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备ID不合法", nil))
    }
}

func (dg *deviceGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dg.registrationId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备ID不能为空", nil))
    }
    dg.ServiceUri = "/v3/devices/" + dg.registrationId

    dg.ReqURI = dg.GetServiceUrl()

    return dg.GetRequest()
}

func NewDeviceGet(key string) *deviceGet {
    dg := &deviceGet{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), ""}
    return dg
}
