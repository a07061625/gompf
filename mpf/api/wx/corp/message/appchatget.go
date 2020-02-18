/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 15:46
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取群聊会话
type appChatGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    chatId   string // 群id
}

func (acg *appChatGet) SetChatId(chatId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, chatId)
    if match {
        acg.chatId = chatId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群id不合法", nil))
    }
}

func (acg *appChatGet) checkData() {
    if len(acg.chatId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群id不能为空", nil))
    }
    acg.ReqData["chatid"] = acg.chatId
}

func (acg *appChatGet) SendRequest() api.APIResult {
    acg.checkData()

    acg.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(acg.corpId, acg.agentTag)
    acg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/get?" + mpf.HTTPCreateParams(acg.ReqData, "none", 1)
    client, req := acg.GetRequest()

    resp, result := acg.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewAppChatGet(corpId, agentTag string) *appChatGet {
    acg := &appChatGet{wx.NewBaseWxCorp(), "", "", ""}
    acg.corpId = corpId
    acg.agentTag = agentTag
    return acg
}
