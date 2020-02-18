/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 11:28
 */
package invoice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/corp"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 更新发票状态
type statusUpdate struct {
    wx.BaseWxCorp
    corpId          string
    agentTag        string
    cardId          string // 发票id
    encryptCode     string // 加密密码
    reimburseStatus string // 报销状态
}

func (su *statusUpdate) SetCardId(cardId string) {
    if len(cardId) > 0 {
        su.cardId = cardId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票id不合法", nil))
    }
}

func (su *statusUpdate) SetEncryptCode(encryptCode string) {
    if len(encryptCode) > 0 {
        su.encryptCode = encryptCode
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "加密密码不合法", nil))
    }
}

func (su *statusUpdate) SetReimburseStatus(reimburseStatus string) {
    _, ok := corp.InvoiceReimburseStatusList[reimburseStatus]
    if ok {
        su.reimburseStatus = reimburseStatus
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "报销状态不合法", nil))
    }
}

func (su *statusUpdate) checkData() {
    if len(su.cardId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票id不能为空", nil))
    }
    if len(su.encryptCode) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "加密密码不能为空", nil))
    }
    if len(su.reimburseStatus) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "报销状态不能为空", nil))
    }
    su.ReqData["card_id"] = su.cardId
    su.ReqData["encrypt_code"] = su.encryptCode
    su.ReqData["reimburse_status"] = su.reimburseStatus
}

func (su *statusUpdate) SendRequest(getType string) api.ApiResult {
    su.checkData()
    reqBody := mpf.JsonMarshal(su.ReqData)

    su.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/card/invoice/reimburse/updateinvoicestatus?access_token=" + wx.NewUtilWx().GetCorpCache(su.corpId, su.agentTag, getType)
    client, req := su.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := su.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewStatusUpdate(corpId, agentTag string) *statusUpdate {
    su := &statusUpdate{wx.NewBaseWxCorp(), "", "", "", "", ""}
    su.corpId = corpId
    su.agentTag = agentTag
    su.ReqContentType = project.HTTPContentTypeJSON
    su.ReqMethod = fasthttp.MethodPost
    return su
}
