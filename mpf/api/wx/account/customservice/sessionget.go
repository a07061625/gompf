/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 23:10
 */
package customservice

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type sessionGet struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
}

func (sg *sessionGet) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        sg.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (sg *sessionGet) checkData() {
    if len(sg.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    sg.ReqData["openid"] = sg.openid
}

func (sg *sessionGet) SendRequest() api.APIResult {
    sg.checkData()

    sg.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(sg.appId)
    sg.ReqURI = "https://api.weixin.qq.com/customservice/kfsession/getsession?" + mpf.HTTPCreateParams(sg.ReqData, "none", 1)
    client, req := sg.GetRequest()

    resp, result := sg.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewSessionGet(appId string) *sessionGet {
    sg := &sessionGet{wx.NewBaseWxAccount(), "", ""}
    sg.appId = appId
    return sg
}
