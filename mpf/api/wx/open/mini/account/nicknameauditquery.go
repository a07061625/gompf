/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 17:18
 */
package account

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

// 小程序改名审核状态查询
type nicknameAuditQuery struct {
    wx.BaseWxOpen
    appId   string // 应用ID
    auditId string // 审核id
}

func (naq *nicknameAuditQuery) SetAuditId(auditId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, auditId)
    if match {
        naq.auditId = auditId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核id不合法", nil))
    }
}

func (naq *nicknameAuditQuery) checkData() {
    if len(naq.auditId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "审核id不能为空", nil))
    }
    naq.ReqData["audit_id"] = naq.auditId
}

func (naq *nicknameAuditQuery) SendRequest() api.APIResult {
    naq.checkData()

    reqBody := mpf.JSONMarshal(naq.ReqData)
    naq.ReqURI = "https://api.weixin.qq.com/wxa/api_wxa_querynickname?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(naq.appId)
    client, req := naq.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := naq.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewNicknameAuditQuery(appId string) *nicknameAuditQuery {
    naq := &nicknameAuditQuery{wx.NewBaseWxOpen(), "", ""}
    naq.appId = appId
    naq.ReqContentType = project.HTTPContentTypeJSON
    naq.ReqMethod = fasthttp.MethodPost
    return naq
}
