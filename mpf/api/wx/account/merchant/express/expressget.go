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

type expressGet struct {
    wx.BaseWxAccount
    appId      string
    templateId string // 模板ID
}

func (eg *expressGet) SetTemplateId(templateId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, templateId)
    if match {
        eg.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板ID不合法", nil))
    }
}

func (eg *expressGet) checkData() {
    if len(eg.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板ID不能为空", nil))
    }
    eg.ReqData["template_id"] = eg.templateId
}

func (eg *expressGet) SendRequest() api.ApiResult {
    eg.checkData()

    reqBody := mpf.JsonMarshal(eg.ReqData)
    eg.ReqUrl = "https://api.weixin.qq.com/merchant/express/getbyid?access_token=" + wx.NewUtilWx().GetSingleAccessToken(eg.appId)
    client, req := eg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := eg.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewExpressGet(appId string) *expressGet {
    eg := &expressGet{wx.NewBaseWxAccount(), "", ""}
    eg.appId = appId
    eg.ReqContentType = project.HttpContentTypeJson
    eg.ReqMethod = fasthttp.MethodPost
    return eg
}
