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

// 获取角色下的员工列表
type simpleList struct {
    dingtalk.BaseCorp
    corpId   string
    agentTag string
    atType   string
    roleId   int // 角色ID
}

func (sl *simpleList) SetRoleId(roleId int) {
    if roleId > 0 {
        sl.roleId = roleId
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不合法", nil))
    }
}

func (sl *simpleList) SetOffset(offset int) {
    if offset >= 0 {
        sl.ExtendData["offset"] = offset
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "偏移量不合法", nil))
    }
}

func (sl *simpleList) SetSize(size int) {
    if (size > 0) && (size <= 200) {
        sl.ExtendData["size"] = size
    } else {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "分页大小不合法", nil))
    }
}

func (sl *simpleList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if sl.roleId <= 0 {
        panic(mperr.NewDingTalkCorp(errorcode.DingTalkCorpParam, "角色ID不能为空", nil))
    }
    sl.ExtendData["role_id"] = sl.roleId

    sl.ReqUrl = dingtalk.UrlService + "/topapi/role/simplelist?access_token=" + dingtalk.NewUtil().GetAccessToken(sl.corpId, sl.agentTag, sl.atType)

    reqBody := mpf.JsonMarshal(sl.ExtendData)
    client, req := sl.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewSimpleList(corpId, agentTag, atType string) *simpleList {
    sl := &simpleList{dingtalk.NewCorp(), "", "", "", 0}
    sl.corpId = corpId
    sl.agentTag = agentTag
    sl.atType = atType
    sl.ExtendData["offset"] = 0
    sl.ExtendData["size"] = 10
    sl.ReqContentType = project.HttpContentTypeJson
    sl.ReqMethod = fasthttp.MethodPost
    return sl
}
