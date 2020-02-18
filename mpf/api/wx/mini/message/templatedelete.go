/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/6 0006
 * Time: 9:10
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

type templateDelete struct {
    wx.BaseWxMini
    appId      string // 应用ID
    templateId string // 模板ID
}

func (td *templateDelete) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        td.templateId = templateId
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "模板ID不合法", nil))
    }
}

func (td *templateDelete) checkData() {
    if len(td.templateId) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "模板ID不能为空", nil))
    }
    td.ReqData["template_id"] = td.templateId
}

func (td *templateDelete) SendRequest(getType string) api.ApiResult {
    td.checkData()
    reqBody := mpf.JsonMarshal(td.ReqData)

    td.ReqUrl = "https://api.weixin.qq.com/cgi-bin/wxopen/template/del?access_token=" + wx.NewUtilWx().GetSingleCache(td.appId, getType)
    client, req := td.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := td.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }
    return result
}

func NewTemplateDelete(appId string) *templateDelete {
    td := &templateDelete{wx.NewBaseWxMini(), "", ""}
    td.appId = appId
    td.ReqContentType = project.HTTPContentTypeJSON
    td.ReqMethod = fasthttp.MethodPost
    return td
}
