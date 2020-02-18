package chat

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询群消息已读人员列表
type readListGet struct {
    dingtalk.BaseCorp
    corpId    string
    agentTag  string
    messageId string // 消息ID
}

func (rlg *readListGet) SetMessageId(messageId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, messageId)
    if match {
        rlg.messageId = messageId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息ID不合法", nil))
    }
}

func (rlg *readListGet) SetCursor(cursor int) {
    if cursor >= 0 {
        rlg.ReqData["cursor"] = strconv.Itoa(cursor)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页游标不合法", nil))
    }
}

func (rlg *readListGet) SetSize(size int) {
    if (size > 0) && (size <= 100) {
        rlg.ReqData["size"] = strconv.Itoa(size)
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (rlg *readListGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(rlg.messageId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "消息ID不能为空", nil))
    }
    rlg.ReqData["messageId"] = rlg.messageId
    rlg.ReqData["access_token"] = dingtalk.NewUtil().GetCorpAccessToken(rlg.corpId, rlg.agentTag)

    rlg.ReqUrl = dingtalk.UrlService + "/chat/getReadList?" + mpf.HTTPCreateParams(rlg.ReqData, "none", 1)

    return rlg.GetRequest()
}

func NewReadListGet(corpId, agentTag string) *readListGet {
    rlg := &readListGet{dingtalk.NewCorp(), "", "", ""}
    rlg.corpId = corpId
    rlg.agentTag = agentTag
    rlg.ReqData["cursor"] = "0"
    rlg.ReqData["size"] = "10"
    return rlg
}
