/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 23:29
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 分组消息发送接口
type groupSend struct {
    alipay.BaseAliPay
    groupId  string                 // 分组ID
    msgType  string                 // 消息类型,text:文本消息 image-text:图文消息
    articles map[string]interface{} // 图文消息内容
    text     map[string]interface{} // 文本消息内容
    image    map[string]interface{} // 图片消息内容
}

func (gs *groupSend) SetGroupId(groupId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,10}$`, groupId)
    if match {
        gs.groupId = groupId
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组ID不合法", nil))
    }
}

func (gs *groupSend) SetMsgType(msgType string) {
    if (msgType == "text") || (msgType == "image-text") || (msgType == "image") {
        gs.msgType = msgType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息类型不合法", nil))
    }
}

func (gs *groupSend) SetArticles(articles map[string]interface{}) {
    if len(articles) > 0 {
        gs.articles = articles
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图文消息内容不合法", nil))
    }
}

func (gs *groupSend) SetText(text map[string]interface{}) {
    if len(text) > 0 {
        gs.text = text
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "文本消息内容不合法", nil))
    }
}

func (gs *groupSend) SetImage(image map[string]interface{}) {
    if len(image) > 0 {
        gs.image = image
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图片消息内容不合法", nil))
    }
}

func (gs *groupSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(gs.groupId) == 0 {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "分组ID不能为空", nil))
    }
    if gs.msgType == "image-text" {
        if len(gs.articles) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图文消息内容不能为空", nil))
        }
        gs.BizContent["articles"] = gs.articles
    } else if gs.msgType == "text" {
        if len(gs.text) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "文本消息内容不能为空", nil))
        }
        gs.BizContent["text"] = gs.text
    } else if gs.msgType == "image" {
        if len(gs.image) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图片消息内容不能为空", nil))
        }
        gs.BizContent["image"] = gs.image
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息类型不能为空", nil))
    }
    gs.BizContent["group_id"] = gs.groupId
    gs.BizContent["msg_type"] = gs.msgType

    return gs.GetRequest()
}

func NewGroupSend(appId string) *groupSend {
    gs := &groupSend{alipay.NewBase(appId), "", "", make(map[string]interface{}), make(map[string]interface{}), make(map[string]interface{})}
    gs.SetMethod("alipay.open.public.message.group.send")
    return gs
}
