package microapp

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取应用的可见范围
type visibleScopesGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    agentId  int // 应用ID
}

func (vsg *visibleScopesGet) SetAgentId(agentId int) {
    if agentId > 0 {
        vsg.agentId = agentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "应用ID不合法", nil))
    }
}

func (vsg *visibleScopesGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if vsg.agentId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "应用ID不能为空", nil))
    }
    vsg.ExtendData["agentId"] = vsg.agentId

    vsg.ReqUrl = dingtalk.UrlService + "/microapp/visible_scopes?access_token=" + dingtalk.NewUtil().GetAccessToken(vsg.corpId, vsg.agentTag, vsg.atType)

    reqBody := mpf.JSONMarshal(vsg.ExtendData)
    client, req := vsg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewVisibleScopesGet(corpId, agentTag, atType string) *visibleScopesGet {
    vsg := &visibleScopesGet{dingtalk.NewCorp(), "", "", "", 0}
    vsg.corpId = corpId
    vsg.agentTag = agentTag
    vsg.atType = atType
    vsg.ReqContentType = project.HTTPContentTypeJSON
    vsg.ReqMethod = fasthttp.MethodPost
    return vsg
}
