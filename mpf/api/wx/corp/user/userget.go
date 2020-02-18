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

// 读取成员
type userGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    userId   string // 用户ID
}

func (ug *userGet) SetUserId(userId string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{1,32}$`, userId)
    if match {
        ug.userId = userId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不合法", nil))
    }
}

func (ug *userGet) checkData() {
    if len(ug.userId) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户ID不能为空", nil))
    }
    ug.ReqData["userid"] = ug.userId
}

func (ug *userGet) SendRequest(getType string) api.ApiResult {
    ug.checkData()

    ug.ReqData["access_token"] = wx.NewUtilWx().GetCorpCache(ug.corpId, ug.agentTag, getType)
    ug.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/get?" + mpf.HTTPCreateParams(ug.ReqData, "none", 1)
    client, req := ug.GetRequest()

    resp, result := ug.SendInner(client, req, errorcode.WxCorpRequestGet)
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

func NewUserGet(corpId, agentTag string) *userGet {
    ug := &userGet{wx.NewBaseWxCorp(), "", "", ""}
    ug.corpId = corpId
    ug.agentTag = agentTag
    return ug
}
