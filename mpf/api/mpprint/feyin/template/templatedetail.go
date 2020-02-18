/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 13:04
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

type templateDetail struct {
    mpprint.BaseFeYin
    templateId string // 模板id
}

func (td *templateDetail) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, templateId)
    if match {
        td.templateId = templateId
        td.ReqURI = mpprint.FeYinServiceDomain + "/template/detail/" + templateId + "?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板id不合法", nil))
    }
}

func (td *templateDetail) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(td.templateId) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板id不能为空", nil))
    }

    td.ReqURI += mpprint.NewUtilPrint().GetFeYinAccessToken(td.GetAppId())

    return td.GetRequest()
}

func (td *templateDetail) SendRequest() api.APIResult {
    client, req := td.checkData()
    resp, result := td.SendInner(client, req, errorcode.PrintFeYinRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.PrintFeYinRequestGet
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewTemplateDetail(appId string) *templateDetail {
    td := &templateDetail{mpprint.NewBaseFeYin(), ""}
    td.SetAppId(appId)
    td.ReqMethod = fasthttp.MethodGet
    td.ReqContentType = project.HTTPContentTypeForm
    return td
}
