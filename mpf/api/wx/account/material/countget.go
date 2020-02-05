/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 23:30
 */
package material

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取素材总数
type countGet struct {
    wx.BaseWxAccount
    appId string
}

func (mc *countGet) SendRequest() api.ApiResult {
    mc.ReqUrl = "https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mc.appId)
    client, req := mc.GetRequest()

    resp, result := mc.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewCountGet(appId string) *countGet {
    mc := &countGet{wx.NewBaseWxAccount(), ""}
    mc.appId = appId
    return mc
}
