/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 18:42
 */
package message

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

// 发送模板消息
type templateSend struct {
    wx.BaseWxAccount
    appId        string
    openid       string                 // 用户openid
    templateId   string                 // 模版ID
    redirectUrl  string                 // 重定向链接地址
    miniProgram  map[string]interface{} // 小程序跳转数据
    templateData map[string]interface{} // 模版数据
}

func (ts *templateSend) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ts.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ts *templateSend) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        ts.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模版ID不能为空", nil))
    }
}

func (ts *templateSend) SetRedirectUrl(redirectUrl string) {
    match, _ := regexp.MatchString(project.RegexUrlHttp, redirectUrl)
    if match {
        ts.redirectUrl = redirectUrl
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "重定向地址不合法", nil))
    }
}

func (ts *templateSend) SetMiniProgram(miniProgram map[string]interface{}) {
    if len(miniProgram) > 0 {
        ts.miniProgram = miniProgram
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "小程序跳转数据不合法", nil))
    }
}

func (ts *templateSend) SetTemplateData(templateData map[string]interface{}) {
    ts.templateData = templateData
}

func (ts *templateSend) checkData() {
    if len(ts.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if len(ts.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模版ID不能为空", nil))
    }
}

func (ts *templateSend) SendRequest() api.ApiResult {
    ts.checkData()

    reqData := make(map[string]interface{})
    reqData["touser"] = ts.openid
    reqData["template_id"] = ts.templateId
    reqData["url"] = ts.redirectUrl
    reqData["data"] = ts.templateData
    if len(ts.miniProgram) > 0 {
        reqData["miniprogram"] = ts.miniProgram
    }
    reqBody := mpf.JsonMarshal(reqData)
    ts.ReqUrl = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ts.appId)
    client, req := ts.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ts.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTemplateSend(appId string) *templateSend {
    ts := &templateSend{wx.NewBaseWxAccount(), "", "", "", "", make(map[string]interface{}), make(map[string]interface{})}
    ts.appId = appId
    ts.ReqContentType = project.HttpContentTypeJson
    ts.ReqMethod = fasthttp.MethodPost
    return ts
}
