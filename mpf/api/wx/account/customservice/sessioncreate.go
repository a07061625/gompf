/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 21:57
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
    "github.com/valyala/fasthttp"
)

type sessionCreate struct {
    wx.BaseWxAccount
    appId     string
    kfAccount string // 客服帐号 格式为: 帐号前缀@公众号微信号
    openid    string // 用户openid
}

func (sc *sessionCreate) SetKfAccount(kfAccount string) {
    if (len(kfAccount) > 0) && (len(kfAccount) <= 30) {
        sc.kfAccount = kfAccount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不合法", nil))
    }
}

func (sc *sessionCreate) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        sc.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (sc *sessionCreate) checkData() {
    if len(sc.kfAccount) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客服帐号不能为空", nil))
    }
    if len(sc.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    sc.ReqData["kf_account"] = sc.kfAccount
    sc.ReqData["openid"] = sc.openid
}

func (sc *sessionCreate) SendRequest() api.ApiResult {
    sc.checkData()

    reqBody := mpf.JsonMarshal(sc.ReqData)
    sc.ReqUrl = "https://api.weixin.qq.com/customservice/kfsession/create?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sc.appId)
    client, req := sc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sc.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewSessionCreate(appId string) *sessionCreate {
    sc := &sessionCreate{wx.NewBaseWxAccount(), "", "", ""}
    sc.appId = appId
    sc.ReqContentType = project.HttpContentTypeJson
    sc.ReqMethod = fasthttp.MethodPost
    return sc
}
