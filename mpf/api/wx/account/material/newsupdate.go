/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 10:36
 */
package material

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改永久图文素材
type newsUpdate struct {
    wx.BaseWxAccount
    appId    string
    mediaId  string                 // 图文消息id
    index    int                    // 文章位置
    articles map[string]interface{} // 文章内容
}

func (nu *newsUpdate) SetMediaId(mediaId string) {
    if len(mediaId) > 0 {
        nu.mediaId = mediaId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图文消息id不合法", nil))
    }
}

func (nu *newsUpdate) SetIndex(index int) {
    if index >= 0 {
        nu.index = index
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文章位置不合法", nil))
    }
}

func (nu *newsUpdate) SetArticles(articles map[string]interface{}) {
    if len(articles) > 0 {
        nu.articles = articles
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文章内容不合法", nil))
    }
}

func (nu *newsUpdate) checkData() {
    if len(nu.mediaId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图文消息id不能为空", nil))
    }
    if nu.index < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文章位置不能为空", nil))
    }
    if len(nu.articles) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文章内容不能为空", nil))
    }
}

func (nu *newsUpdate) SendRequest() api.ApiResult {
    nu.checkData()

    reqData := make(map[string]interface{})
    reqData["media_id"] = nu.mediaId
    reqData["index"] = nu.index
    reqData["articles"] = nu.articles
    reqBody := mpf.JsonMarshal(reqData)
    nu.ReqUrl = "https://api.weixin.qq.com/cgi-bin/material/update_news?access_token=" + wx.NewUtilWx().GetSingleAccessToken(nu.appId)
    client, req := nu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := nu.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewNewsUpdate(appId string) *newsUpdate {
    nu := &newsUpdate{wx.NewBaseWxAccount(), "", "", 0, make(map[string]interface{})}
    nu.appId = appId
    nu.index = -1
    nu.ReqContentType = project.HttpContentTypeJson
    nu.ReqMethod = fasthttp.MethodPost
    return nu
}
