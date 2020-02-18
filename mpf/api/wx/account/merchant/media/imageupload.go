/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 13:12
 */
package media

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type imageUpload struct {
    wx.BaseWxAccount
    appId        string
    imageName    string // 图片名
    imageContent string // 图片内容
}

func (iu *imageUpload) SetImageName(imageName string) {
    if len(imageName) > 0 {
        iu.imageName = imageName
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图片名不合法", nil))
    }
}

func (iu *imageUpload) SetImageContent(imageContent string) {
    if len(imageContent) > 0 {
        iu.imageContent = imageContent
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图片内容不合法", nil))
    }
}

func (iu *imageUpload) checkData() {
    if len(iu.imageName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图片名不能为空", nil))
    }
    if len(iu.imageContent) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "图片内容不能为空", nil))
    }
}

func (iu *imageUpload) SendRequest() api.APIResult {
    iu.checkData()

    iu.ReqData["filename"] = iu.imageName
    iu.ReqData["access_token"] = wx.NewUtilWx().GetSingleAccessToken(iu.appId)
    iu.ReqURI = "https://api.weixin.qq.com/merchant/common/upload_img?" + mpf.HTTPCreateParams(iu.ReqData, "none", 1)
    client, req := iu.GetRequest()
    req.SetBody([]byte(iu.imageContent))

    resp, result := iu.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewImageUpload(appId string) *imageUpload {
    iu := &imageUpload{wx.NewBaseWxAccount(), "", "", ""}
    iu.appId = appId
    iu.ReqContentType = project.HTTPContentTypeForm
    iu.ReqMethod = fasthttp.MethodPost
    return iu
}
