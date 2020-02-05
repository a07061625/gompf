/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 23:46
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 异步单发消息
type customSend struct {
    alipay.BaseAliPay
    userId   string                 // 用户ID
    msgType  string                 // 消息类型,text:文本消息 image-text:图文消息
    articles map[string]interface{} // 图文消息内容
    text     map[string]interface{} // 文本消息内容
}

func (cs *customSend) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, userId)
    if match {
        cs.userId = userId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "用户ID不合法", nil))
    }
}

func (cs *customSend) SetMsgType(msgType string) {
    if (msgType == "text") || (msgType == "image-text") {
        cs.msgType = msgType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息类型不合法", nil))
    }
}

func (cs *customSend) SetArticles(articles map[string]interface{}) {
    if len(articles) > 0 {
        cs.articles = articles
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图文消息内容不合法", nil))
    }
}

func (cs *customSend) SetText(text map[string]interface{}) {
    if len(text) > 0 {
        cs.text = text
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "文本消息内容不合法", nil))
    }
}

func (cs *customSend) SetChatStatus(chatStatus string) {
    if (chatStatus == "0") || (chatStatus == "1") {
        cs.BizContent["chat"] = chatStatus
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "聊天消息状态不合法", nil))
    }
}

func (cs *customSend) SetEventType(eventType string) {
    if (eventType == "follow") || (eventType == "click") || (eventType == "enter_ppchat") {
        cs.BizContent["event_type"] = eventType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "事件类型不合法", nil))
    }
}

func (cs *customSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cs.userId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "用户ID不能为空", nil))
    }
    if cs.msgType == "image-text" {
        if len(cs.articles) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图文消息内容不能为空", nil))
        }
        cs.BizContent["articles"] = cs.articles
    } else if cs.msgType == "text" {
        if len(cs.text) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "文本消息内容不能为空", nil))
        }
        cs.BizContent["text"] = cs.text
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息类型不能为空", nil))
    }
    cs.BizContent["to_user_id"] = cs.userId
    cs.BizContent["msg_type"] = cs.msgType

    return cs.GetRequest()
}

func NewCustomSend(appId string) *customSend {
    cs := &customSend{alipay.NewBase(appId), "", "", make(map[string]interface{}), make(map[string]interface{})}
    cs.BizContent["chat"] = "0"
    cs.SetMethod("alipay.open.public.message.custom.send")
    return cs
}
