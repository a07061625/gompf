/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 8:41
 */
package register

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 第三方平台复用公众号主体快速注册小程序
type miniFast struct {
    wx.BaseWxOpen
    appId  string // 公众号APPID
    ticket string // 授权凭证
}

func (mf *miniFast) SetTicket(ticket string) {
    if len(ticket) > 0 {
        mf.ticket = ticket
    } else {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "授权凭证不合法", nil))
    }
}

func (mf *miniFast) checkData() {
    if len(mf.ticket) == 0 {
        panic(mperr.NewWxOpenAccount(errorcode.WxOpenParam, "授权凭证不能为空", nil))
    }
    mf.ReqData["ticket"] = mf.ticket
}

func (mf *miniFast) SendRequest() api.ApiResult {
    mf.checkData()

    reqBody := mpf.JSONMarshal(mf.ReqData)
    mf.ReqUrl = "https://api.weixin.qq.com/cgi-bin/account/fastregister?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(mf.appId)
    client, req := mf.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := mf.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        wx.NewUtilWx().RefreshOpenAuthorizeInfo(mf.appId, project.WxOpenAuthorizeOperateAuthorized, respData)
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewMiniFast(appId string) *miniFast {
    mf := &miniFast{wx.NewBaseWxOpen(), "", ""}
    mf.appId = appId
    mf.ReqContentType = project.HTTPContentTypeJSON
    mf.ReqMethod = fasthttp.MethodPost
    return mf
}
