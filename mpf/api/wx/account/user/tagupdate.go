/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 8:51
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

type tagUpdate struct {
    wx.BaseWxAccount
    appId   string
    tagId   int    // 标签ID
    tagName string // 标签名
}

func (tu *tagUpdate) SetTagId(tagId int) {
    if tagId > 0 {
        tu.tagId = tagId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不合法", nil))
    }
}

func (tu *tagUpdate) SetTagName(tagName string) {
    if (len(tagName) > 0) && (len(tagName) <= 30) {
        tu.tagName = tagName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签名不合法", nil))
    }
}

func (tu *tagUpdate) checkData() {
    if tu.tagId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不能为空", nil))
    }
    if len(tu.tagName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签名不能为空", nil))
    }
}

func (tu *tagUpdate) SendRequest() api.ApiResult {
    tu.checkData()

    tagInfo := make(map[string]interface{})
    tagInfo["id"] = tu.tagId
    tagInfo["name"] = tu.tagName
    reqData := make(map[string]interface{})
    reqData["tag"] = tagInfo
    reqBody := mpf.JsonMarshal(reqData)
    tu.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/update?access_token=" + wx.NewUtilWx().GetSingleAccessToken(tu.appId)
    client, req := tu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := tu.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTagUpdate(appId string) *tagUpdate {
    tu := &tagUpdate{wx.NewBaseWxAccount(), "", 0, ""}
    tu.appId = appId
    tu.ReqContentType = project.HTTPContentTypeJSON
    tu.ReqMethod = fasthttp.MethodPost
    return tu
}
