/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 15:23
 */
package push

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 消息推送校验
type validatePush struct {
    mppush.BaseJPush
    notification   map[string]interface{} // 通知内容
    messageContent map[string]interface{} // 消息内容
}

func (vp *validatePush) SetNotification(notification map[string]interface{}) {
    if len(notification) > 0 {
        vp.notification = notification
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "通知内容不能为空", nil))
    }
}

func (vp *validatePush) SetMessageContent(messageContent map[string]interface{}) {
    if len(messageContent) > 0 {
        vp.messageContent = messageContent
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息内容不能为空", nil))
    }
}

func (vp *validatePush) SetCid(cid string) {
    if len(cid) > 0 {
        vp.ExtendData["cid"] = cid
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送标识符不合法", nil))
    }
}

func (vp *validatePush) SetPlatformList(platformList []string) {
    if len(platformList) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "平台类型不合法", nil))
    }

    platforms := make([]string, 0)
    for _, v := range platformList {
        _, ok := mppush.JPushPlatformTypes[v]
        if ok {
            platforms = append(platforms, v)
        }
    }
    if len(platforms) > 0 {
        vp.ExtendData["platform"] = platforms
    }
}

func (vp *validatePush) SetAudience(audience map[string]interface{}) {
    if len(audience) > 0 {
        vp.ExtendData["audience"] = audience
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "接收方不能为空", nil))
    }
}

func (vp *validatePush) SetSmsMessage(smsMessage map[string]interface{}) {
    if len(smsMessage) > 0 {
        vp.ExtendData["sms_message"] = smsMessage
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "短信补充内容不能为空", nil))
    }
}

func (vp *validatePush) SetOptions(options map[string]interface{}) {
    if len(options) > 0 {
        vp.ExtendData["options"] = options
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送参数不能为空", nil))
    }
}

func (vp *validatePush) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if (len(vp.notification) == 0) && (len(vp.messageContent) == 0) {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "通知内容和消息内容不能同时为空", nil))
    }
    if len(vp.notification) > 0 {
        vp.ExtendData["notification"] = vp.notification
    }
    if len(vp.messageContent) > 0 {
        vp.ExtendData["message"] = vp.messageContent
    }

    vp.ReqURI = vp.GetServiceUrl()

    reqBody := mpf.JSONMarshal(vp.ExtendData)
    client, req := vp.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewValidatePush(key string) *validatePush {
    vp := &validatePush{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), make(map[string]interface{}), make(map[string]interface{})}
    vp.ServiceUri = "/v3/push/validate"
    vp.ExtendData["platform"] = "all"
    vp.ExtendData["audience"] = "all"
    vp.ReqContentType = project.HTTPContentTypeJSON
    vp.ReqMethod = fasthttp.MethodPost
    return vp
}
