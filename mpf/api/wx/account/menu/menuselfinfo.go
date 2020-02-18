/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 12:20
 */
package menu

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 获取自定义菜单配置
type menuSelfInfo struct {
    wx.BaseWxAccount
    appId string
}

func (msi *menuSelfInfo) SendRequest() api.APIResult {
    msi.ReqURI = "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=" + wx.NewUtilWx().GetSingleAccessToken(msi.appId)
    client, req := msi.GetRequest()

    resp, result := msi.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["selfmenu_info"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewMenuSelfInfo(appId string) *menuSelfInfo {
    msi := &menuSelfInfo{wx.NewBaseWxAccount(), ""}
    msi.appId = appId
    return msi
}
