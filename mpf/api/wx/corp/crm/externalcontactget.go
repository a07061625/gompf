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
)

// 获取外部联系人详情
type externalContactGet struct {
    wx.BaseWxCorp
    corpId         string
    agentTag       string
    externalUserId string // 外部联系人用户ID
}

func (ecg *externalContactGet) SetExternalUserId(externalUserId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, externalUserId)
    if match {
        ecg.externalUserId = externalUserId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "外部联系人用户ID不合法", nil))
    }
}

func (ecg *externalContactGet) checkData() {
    if len(ecg.externalUserId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "外部联系人用户ID不能为空", nil))
    }
    ecg.ReqData["external_userid"] = ecg.externalUserId
}

func (ecg *externalContactGet) SendRequest(getType string) api.APIResult {
    ecg.checkData()

    ecg.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(ecg.corpId, ecg.agentTag, getType)
    ecg.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/crm/get_external_contact?" + mpf.HTTPCreateParams(ecg.ReqData, "none", 1)
    client, req := ecg.GetRequest()

    resp, result := ecg.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewExternalContactGet(corpId, agentTag string) *externalContactGet {
    ecg := &externalContactGet{wx.NewBaseWxCorp(), "", "", ""}
    ecg.corpId = corpId
    ecg.agentTag = agentTag
    return ecg
}
