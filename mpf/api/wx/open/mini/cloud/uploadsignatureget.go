/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/11 0011
 * Time: 9:20
 */
package cloud

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

// 获取上传云函数签名header
type uploadSignatureGet struct {
    wx.BaseWxOpen
    appId         string // 应用ID
    hashedPayload string // 上传签名
}

func (usg *uploadSignatureGet) SetHashedPayload(hashedPayload string) {
    match, _ := regexp.MatchString(`^[0-9a-z]{64}$`, hashedPayload)
    if match {
        usg.hashedPayload = hashedPayload
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "上传签名不合法", nil))
    }
}

func (usg *uploadSignatureGet) checkData() {
    if len(usg.hashedPayload) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "上传签名不能为空", nil))
    }
    usg.ReqData["hashed_payload"] = usg.hashedPayload
}

func (usg *uploadSignatureGet) SendRequest() api.APIResult {
    usg.checkData()

    reqBody := mpf.JSONMarshal(usg.ReqData)
    usg.ReqURI = "https://api.weixin.qq.com/tcb/getuploadsignature?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(usg.appId)
    client, req := usg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := usg.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewUploadSignatureGet(appId string) *uploadSignatureGet {
    usg := &uploadSignatureGet{wx.NewBaseWxOpen(), "", ""}
    usg.appId = appId
    usg.ReqContentType = project.HTTPContentTypeJSON
    usg.ReqMethod = fasthttp.MethodPost
    return usg
}
