/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 14:37
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取小程序的第三方提交代码的页面配置
type pageGet struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (pg *pageGet) SendRequest() api.APIResult {
    pg.ReqURI = "https://api.weixin.qq.com/wxa/get_page?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(pg.appId)
    client, req := pg.GetRequest()

    resp, result := pg.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewPageGet(appId string) *pageGet {
    pg := &pageGet{wx.NewBaseWxOpen(), ""}
    pg.appId = appId
    return pg
}
