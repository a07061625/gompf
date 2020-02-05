/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 19:10
 */
package authorize

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 通过openid获取用户详细信息
type userInfo struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
}

func (ui *userInfo) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ui.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ui *userInfo) checkData() {
    if len(ui.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    ui.ReqData["openid"] = ui.openid
}

func (ui *userInfo) SendRequest() api.ApiResult {
    ui.checkData()

    ui.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(ui.appId)
    ui.ReqUrl = "https://api.weixin.qq.com/cgi-bin/user/info?" + mpf.HttpCreateParams(ui.ReqData, "none", 1)
    client, req := ui.GetRequest()

    resp, result := ui.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["openid"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewUserInfo(appId string) *userInfo {
    ui := &userInfo{wx.NewBaseWxAccount(), "", ""}
    ui.appId = appId
    ui.ReqData["lang"] = "zh_CN"
    return ui
}
