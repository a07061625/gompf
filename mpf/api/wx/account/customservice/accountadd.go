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

type accountAdd struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
    nickname  string // 客服昵称
}

func (aa *accountAdd) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        aa.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (aa *accountAdd) SetNickname(nickname string) {
    if len(nickname) > 0 {
        trueName := []rune(nickname)
        aa.nickname = string(trueName[:16])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服昵称不合法", nil))
    }
}

func (aa *accountAdd) checkData() {
    if len(aa.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    if len(aa.nickname) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服昵称不能为空", nil))
    }
    aa.ReqData["kf_account"] = aa.kfAccount
    aa.ReqData["nickname"] = aa.nickname
}

func (aa *accountAdd) SendRequest() api.ApiResult {
    aa.checkData()

    reqBody := mpf.JSONMarshal(aa.ReqData)
    aa.ReqUrl = "https://api.weixin.qq.com/customservice/kfaccount/add?access_token=" + wx.NewUtilWx().GetSingleAccessToken(aa.appId)
    client, req := aa.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := aa.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewAccountAdd(appId string) *accountAdd {
    aa := &accountAdd{wx.NewBaseWxAccount(), "", "", ""}
    aa.appId = appId
    aa.ReqContentType = project.HTTPContentTypeJSON
    aa.ReqMethod = fasthttp.MethodPost
    return aa
}
