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

// 判断设备与标签绑定关系
type tagRegistration struct {
    mppush.BaseJPush
    tag            string // 标签名
    registrationId string // 设备ID
}

func (tr *tagRegistration) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        tr.tag = tag
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签名不合法", nil))
    }
}

func (tr *tagRegistration) SetRegistrationId(registrationId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, registrationId)
    if match {
        tr.registrationId = registrationId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备ID不合法", nil))
    }
}

func (tr *tagRegistration) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tr.tag) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签名不能为空", nil))
    }
    if len(tr.registrationId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备ID不能为空", nil))
    }
    tr.ServiceUri = "/v3/tags/" + tr.tag + "/registration_ids/" + tr.registrationId

    tr.ReqUrl = tr.GetServiceUrl()

    return tr.GetRequest()
}

func NewTagRegistration(key string) *tagRegistration {
    tr := &tagRegistration{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), "", ""}
    return tr
}
