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

// 获取外部联系人列表
type externalContactListGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    userId   string // 成员用户ID
}

func (ecg *externalContactListGet) SetUserId(userId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, userId)
    if match {
        ecg.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员用户ID不合法", nil))
    }
}

func (ecg *externalContactListGet) checkData() {
    if len(ecg.userId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "成员用户ID不能为空", nil))
    }
    ecg.ReqData["userid"] = ecg.userId
}

func (ecg *externalContactListGet) SendRequest(getType string) api.ApiResult {
    ecg.checkData()

    ecg.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(ecg.corpId, ecg.agentTag, getType)
    ecg.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/crm/get_external_contact_list?" + mpf.HTTPCreateParams(ecg.ReqData, "none", 1)
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

func NewExternalContactListGet(corpId, agentTag string) *externalContactListGet {
    ecg := &externalContactListGet{wx.NewBaseWxCorp(), "", "", ""}
    ecg.corpId = corpId
    ecg.agentTag = agentTag
    return ecg
}
