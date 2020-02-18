package microapp

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 设置应用的可见范围
type visibleScopesSet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    agentId  int // 应用ID
}

func (vss *visibleScopesSet) SetAgentId(agentId int) {
    if agentId > 0 {
        vss.agentId = agentId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "应用ID不合法", nil))
    }
}

func (vss *visibleScopesSet) SetIsHidden(isHidden bool) {
    vss.ExtendData["isHidden"] = isHidden
}

func (vss *visibleScopesSet) SetDeptVisibleScopes(deptVisibleScopes []int) {
    departList := make([]int, 0)
    for _, v := range deptVisibleScopes {
        if v > 0 {
            departList = append(departList, v)
        }
    }
    vss.ExtendData["deptVisibleScopes"] = departList
}

func (vss *visibleScopesSet) SetUserVisibleScopes(userVisibleScopes []string) {
    userList := make([]string, 0)
    for _, v := range userVisibleScopes {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, v)
        if match {
            userList = append(userList, v)
        }
    }
    vss.ExtendData["userVisibleScopes"] = userList
}

func (vss *visibleScopesSet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if vss.agentId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "应用ID不能为空", nil))
    }
    vss.ExtendData["agentId"] = vss.agentId

    vss.ReqUrl = dingtalk.UrlService + "/microapp/set_visible_scopes?access_token=" + dingtalk.NewUtil().GetAccessToken(vss.corpId, vss.agentTag, vss.atType)

    reqBody := mpf.JSONMarshal(vss.ExtendData)
    client, req := vss.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewVisibleScopesSet(corpId, agentTag, atType string) *visibleScopesSet {
    vss := &visibleScopesSet{dingtalk.NewCorp(), "", "", "", 0}
    vss.corpId = corpId
    vss.agentTag = agentTag
    vss.atType = atType
    vss.ExtendData["isHidden"] = false
    vss.ExtendData["deptVisibleScopes"] = make([]int, 0)
    vss.ExtendData["userVisibleScopes"] = make([]string, 0)
    vss.ReqContentType = project.HTTPContentTypeJSON
    vss.ReqMethod = fasthttp.MethodPost
    return vss
}
