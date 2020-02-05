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

// 设置设备的别名与标签
type deviceSet struct {
    mppush.BaseJPush
    registrationId string // 设备ID
}

func (ds *deviceSet) SetRegistrationId(registrationId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, registrationId)
    if match {
        ds.registrationId = registrationId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备ID不合法", nil))
    }
}

func (ds *deviceSet) SetTags(tags map[string]interface{}) {
    if len(tags) > 0 {
        ds.ExtendData["tags"] = tags
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签列表不合法", nil))
    }
}

func (ds *deviceSet) SetAliasName(aliasName string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, aliasName)
    if match {
        ds.ExtendData["alias"] = aliasName
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备别名不合法", nil))
    }
}

func (ds *deviceSet) SetMobile(mobile string) {
    match, _ := regexp.MatchString(project.RegexPhone, mobile)
    if match {
        ds.ExtendData["mobile"] = mobile
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "关联手机号码不合法", nil))
    }
}

func (ds *deviceSet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ds.registrationId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备ID不能为空", nil))
    }
    ds.ServiceUri = "/v3/devices/" + ds.registrationId

    ds.ReqUrl = ds.GetServiceUrl()

    reqBody := mpf.JsonMarshal(ds.ExtendData)
    client, req := ds.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewDeviceSet(key string) *deviceSet {
    ds := &deviceSet{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), ""}
    ds.ExtendData["tags"] = ""
    ds.ExtendData["alias"] = ""
    ds.ExtendData["mobile"] = ""
    ds.ReqContentType = project.HttpContentTypeJson
    ds.ReqMethod = fasthttp.MethodPost
    return ds
}
