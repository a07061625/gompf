/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 11:16
 */
package media

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type newsUpload struct {
    wx.BaseWxAccount
    appId    string
    articles []map[string]interface{} // 图文消息列表
}

func (nu *newsUpload) SetArticles(articles []map[string]interface{}) {
    if len(articles) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图文消息不能为空", nil))
    } else if len(articles) > 8 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图文消息不能超过8个", nil))
    }
    nu.articles = articles
}

func (nu *newsUpload) checkData() {
    if len(nu.articles) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图文消息不能为空", nil))
    }
}

func (nu *newsUpload) SendRequest() api.APIResult {
    nu.checkData()

    reqData := make(map[string]interface{})
    reqData["articles"] = nu.articles
    reqBody := mpf.JSONMarshal(reqData)
    nu.ReqURI = "https://api.weixin.qq.com/cgi-bin/media/uploadnews?access_token=" + wx.NewUtilWx().GetSingleAccessToken(nu.appId)
    client, req := nu.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := nu.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["media_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewNewsUpload(appId string) *newsUpload {
    nu := &newsUpload{wx.NewBaseWxAccount(), "", make([]map[string]interface{}, 0)}
    nu.appId = appId
    nu.ReqContentType = project.HTTPContentTypeJSON
    nu.ReqMethod = fasthttp.MethodPost
    return nu
}
