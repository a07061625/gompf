/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 17:28
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

// 推送订阅模板消息
type subscribeSend struct {
    wx.BaseWxAccount
    appId       string
    openid      string                 // 用户openid
    templateId  string                 // 模版ID
    redirectUrl string                 // 重定向地址
    miniProgram map[string]interface{} // 小程序跳转数据
    scene       int                    // 订阅场景值
    msgTitle    string                 // 消息标题
    msgData     map[string]interface{} // 消息内容
}

func (ss *subscribeSend) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ss.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ss *subscribeSend) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        ss.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模版ID不能为空", nil))
    }
}

func (ss *subscribeSend) SetRedirectUrl(redirectUrl string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, redirectUrl)
    if match {
        ss.redirectUrl = redirectUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向地址不合法", nil))
    }
}

func (ss *subscribeSend) SetMiniProgram(miniProgram map[string]interface{}) {
    if len(miniProgram) > 0 {
        ss.miniProgram = miniProgram
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "小程序跳转数据不合法", nil))
    }
}

func (ss *subscribeSend) SetScene(scene int) {
    if (scene >= 0) && (scene <= 10000) {
        ss.scene = scene
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订阅场景值不合法", nil))
    }
}

func (ss *subscribeSend) SetMsgTitle(msgTitle string) {
    if len(msgTitle) > 0 {
        trueTitle := []rune(msgTitle)
        ss.msgTitle = string(trueTitle[:15])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息标题不合法", nil))
    }
}

func (ss *subscribeSend) SetMsgData(msgData map[string]interface{}) {
    if len(msgData) > 0 {
        ss.msgData = msgData
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息内容不合法", nil))
    }
}

func (ss *subscribeSend) checkData() {
    if len(ss.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if len(ss.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模版ID不能为空", nil))
    }
    if ss.scene < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订阅场景值不能为空", nil))
    }
    if len(ss.msgTitle) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息标题不能为空", nil))
    }
    if len(ss.msgData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息内容不能为空", nil))
    }
}

func (ss *subscribeSend) SendRequest() api.APIResult {
    ss.checkData()

    reqData := make(map[string]interface{})
    reqData["touser"] = ss.openid
    reqData["template_id"] = ss.templateId
    reqData["scene"] = ss.scene
    reqData["title"] = ss.msgTitle
    reqData["data"] = ss.msgData
    if len(ss.miniProgram) > 0 {
        reqData["miniprogram"] = ss.miniProgram
    }
    if len(ss.redirectUrl) > 0 {
        reqData["url"] = ss.redirectUrl
    }
    reqBody := mpf.JSONMarshal(reqData)
    ss.ReqURI = "https://api.weixin.qq.com/cgi-bin/message/template/subscribe?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ss.appId)
    client, req := ss.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ss.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewSubscribeSend(appId string) *subscribeSend {
    ss := &subscribeSend{wx.NewBaseWxAccount(), "", "", "", "", make(map[string]interface{}), 0, "", make(map[string]interface{})}
    ss.appId = appId
    ss.scene = -1
    ss.ReqContentType = project.HTTPContentTypeJSON
    ss.ReqMethod = fasthttp.MethodPost
    return ss
}
