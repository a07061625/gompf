/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 9:07
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type authorizeUrl struct {
    wx.BaseWxOpen
    urlAuthCallback string
}

func (au *authorizeUrl) SendRequest() api.ApiResult {
    reqBody := mpf.JsonMarshal(au.ReqData)
    au.ReqUrl = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := au.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := au.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    res := make(map[string]string)
    res["url"] = ""
    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    authCode, ok := respData["pre_auth_code"]
    if ok {
        urlParams := make(map[string]string)
        urlParams["component_appid"] = au.ReqData["component_appid"]
        urlParams["pre_auth_code"] = authCode.(string)
        urlParams["redirect_uri"] = au.urlAuthCallback
        res["url"] = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?" + mpf.HttpCreateParams(urlParams, "none", 1)
    }
    result.Data = res

    return result
}

func NewAuthorizeUrl(appId string) *authorizeUrl {
    conf := wx.NewConfig().GetOpen()
    au := &authorizeUrl{wx.NewBaseWxOpen(), ""}
    au.ReqData["component_appid"] = conf.GetAppId()
    au.urlAuthCallback = conf.GetUrlAuthCallback()
    au.ReqContentType = project.HTTPContentTypeJSON
    au.ReqMethod = fasthttp.MethodPost
    return au
}
