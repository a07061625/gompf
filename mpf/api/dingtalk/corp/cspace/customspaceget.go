package cspace

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取企业下的自定义空间
type customSpaceGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    domain   string // 域名
}

func (csg *customSpaceGet) SetDomain(domain string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,10}$`, domain)
    if match {
        csg.domain = domain
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "域名不合法", nil))
    }
}

func (csg *customSpaceGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(csg.domain) > 0 {
        csg.ReqData["domain"] = csg.domain
        csg.ReqData["access_token"] = dingtalk.NewUtil().GetCorpAccessToken(csg.corpId, csg.agentTag)
    } else {
        agentInfo := dingtalk.NewConfig().GetCorp(csg.corpId).GetAgentInfo(csg.agentTag)
        csg.ReqData["agent_id"] = agentInfo["id"]
        csg.ReqData["access_token"] = dingtalk.NewUtil().GetProviderAuthorizeAccessToken(csg.corpId)
    }
    csg.ReqURI = dingtalk.UrlService + "/cspace/get_custom_space?" + mpf.HTTPCreateParams(csg.ReqData, "none", 1)

    return csg.GetRequest()
}

func NewCustomSpaceGet(corpId, agentTag string) *customSpaceGet {
    csg := &customSpaceGet{dingtalk.NewCorp(), "", "", ""}
    csg.corpId = corpId
    csg.agentTag = agentTag
    return csg
}
