/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/28 0028
 * Time: 21:50
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

type deviceBind struct {
    mpprint.BaseFeYin
    deviceNo string // 机器编号
}

func (db *deviceBind) SetDeviceNo(deviceNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, deviceNo)
    if match {
        db.deviceNo = deviceNo
        db.ReqUrl = mpprint.FeYinServiceDomain + "/device/" + db.deviceNo + "/bind?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不合法", nil))
    }
}

func (db *deviceBind) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(db.deviceNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不能为空", nil))
    }

    db.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(db.GetAppId())
    client, req := db.GetRequest()
    req.SetBody([]byte("[]"))

    return client, req
}

func (db *deviceBind) SendRequest() api.ApiResult {
    client, req := db.checkData()
    resp, result := db.SendInner(client, req, errorcode.PrintFeYinRequestPost)
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

func NewDeviceBind(appId string) *deviceBind {
    db := &deviceBind{mpprint.NewBaseFeYin(), ""}
    db.SetAppId(appId)
    db.ReqMethod = fasthttp.MethodPost
    db.ReqContentType = project.HTTPContentTypeJSON
    return db
}
