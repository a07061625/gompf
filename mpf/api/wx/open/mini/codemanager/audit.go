/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 10:18
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 将第三方提交的代码包提交审核
type audit struct {
    wx.BaseWxOpen
    appId     string                   // 应用ID
    auditList []map[string]interface{} // 审核列表
}

func (ca *audit) SetAuditList(auditList []map[string]interface{}) {
    if len(auditList) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核列表不能为空", nil))
    } else if len(auditList) > 5 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核列表数量不能超过5个", nil))
    }
    ca.auditList = auditList
}

func (ca *audit) checkData() {
    if len(ca.auditList) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核列表不能为空", nil))
    }
}

func (ca *audit) SendRequest() api.ApiResult {
    ca.checkData()

    reqData := make(map[string]interface{})
    reqData["item_list"] = ca.auditList
    reqBody := mpf.JsonMarshal(ca.ReqData)
    ca.ReqUrl = "https://api.weixin.qq.com/wxa/submit_audit?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ca.appId)
    client, req := ca.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ca.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAudit(appId string) *audit {
    ca := &audit{wx.NewBaseWxOpen(), "", make([]map[string]interface{}, 0)}
    ca.appId = appId
    ca.ReqContentType = project.HttpContentTypeJson
    ca.ReqMethod = fasthttp.MethodPost
    return ca
}
