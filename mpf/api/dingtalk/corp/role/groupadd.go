package role

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建角色组
type groupAdd struct {
    dingtalk.BaseCorp
    corpId    string
    agentTag  string
    atType    string
    groupName string // 角色组名称
}

func (ga *groupAdd) SetGroupName(groupName string) {
    if len(groupName) > 0 {
        ga.groupName = groupName
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色组名称不合法", nil))
    }
}

func (ga *groupAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ga.groupName) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色组名称不能为空", nil))
    }
    ga.ExtendData["name"] = ga.groupName

    ga.ReqUrl = dingtalk.UrlService + "/role/add_role_group?access_token=" + dingtalk.NewUtil().GetAccessToken(ga.corpId, ga.agentTag, ga.atType)

    reqBody := mpf.JsonMarshal(ga.ExtendData)
    client, req := ga.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewGroupAdd(corpId, agentTag, atType string) *groupAdd {
    ga := &groupAdd{dingtalk.NewCorp(), "", "", "", ""}
    ga.corpId = corpId
    ga.agentTag = agentTag
    ga.atType = atType
    ga.ReqContentType = project.HTTPContentTypeJSON
    ga.ReqMethod = fasthttp.MethodPost
    return ga
}
