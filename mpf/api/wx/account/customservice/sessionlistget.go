/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 23:15
 */
package customservice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

type sessionListGet struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
}

func (slg *sessionListGet) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        slg.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (slg *sessionListGet) checkData() {
    if len(slg.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    slg.ReqData["kf_account"] = slg.kfAccount
}

func (slg *sessionListGet) SendRequest() api.APIResult {
    slg.checkData()

    slg.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(slg.appId)
    slg.ReqURI = "https://api.weixin.qq.com/customservice/kfsession/getsessionlist?" + mpf.HTTPCreateParams(slg.ReqData, "none", 1)
    client, req := slg.GetRequest()

    resp, result := slg.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["sessionlist"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewSessionListGet(appId string) *sessionListGet {
    slg := &sessionListGet{wx.NewBaseWxAccount(), "", ""}
    slg.appId = appId
    return slg
}
