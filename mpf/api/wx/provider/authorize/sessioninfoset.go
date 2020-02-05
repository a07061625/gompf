/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 20:09
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

// 设置授权配置
type sessionInfoSet struct {
    wx.BaseWxProvider
    preAuthCode string // 预授权码
    authType    int    // 授权类型,默认值为0 0:正式授权 1:测试授权
}

func (sis *sessionInfoSet) SetPreAuthCode(preAuthCode string) {
    if len(preAuthCode) > 0 {
        sis.preAuthCode = preAuthCode
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "预授权码不合法", nil))
    }
}

func (sis *sessionInfoSet) SetAuthType(authType int) {
    if (authType == 0) || (authType == 1) {
        sis.authType = authType
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权类型不合法", nil))
    }
}

func (sis *sessionInfoSet) checkData() {
    if len(sis.preAuthCode) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "预授权码不能为空", nil))
    }
}

func (sis *sessionInfoSet) SendRequest() api.ApiResult {
    sis.checkData()

    sessionInfo := make(map[string]interface{})
    sessionInfo["auth_type"] = sis.authType
    reqData := make(map[string]interface{})
    reqData["pre_auth_code"] = sis.preAuthCode
    reqData["session_info"] = sessionInfo
    reqBody := mpf.JsonMarshal(reqData)
    sis.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/set_session_info?suite_access_token=" + wx.NewUtilWx().GetProviderSuiteToken()
    client, req := sis.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sis.SendInner(client, req, errorcode.WxProviderRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewSessionInfoSet() *sessionInfoSet {
    sis := &sessionInfoSet{wx.NewBaseWxProvider(), "", 0}
    sis.authType = 0
    sis.ReqContentType = project.HttpContentTypeJson
    sis.ReqMethod = fasthttp.MethodPost
    return sis
}
