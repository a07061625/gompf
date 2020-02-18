/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 17:41
 */
package shelf

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type shelfAll struct {
    wx.BaseWxAccount
    appId string
}

func (sa *shelfAll) SendRequest() api.ApiResult {
    sa.ReqUrl = "https://api.weixin.qq.com/merchant/shelf/getall?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sa.appId)
    client, req := sa.GetRequest()

    resp, result := sa.SendInner(client, req, errorcode.WxAccountRequestGet)
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

func NewShelfAll(appId string) *shelfAll {
    sa := &shelfAll{wx.NewBaseWxAccount(), ""}
    sa.appId = appId
    return sa
}
