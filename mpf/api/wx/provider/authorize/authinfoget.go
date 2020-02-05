/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 19:54
 */
package authorize

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

// 获取企业授权信息
type authInfoGet struct {
    wx.BaseWxProvider
    authCorpId    string // 授权企业ID
    permanentCode string // 永久授权码
}

func (aig *authInfoGet) SetAuthCorpId(authCorpId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, authCorpId)
    if match {
        aig.authCorpId = authCorpId
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权企业ID不合法", nil))
    }
}

func (aig *authInfoGet) SetPermanentCode(permanentCode string) {
    if len(permanentCode) > 0 {
        aig.permanentCode = permanentCode
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "永久授权码不合法", nil))
    }
}

func (aig *authInfoGet) checkData() {
    if len(aig.authCorpId) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "授权企业ID不能为空", nil))
    }
    if len(aig.permanentCode) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "应用ID不能为空", nil))
    }
    aig.ReqData["auth_corpid"] = aig.authCorpId
    aig.ReqData["permanent_code"] = aig.permanentCode
}

func (aig *authInfoGet) SendRequest() api.ApiResult {
    aig.checkData()

    reqBody := mpf.JsonMarshal(aig.ReqData)
    aig.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_auth_info?suite_access_token=" + wx.NewUtilWx().GetProviderSuiteToken()
    client, req := aig.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := aig.SendInner(client, req, errorcode.WxProviderRequestPost)
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

func NewAuthInfoGet() *authInfoGet {
    aig := &authInfoGet{wx.NewBaseWxProvider(), "", ""}
    aig.ReqContentType = project.HttpContentTypeJson
    aig.ReqMethod = fasthttp.MethodPost
    return aig
}
