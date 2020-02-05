package auth

import (
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/valyala/fasthttp"
)

// 获取通讯录权限范围
type scopes struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
}

func (s *scopes) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    s.ReqUrl = dingtalk.UrlService + "/auth/scopes?access_token=" + dingtalk.NewUtil().GetAccessToken(s.corpId, s.agentTag, s.atType)

    return s.GetRequest()
}

func NewScopes(corpId, agentTag, atType string) *scopes {
    s := &scopes{dingtalk.NewCorp(), "", "", ""}
    s.corpId = corpId
    s.agentTag = agentTag
    s.atType = atType
    return s
}
