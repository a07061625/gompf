package role

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建角色
type roleAdd struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleName string // 角色名称
    groupId  int    // 角色组ID
}

func (ra *roleAdd) SetRoleName(roleName string) {
    if len(roleName) > 0 {
        ra.roleName = roleName
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色名称不合法", nil))
    }
}

func (ra *roleAdd) SetGroupId(groupId int) {
    if groupId > 0 {
        ra.groupId = groupId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色组ID不合法", nil))
    }
}

func (ra *roleAdd) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ra.roleName) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色名称不能为空", nil))
    }
    if ra.groupId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色组ID不能为空", nil))
    }
    ra.ExtendData["roleName"] = ra.roleName
    ra.ExtendData["groupId"] = ra.groupId

    ra.ReqUrl = dingtalk.UrlService + "/role/add_role?access_token=" + dingtalk.NewUtil().GetAccessToken(ra.corpId, ra.agentTag, ra.atType)

    reqBody := mpf.JsonMarshal(ra.ExtendData)
    client, req := ra.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRoleAdd(corpId, agentTag, atType string) *roleAdd {
    ra := &roleAdd{dingtalk.NewCorp(), "", "", "", "", 0}
    ra.corpId = corpId
    ra.agentTag = agentTag
    ra.atType = atType
    ra.ReqContentType = project.HttpContentTypeJson
    ra.ReqMethod = fasthttp.MethodPost
    return ra
}
