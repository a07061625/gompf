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

// 根据标签进行群发
type massSendAll struct {
    wx.BaseWxAccount
    appId       string
    msgType     string                 // 消息类型
    msgData     map[string]interface{} // 消息数据
    filter      map[string]interface{} // 接收者数据
    sendReprint int                    // 转载群发标识
    msgId       string                 // 群发消息ID
}

func (msa *massSendAll) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := account.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不支持", nil))
    } else if len(msgData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息数据不能为空", nil))
    }
    msa.msgType = msgType
    msa.msgData = msgData
}

func (msa *massSendAll) SetSendReprint(sendReprint int) {
    if (sendReprint == 0) || (sendReprint == 1) {
        msa.sendReprint = sendReprint
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "转载群发标识不合法", nil))
    }
}

func (msa *massSendAll) SetFilter(filter map[string]interface{}) {
    if len(filter) > 0 {
        msa.filter = filter
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "接收者数据不合法", nil))
    }
}

func (msa *massSendAll) SetMsgId(msgId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,64}$`, msgId)
    if match {
        msa.msgId = msgId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "群发消息ID不合法", nil))
    }
}

func (msa *massSendAll) checkData() {
    if len(msa.msgType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不能为空", nil))
    }
    if (msa.msgType == account.MessageTypeMpNews) && (msa.sendReprint < 0) {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "转载群发标识不能为空", nil))
    }
    if len(msa.filter) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "接收者数据不能为空", nil))
    }
    if len(msa.msgId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "群发消息ID不能为空", nil))
    }
}

func (msa *massSendAll) SendRequest() api.ApiResult {
    msa.checkData()

    reqData := make(map[string]interface{})
    reqData["msgtype"] = msa.msgType
    reqData[msa.msgType] = msa.msgData
    reqData["filter"] = msa.filter
    reqData["clientmsgid"] = msa.msgId
    if msa.msgType == account.MessageTypeMpNews {
        reqData["send_ignore_reprint"] = msa.sendReprint
    }
    reqBody := mpf.JsonMarshal(reqData)
    msa.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=" + wx.NewUtilWx().GetSingleAccessToken(msa.appId)
    client, req := msa.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := msa.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewMassSendAll(appId string) *massSendAll {
    msa := &massSendAll{wx.NewBaseWxAccount(), "", "", make(map[string]interface{}), make(map[string]interface{}), 0, ""}
    msa.appId = appId
    msa.sendReprint = -1
    msa.ReqContentType = project.HttpContentTypeJson
    msa.ReqMethod = fasthttp.MethodPost
    return msa
}
