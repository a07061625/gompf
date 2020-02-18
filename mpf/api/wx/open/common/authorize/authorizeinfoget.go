/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 11:08
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

type authorizeInfoGet struct {
    wx.BaseWxOpen
    appId string // 授权公众号或小程序的appid
}

func (aig *authorizeInfoGet) SendRequest() api.ApiResult {
    reqBody := mpf.JSONMarshal(aig.ReqData)
    aig.ReqUrl = "https://api.weixin.qq.com/cgi-bin/component/api_get_authorizer_info?component_access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := aig.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := aig.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["authorizer_info"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAuthorizeInfoGet(appId string) *authorizeInfoGet {
    conf := wx.NewConfig().GetOpen()
    aig := &authorizeInfoGet{wx.NewBaseWxOpen(), ""}
    aig.ReqData["component_appid"] = conf.GetAppId()
    aig.ReqData["authorizer_appid"] = appId
    aig.ReqContentType = project.HTTPContentTypeJSON
    aig.ReqMethod = fasthttp.MethodPost
    return aig
}
