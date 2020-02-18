/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 18:45
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 通过授权码获取用户详细信息
type userInfoCode struct {
    wx.BaseWxAccount
    appId    string
    authCode string // 授权码
}

func (uic *userInfoCode) SetAuthCode(authCode string) {
    if len(authCode) > 0 {
        uic.authCode = authCode
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不合法", nil))
    }
}

func (uic *userInfoCode) checkData() {
    if len(uic.authCode) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不能为空", nil))
    }
}

func (uic *userInfoCode) SendRequest() api.APIResult {
    uic.checkData()

    userBase := NewUserBase(uic.appId)
    userBase.SetAuthCode(uic.authCode)
    baseRes := userBase.SendRequest()
    if baseRes.Code > 0 {
        return baseRes
    }

    baseData := baseRes.Data.(map[string]interface{})
    uic.ReqData["access_token"] = baseData["access_token"].(string)
    uic.ReqData["openid"] = baseData["openid"].(string)
    uic.ReqURI = "https://api.weixin.qq.com/sns/userinfo?" + mpf.HTTPCreateParams(uic.ReqData, "none", 1)
    client, req := uic.GetRequest()

    resp, result := uic.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["openid"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewUserInfoCode(appId string) *userInfoCode {
    uic := &userInfoCode{wx.NewBaseWxAccount(), "", ""}
    uic.appId = appId
    uic.ReqData["lang"] = "zh_CN"
    return uic
}
