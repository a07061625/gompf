/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 23:15
 */
package tool

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type callbackIpGet struct {
    wx.BaseWxAccount
    appId string
}

func (cig *callbackIpGet) SendRequest() api.ApiResult {
    cig.ReqUrl = "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=" + wx.NewUtilWx().GetSingleAccessToken(cig.appId)
    client, req := cig.GetRequest()

    resp, result := cig.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["ip_list"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewCallbackIpGet(appId string) *callbackIpGet {
    cig := &callbackIpGet{wx.NewBaseWxAccount(), ""}
    cig.appId = appId
    return cig
}
