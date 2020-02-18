/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 22:42
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取展示的公众号信息
type showItemGet struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (sig *showItemGet) SendRequest() api.APIResult {
    sig.ReqURI = "https://api.weixin.qq.com/wxa/getshowwxaitem?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(sig.appId)
    client, req := sig.GetRequest()

    resp, result := sig.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewShowItemGet(appId string) *showItemGet {
    sig := &showItemGet{wx.NewBaseWxOpen(), ""}
    sig.appId = appId
    return sig
}
