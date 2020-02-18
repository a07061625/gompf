/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 23:21
 */
package customservice

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

type waitCaseGet struct {
    wx.BaseWxAccount
    appId string
}

func (wcg *waitCaseGet) SendRequest() api.ApiResult {
    wcg.ReqUrl = "https://api.weixin.qq.com/customservice/kfsession/getwaitcase?access_token=" + wx.NewUtilWx().GetSingleAccessToken(wcg.appId)
    client, req := wcg.GetRequest()

    resp, result := wcg.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["waitcaselist"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewWaitCaseGet(appId string) *waitCaseGet {
    wcg := &waitCaseGet{wx.NewBaseWxAccount(), ""}
    wcg.appId = appId
    return wcg
}
