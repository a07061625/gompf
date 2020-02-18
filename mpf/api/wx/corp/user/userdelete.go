/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/9 0009
 * Time: 16:54
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

// 删除成员
type userDelete struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    userId   string // 用户ID
}

func (ud *userDelete) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{1,32}$`, userId)
    if match {
        ud.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不合法", nil))
    }
}

func (ud *userDelete) checkData() {
    if len(ud.userId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不能为空", nil))
    }
    ud.ReqData["userid"] = ud.userId
}

func (ud *userDelete) SendRequest(getType string) api.APIResult {
    ud.checkData()

    ud.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(ud.corpId, ud.agentTag, getType)
    ud.ReqURI = "https://qyapi.weixin.qq.com/cgi-bin/user/delete?" + mpf.HTTPCreateParams(ud.ReqData, "none", 1)
    client, req := ud.GetRequest()

    resp, result := ud.SendInner(client, req, errorcode.WxCorpRequestGet)
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

func NewUserDelete(corpId, agentTag string) *userDelete {
    ud := &userDelete{wx.NewBaseWxCorp(), "", "", ""}
    ud.corpId = corpId
    ud.agentTag = agentTag
    return ud
}
