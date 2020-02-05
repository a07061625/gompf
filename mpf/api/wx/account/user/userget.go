/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 9:09
 */
package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type userGet struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
}

func (ug *userGet) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ug.ReqData["next_openid"] = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ug *userGet) SendRequest() api.ApiResult {
    ug.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(ug.appId)
    ug.ReqUrl = "https://api.weixin.qq.com/cgi-bin/user/get?" + mpf.HttpCreateParams(ug.ReqData, "none", 1)
    client, req := ug.GetRequest()

    resp, result := ug.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewUserGet(appId string) *userGet {
    ug := &userGet{wx.NewBaseWxAccount(), "", ""}
    ug.appId = appId
    ug.ReqData["next_openid"] = ""
    return ug
}
