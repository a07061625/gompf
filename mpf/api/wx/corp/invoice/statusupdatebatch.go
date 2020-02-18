/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/7 0007
 * Time: 12:37
 */
package invoice

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/corp"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 批量更新发票状态
type statusUpdateBatch struct {
    wx.BaseWxCorp
    corpId          string
    agentTag        string
    openid          string              // 用户openid
    reimburseStatus string              // 报销状态
    invoiceList     []map[string]string // 发票列表
}

func (sub *statusUpdateBatch) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        sub.openid = openid
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不合法", nil))
    }
}

func (sub *statusUpdateBatch) SetReimburseStatus(reimburseStatus string) {
    _, ok := corp.InvoiceReimburseStatusList[reimburseStatus]
    if ok {
        sub.reimburseStatus = reimburseStatus
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "报销状态不合法", nil))
    }
}

func (sub *statusUpdateBatch) SetInvoiceList(invoiceList []map[string]string) {
    if len(invoiceList) > 0 {
        sub.invoiceList = invoiceList
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票列表不合法", nil))
    }
}

func (sub *statusUpdateBatch) checkData() {
    if len(sub.openid) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不能为空", nil))
    }
    if len(sub.reimburseStatus) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "报销状态不能为空", nil))
    }
    if len(sub.invoiceList) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发票列表不能为空", nil))
    }
}

func (sub *statusUpdateBatch) SendRequest(getType string) api.APIResult {
    sub.checkData()
    reqData := make(map[string]interface{})
    reqData["openid"] = sub.openid
    reqData["reimburse_status"] = sub.reimburseStatus
    reqData["invoice_list"] = sub.invoiceList
    reqBody := mpf.JSONMarshal(reqData)

    sub.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/card/invoice/reimburse/updatestatusbatch?access_token=" + wx.NewUtilWx().GetCorpCache(sub.corpId, sub.agentTag, getType)
    client, req := sub.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sub.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewStatusUpdateBatch(corpId, agentTag string) *statusUpdateBatch {
    sub := &statusUpdateBatch{wx.NewBaseWxCorp(), "", "", "", "", make([]map[string]string, 0)}
    sub.corpId = corpId
    sub.agentTag = agentTag
    sub.ReqContentType = project.HTTPContentTypeJSON
    sub.ReqMethod = fasthttp.MethodPost
    return sub
}
