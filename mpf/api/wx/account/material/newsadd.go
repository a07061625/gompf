/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 10:28
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

// 新增永久图文素材
type newsAdd struct {
    wx.BaseWxAccount
    appId    string
    articles []map[string]interface{} // 文章列表
}

func (na *newsAdd) SetArticles(articles []map[string]interface{}) {
    if len(articles) == 0 {
        na.articles = articles
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文章不能为空", nil))
    }
}

func (na *newsAdd) checkData() {
    if len(na.articles) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "文章不能为空", nil))
    }
}

func (na *newsAdd) SendRequest() api.APIResult {
    na.checkData()

    reqData := make(map[string]interface{})
    reqData["articles"] = na.articles
    reqBody := mpf.JSONMarshal(reqData)
    na.ReqURI = "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=" + wx.NewUtilWx().GetSingleAccessToken(na.appId)
    client, req := na.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := na.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewNewsAdd(appId string) *newsAdd {
    na := &newsAdd{wx.NewBaseWxAccount(), "", make([]map[string]interface{}, 0)}
    na.appId = appId
    na.ReqContentType = project.HTTPContentTypeJSON
    na.ReqMethod = fasthttp.MethodPost
    return na
}
