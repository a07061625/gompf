/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 17:25
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/corp"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送互联企业消息
type linkedMessageSend struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    userList  []string               // 成员ID列表
    partyList []string               // 部门ID列表
    tagList   []uint                 // 标签ID列表
    sendFlag  int                    // 发送消息标识,默认0 0:不发送给应用可见范围内的所有人 1:发送给应用可见范围内的所有人
    safeFlag  int                    // 保密消息标识,默认0 0:否 1:是
    msgType   string                 // 消息类型
    msgData   map[string]interface{} // 消息数据
}

func (lms *linkedMessageSend) SetUserList(userList []string) {
    lms.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            lms.userList = append(lms.userList, v)
        }
    }
    if len(lms.userList) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员ID不能超过1000个", nil))
    }
}

func (lms *linkedMessageSend) SetPartyList(partyList []string) {
    lms.partyList = make([]string, 0)
    for _, v := range partyList {
        if len(v) > 0 {
            lms.partyList = append(lms.partyList, v)
        }
    }
    if len(lms.partyList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门ID不能超过100个", nil))
    }
}

func (lms *linkedMessageSend) SetTagList(tagList []uint) {
    lms.tagList = make([]uint, 0)
    for _, v := range tagList {
        if v > 0 {
            lms.tagList = append(lms.tagList, v)
        }
    }
    if len(lms.tagList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签ID不能超过100个", nil))
    }
}

func (lms *linkedMessageSend) SetSendFlag(sendFlag int) {
    if (sendFlag == 0) || (sendFlag == 1) {
        lms.sendFlag = sendFlag
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发送消息标识不合法", nil))
    }
}

func (lms *linkedMessageSend) SetSafeFlag(safeFlag int) {
    if (safeFlag == 0) || (safeFlag == 1) {
        lms.safeFlag = safeFlag
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "保密消息标识不合法", nil))
    }
}

func (lms *linkedMessageSend) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := corp.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不支持", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息数据不能为空", nil))
    }
    lms.msgType = msgType
    lms.msgData = msgData
}

func (lms *linkedMessageSend) checkData() {
    if len(lms.msgType) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不能为空", nil))
    }
}

func (lms *linkedMessageSend) SendRequest() api.APIResult {
    lms.checkData()

    agentInfo := wx.NewConfig().GetCorp(lms.corpId).GetAgentInfo(lms.agentTag)
    reqData := make(map[string]interface{})
    reqData["agentid"] = agentInfo["id"]
    reqData["toall"] = lms.sendFlag
    reqData["safe"] = lms.safeFlag
    if len(lms.userList) > 0 {
        reqData["touser"] = lms.userList
    }
    if len(lms.partyList) > 0 {
        reqData["toparty"] = lms.partyList
    }
    if len(lms.tagList) > 0 {
        reqData["totag"] = lms.tagList
    }
    reqData["msgtype"] = lms.msgType
    reqData[lms.msgType] = lms.msgData
    reqBody := mpf.JSONMarshal(reqData)

    lms.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/linkedcorp/message/send?access_token=" + wx.NewUtilWx().GetCorpAccessToken(lms.corpId, lms.agentTag)
    client, req := lms.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := lms.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewLinkedMessageSend(corpId, agentTag string) *linkedMessageSend {
    lms := &linkedMessageSend{wx.NewBaseWxCorp(), "", "", make([]string, 0), make([]string, 0), make([]uint, 0), 0, 0, "", make(map[string]interface{})}
    lms.corpId = corpId
    lms.agentTag = agentTag
    lms.sendFlag = 0
    lms.safeFlag = 0
    lms.ReqContentType = project.HTTPContentTypeJSON
    lms.ReqMethod = fasthttp.MethodPost
    return lms
}
