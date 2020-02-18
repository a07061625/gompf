/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 16:01
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改功能介绍
type signatureModify struct {
    wx.BaseWxOpen
    appId     string // 应用ID
    signature string // 功能介绍
}

func (sm *signatureModify) SetSignature(signature string) {
    if len(signature) > 0 {
        sm.signature = signature
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "功能介绍不合法", nil))
    }
}

func (sm *signatureModify) checkData() {
    if len(sm.signature) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "功能介绍不能为空", nil))
    }
    sm.ReqData["signature"] = sm.signature
}

func (sm *signatureModify) SendRequest() api.APIResult {
    sm.checkData()

    reqBody := mpf.JSONMarshal(sm.ReqData)
    sm.ReqURI = "https://api.weixin.qq.com/cgi-bin/account/modifysignature?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(sm.appId)
    client, req := sm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sm.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewSignatureModify(appId string) *signatureModify {
    sm := &signatureModify{wx.NewBaseWxOpen(), "", ""}
    sm.appId = appId
    sm.ReqContentType = project.HTTPContentTypeJSON
    sm.ReqMethod = fasthttp.MethodPost
    return sm
}
