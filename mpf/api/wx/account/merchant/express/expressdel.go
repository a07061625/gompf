/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 18:02
 */
package express

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

type expressDel struct {
    wx.BaseWxAccount
    appId      string
    templateId string // 模板ID
}

func (ed *expressDel) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, templateId)
    if match {
        ed.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板ID不合法", nil))
    }
}

func (ed *expressDel) checkData() {
    if len(ed.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板ID不能为空", nil))
    }
    ed.ReqData["template_id"] = ed.templateId
}

func (ed *expressDel) SendRequest() api.ApiResult {
    ed.checkData()

    reqBody := mpf.JsonMarshal(ed.ReqData)
    ed.ReqUrl = "https://api.weixin.qq.com/merchant/express/del?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ed.appId)
    client, req := ed.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ed.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewExpressDel(appId string) *expressDel {
    ed := &expressDel{wx.NewBaseWxAccount(), "", ""}
    ed.appId = appId
    ed.ReqContentType = project.HttpContentTypeJson
    ed.ReqMethod = fasthttp.MethodPost
    return ed
}
