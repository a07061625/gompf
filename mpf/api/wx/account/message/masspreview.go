/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 15:37
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

// 预览接口
type massPreview struct {
    wx.BaseWxAccount
    appId   string
    msgType string                 // 消息类型
    msgData map[string]interface{} // 消息数据
    openid  string                 // 用户openid
    wxName  string                 // 公众号名称
}

func (mp *massPreview) SetMsgInfo(msgType string, msgData map[string]interface{}) {
    _, ok := account.MessageTypes[msgType]
    if !ok {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不支持", nil))
    } else if len(msgData) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息数据不能为空", nil))
    }
    mp.msgType = msgType
    mp.msgData = msgData
}

func (mp *massPreview) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        mp.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (mp *massPreview) SetWxName(wxName string) {
    if len(wxName) > 0 {
        mp.wxName = wxName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "公众号名称不合法", nil))
    }
}

func (mp *massPreview) checkData() {
    if len(mp.msgType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "消息类型不能为空", nil))
    }
    if (len(mp.openid) == 0) && (len(mp.wxName) == 0) {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid和公众号名称不能都为空", nil))
    }
}

func (mp *massPreview) SendRequest() api.ApiResult {
    mp.checkData()

    reqData := make(map[string]interface{})
    reqData["msgtype"] = mp.msgType
    reqData[mp.msgType] = mp.msgData
    if len(mp.openid) > 0 {
        reqData["touser"] = mp.openid
    }
    if len(mp.wxName) > 0 {
        reqData["towxname"] = mp.wxName
    }
    reqBody := mpf.JsonMarshal(reqData)
    mp.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/mass/preview?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mp.appId)
    client, req := mp.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mp.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewMassPreview(appId string) *massPreview {
    mp := &massPreview{wx.NewBaseWxAccount(), "", "", make(map[string]interface{}), "", ""}
    mp.appId = appId
    mp.ReqContentType = project.HttpContentTypeJson
    mp.ReqMethod = fasthttp.MethodPost
    return mp
}
