/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 0:24
 */
package user

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type tagCreate struct {
    wx.BaseWxAccount
    appId   string
    tagName string // 标签名
}

func (tc *tagCreate) SetTagName(tagName string) {
    if (len(tagName) > 0) && (len(tagName) <= 30) {
        tc.tagName = tagName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签名不合法", nil))
    }
}

func (tc *tagCreate) checkData() {
    if len(tc.tagName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签名不能为空", nil))
    }
}

func (tc *tagCreate) SendRequest() api.ApiResult {
    tc.checkData()

    tagInfo := make(map[string]interface{})
    tagInfo["name"] = tc.tagName
    reqData := make(map[string]interface{})
    reqData["tag"] = tagInfo
    reqBody := mpf.JsonMarshal(reqData)
    tc.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/create?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tc.appId)
    client, req := tc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tc.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["tag"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTagCreate(appId string) *tagCreate {
    tc := &tagCreate{wx.NewBaseWxAccount(), "", ""}
    tc.appId = appId
    tc.ReqContentType = project.HttpContentTypeJson
    tc.ReqMethod = fasthttp.MethodPost
    return tc
}
