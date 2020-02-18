/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 9:07
 */
package crm

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 离职成员的外部联系人再分配
type externalContactTransfer struct {
    wx.BaseWxCorp
    corpId         string
    agentTag       string
    externalUserId string // 外部联系人用户ID
    handoverUserId string // 离职成员用户ID
    takeoverUserId string // 接替成员用户ID
}

func (ect *externalContactTransfer) SetExternalUserId(externalUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, externalUserId)
    if match {
        ect.externalUserId = externalUserId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "外部联系人用户ID不合法", nil))
    }
}

func (ect *externalContactTransfer) SetHandoverUserId(handoverUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, handoverUserId)
    if match {
        ect.handoverUserId = handoverUserId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "离职成员用户ID不合法", nil))
    }
}

func (ect *externalContactTransfer) SetTakeoverUserId(takeoverUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, takeoverUserId)
    if match {
        ect.takeoverUserId = takeoverUserId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "接替成员用户ID不合法", nil))
    }
}

func (ect *externalContactTransfer) checkData() {
    if len(ect.externalUserId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "外部联系人用户ID不能为空", nil))
    }
    if len(ect.externalUserId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "外部联系人用户ID不能为空", nil))
    }
    if len(ect.externalUserId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "接替成员用户ID不能为空", nil))
    }
    ect.ReqData["external_userid"] = ect.externalUserId
    ect.ReqData["handover_userid"] = ect.handoverUserId
    ect.ReqData["takeover_userid"] = ect.takeoverUserId
}

func (ect *externalContactTransfer) SendRequest(getType string) api.ApiResult {
    ect.checkData()
    reqBody := mpf.JsonMarshal(ect.ReqData)

    ect.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/crm/transfer_external_contact?access_token=" + wx.NewUtilWx().GetCorpCache(ect.corpId, ect.agentTag, getType)
    client, req := ect.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ect.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewExternalContactTransfer(corpId, agentTag string) *externalContactTransfer {
    ect := &externalContactTransfer{wx.NewBaseWxCorp(), "", "", "", "", ""}
    ect.corpId = corpId
    ect.agentTag = agentTag
    ect.ReqContentType = project.HTTPContentTypeJSON
    ect.ReqMethod = fasthttp.MethodPost
    return ect
}
