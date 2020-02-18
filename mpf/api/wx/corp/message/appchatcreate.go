/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 15:46
 */
package message

import (
    "regexp"
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 创建群聊会话
type appChatCreate struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    name     string   // 群聊名
    owner    string   // 群主id
    userList []string // 群成员列表
    chatId   string   // 群id
}

func (acc *appChatCreate) SetName(name string) {
    if len(name) > 0 {
        trueName := []rune(name)
        acc.name = string(trueName[:25])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群聊名不合法", nil))
    }
}

func (acc *appChatCreate) SetOwner(owner string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, owner)
    if match {
        acc.owner = owner
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群主id不合法", nil))
    }
}

func (acc *appChatCreate) SetUserList(userList []string) {
    acc.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            acc.userList = append(acc.userList, v)
        }
    }
    if len(acc.userList) < 2 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群成员不能少于2个", nil))
    }
    if len(acc.userList) > 500 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群成员不能超过500个", nil))
    }
}

func (acc *appChatCreate) SetChatId(chatId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, chatId)
    if match {
        acc.chatId = chatId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群id不合法", nil))
    }
}

func (acc *appChatCreate) checkData() {
    if len(acc.name) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群聊名不能为空", nil))
    }
    if len(acc.userList) < 2 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群成员不能少于2个", nil))
    }
}

func (acc *appChatCreate) SendRequest() api.ApiResult {
    acc.checkData()

    reqData := make(map[string]interface{})
    reqData["name"] = acc.name
    if len(acc.owner) > 0 {
        reqData["owner"] = acc.owner
    }
    reqData["chatid"] = acc.chatId
    reqData["userlist"] = acc.userList
    reqBody := mpf.JsonMarshal(reqData)

    acc.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/appchat/create?access_token=" + wx.NewUtilWx().GetCorpAccessToken(acc.corpId, acc.agentTag)
    client, req := acc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := acc.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewAppChatCreate(corpId, agentTag string) *appChatCreate {
    acc := &appChatCreate{wx.NewBaseWxCorp(), "", "", "", "", make([]string, 0), ""}
    acc.corpId = corpId
    acc.agentTag = agentTag
    acc.chatId = mpf.ToolCreateNonceStr(8, "numlower") + strconv.FormatInt(time.Now().Unix(), 10)
    acc.ReqContentType = project.HTTPContentTypeJSON
    acc.ReqMethod = fasthttp.MethodPost
    return acc
}
