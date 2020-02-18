/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 15:17
 */
package message

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/corp"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送企业消息
type messageSend struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    userList  []string               // 成员ID列表
    partyList []uint                 // 部门ID列表
    tagList   []uint                 // 标签ID列表
    msgType   string                 // 消息类型
    msgData   map[string]interface{} // 消息数据
    safeTag   int                    // 保密消息标识,默认0 0:否 1:是
}

func (ms *messageSend) SetUserList(userList []string) {
    ms.userList = make([]string, 0)
    for _, v := range userList {
        match, _ := regexp.MatchString(project.RegexDigitLower, v)
        if match {
            ms.userList = append(ms.userList, v)
        }
    }
    if len(ms.userList) > 1000 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员ID不能超过1000个", nil))
    }
}

func (ms *messageSend) SetPartyList(partyList []uint) {
    ms.partyList = make([]uint, 0)
    for _, v := range partyList {
        if v > 0 {
            ms.partyList = append(ms.partyList, v)
        }
    }
    if len(ms.partyList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "部门ID不能超过100个", nil))
    }
}

func (ms *messageSend) SetTagList(tagList []uint) {
    ms.tagList = make([]uint, 0)
    for _, v := range tagList {
        if v > 0 {
            ms.tagList = append(ms.tagList, v)
        }
    }
    if len(ms.tagList) > 100 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "标签ID不能超过100个", nil))
    }
}

func (ms *messageSend) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := corp.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不支持", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息数据不能为空", nil))
    }
    ms.msgType = msgType
    ms.msgData = msgData
}

func (ms *messageSend) SetSafeTag(safeTag int) {
    if (safeTag == 0) || (safeTag == 1) {
        ms.safeTag = safeTag
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "保密消息标识不合法", nil))
    }
}

func (ms *messageSend) checkData() {
    if len(ms.msgType) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "消息类型不能为空", nil))
    }
}

func (ms *messageSend) SendRequest(getType string) api.ApiResult {
    ms.checkData()

    agentInfo := wx.NewConfig().GetCorp(ms.corpId).GetAgentInfo(ms.agentTag)
    reqData := make(map[string]interface{})
    reqData["agentid"] = agentInfo["id"]
    reqData["safe"] = strconv.Itoa(ms.safeTag)
    if len(ms.userList) == 0 {
        reqData["touser"] = "@all"
    } else {
        reqData["touser"] = strings.Join(ms.userList, "|")
    }
    partyStr := ""
    for _, partyId := range ms.partyList {
        partyStr += "|" + strconv.Itoa(int(partyId))
    }
    if len(partyStr) == 0 {
        reqData["toparty"] = ""
    } else {
        reqData["toparty"] = partyStr[1:]
    }
    tagStr := ""
    for _, tagId := range ms.tagList {
        tagStr += "|" + strconv.Itoa(int(tagId))
    }
    if len(tagStr) == 0 {
        reqData["totag"] = ""
    } else {
        reqData["totag"] = tagStr[1:]
    }
    reqData["msgtype"] = ms.msgType
    reqData[ms.msgType] = ms.msgData
    reqBody := mpf.JSONMarshal(reqData)

    ms.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + wx.NewUtilWx().GetCorpCache(ms.corpId, ms.agentTag, getType)
    client, req := ms.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ms.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewMessageSend(corpId, agentTag string) *messageSend {
    ms := &messageSend{wx.NewBaseWxCorp(), "", "", make([]string, 0), make([]uint, 0), make([]uint, 0), "", make(map[string]interface{}), 0}
    ms.corpId = corpId
    ms.agentTag = agentTag
    ms.ReqContentType = project.HTTPContentTypeJSON
    ms.ReqMethod = fasthttp.MethodPost
    return ms
}
