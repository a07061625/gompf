/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 9:25
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

type authorizeInfo struct {
    wx.BaseWxOpen
    authCode string // 授权码
}

func (ai *authorizeInfo) SetAuthCode(authCode string) {
    if len(authCode) > 0 {
        ai.authCode = authCode
    } else {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "授权码不合法", nil))
    }
}

func (ai *authorizeInfo) checkData() {
    if len(ai.authCode) == 0 {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "授权码不能为空", nil))
    }
    ai.ReqData["authorization_code"] = ai.authCode
}

func (ai *authorizeInfo) SendRequest() api.ApiResult {
    ai.checkData()

    reqBody := mpf.JsonMarshal(ai.ReqData)
    ai.ReqUrl = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := ai.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ai.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["authorization_info"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = "授权失败,请重新授权"
    }

    return result
}

func NewAuthorizeInfo() *authorizeInfo {
    conf := wx.NewConfig().GetOpen()
    ai := &authorizeInfo{wx.NewBaseWxOpen(), ""}
    ai.ReqData["component_appid"] = conf.GetAppId()
    ai.ReqContentType = project.HttpContentTypeJson
    ai.ReqMethod = fasthttp.MethodPost
    return ai
}
