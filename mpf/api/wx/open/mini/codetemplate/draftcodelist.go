/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 14:57
 */
package codetemplate

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取草稿箱内的所有临时代码草稿
type draftCodeList struct {
    wx.BaseWxOpen
}

func (dcl *draftCodeList) SendRequest() api.ApiResult {
    dcl.ReqUrl = "https://api.weixin.qq.com/wxa/gettemplatedraftlist?access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := dcl.GetRequest()

    resp, result := dcl.SendInner(client, req, errorcode.WxOpenRequestGet)
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

func NewDraftCodeList() *draftCodeList {
    return &draftCodeList{wx.NewBaseWxOpen()}
}
