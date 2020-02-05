/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:24
 */
package member

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

type memberInfo struct {
    mpprint.BaseFeYin
    uid string // 商户id
}

func (mi *memberInfo) SetUid(uid string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, uid)
    if match {
        mi.uid = uid
        mi.ReqUrl = mpprint.FeYinServiceDomain + "/member/" + uid + "?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "商户id不合法", nil))
    }
}

func (mi *memberInfo) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mi.uid) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "商户id不能为空", nil))
    }

    mi.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(mi.GetAppId())

    return mi.GetRequest()
}

func (mi *memberInfo) SendRequest() api.ApiResult {
    client, req := mi.checkData()
    resp, result := mi.SendInner(client, req, errorcode.PrintFeYinRequestGet)
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

func NewMemberInfo(appId string) *memberInfo {
    mi := &memberInfo{mpprint.NewBaseFeYin(), ""}
    mi.SetAppId(appId)
    mi.ReqMethod = fasthttp.MethodGet
    mi.ReqContentType = project.HttpContentTypeForm
    return mi
}
