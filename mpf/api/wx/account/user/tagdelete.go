/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 8:44
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

type tagDelete struct {
    wx.BaseWxAccount
    appId string
    tagId int // 标签ID
}

func (td *tagDelete) SetTagId(tagId int) {
    if tagId > 0 {
        td.tagId = tagId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不合法", nil))
    }
}

func (td *tagDelete) checkData() {
    if td.tagId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "标签ID不能为空", nil))
    }
}

func (td *tagDelete) SendRequest() api.ApiResult {
    td.checkData()

    tagInfo := make(map[string]interface{})
    tagInfo["id"] = td.tagId
    reqData := make(map[string]interface{})
    reqData["tag"] = tagInfo
    reqBody := mpf.JSONMarshal(reqData)
    td.ReqUrl = "https://api.weixin.qq.com/cgi-bin/tags/delete?access_token=" + wx.NewUtilWx().GetSingleAccessToken(td.appId)
    client, req := td.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := td.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewTagDelete(appId string) *tagDelete {
    td := &tagDelete{wx.NewBaseWxAccount(), "", 0}
    td.appId = appId
    td.ReqContentType = project.HTTPContentTypeJSON
    td.ReqMethod = fasthttp.MethodPost
    return td
}
