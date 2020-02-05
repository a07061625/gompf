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

// 删除模板消息
type templateDel struct {
    wx.BaseWxAccount
    appId      string
    templateId string // 模板消息ID
}

func (ta *templateDel) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        ta.templateId = templateId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板消息ID不合法", nil))
    }
}

func (ta *templateDel) checkData() {
    if len(ta.templateId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "模板消息ID不能为空", nil))
    }
    ta.ReqData["template_id"] = ta.templateId
}

func (ta *templateDel) SendRequest() api.ApiResult {
    ta.checkData()

    reqBody := mpf.JsonMarshal(ta.ReqData)
    ta.ReqUrl = "https://api.weixin.qq.com/cgi-bin/template/del_private_template?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ta.appId)
    client, req := ta.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ta.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewTemplateDel(appId string) *templateDel {
    ta := &templateDel{wx.NewBaseWxAccount(), "", ""}
    ta.appId = appId
    ta.ReqContentType = project.HttpContentTypeJson
    ta.ReqMethod = fasthttp.MethodPost
    return ta
}
