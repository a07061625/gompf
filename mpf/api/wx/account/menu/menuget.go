/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 11:36
 */
package menu

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 查询菜单
type menuGet struct {
    wx.BaseWxAccount
    appId string
}

func (mg *menuGet) SendRequest() api.APIResult {
    mg.ReqURI = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=" + wx.NewUtilWx().GetSingleAccessToken(mg.appId)
    client, req := mg.GetRequest()

    resp, result := mg.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["menu"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMenuGet(appId string) *menuGet {
    mg := &menuGet{wx.NewBaseWxAccount(), ""}
    mg.appId = appId
    return mg
}
