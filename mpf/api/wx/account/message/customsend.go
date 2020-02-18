/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 13:20
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

// 发送客服消息
type customSend struct {
    wx.BaseWxAccount
    appId       string
    accessToken string                 // 令牌
    openid      string                 // 用户openid
    msgType     string                 // 消息类型
    msgData     map[string]interface{} // 消息数据
}

func (cs *customSend) SetAccessToken(accessToken string) {
    if len(accessToken) > 0 {
        cs.accessToken = accessToken
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "令牌不合法", nil))
    }
}

func (cs *customSend) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        cs.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (cs *customSend) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := msgCustomTypes[msgType]
    if !ok {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不合法", nil))
    }
    if len(msgData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息数据不能为空", nil))
    }
    cs.msgType = msgType
    cs.msgData = msgData
}

func (cs *customSend) checkData() {
    if len(cs.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if len(cs.msgType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不能为空", nil))
    }
}

func (cs *customSend) SendRequest(getType string) api.ApiResult {
    cs.checkData()

    reqData := make(map[string]interface{})
    reqData["touser"] = cs.openid
    reqData["msgtype"] = cs.msgType
    reqData[cs.msgType] = cs.msgData
    reqBody := mpf.JSONMarshal(reqData)
    if len(cs.accessToken) > 0 {
        cs.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + cs.accessToken
    } else {
        cs.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + wx.NewUtilWx().GetSingleCache(cs.appId, getType)
    }
    cs.ReqHeader["Expect"] = ""
    client, req := cs.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cs.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewCustomSend(appId string) *customSend {
    cs := &customSend{wx.NewBaseWxAccount(), "", "", "", "", make(map[string]interface{})}
    cs.appId = appId
    cs.ReqContentType = project.HTTPContentTypeJSON
    cs.ReqMethod = fasthttp.MethodPost
    return cs
}

const (
    MsgCustomTypeText   = "text"
    MsgCustomTypeImage  = "image"
    MsgCustomTypeVoice  = "voice"
    MsgCustomTypeVideo  = "video"
    MsgCustomTypeMusic  = "music"
    MsgCustomTypeNews   = "news"
    MsgCustomTypeMpNews = "mpnews"
    MsgCustomTypeMenu   = "msgmenu"
)

var (
    msgCustomTypes map[string]string
)

func init() {
    msgCustomTypes = make(map[string]string)
    msgCustomTypes[MsgCustomTypeText] = "文本"
    msgCustomTypes[MsgCustomTypeImage] = "图片"
    msgCustomTypes[MsgCustomTypeVoice] = "语音"
    msgCustomTypes[MsgCustomTypeVideo] = "视频"
    msgCustomTypes[MsgCustomTypeMusic] = "音乐"
    msgCustomTypes[MsgCustomTypeNews] = "图文(跳转到外链)"
    msgCustomTypes[MsgCustomTypeMpNews] = "图文(跳转到图文消息页面)"
    msgCustomTypes[MsgCustomTypeMenu] = "菜单"
}
