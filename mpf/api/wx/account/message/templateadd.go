/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/12 0012
 * Time: 18:23
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

// 创建模板消息
type templateAdd struct {
    wx.BaseWxAccount
    appId           string
    templateShortId string // 模板编号
}

func (ta *templateAdd) SetTemplateShortId(templateShortId string) {
    if len(templateShortId) > 0 {
        ta.templateShortId = templateShortId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板编号不合法", nil))
    }
}

func (ta *templateAdd) checkData() {
    if len(ta.templateShortId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板编号不能为空", nil))
    }
    ta.ReqData["template_id_short"] = ta.templateShortId
}

func (ta *templateAdd) SendRequest() api.APIResult {
    ta.checkData()

    reqBody := mpf.JSONMarshal(ta.ReqData)
    ta.ReqURI = "https://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ta.appId)
    client, req := ta.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ta.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewTemplateAdd(appId string) *templateAdd {
    ta := &templateAdd{wx.NewBaseWxAccount(), "", ""}
    ta.appId = appId
    ta.ReqContentType = project.HTTPContentTypeJSON
    ta.ReqMethod = fasthttp.MethodPost
    return ta
}
