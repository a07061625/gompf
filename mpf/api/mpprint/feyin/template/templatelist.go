/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 13:10
 */
package template

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/mpprint"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type templateList struct {
    mpprint.BaseFeYin
}

func (tl *templateList) checkData() (*fasthttp.Client, *fasthttp.Request) {
    tl.ReqUrl = mpprint.FeYinServiceDomain + "/templates?access_token=" + mpprint.NewUtilPrint().GetFeYinAccessToken(tl.GetAppId())

    return tl.GetRequest()
}

func (tl *templateList) SendRequest() api.ApiResult {
    client, req := tl.checkData()
    resp, result := tl.SendInner(client, req, errorcode.PrintFeYinRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.PrintFeYinRequestGet
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewTemplateList(appId string) *templateList {
    tl := &templateList{mpprint.NewBaseFeYin()}
    tl.SetAppId(appId)
    tl.ReqMethod = fasthttp.MethodGet
    tl.ReqContentType = project.HttpContentTypeForm
    return tl
}
