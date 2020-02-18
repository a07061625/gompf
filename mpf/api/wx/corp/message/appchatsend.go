/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 17:08
 */
package message

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/corp"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送群聊会话消息
type appChatSend struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    chatId   string                 // 群id
    safeFlag int                    // 保密消息标识,默认0 0:否 1:是
    msgType  string                 // 消息类型
    msgData  map[string]interface{} // 消息数据
}

func (acs *appChatSend) SetChatId(chatId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, chatId)
    if match {
        acs.chatId = chatId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群id不合法", nil))
    }
}

func (acs *appChatSend) SetSafeFlag(safeFlag int) {
    if (safeFlag == 0) || (safeFlag == 1) {
        acs.safeFlag = safeFlag
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "保密消息标识不合法", nil))
    }
}

func (acs *appChatSend) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    if msgType == corp.MessageTypeMiniNotice {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不支持", nil))
    }
    _, ok := corp.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不支持", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息数据不能为空", nil))
    }
    acs.msgType = msgType
    acs.msgData = msgData
}

func (acs *appChatSend) checkData() {
    if len(acs.chatId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "群id不能为空", nil))
    }
    if len(acs.msgData) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不能为空", nil))
    }
}

func (acs *appChatSend) SendRequest() api.APIResult {
    acs.checkData()

    reqData := make(map[string]interface{})
    reqData["chatid"] = acs.chatId
    reqData["safe"] = strconv.Itoa(acs.safeFlag)
    reqData["msgtype"] = acs.msgType
    reqData[acs.msgType] = acs.msgData
    reqBody := mpf.JSONMarshal(reqData)

    acs.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token=" + wx.NewUtilWx().GetCorpAccessToken(acs.corpId, acs.agentTag)
    client, req := acs.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := acs.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewAppChatSend(corpId, agentTag string) *appChatSend {
    acs := &appChatSend{wx.NewBaseWxCorp(), "", "", "", 0, "", make(map[string]interface{})}
    acs.corpId = corpId
    acs.agentTag = agentTag
    acs.safeFlag = 0
    acs.ReqContentType = project.HTTPContentTypeJSON
    acs.ReqMethod = fasthttp.MethodPost
    return acs
}
