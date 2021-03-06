/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 19:15
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取登录用户信息
type loginInfoGet struct {
    wx.BaseWxProvider
    authCode string // 授权码
}

func (lig *loginInfoGet) SetAuthCode(authCode string) {
    if len(authCode) > 0 {
        lig.authCode = authCode
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权码不合法", nil))
    }
}

func (lig *loginInfoGet) checkData() {
    if len(lig.authCode) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权码不能为空", nil))
    }
    lig.ReqData["auth_code"] = lig.authCode
}

func (lig *loginInfoGet) SendRequest() api.APIResult {
    lig.checkData()

    reqBody := mpf.JSONMarshal(lig.ReqData)
    lig.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/service/get_login_info?access_token=" + wx.NewUtilWx().GetProviderToken()
    client, req := lig.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := lig.SendInner(client, req, errorcode.WxProviderRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewLoginInfoGet() *loginInfoGet {
    lig := &loginInfoGet{wx.NewBaseWxProvider(), ""}
    lig.ReqContentType = project.HTTPContentTypeJSON
    lig.ReqMethod = fasthttp.MethodPost
    return lig
}
