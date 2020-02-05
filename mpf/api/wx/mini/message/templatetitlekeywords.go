/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 9:49
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

type templateTitleKeywords struct {
    wx.BaseWxMini
    appId   string // 应用ID
    titleId string // 模板标题id
}

func (ttk *templateTitleKeywords) SetTitleId(titleId string) {
    if len(titleId) > 0 {
        ttk.titleId = titleId
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "模板标题id不合法", nil))
    }
}

func (ttk *templateTitleKeywords) checkData() {
    if len(ttk.titleId) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "模板标题id不能为空", nil))
    }
    ttk.ReqData["id"] = ttk.titleId
}

func (ttk *templateTitleKeywords) SendRequest(getType string) api.ApiResult {
    ttk.checkData()
    reqBody := mpf.JsonMarshal(ttk.ReqData)

    ttk.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/template/library/get?access_token=" + wx.NewUtilWx().GetSingleCache(ttk.appId, getType)
    client, req := ttk.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ttk.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    _, ok := respData["id"]
    if ok {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTemplateTitleKeywords(appId string) *templateTitleKeywords {
    ttk := &templateTitleKeywords{wx.NewBaseWxMini(), "", ""}
    ttk.appId = appId
    ttk.ReqContentType = project.HttpContentTypeJson
    ttk.ReqMethod = fasthttp.MethodPost
    return ttk
}
