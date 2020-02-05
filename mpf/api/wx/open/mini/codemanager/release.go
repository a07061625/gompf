/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 12:58
 */
package codemanager

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

// 发布已通过审核的小程序
type release struct {
    wx.BaseWxOpen
    appId string // 应用ID
}

func (r *release) SendRequest() api.ApiResult {
    r.ReqUrl = "https://api.weixin.qq.com/wxa/release?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(r.appId)
    client, req := r.GetRequest()
    req.SetBody([]byte("{}"))

    resp, result := r.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewRelease(appId string) *release {
    r := &release{wx.NewBaseWxOpen(), ""}
    r.appId = appId
    r.ReqContentType = project.HttpContentTypeJson
    r.ReqMethod = fasthttp.MethodPost
    return r
}
