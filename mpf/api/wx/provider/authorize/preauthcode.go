/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 20:05
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取预授权码
type preAuthCode struct {
    wx.BaseWxProvider
}

func (pac *preAuthCode) SendRequest() api.ApiResult {
    pac.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/service/get_pre_auth_code?suite_access_token=" + wx.NewUtilWx().GetProviderSuiteToken()
    client, req := pac.GetRequest()

    resp, result := pac.SendInner(client, req, errorcode.WxProviderRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["pre_auth_code"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxProviderRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewPreAuthCode() *preAuthCode {
    return &preAuthCode{wx.NewBaseWxProvider()}
}
