/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 12:53
 */
package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 二次验证成员加入
type authSuccess struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    userId   string // 用户ID
}

func (as *authSuccess) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, userId)
    if match {
        as.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不合法", nil))
    }
}

func (as *authSuccess) checkData() {
    if len(as.userId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不能为空", nil))
    }
    as.ReqData["userid"] = as.userId
}

func (as *authSuccess) SendRequest(getType string) api.APIResult {
    as.checkData()

    as.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(as.corpId, as.agentTag, getType)
    as.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/user/authsucc?" + mpf.HTTPCreateParams(as.ReqData, "none", 1)
    client, req := as.GetRequest()

    resp, result := as.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewAuthSuccess(corpId, agentTag string) *authSuccess {
    as := &authSuccess{wx.NewBaseWxCorp(), "", "", ""}
    as.corpId = corpId
    as.agentTag = agentTag
    return as
}
