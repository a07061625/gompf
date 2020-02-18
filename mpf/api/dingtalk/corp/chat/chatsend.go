package chat

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送群消息
type chatSend struct {
    dingtalk.BaseCorp
    corpId     string
    agentTag   string
    chatId     string                 // 会话ID
    msgContent map[string]interface{} // 消息内容
}

func (cs *chatSend) SetChatId(chatId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, chatId)
    if match {
        cs.chatId = chatId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "会话ID不合法", nil))
    }
}

func (cs *chatSend) SetMsgContent(msgType string, msgData map[string]interface{}) {
    _, ok := dingtalk.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息类型不支持", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息数据不能为空", nil))
    }
    cs.msgContent["msgtype"] = msgType
    cs.msgContent[msgType] = msgData
}

func (cs *chatSend) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cs.chatId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "会话ID不能为空", nil))
    }
    if len(cs.msgContent) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息内容不能为空", nil))
    }
    cs.ExtendData["chatid"] = cs.chatId
    cs.ExtendData["msg"] = cs.msgContent

    cs.ReqUrl = dingtalk.UrlService + "/chat/send?access_token=" + dingtalk.NewUtil().GetCorpAccessToken(cs.corpId, cs.agentTag)

    reqBody := mpf.JsonMarshal(cs.ExtendData)
    client, req := cs.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewChatSend(corpId, agentTag string) *chatSend {
    cs := &chatSend{dingtalk.NewCorp(), "", "", "", make(map[string]interface{})}
    cs.corpId = corpId
    cs.agentTag = agentTag
    cs.ReqContentType = project.HTTPContentTypeJSON
    cs.ReqMethod = fasthttp.MethodPost
    return cs
}
