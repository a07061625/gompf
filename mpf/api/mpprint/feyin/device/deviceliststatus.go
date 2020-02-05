/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/28 0028
 * Time: 22:25
 */
package device

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/mpprint"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type deviceListStatus struct {
    mpprint.BaseFeYin
}

func (dls *deviceListStatus) checkData() (*fasthttp.Client, *fasthttp.Request) {
    dls.ReqUrl = mpprint.FeYinServiceDomain + "/devices?access_token=" + mpprint.NewUtilPrint().GetFeYinAccessToken(dls.GetAppId())
    return dls.GetRequest()
}

func (dls *deviceListStatus) SendRequest() api.ApiResult {
    client, req := dls.checkData()
    resp, result := dls.SendInner(client, req, errorcode.PrintFeYinRequestGet)
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

func NewDeviceListStatus(appId string) *deviceListStatus {
    dls := &deviceListStatus{mpprint.NewBaseFeYin()}
    dls.SetAppId(appId)
    dls.ReqMethod = fasthttp.MethodGet
    dls.ReqContentType = project.HttpContentTypeForm
    return dls
}
