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

// 获取角色详情
type roleGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleId   int // 角色ID
}

func (rg *roleGet) SetRoleId(roleId int) {
    if roleId > 0 {
        rg.roleId = roleId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不合法", nil))
    }
}

func (rg *roleGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if rg.roleId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不能为空", nil))
    }
    rg.ExtendData["roleId"] = rg.roleId

    rg.ReqUrl = dingtalk.UrlService + "/topapi/role/getrole?access_token=" + dingtalk.NewUtil().GetAccessToken(rg.corpId, rg.agentTag, rg.atType)

    reqBody := mpf.JsonMarshal(rg.ExtendData)
    client, req := rg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRoleGet(corpId, agentTag, atType string) *roleGet {
    rg := &roleGet{dingtalk.NewCorp(), "", "", "", 0}
    rg.corpId = corpId
    rg.agentTag = agentTag
    rg.atType = atType
    rg.ReqContentType = project.HttpContentTypeJson
    rg.ReqMethod = fasthttp.MethodPost
    return rg
}
