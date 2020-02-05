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

type chatGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    chatId   string // 会话ID
}

func (cg *chatGet) SetChatId(chatId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, chatId)
    if match {
        cg.chatId = chatId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "会话ID不合法", nil))
    }
}

func (cg *chatGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cg.chatId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "会话ID不能为空", nil))
    }
    cg.ReqData["chatid"] = cg.chatId
    cg.ReqData["access_token"] = dingtalk.NewUtil().GetCorpAccessToken(cg.corpId, cg.agentTag)

    cg.ReqUrl = dingtalk.UrlService + "/chat/get?" + mpf.HttpCreateParams(cg.ReqData, "none", 1)

    return cg.GetRequest()
}

func NewChatGet(corpId, agentTag string) *chatGet {
    cg := &chatGet{dingtalk.NewCorp(), "", "", ""}
    cg.corpId = corpId
    cg.agentTag = agentTag
    return cg
}
