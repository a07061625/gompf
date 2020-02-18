/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 12:06
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 小程序审核撤回
type auditUndo struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (au *auditUndo) SendRequest() api.ApiResult {
    au.ReqUrl = "https://api.weixin.qq.com/wxa/undocodeaudit?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(au.appId)
    client, req := au.GetRequest()

    resp, result := au.SendInner(client, req, errorcode.WxOpenRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAuditUndo(appId string) *auditUndo {
    au := &auditUndo{wx.NewBaseWxOpen(), ""}
    au.appId = appId
    au.ReqContentType = project.HTTPContentTypeJSON
    au.ReqMethod = fasthttp.MethodPost
    return au
}
