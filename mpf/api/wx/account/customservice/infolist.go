/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 22:50
 */
package customservice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type infoList struct {
    wx.BaseWxAccount
    appId string
}

func (il *infoList) SendRequest() api.APIResult {
    il.ReqURI = "https://api.weixin.qq.com/cgi-bin/customservice/getkflist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(il.appId)
    client, req := il.GetRequest()

    resp, result := il.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["kf_list"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewInfoList(appId string) *infoList {
    il := &infoList{wx.NewBaseWxAccount(), ""}
    il.appId = appId
    return il
}
