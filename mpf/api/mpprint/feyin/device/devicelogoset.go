/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 10:33
 */
package device

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/mpprint"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type deviceLogoSet struct {
    mpprint.BaseFeYin
    deviceNo  string // 机器编号
    path      string // LOGO图片链接
    threshold int    // 图片灰度值
}

func (dls *deviceLogoSet) SetDeviceNo(deviceNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, deviceNo)
    if match {
        dls.deviceNo = deviceNo
        dls.ReqUrl = mpprint.FeYinServiceDomain + "/device/" + deviceNo + "/setting/logo?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不合法", nil))
    }
}

func (dls *deviceLogoSet) SetPath(path string) {
    match, _ := regexp.MatchString(project.RegexURLHTTP, path)
    if match {
        dls.ReqData["path"] = path
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "LOGO图片链接不合法", nil))
    }
}

func (dls *deviceLogoSet) SetThreshold(threshold int) {
    if (threshold > 0) && (threshold <= 255) {
        dls.ReqData["threshold"] = strconv.Itoa(threshold)
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "图片灰度值不合法", nil))
    }
}

func (dls *deviceLogoSet) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(dls.deviceNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不能为空", nil))
    }
    _, ok := dls.ReqData["path"]
    if !ok {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "LOGO图片链接不能为空", nil))
    }

    dls.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(dls.GetAppId())
    client, req := dls.GetRequest()
    reqBody := mpf.JSONMarshal(dls.ReqData)
    req.SetBody([]byte(reqBody))

    return client, req
}

func (dls *deviceLogoSet) SendRequest() api.ApiResult {
    client, req := dls.checkData()
    resp, result := dls.SendInner(client, req, errorcode.PrintFeYinRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewDeviceLogoSet(appId string) *deviceLogoSet {
    db := &deviceLogoSet{mpprint.NewBaseFeYin(), "", "", 0}
    db.SetAppId(appId)
    db.ReqData["threshold"] = "200"
    db.ReqMethod = fasthttp.MethodPost
    db.ReqContentType = project.HTTPContentTypeJSON
    return db
}
