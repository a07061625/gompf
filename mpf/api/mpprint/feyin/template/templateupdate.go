/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 12:53
 */
package template

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/mpprint"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type templateUpdate struct {
    mpprint.BaseFeYin
    templateId string // 模板id
    name       string // 模板名称
    content    string // 模板内容
    catalog    string // 模板归类
    desc       string // 模板说明
}

func (tu *templateUpdate) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, templateId)
    if match {
        tu.templateId = templateId
        tu.ReqUrl = mpprint.FeYinServiceDomain + "/template/" + templateId + "?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板id不合法", nil))
    }
}

func (tu *templateUpdate) SetName(name string) {
    if (len(name) > 0) && (len(name) <= 30) {
        tu.ReqData["name"] = name
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板名称不合法", nil))
    }
}

func (tu *templateUpdate) SetContent(content string) {
    if len(content) > 0 {
        tu.ReqData["content"] = content
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板内容不合法", nil))
    }
}

func (tu *templateUpdate) SetCatalog(catalog string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, catalog)
    if match {
        tu.ReqData["catalog"] = catalog
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板归类不合法", nil))
    }
}

func (tu *templateUpdate) SetDesc(desc string) {
    if len(desc) > 0 {
        tu.ReqData["desc"] = desc
    }
}

func (tu *templateUpdate) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tu.templateId) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板id不能为空", nil))
    }
    if len(tu.ReqData) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板数据不能为空", nil))
    }

    tu.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(tu.GetAppId())
    client, req := tu.GetRequest()
    reqBody := mpf.JsonMarshal(tu.ReqData)
    req.SetBody([]byte(reqBody))

    return client, req
}

func (tu *templateUpdate) SendRequest() api.ApiResult {
    client, req := tu.checkData()
    resp, result := tu.SendInner(client, req, errorcode.PrintFeYinRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTemplateUpdate(appId string) *templateUpdate {
    tc := &templateUpdate{mpprint.NewBaseFeYin(), "", "", "", "", ""}
    tc.SetAppId(appId)
    tc.ReqMethod = fasthttp.MethodPost
    tc.ReqContentType = project.HttpContentTypeJson
    return tc
}
