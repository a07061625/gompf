/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 15:41
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取帐号基本信息
type basicInfoGet struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (big *basicInfoGet) SendRequest() api.ApiResult {
    big.ReqUrl = "https://api.weixin.qq.com/cgi-bin/account/getaccountbasicinfo?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(big.appId)
    client, req := big.GetRequest()

    resp, result := big.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewBasicInfoGet(appId string) *basicInfoGet {
    big := &basicInfoGet{wx.NewBaseWxOpen(), ""}
    big.appId = appId
    return big
}
