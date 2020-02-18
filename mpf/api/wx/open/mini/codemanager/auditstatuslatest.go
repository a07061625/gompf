/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 12:01
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 查询最新一次提交的审核状态
type auditStatusLatest struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (asl *auditStatusLatest) SendRequest() api.ApiResult {
    asl.ReqUrl = "https://api.weixin.qq.com/wxa/get_latest_auditstatus?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(asl.appId)
    client, req := asl.GetRequest()

    resp, result := asl.SendInner(client, req, errorcode.WxOpenRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewAuditStatusLatest(appId string) *auditStatusLatest {
    asl := &auditStatusLatest{wx.NewBaseWxOpen(), ""}
    asl.appId = appId
    return asl
}
