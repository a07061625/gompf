/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 8:49
 */
package user

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/valyala/fasthttp"
)

type tagGet struct {
    wx.BaseWxAccount
    appId string
}

func (tg *tagGet) SendRequest() api.ApiResult {
    tg.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/get?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tg.appId)
    client, req := tg.GetRequest()

    resp, result := tg.SendInner(client, req, errorcode.WxAccountRequestGet)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["tags"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestGet
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTagGet(appId string) *tagGet {
    tg := &tagGet{wx.NewBaseWxAccount(), ""}
    tg.appId = appId
    tg.ReqContentType = project.HTTPContentTypeJSON
    tg.ReqMethod = fasthttp.MethodPost
    return tg
}
