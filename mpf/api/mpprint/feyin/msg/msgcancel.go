/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:35
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

type msgCancel struct {
    mpprint.BaseFeYin
    msgNo string // 消息ID
}

func (mc *msgCancel) SetMsgNo(msgNo string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, msgNo)
    if match {
        mc.msgNo = msgNo
        mc.ReqURI = mpprint.FeYinServiceDomain + "/msg/" + msgNo + "/cancel?access_token="
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息ID不合法", nil))
    }
}

func (mc *msgCancel) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mc.msgNo) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息ID不能为空", nil))
    }

    mc.ReqURI += mpprint.NewUtilPrint().GetFeYinAccessToken(mc.GetAppId())
    client, req := mc.GetRequest()
    req.SetBody([]byte("[]"))

    return client, req
}

func (mc *msgCancel) SendRequest() api.APIResult {
    client, req := mc.checkData()
    resp, result := mc.SendInner(client, req, errorcode.PrintFeYinRequestPost)
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

func NewMsgCancel(appId string) *msgCancel {
    mc := &msgCancel{mpprint.NewBaseFeYin(), ""}
    mc.SetAppId(appId)
    mc.ReqMethod = fasthttp.MethodPost
    mc.ReqContentType = project.HTTPContentTypeJSON
    return mc
}
