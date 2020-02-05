/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 0:33
 */
package tag

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

// 增加标签成员
type tagUsersAdd struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    tagId     string   // 标签ID
    userList  []string // 成员ID列表
    partyList []int    // 部门ID列表
}

func (tua *tagUsersAdd) SetTagId(tagId string) {
    if len(tagId) > 0 {
        tua.tagId = tagId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不合法", nil))
    }
}

func (tua *tagUsersAdd) SetUserList(userList []string) {
    tua.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            tua.userList = append(tua.userList, v)
        }
    }
    if len(tua.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员ID列表不能为空", nil))
    } else if len(tua.userList) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员ID列表不能超过1000个", nil))
    }
}

func (tua *tagUsersAdd) SetPartyList(partyList []int) {
    tua.partyList = make([]int, 0)
    for _, v := range partyList {
        if v > 0 {
            tua.partyList = append(tua.partyList, v)
        }
    }
    if len(tua.partyList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门ID列表不能超过100个", nil))
    }
}

func (tua *tagUsersAdd) checkData() {
    if len(tua.tagId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不能为空", nil))
    }
    if (len(tua.userList) == 0) && (len(tua.partyList) == 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员列表和部门列表不能同时为空", nil))
    }
}

func (tua *tagUsersAdd) SendRequest(getType string) api.ApiResult {
    tua.checkData()

    reqData := make(map[string]interface{})
    reqData["tagid"] = tua.tagId
    reqData["userlist"] = tua.userList
    reqData["partylist"] = tua.partyList
    reqBody := mpf.JsonMarshal(reqData)
    tua.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/tag/addtagusers?access_token=" + wx.NewUtilWx().GetCorpCache(tua.corpId, tua.agentTag, getType)
    client, req := tua.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tua.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewTagUsersAdd(corpId, agentTag string) *tagUsersAdd {
    tua := &tagUsersAdd{wx.NewBaseWxCorp(), "", "", "", make([]string, 0), make([]int, 0)}
    tua.corpId = corpId
    tua.agentTag = agentTag
    tua.ReqContentType = project.HttpContentTypeJson
    tua.ReqMethod = fasthttp.MethodPost
    return tua
}
