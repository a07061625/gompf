package health

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户钉钉运动开启状态
type stepInfoUserStatus struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
}

func (ius *stepInfoUserStatus) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ius.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (ius *stepInfoUserStatus) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ius.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    ius.ExtendData["userid"] = ius.userId

    ius.ReqUrl = dingtalk.UrlService + "/topapi/health/stepinfo/getuserstatus?access_token=" + dingtalk.NewUtil().GetAccessToken(ius.corpId, ius.agentTag, ius.atType)

    reqBody := mpf.JSONMarshal(ius.ExtendData)
    client, req := ius.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewStepInfoUserStatus(corpId, agentTag, atType string) *stepInfoUserStatus {
    ius := &stepInfoUserStatus{dingtalk.NewCorp(), "", "", "", ""}
    ius.corpId = corpId
    ius.agentTag = agentTag
    ius.atType = atType
    ius.ReqContentType = project.HTTPContentTypeJSON
    ius.ReqMethod = fasthttp.MethodPost
    return ius
}
