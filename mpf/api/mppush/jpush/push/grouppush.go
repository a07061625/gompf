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

// 消息分组推送
type groupPush struct {
    mppush.BaseJPush
    notification   map[string]interface{} // 通知内容
    messageContent map[string]interface{} // 消息内容
}

func (gp *groupPush) SetNotification(notification map[string]interface{}) {
    if len(notification) > 0 {
        gp.notification = notification
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "通知内容不能为空", nil))
    }
}

func (gp *groupPush) SetMessageContent(messageContent map[string]interface{}) {
    if len(messageContent) > 0 {
        gp.messageContent = messageContent
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息内容不能为空", nil))
    }
}

func (gp *groupPush) SetCid(cid string) {
    if len(cid) > 0 {
        gp.ExtendData["cid"] = cid
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送标识符不合法", nil))
    }
}

func (gp *groupPush) SetPlatformList(platformList []string) {
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
        gp.ExtendData["platform"] = platforms
    }
}

func (gp *groupPush) SetAudience(audience map[string]interface{}) {
    if len(audience) > 0 {
        gp.ExtendData["audience"] = audience
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "接收方不能为空", nil))
    }
}

func (gp *groupPush) SetSmsMessage(smsMessage map[string]interface{}) {
    if len(smsMessage) > 0 {
        gp.ExtendData["sms_message"] = smsMessage
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "短信补充内容不能为空", nil))
    }
}

func (gp *groupPush) SetOptions(options map[string]interface{}) {
    if len(options) > 0 {
        gp.ExtendData["options"] = options
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送参数不能为空", nil))
    }
}

func (gp *groupPush) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if (len(gp.notification) == 0) && (len(gp.messageContent) == 0) {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "通知内容和消息内容不能同时为空", nil))
    }
    if len(gp.notification) > 0 {
        gp.ExtendData["notification"] = gp.notification
    }
    if len(gp.messageContent) > 0 {
        gp.ExtendData["message"] = gp.messageContent
    }

    gp.ReqUrl = gp.GetServiceUrl()

    reqBody := mpf.JSONMarshal(gp.ExtendData)
    client, req := gp.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewGroupPush(key string) *groupPush {
    gp := &groupPush{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "group"), make(map[string]interface{}), make(map[string]interface{})}
    gp.ServiceUri = "/v3/grouppush"
    gp.ExtendData["platform"] = "all"
    gp.ExtendData["audience"] = "all"
    gp.ReqContentType = project.HTTPContentTypeJSON
    gp.ReqMethod = fasthttp.MethodPost
    return gp
}
