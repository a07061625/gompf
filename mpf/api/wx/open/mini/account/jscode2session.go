/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 12:53
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取第三方平台 session_key 和 openid
type jsCode2Session struct {
    wx.BaseWxOpen
    jsCode string // 登录凭证
}

func (jcs *jsCode2Session) SetJsCode(jsCode string) {
    if len(jsCode) > 0 {
        jcs.jsCode = jsCode
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "登录凭证不合法", nil))
    }
}

func (jcs *jsCode2Session) checkData() {
    if len(jcs.jsCode) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "登录凭证不能为空", nil))
    }
    jcs.ReqData["js_code"] = jcs.jsCode
}

func (jcs *jsCode2Session) SendRequest() api.ApiResult {
    jcs.checkData()

    jcs.ReqData["component_access_token"] = wx.NewUtilWx().GetOpenAccessToken()
    jcs.ReqUrl = "https://api.weixin.qq.com/sns/component/jscode2session?" + mpf.HttpCreateParams(jcs.ReqData, "none", 1)
    client, req := jcs.GetRequest()

    resp, result := jcs.SendInner(client, req, errorcode.WxOpenRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["openid"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewJsCode2Session(appId string) *jsCode2Session {
    conf := wx.NewConfig().GetOpen()
    jcs := &jsCode2Session{wx.NewBaseWxOpen(), ""}
    jcs.ReqData["component_appid"] = conf.GetAppId()
    jcs.ReqData["appid"] = appId
    jcs.ReqData["grant_type"] = "authorization_code"
    return jcs
}
