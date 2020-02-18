/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 14:40
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改小程序线上代码的可见状态
type visitStatusChange struct {
    wx.BaseWxOpen
    appId       string // 应用ID
    visitStatus string // 访问状态
}

func (vsc *visitStatusChange) SetVisitStatus(visitStatus string) {
    if (visitStatus == "open") || (visitStatus == "close") {
        vsc.visitStatus = visitStatus
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "访问状态不支持", nil))
    }
}

func (vsc *visitStatusChange) checkData() {
    if len(vsc.visitStatus) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "访问状态不能为空", nil))
    }
    vsc.ReqData["action"] = vsc.visitStatus
}

func (vsc *visitStatusChange) SendRequest() api.ApiResult {
    vsc.checkData()

    reqBody := mpf.JSONMarshal(vsc.ReqData)
    vsc.ReqUrl = "https://api.weixin.qq.com/wxa/change_visitstatus?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(vsc.appId)
    client, req := vsc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := vsc.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewVisitStatusChange(appId string) *visitStatusChange {
    vsc := &visitStatusChange{wx.NewBaseWxOpen(), "", ""}
    vsc.appId = appId
    vsc.ReqContentType = project.HTTPContentTypeJSON
    vsc.ReqMethod = fasthttp.MethodPost
    return vsc
}
