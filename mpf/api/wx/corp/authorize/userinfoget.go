/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 22:32
 */
package authorize

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
)

// 获取访问用户身份
type userInfoGet struct {
    wx.BaseWxCorp
    corpId   string
    agentTag string
    authCode string // 授权码
}

func (uig *userInfoGet) SetAuthCode(authCode string) {
    if len(authCode) > 0 {
        uig.authCode = authCode
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "授权码不合法", nil))
    }
}

func (uig *userInfoGet) checkData() {
    if len(uig.authCode) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "授权码不能为空", nil))
    }
    uig.ReqData["code"] = uig.authCode
}

func (uig *userInfoGet) SendRequest() api.ApiResult {
    uig.checkData()
    uig.ReqData["access_token"] = wx.NewUtilWx().GetCorpAccessToken(uig.corpId, uig.agentTag)
    uig.ReqUrl = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?" + mpf.HttpCreateParams(uig.ReqData, "none", 1)

    client, req := uig.GetRequest()
    resp, result := uig.SendInner(client, req, errorcode.WxCorpRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxCorpRequestGet
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewUserInfoGet(corpId, agentTag string) *userInfoGet {
    uig := &userInfoGet{wx.NewBaseWxCorp(), "", "", ""}
    uig.corpId = corpId
    uig.agentTag = agentTag
    return uig
}
