/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 11:58
 */
package codemanager

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

// 查询某个指定版本的审核状态
type auditStatus struct {
    wx.BaseWxOpen
    appId   string // 应用ID
    auditId string // 审核id
}

func (cas *auditStatus) SetAuditId(auditId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, auditId)
    if match {
        cas.auditId = auditId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核id不合法", nil))
    }
}

func (cas *auditStatus) checkData() {
    if len(cas.auditId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核id不能为空", nil))
    }
    cas.ReqData["audit_id"] = cas.auditId
}

func (cas *auditStatus) SendRequest() api.ApiResult {
    cas.checkData()

    reqBody := mpf.JsonMarshal(cas.ReqData)
    cas.ReqUrl = "https://api.weixin.qq.com/wxa/get_auditstatus?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(cas.appId)
    client, req := cas.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cas.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewCodeAuditStatus(appId string) *auditStatus {
    cas := &auditStatus{wx.NewBaseWxOpen(), "", ""}
    cas.appId = appId
    cas.ReqContentType = project.HttpContentTypeJson
    cas.ReqMethod = fasthttp.MethodPost
    return cas
}
