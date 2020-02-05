/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:30
 */
package member

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/mpprint"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type memberList struct {
    mpprint.BaseFeYin
}

func (ml *memberList) checkData() (*fasthttp.Client, *fasthttp.Request) {
    ml.ReqUrl = mpprint.FeYinServiceDomain + "/app/" + ml.GetAppId() + "/members?access_token=" + mpprint.NewUtilPrint().GetFeYinAccessToken(ml.GetAppId())

    return ml.GetRequest()
}

func (ml *memberList) SendRequest() api.ApiResult {
    client, req := ml.checkData()
    resp, result := ml.SendInner(client, req, errorcode.PrintFeYinRequestGet)
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

func NewMemberList(appId string) *memberList {
    ml := &memberList{mpprint.NewBaseFeYin()}
    ml.SetAppId(appId)
    ml.ReqMethod = fasthttp.MethodGet
    ml.ReqContentType = project.HttpContentTypeForm
    return ml
}
