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

type accountInviteWorker struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
    inviteWx  string // 客服微信号
}

func (aiw *accountInviteWorker) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        aiw.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (aiw *accountInviteWorker) SetInviteWx(inviteWx string) {
    if len(inviteWx) > 0 {
        aiw.inviteWx = inviteWx
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服微信号不合法", nil))
    }
}

func (aiw *accountInviteWorker) checkData() {
    if len(aiw.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    if len(aiw.inviteWx) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服微信号不能为空", nil))
    }
    aiw.ReqData["kf_account"] = aiw.kfAccount
    aiw.ReqData["invite_wx"] = aiw.inviteWx
}

func (aiw *accountInviteWorker) SendRequest() api.ApiResult {
    aiw.checkData()

    reqBody := mpf.JsonMarshal(aiw.ReqData)
    aiw.ReqUrl = "https://api.weixin.qq.com/customservice/kfaccount/inviteworker?access_token=" + wx.NewUtilWx().GetSingleAccessToken(aiw.appId)
    client, req := aiw.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := aiw.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewAccountInviteWorker(appId string) *accountInviteWorker {
    aiw := &accountInviteWorker{wx.NewBaseWxAccount(), "", "", ""}
    aiw.appId = appId
    aiw.ReqContentType = project.HTTPContentTypeJSON
    aiw.ReqMethod = fasthttp.MethodPost
    return aiw
}
