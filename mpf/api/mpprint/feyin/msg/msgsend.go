/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/29 0029
 * Time: 11:47
 */
package msg

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/mpprint"
    "github.com/a07061625/gompf/mpf/mpcache"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type msgSend struct {
    mpprint.BaseFeYin
    devices      []string          // 机器编号列表
    msgNo        string            // 消息ID
    msgContent   string            // 消息内容
    templateId   string            // 模板id
    templateData map[string]string // 模板数据
}

func (ms *msgSend) SetDevices(devices []string) {
    ms.devices = make([]string, 0)
    for _, v := range devices {
        match, _ := regexp.MatchString(project.RegexDigit, v)
        if match {
            ms.devices = append(ms.devices, v)
        }
    }
}

func (ms *msgSend) SetMsgNo(msgNo string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, msgNo)
    if match {
        ms.ReqData["msg_no"] = msgNo
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息ID不合法", nil))
    }
}

func (ms *msgSend) SetMsgContent(msgContent string) {
    if len(msgContent) > 0 {
        ms.ReqData["msg_content"] = msgContent
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息内容不合法", nil))
    }
}

func (ms *msgSend) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, templateId)
    if match {
        ms.ReqData["template_id"] = templateId
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板id不合法", nil))
    }
}

func (ms *msgSend) SetTemplateData(templateData map[string]string) {
    if len(templateData) > 0 {
        ms.templateData = templateData
    } else {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板数据不合法", nil))
    }
}

func (ms *msgSend) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ms.devices) == 0 {
        panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "机器编号不能为空", nil))
    }

    newData := make(map[string]interface{})
    for k, v := range ms.ReqData {
        newData[k] = v
    }
    newData["device_no"] = strings.Join(ms.devices, ",")

    _, ok := ms.ReqData["template_id"]
    if ok {
        if len(ms.templateData) == 0 {
            panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "模板数据不能为空", nil))
        }
        newData["template_data"] = ms.templateData
        delete(newData, "msg_content")
    } else {
        _, ok = ms.ReqData["msg_content"]
        if !ok {
            panic(mperr.NewPrintFeYin(errorcode.PrintFeYinParam, "消息内容不能为空", nil))
        }
    }
    reqBody := mpf.JSONMarshal(newData)

    ms.ReqUrl = mpprint.FeYinServiceDomain + "/msg?access_token=" + mpprint.NewUtilPrint().GetFeYinAccessToken(ms.GetAppId())

    client, req := ms.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func (ms *msgSend) SendRequest() api.ApiResult {
    client, req := ms.checkData()
    resp, result := ms.SendInner(client, req, errorcode.PrintFeYinRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["msg_no"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.PrintFeYinRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMsgSend(appId string) *msgSend {
    ms := &msgSend{mpprint.NewBaseFeYin(), make([]string, 0), "", "", "", make(map[string]string)}
    ms.SetAppId(appId)
    ms.ReqData["appid"] = appId
    ms.ReqData["msg_no"] = mpcache.NewUtilCache().CreateUniqueId()
    ms.ReqMethod = fasthttp.MethodPost
    ms.ReqContentType = project.HTTPContentTypeJSON
    return ms
}
