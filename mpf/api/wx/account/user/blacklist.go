/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 9:19
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
    "github.com/valyala/fasthttp"
)

// 获取黑名单列表
type blackList struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
}

func (bl *blackList) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        bl.ReqData["begin_openid"] = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (bl *blackList) SendRequest() api.APIResult {
    reqBody := mpf.JSONMarshal(bl.ReqData)
    bl.ReqURI = "https://api.weixin.qq.com/cgi-bin/tags/members/getblacklist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(bl.appId)
    client, req := bl.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := bl.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewBlackList(appId string) *blackList {
    bl := &blackList{wx.NewBaseWxAccount(), "", ""}
    bl.appId = appId
    bl.ReqData["begin_openid"] = ""
    bl.ReqContentType = project.HTTPContentTypeJSON
    bl.ReqMethod = fasthttp.MethodPost
    return bl
}
