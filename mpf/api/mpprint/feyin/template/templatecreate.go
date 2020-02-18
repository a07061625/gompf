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

type templateCreate struct {
    mpprint.BaseFeYin
    name    string // 模板名称
    content string // 模板内容
    catalog string // 模板归类
    desc    string // 模板说明
}

func (tc *templateCreate) SetName(name string) {
    if (len(name) > 0) && (len(name) <= 30) {
        tc.ReqData["name"] = name
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板名称不合法", nil))
    }
}

func (tc *templateCreate) SetContent(content string) {
    if len(content) > 0 {
        tc.ReqData["content"] = content
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板内容不合法", nil))
    }
}

func (tc *templateCreate) SetCatalog(catalog string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, catalog)
    if match {
        tc.ReqData["catalog"] = catalog
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板归类不合法", nil))
    }
}

func (tc *templateCreate) SetDesc(desc string) {
    if len(desc) > 0 {
        tc.ReqData["desc"] = desc
    }
}

func (tc *templateCreate) checkData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := tc.ReqData["name"]
    if !ok {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板名称不能为空", nil))
    }
    _, ok = tc.ReqData["content"]
    if !ok {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板内容不能为空", nil))
    }
    _, ok = tc.ReqData["catalog"]
    if !ok {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板归类不能为空", nil))
    }

    tc.ReqUrl = mpprint.FeYinServiceDomain + "/template?access_token=" + mpprint.NewUtilPrint().GetFeYinAccessToken(tc.GetAppId())
    client, req := tc.GetRequest()
    reqBody := mpf.JsonMarshal(tc.ReqData)
    req.SetBody([]byte(reqBody))

    return client, req
}

func (tc *templateCreate) SendRequest() api.ApiResult {
    client, req := tc.checkData()
    resp, result := tc.SendInner(client, req, errorcode.PrintFeYinRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["template_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTemplateCreate(appId string) *templateCreate {
    tc := &templateCreate{mpprint.NewBaseFeYin(), "", "", "", ""}
    tc.SetAppId(appId)
    tc.ReqMethod = fasthttp.MethodPost
    tc.ReqContentType = project.HTTPContentTypeJSON
    return tc
}
