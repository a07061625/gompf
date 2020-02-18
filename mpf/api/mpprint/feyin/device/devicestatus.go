/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:12
 */
package device

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

type deviceStatus struct {
    mpprint.BaseFeYin
    deviceNo string // 机器编号
}

func (ds *deviceStatus) SetDeviceNo(deviceNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, deviceNo)
    if match {
        ds.deviceNo = deviceNo
        ds.ReqUrl = mpprint.FeYinServiceDomain + "/device/" + deviceNo + "/status?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不合法", nil))
    }
}

func (ds *deviceStatus) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ds.deviceNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不能为空", nil))
    }

    ds.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(ds.GetAppId())

    return ds.GetRequest()
}

func (ds *deviceStatus) SendRequest() api.ApiResult {
    client, req := ds.checkData()
    resp, result := ds.SendInner(client, req, errorcode.PrintFeYinRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["device_no"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewDeviceStatus(appId string) *deviceStatus {
    db := &deviceStatus{mpprint.NewBaseFeYin(), ""}
    db.SetAppId(appId)
    db.ReqMethod = fasthttp.MethodGet
    db.ReqContentType = project.HTTPContentTypeForm
    return db
}
