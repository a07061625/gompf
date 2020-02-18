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

// 获取企业外部联系人详情
type extContactGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户id
}

func (ecg *extContactGet) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ecg.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不合法", nil))
    }
}

func (ecg *extContactGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ecg.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户id不能为空", nil))
    }
    ecg.ExtendData["user_id"] = ecg.userId

    ecg.ReqURI = dingtalk.UrlService + "/topapi/extcontact/get?access_token=" + dingtalk.NewUtil().GetAccessToken(ecg.corpId, ecg.agentTag, ecg.atType)

    reqBody := mpf.JSONMarshal(ecg.ExtendData)
    client, req := ecg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewExtContactGet(corpId, agentTag, atType string) *extContactGet {
    ecg := &extContactGet{dingtalk.NewCorp(), "", "", "", ""}
    ecg.corpId = corpId
    ecg.agentTag = agentTag
    ecg.atType = atType
    ecg.ReqContentType = project.HTTPContentTypeJSON
    ecg.ReqMethod = fasthttp.MethodPost
    return ecg
}
