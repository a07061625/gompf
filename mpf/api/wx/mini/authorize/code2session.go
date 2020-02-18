/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 0:03
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type code2Session struct {
    wx.BaseWxMini
    jsCode string // 授权码
}

func (cs *code2Session) SetCode(jsCode string) {
    if len(jsCode) > 0 {
        cs.ReqData["js_code"] = jsCode
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "授权码不合法", nil))
    }
}

func (cs *code2Session) checkData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := cs.ReqData["js_code"]
    if !ok {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "授权码不能为空", nil))
    }

    cs.ReqURI = "https://api.weixin.qq.com/sns/jscode2session?" + mpf.HTTPCreateParams(cs.ReqData, "none", 1)

    return cs.GetRequest()
}

func (cs *code2Session) SendRequest() api.APIResult {
    client, req := cs.checkData()
    resp, result := cs.SendInner(client, req, errorcode.WxMiniRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["openid"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewCode2Session(appId string) *code2Session {
    conf := wx.NewConfig().GetAccount(appId)
    cs := &code2Session{wx.NewBaseWxMini(), ""}
    cs.ReqData["appid"] = conf.GetAppId()
    cs.ReqData["secret"] = conf.GetSecret()
    cs.ReqData["grant_type"] = "authorization_code"
    return cs
}
