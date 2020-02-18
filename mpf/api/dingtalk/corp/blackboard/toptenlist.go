package blackboard

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取用户公告数据
type topTenList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    userId   string // 用户ID
}

func (ttl *topTenList) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ttl.userId = userId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不合法", nil))
    }
}

func (ttl *topTenList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ttl.userId) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "用户ID不能为空", nil))
    }
    ttl.ExtendData["userid"] = ttl.userId

    ttl.ReqUrl = dingtalk.UrlService + "/topapi/blackboard/listtopten?access_token=" + dingtalk.NewUtil().GetAccessToken(ttl.corpId, ttl.agentTag, ttl.atType)

    reqBody := mpf.JsonMarshal(ttl.ExtendData)
    client, req := ttl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewTopTenList(corpId, agentTag, atType string) *topTenList {
    ttl := &topTenList{dingtalk.NewCorp(), "", "", "", ""}
    ttl.corpId = corpId
    ttl.agentTag = agentTag
    ttl.atType = atType
    ttl.ReqContentType = project.HTTPContentTypeJSON
    ttl.ReqMethod = fasthttp.MethodPost
    return ttl
}
