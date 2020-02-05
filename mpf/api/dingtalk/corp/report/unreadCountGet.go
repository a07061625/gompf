package report

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户日志未读数
type unreadCountGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
}

func (ucg *unreadCountGet) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ucg.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (ucg *unreadCountGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ucg.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    ucg.ExtendData["userid"] = ucg.userId

    ucg.ReqUrl = dingtalk.UrlService + "/topapi/report/getunreadcount?access_token=" + dingtalk.NewUtil().GetAccessToken(ucg.corpId, ucg.agentTag, ucg.atType)

    reqBody := mpf.JsonMarshal(ucg.ExtendData)
    client, req := ucg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewUnreadCountGet(corpId, agentTag, atType string) *unreadCountGet {
    ucg := &unreadCountGet{dingtalk.NewCorp(), "", "", "", ""}
    ucg.corpId = corpId
    ucg.agentTag = agentTag
    ucg.atType = atType
    ucg.ReqContentType = project.HttpContentTypeJson
    ucg.ReqMethod = fasthttp.MethodPost
    return ucg
}
