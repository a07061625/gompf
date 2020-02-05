/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 23:29
 */
package message

import (
    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 群发消息
type totalSend struct {
    alipay.BaseAliPay
    msgType  string                 // 消息类型,text:文本消息 image-text:图文消息
    articles map[string]interface{} // 图文消息内容
    text     map[string]interface{} // 文本消息内容
}

func (ts *totalSend) SetMsgType(msgType string) {
    if (msgType == "text") || (msgType == "image-text") {
        ts.msgType = msgType
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息类型不合法", nil))
    }
}

func (ts *totalSend) SetArticles(articles map[string]interface{}) {
    if len(articles) > 0 {
        ts.articles = articles
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图文消息内容不合法", nil))
    }
}

func (ts *totalSend) SetText(text map[string]interface{}) {
    if len(text) > 0 {
        ts.text = text
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "文本消息内容不合法", nil))
    }
}

func (ts *totalSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if ts.msgType == "image-text" {
        if len(ts.articles) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "图文消息内容不能为空", nil))
        }
        ts.BizContent["articles"] = ts.articles
    } else if ts.msgType == "text" {
        if len(ts.text) == 0 {
            panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "文本消息内容不能为空", nil))
        }
        ts.BizContent["text"] = ts.text
    } else {
        panic(mperr.NewAliPayLife(errorcode.AliPayLifeParam, "消息类型不能为空", nil))
    }
    ts.BizContent["msg_type"] = ts.msgType

    return ts.GetRequest()
}

func NewTotalSend(appId string) *totalSend {
    ts := &totalSend{alipay.NewBase(appId), "", make(map[string]interface{}), make(map[string]interface{})}
    ts.SetMethod("alipay.open.public.message.total.send")
    return ts
}
