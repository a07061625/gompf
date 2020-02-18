/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 18:13
 */
package express

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type expressAll struct {
    wx.BaseWxAccount
    appId string
}

func (ea *expressAll) SendRequest() api.APIResult {
    ea.ReqURI = "https://api.weixin.qq.com/merchant/express/getall?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ea.appId)
    client, req := ea.GetRequest()

    resp, result := ea.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewExpressAll(appId string) *expressAll {
    ea := &expressAll{wx.NewBaseWxAccount(), ""}
    ea.appId = appId
    return ea
}
