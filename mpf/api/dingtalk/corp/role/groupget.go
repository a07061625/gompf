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

// 获取角色组
type groupGet struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    groupId  int // 角色组ID
}

func (gg *groupGet) SetRoleId(groupId int) {
    if groupId > 0 {
        gg.groupId = groupId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色组ID不合法", nil))
    }
}

func (gg *groupGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if gg.groupId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色组ID不能为空", nil))
    }
    gg.ExtendData["group_id"] = gg.groupId

    gg.ReqURI = dingtalk.UrlService + "/topapi/role/getrolegroup?access_token=" + dingtalk.NewUtil().GetAccessToken(gg.corpId, gg.agentTag, gg.atType)

    reqBody := mpf.JSONMarshal(gg.ExtendData)
    client, req := gg.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewGroupGet(corpId, agentTag, atType string) *groupGet {
    gg := &groupGet{dingtalk.NewCorp(), "", "", "", 0}
    gg.corpId = corpId
    gg.agentTag = agentTag
    gg.atType = atType
    gg.ReqContentType = project.HTTPContentTypeJSON
    gg.ReqMethod = fasthttp.MethodPost
    return gg
}
