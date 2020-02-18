/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 13:00
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 小程序版本回退
type rollback struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (r *rollback) SendRequest() api.APIResult {
    r.ReqURI = "https://api.weixin.qq.com/wxa/revertcoderelease?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(r.appId)
    client, req := r.GetRequest()

    resp, result := r.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewRollback(appId string) *rollback {
    r := &rollback{wx.NewBaseWxOpen(), ""}
    r.appId = appId
    return r
}
