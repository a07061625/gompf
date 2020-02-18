/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 22:46
 */
package batch

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 邀请成员
type inviteUser struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    users     []string // 用户ID列表
    partyList []uint   // 部门ID列表
    tags      []uint   // 标签ID列表
}

func (iu *inviteUser) SetUsers(users []string) {
    iu.users = make([]string, 0)
    for _, v := range users {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            iu.users = append(iu.users, v)
        }
    }
    if len(iu.users) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID列表不能为空", nil))
    } else if len(iu.users) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID列表不能超过1000个", nil))
    }
}

func (iu *inviteUser) SetPartyList(partyList []uint) {
    iu.partyList = make([]uint, 0)
    for _, v := range partyList {
        if v > 0 {
            iu.partyList = append(iu.partyList, v)
        }
    }
    if len(iu.partyList) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门ID列表不能超过1000个", nil))
    }
}

func (iu *inviteUser) SetTags(tags []uint) {
    iu.tags = make([]uint, 0)
    for _, v := range tags {
        if v > 0 {
            iu.tags = append(iu.tags, v)
        }
    }
    if len(iu.tags) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签ID列表不能超过1000个", nil))
    }
}

func (iu *inviteUser) checkData() {
    if (len(iu.users) == 0) && (len(iu.partyList) == 0) && (len(iu.tags) == 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户列表,部门列表和标签列表不能同时为空", nil))
    }
}

func (iu *inviteUser) SendRequest() api.ApiResult {
    iu.checkData()
    reqData := make(map[string]interface{})
    reqData["user"] = iu.users
    reqData["party"] = iu.partyList
    reqData["tag"] = iu.tags
    reqBody := mpf.JsonMarshal(reqData)

    iu.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/batch/invite?access_token=" + wx.NewUtilWx().GetCorpAccessToken(iu.corpId, iu.agentTag)
    client, req := iu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := iu.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewInviteUser(corpId, agentTag string) *inviteUser {
    iu := &inviteUser{wx.NewBaseWxCorp(), "", "", make([]string, 0), make([]uint, 0), make([]uint, 0)}
    iu.corpId = corpId
    iu.agentTag = agentTag
    iu.ReqContentType = project.HTTPContentTypeJSON
    iu.ReqMethod = fasthttp.MethodPost
    return iu
}
