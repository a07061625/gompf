/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 15:07
 */
package codetemplate

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取代码模版库中的所有小程序代码模版
type codeList struct {
    wx.BaseWxOpen
}

func (cl *codeList) SendRequest() api.ApiResult {
    cl.ReqUrl = "https://api.weixin.qq.com/wxa/gettemplatelist?access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := cl.GetRequest()

    resp, result := cl.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewCodeList() *codeList {
    return &codeList{wx.NewBaseWxOpen()}
}
