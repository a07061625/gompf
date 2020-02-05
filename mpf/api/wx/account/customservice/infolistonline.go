/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 22:54
 */
package customservice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type infoListOnline struct {
    wx.BaseWxAccount
    appId string
}

func (ilo *infoListOnline) SendRequest() api.ApiResult {
    ilo.ReqUrl = "https://api.weixin.qq.com/cgi-bin/customservice/getonlinekflist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ilo.appId)
    client, req := ilo.GetRequest()

    resp, result := ilo.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["kf_online_list"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewInfoListOnline(appId string) *infoListOnline {
    ilo := &infoListOnline{wx.NewBaseWxAccount(), ""}
    ilo.appId = appId
    return ilo
}
