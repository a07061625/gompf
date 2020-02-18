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

// 删除标签成员
type tagUsersDel struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    tagId     string   // 标签ID
    userList  []string // 成员ID列表
    partyList []int    // 部门ID列表
}

func (tud *tagUsersDel) SetTagId(tagId string) {
    if len(tagId) > 0 {
        tud.tagId = tagId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不合法", nil))
    }
}

func (tud *tagUsersDel) SetUserList(userList []string) {
    tud.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            tud.userList = append(tud.userList, v)
        }
    }
    if len(tud.userList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员ID列表不能为空", nil))
    } else if len(tud.userList) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员ID列表不能超过1000个", nil))
    }
}

func (tud *tagUsersDel) SetPartyList(partyList []int) {
    tud.partyList = make([]int, 0)
    for _, v := range partyList {
        if v > 0 {
            tud.partyList = append(tud.partyList, v)
        }
    }
    if len(tud.partyList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门ID列表不能超过100个", nil))
    }
}

func (tud *tagUsersDel) checkData() {
    if len(tud.tagId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签id不能为空", nil))
    }
    if (len(tud.userList) == 0) && (len(tud.partyList) == 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员列表和部门列表不能同时为空", nil))
    }
}

func (tud *tagUsersDel) SendRequest(getType string) api.APIResult {
    tud.checkData()

    reqData := make(map[string]interface{})
    reqData["tagid"] = tud.tagId
    reqData["userlist"] = tud.userList
    reqData["partylist"] = tud.partyList
    reqBody := mpf.JSONMarshal(reqData)
    tud.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/tag/deltagusers?access_token=" + wx.NewUtilWx().GetCorpCache(tud.corpId, tud.agentTag, getType)
    client, req := tud.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tud.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTagUsersDel(corpId, agentTag string) *tagUsersDel {
    tud := &tagUsersDel{wx.NewBaseWxCorp(), "", "", "", make([]string, 0), make([]int, 0)}
    tud.corpId = corpId
    tud.agentTag = agentTag
    tud.ReqContentType = project.HTTPContentTypeJSON
    tud.ReqMethod = fasthttp.MethodPost
    return tud
}
