/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:18
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

type deviceUnbind struct {
    mpprint.BaseFeYin
    deviceNo string // 机器编号
}

func (du *deviceUnbind) SetDeviceNo(deviceNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, deviceNo)
    if match {
        du.deviceNo = deviceNo
        du.ReqUrl = mpprint.FeYinServiceDomain + "/device/" + deviceNo + "/unbind?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不合法", nil))
    }
}

func (du *deviceUnbind) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(du.deviceNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不能为空", nil))
    }

    du.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(du.GetAppId())
    client, req := du.GetRequest()
    req.SetBody([]byte("[]"))

    return client, req
}

func (du *deviceUnbind) SendRequest() api.ApiResult {
    client, req := du.checkData()
    resp, result := du.SendInner(client, req, errorcode.PrintFeYinRequestPost)
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

func NewDeviceUnbind(appId string) *deviceUnbind {
    db := &deviceUnbind{mpprint.NewBaseFeYin(), ""}
    db.SetAppId(appId)
    db.ReqMethod = fasthttp.MethodPost
    db.ReqContentType = project.HttpContentTypeJson
    return db
}
