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

// 删除角色
type roleDelete struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleId   int // 角色ID
}

func (rd *roleDelete) SetRoleId(roleId int) {
    if roleId > 0 {
        rd.roleId = roleId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不合法", nil))
    }
}

func (rd *roleDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if rd.roleId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不能为空", nil))
    }
    rd.ExtendData["role_id"] = rd.roleId

    rd.ReqUrl = dingtalk.UrlService + "/topapi/role/deleterole?access_token=" + dingtalk.NewUtil().GetAccessToken(rd.corpId, rd.agentTag, rd.atType)

    reqBody := mpf.JsonMarshal(rd.ExtendData)
    client, req := rd.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewRoleDelete(corpId, agentTag, atType string) *roleDelete {
    rd := &roleDelete{dingtalk.NewCorp(), "", "", "", 0}
    rd.corpId = corpId
    rd.agentTag = agentTag
    rd.atType = atType
    rd.ReqContentType = project.HttpContentTypeJson
    rd.ReqMethod = fasthttp.MethodPost
    return rd
}
