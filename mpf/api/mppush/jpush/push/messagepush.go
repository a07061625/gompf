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

// 消息推送
type messagePush struct {
    mppush.BaseJPush
    notification   map[string]interface{} // 通知内容
    messageContent map[string]interface{} // 消息内容
}

func (mp *messagePush) SetNotification(notification map[string]interface{}) {
    if len(notification) > 0 {
        mp.notification = notification
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "通知内容不能为空", nil))
    }
}

func (mp *messagePush) SetMessageContent(messageContent map[string]interface{}) {
    if len(messageContent) > 0 {
        mp.messageContent = messageContent
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "消息内容不能为空", nil))
    }
}

func (mp *messagePush) SetCid(cid string) {
    if len(cid) > 0 {
        mp.ExtendData["cid"] = cid
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送标识符不合法", nil))
    }
}

func (mp *messagePush) SetPlatformList(platformList []string) {
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
        mp.ExtendData["platform"] = platforms
    }
}

func (mp *messagePush) SetAudience(audience map[string]interface{}) {
    if len(audience) > 0 {
        mp.ExtendData["audience"] = audience
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "接收方不能为空", nil))
    }
}

func (mp *messagePush) SetSmsMessage(smsMessage map[string]interface{}) {
    if len(smsMessage) > 0 {
        mp.ExtendData["sms_message"] = smsMessage
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "短信补充内容不能为空", nil))
    }
}

func (mp *messagePush) SetOptions(options map[string]interface{}) {
    if len(options) > 0 {
        mp.ExtendData["options"] = options
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送参数不能为空", nil))
    }
}

func (mp *messagePush) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if (len(mp.notification) == 0) && (len(mp.messageContent) == 0) {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "通知内容和消息内容不能同时为空", nil))
    }
    if len(mp.notification) > 0 {
        mp.ExtendData["notification"] = mp.notification
    }
    if len(mp.messageContent) > 0 {
        mp.ExtendData["message"] = mp.messageContent
    }

    mp.ReqUrl = mp.GetServiceUrl()

    reqBody := mpf.JSONMarshal(mp.ExtendData)
    client, req := mp.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewMessagePush(key string) *messagePush {
    mp := &messagePush{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), make(map[string]interface{}), make(map[string]interface{})}
    mp.ServiceUri = "/v3/push"
    mp.ExtendData["platform"] = "all"
    mp.ExtendData["audience"] = "all"
    mp.ReqContentType = project.HTTPContentTypeJSON
    mp.ReqMethod = fasthttp.MethodPost
    return mp
}
