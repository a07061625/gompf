/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 0:19
 */
package message

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/mini"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 发送客服消息给用户
type customSend struct {
    wx.BaseWxMini
    appId   string                 // 应用ID
    touser  string                 // 用户openid
    msgType string                 // 消息类型
    msgData map[string]interface{} // 消息数据
}

func (cs *customSend) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        cs.touser = openid
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "用户openid不合法", nil))
    }
}

func (cs *customSend) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := mini.MessageCustomTypes[msgType]
    if !ok {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "消息类型不合法", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "消息数据不能为空", nil))
    }
    cs.msgType = msgType
    cs.msgData = msgData
}

func (cs *customSend) checkData() {
    if len(cs.touser) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "用户openid不能为空", nil))
    }
    if len(cs.msgType) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "消息类型不能为空", nil))
    }
}

func (cs *customSend) SendRequest(getType string) api.ApiResult {
    cs.checkData()
    reqData := make(map[string]interface{})
    reqData["touser"] = cs.touser
    reqData["msgtype"] = cs.msgType
    reqData[cs.msgType] = cs.msgData
    reqBody := mpf.JSONMarshal(reqData)

    cs.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + wx.NewUtilWx().GetSingleCache(cs.appId, getType)
    client, req := cs.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cs.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewCustomSend(appId string) *customSend {
    cs := &customSend{wx.NewBaseWxMini(), "", "", "", make(map[string]interface{})}
    cs.appId = appId
    cs.ReqContentType = project.HTTPContentTypeJSON
    cs.ReqMethod = fasthttp.MethodPost
    cs.ReqHeader["Expect"] = ""
    return cs
}
