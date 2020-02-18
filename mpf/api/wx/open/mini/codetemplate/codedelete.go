/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 15:03
 */
package codetemplate

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 删除指定小程序代码模版
type codeDelete struct {
    wx.BaseWxOpen
    templateId string // 模板ID
}

func (cd *codeDelete) SetTemplateId(templateId string) {
    if len(templateId) > 0 {
        cd.templateId = templateId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "模板ID不合法", nil))
    }
}

func (cd *codeDelete) checkData() {
    if len(cd.templateId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "模板ID不能为空", nil))
    }
    cd.ReqData["template_id"] = cd.templateId
}

func (cd *codeDelete) SendRequest() api.ApiResult {
    cd.checkData()

    reqBody := mpf.JsonMarshal(cd.ReqData)
    cd.ReqUrl = "https://api.weixin.qq.com/wxa/deletetemplate?access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := cd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := cd.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JsonUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewCodeDelete() *codeDelete {
    cd := &codeDelete{wx.NewBaseWxOpen(), ""}
    cd.ReqContentType = project.HTTPContentTypeJSON
    cd.ReqMethod = fasthttp.MethodPost
    return cd
}
