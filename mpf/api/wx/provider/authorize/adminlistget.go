/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 19:45
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

// 获取应用的管理员列表
type adminListGet struct {
    wx.BaseWxProvider
    authCorpId string // 授权企业ID
    agentId    string // 应用ID
}

func (alg *adminListGet) SetAgentId(agentId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, agentId)
    if match {
        alg.agentId = agentId
    } else {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "应用ID不合法", nil))
    }
}

func (alg *adminListGet) checkData() {
    if len(alg.agentId) == 0 {
        panic(mperr.NewWxProvider(errorcode.WxProviderParam, "应用ID不能为空", nil))
    }
    alg.ReqData["agentid"] = alg.agentId
}

func (alg *adminListGet) SendRequest() api.APIResult {
    alg.checkData()

    reqBody := mpf.JSONMarshal(alg.ReqData)
    alg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/service/get_admin_list?suite_access_token=" + wx.NewUtilWx().GetProviderSuiteToken()
    client, req := alg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := alg.SendInner(client, req, errorcode.WxProviderRequestPost)
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

func NewAdminListGet(corpId string) *adminListGet {
    alg := &adminListGet{wx.NewBaseWxProvider(), "", ""}
    alg.ReqData["auth_corpid"] = corpId
    alg.ReqContentType = project.HTTPContentTypeJSON
    alg.ReqMethod = fasthttp.MethodPost
    return alg
}
