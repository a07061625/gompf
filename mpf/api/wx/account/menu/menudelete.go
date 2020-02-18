/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 12:18
 */
package menu

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
)

// 删除菜单
type menuDelete struct {
    wx.BaseWxAccount
    appId string
}

func (md *menuDelete) SendRequest() api.ApiResult {
    md.ReqUrl = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=" + wx.NewUtilWx().GetSingleAccessToken(md.appId)
    client, req := md.GetRequest()

    resp, result := md.SendInner(client, req, errorcode.WxAccountRequestGet)
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

func NewMenuDelete(appId string) *menuDelete {
    md := &menuDelete{wx.NewBaseWxAccount(), ""}
    md.appId = appId
    return md
}
