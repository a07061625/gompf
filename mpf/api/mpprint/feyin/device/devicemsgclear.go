/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 10:51
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

type deviceMsgClear struct {
    mpprint.BaseFeYin
    deviceNo string // 机器编号
}

func (dmc *deviceMsgClear) SetDeviceNo(deviceNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, deviceNo)
    if match {
        dmc.deviceNo = deviceNo
        dmc.ReqUrl = mpprint.FeYinServiceDomain + "/device/" + deviceNo + "/msg/clear?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不合法", nil))
    }
}

func (dmc *deviceMsgClear) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dmc.deviceNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不能为空", nil))
    }

    dmc.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(dmc.GetAppId())
    client, req := dmc.GetRequest()
    req.SetBody([]byte("[]"))

    return client, req
}

func (dmc *deviceMsgClear) SendRequest() api.ApiResult {
    client, req := dmc.checkData()
    resp, result := dmc.SendInner(client, req, errorcode.PrintFeYinRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["clear_cnt"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewDeviceMsgClear(appId string) *deviceMsgClear {
    db := &deviceMsgClear{mpprint.NewBaseFeYin(), ""}
    db.SetAppId(appId)
    db.ReqMethod = fasthttp.MethodPost
    db.ReqContentType = project.HttpContentTypeJson
    return db
}
