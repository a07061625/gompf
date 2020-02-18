/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 9:31
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

type templateSend struct {
    wx.BaseWxMini
    appId           string                 // 应用ID
    openid          string                 // 用户openid
    templateId      string                 // 模板ID
    templateData    map[string]interface{} // 模板内容
    redirectUrl     string                 // 跳转页面
    formId          string                 // 表单ID
    emphasisKeyword string                 // 放大的关键词
}

func (ts *templateSend) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ts.openid = openid
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "用户openid不合法", nil))
    }
}

func (ts *templateSend) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        ts.templateId = templateId
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "模板ID不能为空", nil))
    }
}

func (ts *templateSend) SetRedirectUrl(redirectUrl string) {
    if len(redirectUrl) > 0 {
        ts.redirectUrl = redirectUrl
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "跳转页面不能为空", nil))
    }
}

func (ts *templateSend) SetFormId(formId string) {
    if len(formId) > 0 {
        ts.formId = formId
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "表单ID不能为空", nil))
    }
}

func (ts *templateSend) SetTemplateData(templateData map[string]interface{}) {
    ts.templateData = templateData
}

func (ts *templateSend) SetEmphasisKeyword(emphasisKeyword string) {
    if len(emphasisKeyword) > 0 {
        ts.emphasisKeyword = emphasisKeyword
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "关键词不能为空", nil))
    }
}

func (ts *templateSend) checkData() {
    if len(ts.openid) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "用户openid不能为空", nil))
    }
    if len(ts.templateId) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "模板ID不能为空", nil))
    }
    if len(ts.formId) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "表单ID不能为空", nil))
    }
}

func (ts *templateSend) SendRequest(getType string) api.APIResult {
    ts.checkData()
    reqData := make(map[string]interface{})
    reqData["touser"] = ts.openid
    reqData["template_id"] = ts.templateId
    reqData["page"] = ts.redirectUrl
    reqData["form_id"] = ts.formId
    reqData["data"] = ts.templateData
    if len(ts.emphasisKeyword) > 0 {
        reqData["emphasis_keyword"] = ts.emphasisKeyword
    }
    reqBody := mpf.JSONMarshal(reqData)

    ts.ReqURI = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=" + wx.NewUtilWx().GetSingleCache(ts.appId, getType)
    client, req := ts.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ts.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTemplateSend(appId string) *templateSend {
    ts := &templateSend{wx.NewBaseWxMini(), "", "", "", make(map[string]interface{}), "", "", ""}
    ts.appId = appId
    ts.ReqContentType = project.HTTPContentTypeJSON
    ts.ReqMethod = fasthttp.MethodPost
    return ts
}
