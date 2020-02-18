/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 21:57
 */
package customservice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

type accountDel struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
}

func (ad *accountDel) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        ad.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (ad *accountDel) checkData() {
    if len(ad.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    ad.ReqData["kf_account"] = ad.kfAccount
}

func (ad *accountDel) SendRequest() api.APIResult {
    ad.checkData()

    ad.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(ad.appId)
    ad.ReqURI = "https://api.weixin.qq.com/customservice/kfaccount/del?" + mpf.HTTPCreateParams(ad.ReqData, "none", 1)
    client, req := ad.GetRequest()

    resp, result := ad.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAccountDel(appId string) *accountDel {
    ad := &accountDel{wx.NewBaseWxAccount(), "", ""}
    ad.appId = appId
    return ad
}
