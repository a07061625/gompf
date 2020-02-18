/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 16:51
 */
package message

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

// 修改群聊会话
type appChatUpdate struct {
    wx.BaseWxCorp
    corpId      string
    agentTag    string
    chatId      string   // 群id
    name        string   // 群聊名
    owner       string   // 群主id
    addUserList []string // 添加成员列表
    delUserList []string // 踢出成员列表
}

func (acu *appChatUpdate) SetChatId(chatId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, chatId)
    if match {
        acu.chatId = chatId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群id不合法", nil))
    }
}

func (acu *appChatUpdate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        acu.name = string(trueName[:25])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群聊名不合法", nil))
    }
}

func (acu *appChatUpdate) SetOwner(owner string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, owner)
    if match {
        acu.owner = owner
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群主id不合法", nil))
    }
}

func (acu *appChatUpdate) SetAddUserList(addUserList []string) {
    acu.addUserList = make([]string, 0)
    for _, v := range addUserList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            acu.addUserList = append(acu.addUserList, v)
        }
    }
}

func (acu *appChatUpdate) SetDelUserList(delUserList []string) {
    acu.delUserList = make([]string, 0)
    for _, v := range delUserList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            acu.delUserList = append(acu.delUserList, v)
        }
    }
}

func (acu *appChatUpdate) checkData() {
    if len(acu.chatId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群聊名不能为空", nil))
    }
}

func (acu *appChatUpdate) SendRequest() api.ApiResult {
    acu.checkData()

    reqData := make(map[string]interface{})
    reqData["chatid"] = acu.chatId
    if len(acu.name) > 0 {
        reqData["name"] = acu.name
    }
    if len(acu.owner) > 0 {
        reqData["owner"] = acu.owner
    }
    if len(acu.addUserList) > 0 {
        reqData["add_user_list"] = acu.addUserList
    }
    if len(acu.delUserList) > 0 {
        reqData["del_user_list"] = acu.delUserList
    }
    reqBody := mpf.JsonMarshal(reqData)

    acu.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/appchat/update?access_token=" + wx.NewUtilWx().GetCorpAccessToken(acu.corpId, acu.agentTag)
    client, req := acu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := acu.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewAppChatUpdate(corpId, agentTag string) *appChatUpdate {
    acu := &appChatUpdate{wx.NewBaseWxCorp(), "", "", "", "", "", make([]string, 0), make([]string, 0)}
    acu.corpId = corpId
    acu.agentTag = agentTag
    acu.ReqContentType = project.HTTPContentTypeJSON
    acu.ReqMethod = fasthttp.MethodPost
    return acu
}
