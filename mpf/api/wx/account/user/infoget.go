/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 0:17
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

type infoGet struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
}

func (ig *infoGet) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ig.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ig *infoGet) checkData() {
    if len(ig.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    ig.ReqData["openid"] = ig.openid
}

func (ig *infoGet) SendRequest() api.ApiResult {
    ig.checkData()

    ig.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(ig.appId)
    ig.ReqUrl = "https://api.weixin.qq.com/cgi-bin/user/info?" + mpf.HttpCreateParams(ig.ReqData, "none", 1)
    client, req := ig.GetRequest()

    resp, result := ig.SendInner(client, req, errorcode.WxAccountRequestGet)
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

func NewInfoGet(appId string) *infoGet {
    ig := &infoGet{wx.NewBaseWxAccount(), "", ""}
    ig.appId = appId
    ig.ReqData["lang"] = "zh_CN"
    return ig
}
