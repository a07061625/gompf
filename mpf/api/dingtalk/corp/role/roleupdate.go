/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/31 0031
 * Time: 22:33
 */
package role

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/dingtalk"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新角色
type roleUpdate struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleId   int    // 角色ID
    roleName string // 角色名称
}

func (ru *roleUpdate) SetRoleId(roleId int) {
    if roleId > 0 {
        ru.roleId = roleId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不合法", nil))
    }
}

func (ru *roleUpdate) SetRoleName(roleName string) {
    if len(roleName) > 0 {
        ru.roleName = roleName
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色名称不合法", nil))
    }
}

func (ru *roleUpdate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if ru.roleId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不能为空", nil))
    }
    if len(ru.roleName) == 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色名称不能为空", nil))
    }
    ru.ExtendData["roleId"] = ru.roleId
    ru.ExtendData["roleName"] = ru.roleName

    ru.ReqUrl = dingtalk.UrlService + "/role/update_role?access_token=" + dingtalk.NewUtil().GetAccessToken(ru.corpId, ru.agentTag, ru.atType)

    reqBody := mpf.JSONMarshal(ru.ExtendData)
    client, req := ru.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRoleUpdate(corpId, agentTag, atType string) *roleUpdate {
    ru := &roleUpdate{dingtalk.NewCorp(), "", "", "", 0, ""}
    ru.corpId = corpId
    ru.agentTag = agentTag
    ru.atType = atType
    ru.ReqContentType = project.HTTPContentTypeJSON
    ru.ReqMethod = fasthttp.MethodPost
    return ru
}
