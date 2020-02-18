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
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type accountUpdate struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
    nickname  string // 客服昵称
}

func (au *accountUpdate) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        au.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (au *accountUpdate) SetNickname(nickname string) {
    if len(nickname) > 0 {
        trueName := []rune(nickname)
        au.nickname = string(trueName[:16])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服昵称不合法", nil))
    }
}

func (au *accountUpdate) checkData() {
    if len(au.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    if len(au.nickname) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服昵称不能为空", nil))
    }
    au.ReqData["kf_account"] = au.kfAccount
    au.ReqData["nickname"] = au.nickname
}

func (au *accountUpdate) SendRequest() api.ApiResult {
    au.checkData()

    reqBody := mpf.JsonMarshal(au.ReqData)
    au.ReqUrl = "https://api.weixin.qq.com/customservice/kfaccount/update?access_token=" + wx.NewUtilWx().GetSingleAccessToken(au.appId)
    client, req := au.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := au.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewAccountUpdate(appId string) *accountUpdate {
    au := &accountUpdate{wx.NewBaseWxAccount(), "", "", ""}
    au.appId = appId
    au.ReqContentType = project.HTTPContentTypeJSON
    au.ReqMethod = fasthttp.MethodPost
    return au
}
