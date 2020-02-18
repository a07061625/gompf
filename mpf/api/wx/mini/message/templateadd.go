/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 8:56
 */
package message

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type templateAdd struct {
    wx.BaseWxMini
    appId      string // 应用ID
    titleId    string // 标题ID
    keywordIds []uint // 关键词ID列表
}

func (ta *templateAdd) SetTitleId(titleId string) {
    if len(titleId) > 0 {
        ta.titleId = titleId
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "标题ID不合法", nil))
    }
}

func (ta *templateAdd) SetKeywordIds(keywordIds []uint) {
    if len(keywordIds) > 10 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "关键词ID不能超过10个", nil))
    }

    ta.keywordIds = keywordIds
}

func (ta *templateAdd) checkData() {
    if len(ta.titleId) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "标题ID不能为空", nil))
    }
}

func (ta *templateAdd) SendRequest(getType string) api.ApiResult {
    ta.checkData()
    reqData := make(map[string]interface{})
    reqData["id"] = ta.titleId
    reqData["keyword_id_list"] = ta.keywordIds
    reqBody := mpf.JSONMarshal(reqData)

    ta.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/template/add?access_token=" + wx.NewUtilWx().GetSingleCache(ta.appId, getType)
    client, req := ta.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ta.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    _, ok := respData["template_id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTemplateAdd(appId string) *templateAdd {
    ta := &templateAdd{wx.NewBaseWxMini(), "", "", make([]uint, 0)}
    ta.appId = appId
    ta.ReqContentType = project.HTTPContentTypeJSON
    ta.ReqMethod = fasthttp.MethodPost
    return ta
}
