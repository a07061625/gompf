/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 8:58
 */
package user

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type tagListGet struct {
    wx.BaseWxAccount
    appId  string
    openid string // 用户openid
}

func (tlg *tagListGet) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        tlg.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (tlg *tagListGet) checkData() {
    if len(tlg.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    tlg.ReqData["openid"] = tlg.openid
}

func (tlg *tagListGet) SendRequest() api.ApiResult {
    tlg.checkData()

    reqBody := mpf.JsonMarshal(tlg.ReqData)
    tlg.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tlg.appId)
    client, req := tlg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tlg.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewTagListGet(appId string) *tagListGet {
    tlg := &tagListGet{wx.NewBaseWxAccount(), "", ""}
    tlg.appId = appId
    tlg.ReqContentType = project.HTTPContentTypeJSON
    tlg.ReqMethod = fasthttp.MethodPost
    return tlg
}
