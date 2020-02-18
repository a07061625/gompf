/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:42
 */
package msg

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

type msgStatus struct {
    mpprint.BaseFeYin
    msgNo string // 消息ID
}

func (ms *msgStatus) SetMsgNo(msgNo string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, msgNo)
    if match {
        ms.msgNo = msgNo
        ms.ReqUrl = mpprint.FeYinServiceDomain + "/msg/" + msgNo + "/status?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息ID不合法", nil))
    }
}

func (ms *msgStatus) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ms.msgNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息ID不能为空", nil))
    }

    ms.ReqUrl += mpprint.NewUtilPrint().GetFeYinAccessToken(ms.GetAppId())
    return ms.GetRequest()
}

func (ms *msgStatus) SendRequest() api.ApiResult {
    client, req := ms.checkData()
    resp, result := ms.SendInner(client, req, errorcode.PrintFeYinRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["msg_no"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMsgStatus(appId string) *msgStatus {
    ms := &msgStatus{mpprint.NewBaseFeYin(), ""}
    ms.SetAppId(appId)
    ms.ReqMethod = fasthttp.MethodGet
    ms.ReqContentType = project.HTTPContentTypeForm
    return ms
}
