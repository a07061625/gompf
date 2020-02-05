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

// 更新标签
type tagUpdate struct {
    mppush.BaseJPush
    tag              string                 // 标签名
    registrationList map[string]interface{} // 设备列表
}

func (tu *tagUpdate) SetTag(tag string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, tag)
    if match {
        tu.tag = tag
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签名不合法", nil))
    }
}

func (tu *tagUpdate) SetRegistrationList(registrationList map[string]interface{}) {
    if len(registrationList) > 0 {
        tu.registrationList = registrationList
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能为空", nil))
    }
}

func (tu *tagUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tu.tag) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "标签名不能为空", nil))
    }
    if len(tu.registrationList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "设备列表不能为空", nil))
    }
    tu.ServiceUri = "/v3/tags/" + tu.tag
    tu.ExtendData["registration_ids"] = tu.registrationList

    tu.ReqUrl = tu.GetServiceUrl()

    reqBody := mpf.JsonMarshal(tu.ExtendData)
    client, req := tu.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewTagUpdate(key string) *tagUpdate {
    tu := &tagUpdate{mppush.NewBaseJPush(mppush.JPushServiceDomainDevice, key, "app"), "", make(map[string]interface{})}
    tu.ReqContentType = project.HttpContentTypeJson
    tu.ReqMethod = fasthttp.MethodPost
    return tu
}
