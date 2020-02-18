/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 15:00
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

// 将草稿箱的草稿选为小程序代码模版
type codeAdd struct {
    wx.BaseWxOpen
    draftId string // 草稿ID
}

func (ca *codeAdd) SetDraftId(draftId string) {
    if len(draftId) > 0 {
        ca.draftId = draftId
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "草稿ID不合法", nil))
    }
}

func (ca *codeAdd) checkData() {
    if len(ca.draftId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "草稿ID不能为空", nil))
    }
    ca.ReqData["draft_id"] = ca.draftId
}

func (ca *codeAdd) SendRequest() api.APIResult {
    ca.checkData()

    reqBody := mpf.JSONMarshal(ca.ReqData)
    ca.ReqURI = "https://api.weixin.qq.com/wxa/addtotemplate?access_token=" + wx.NewUtilWx().GetOpenAccessToken()
    client, req := ca.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ca.SendInner(client, req, errorcode.WxOpenRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, _ := mpf.JSONUnmarshalMap(resp.Content)
    errCode, ok := respData["errcode"]
    if ok && (errCode.(int) == 0) {
        result.Data = respData
    } else {
        result.Code = errorcode.WxOpenRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewCodeAdd() *codeAdd {
    ca := &codeAdd{wx.NewBaseWxOpen(), ""}
    ca.ReqContentType = project.HTTPContentTypeJSON
    ca.ReqMethod = fasthttp.MethodPost
    return ca
}
