/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 12:11
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

// 小程序分阶段发布
type grayRelease struct {
    wx.BaseWxOpen
    appId      string // 应用ID
    percentage int    // 灰度百分比
}

func (gr *grayRelease) SetPercentage(percentage int) {
    if (percentage > 0) && (percentage <= 100) {
        gr.percentage = percentage
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "灰度百分比不合法", nil))
    }
}

func (gr *grayRelease) checkData() {
    if gr.percentage <= 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "灰度百分比不能为空", nil))
    }
}

func (gr *grayRelease) SendRequest() api.ApiResult {
    gr.checkData()

    reqData := make(map[string]interface{})
    reqData["gray_percentage"] = gr.percentage
    reqBody := mpf.JSONMarshal(gr.ReqData)
    gr.ReqUrl = "https://api.weixin.qq.com/wxa/grayrelease?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(gr.appId)
    client, req := gr.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := gr.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewGrayRelease(appId string) *grayRelease {
    gr := &grayRelease{wx.NewBaseWxOpen(), "", 0}
    gr.appId = appId
    gr.percentage = 0
    gr.ReqContentType = project.HTTPContentTypeJSON
    gr.ReqMethod = fasthttp.MethodPost
    return gr
}
