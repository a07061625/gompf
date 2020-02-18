/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:01
 */
package cloud

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 开通云开发权限的帐号
type userCreate struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (uc *userCreate) SendRequest() api.ApiResult {
    uc.ReqUrl = "https://api.weixin.qq.com/tcb/createclouduser?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(uc.appId)
    client, req := uc.GetRequest()
    req.SetBody([]byte("{}"))

    resp, result := uc.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewUserCreate(appId string) *userCreate {
    uc := &userCreate{wx.NewBaseWxOpen(), ""}
    uc.appId = appId
    uc.ReqContentType = project.HTTPContentTypeJSON
    uc.ReqMethod = fasthttp.MethodPost
    return uc
}
