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

// 通过授权码获取用户基本信息
type userBase struct {
    wx.BaseWxAccount
    appId    string
    authCode string // 授权码
}

func (ub *userBase) SetAuthCode(authCode string) {
    if len(authCode) > 0 {
        ub.authCode = authCode
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不合法", nil))
    }
}

func (ub *userBase) checkData() {
    if len(ub.authCode) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不能为空", nil))
    }
    ub.ReqData["code"] = ub.authCode
}

func (ub *userBase) SendRequest() api.APIResult {
    ub.checkData()

    ub.ReqURI = "https://api.weixin.qq.com/sns/oauth2/access_token?" + mpf.HTTPCreateParams(ub.ReqData, "none", 1)
    client, req := ub.GetRequest()

    resp, result := ub.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["access_token"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewUserBase(appId string) *userBase {
    conf := wx.NewConfig().GetAccount(appId)
    ub := &userBase{wx.NewBaseWxAccount(), "", ""}
    ub.appId = appId
    ub.ReqData["appid"] = conf.GetAppId()
    ub.ReqData["secret"] = conf.GetSecret()
    ub.ReqData["grant_type"] = "authorization_code"
    return ub
}
