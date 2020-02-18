/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 17:27
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 微信认证名称检测
type nicknameCheck struct {
    wx.BaseWxOpen
    appId    string // 应用ID
    nickname string // 昵称
}

func (ac *nicknameCheck) SetNickname(nickname string) {
    if len(nickname) > 0 {
        ac.nickname = nickname
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "昵称不合法", nil))
    }
}

func (ac *nicknameCheck) checkData() {
    if len(ac.nickname) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "昵称不能为空", nil))
    }
    ac.ReqData["nick_name"] = ac.nickname
}

func (ac *nicknameCheck) SendRequest() api.ApiResult {
    ac.checkData()

    reqBody := mpf.JsonMarshal(ac.ReqData)
    ac.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxverify/checkwxverifynickname?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(ac.appId)
    client, req := ac.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ac.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewNicknameCheck(appId string) *nicknameCheck {
    ac := &nicknameCheck{wx.NewBaseWxOpen(), "", ""}
    ac.appId = appId
    ac.ReqContentType = project.HTTPContentTypeJSON
    ac.ReqMethod = fasthttp.MethodPost
    return ac
}
