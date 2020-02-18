/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 9:27
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

type userGetByTag struct {
    wx.BaseWxAccount
    appId  string
    tagId  int    // 标签ID
    openid string // 用户openid
}

func (ugt *userGetByTag) SetTagId(tagId int) {
    if tagId > 0 {
        ugt.tagId = tagId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不合法", nil))
    }
}

func (ugt *userGetByTag) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ugt.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ugt *userGetByTag) checkData() {
    if ugt.tagId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不能为空", nil))
    }
}

func (ugt *userGetByTag) SendRequest() api.ApiResult {
    ugt.checkData()

    reqData := make(map[string]interface{})
    reqData["tagid"] = ugt.tagId
    reqData["next_openid"] = ugt.openid
    reqBody := mpf.JSONMarshal(reqData)
    ugt.ReqUrl = "https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ugt.appId)
    client, req := ugt.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ugt.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["errcode"]
    if ok {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    } else {
        result.Data = respData
    }

    return result
}

func NewUserGetByTag(appId string) *userGetByTag {
    ugt := &userGetByTag{wx.NewBaseWxAccount(), "", 0, ""}
    ugt.appId = appId
    ugt.ReqContentType = project.HTTPContentTypeJSON
    ugt.ReqMethod = fasthttp.MethodPost
    return ugt
}
