package extcontact

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除外部联系人
type extContactDelete struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户id
}

func (ecd *extContactDelete) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ecd.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (ecd *extContactDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ecd.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    ecd.ExtendData["user_id"] = ecd.userId

    ecd.ReqUrl = dingtalk.UrlService + "/topapi/extcontact/delete?access_token=" + dingtalk.NewUtil().GetAccessToken(ecd.corpId, ecd.agentTag, ecd.atType)

    reqBody := mpf.JSONMarshal(ecd.ExtendData)
    client, req := ecd.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewExtContactDelete(corpId, agentTag, atType string) *extContactDelete {
    ecd := &extContactDelete{dingtalk.NewCorp(), "", "", "", ""}
    ecd.corpId = corpId
    ecd.agentTag = agentTag
    ecd.atType = atType
    ecd.ReqContentType = project.HTTPContentTypeJSON
    ecd.ReqMethod = fasthttp.MethodPost
    return ecd
}
