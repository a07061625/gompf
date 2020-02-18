/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 15:50
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/account"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 根据OpenID列表群发
type massSend struct {
    wx.BaseWxAccount
    appId       string
    msgType     string                 // 消息类型
    msgData     map[string]interface{} // 消息数据
    openidList  []string               // 用户openid列表
    sendReprint int                    // 转载群发标识
}

func (ms *massSend) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := account.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不支持", nil))
    } else if len(msgData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息数据不能为空", nil))
    }
    ms.msgType = msgType
    ms.msgData = msgData
}

func (ms *massSend) SetSendReprint(sendReprint int) {
    if (sendReprint == 0) || (sendReprint == 1) {
        ms.sendReprint = sendReprint
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "转载群发标识不合法", nil))
    }
}

func (ms *massSend) SetOpenidList(openidList []string) {
    ms.openidList = make([]string, 0)
    for _, v := range openidList {
        match, _ := regexp.MatchString(project.RegexWxOpenid, v)
        if match {
            ms.openidList = append(ms.openidList, v)
        }
    }
}

func (ms *massSend) checkData() {
    if len(ms.msgType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不能为空", nil))
    }
    if (ms.msgType == account.MessageTypeMpNews) && (ms.sendReprint < 0) {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "转载群发标识不能为空", nil))
    }
    if len(ms.openidList) < 2 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能少于2个", nil))
    } else if len(ms.openidList) > 10000 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid列表不能超过10000个", nil))
    }
}

func (ms *massSend) SendRequest() api.ApiResult {
    ms.checkData()

    reqData := make(map[string]interface{})
    reqData["msgtype"] = ms.msgType
    reqData[ms.msgType] = ms.msgData
    reqData["touser"] = ms.openidList
    if ms.msgType == account.MessageTypeMpNews {
        reqData["send_ignore_reprint"] = ms.sendReprint
    }
    reqBody := mpf.JsonMarshal(reqData)
    ms.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ms.appId)
    client, req := ms.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ms.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMassSend(appId string) *massSend {
    ms := &massSend{wx.NewBaseWxAccount(), "", "", make(map[string]interface{}), make([]string, 0), 0}
    ms.appId = appId
    ms.sendReprint = -1
    ms.ReqContentType = project.HTTPContentTypeJSON
    ms.ReqMethod = fasthttp.MethodPost
    return ms
}
