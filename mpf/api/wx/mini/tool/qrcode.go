/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/5 0005
 * Time: 12:59
 */
package tool

import (
    "encoding/base64"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type qrCode struct {
    wx.BaseWxMini
    appId     string            // 应用ID
    scene     string            // 场景
    page      string            // 页面地址
    width     int               // 二维码宽度
    autoColor bool              // 线条颜色配置
    lineColor map[string]string // 线条rgb颜色
    isHyaLine bool              // 透明底色标识
}

func (qc *qrCode) SetScene(scene string) {
    trueScene := strings.TrimSpace(scene)
    length := len(trueScene)
    if length == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "场景标识不能为空", nil))
    } else if length > 32 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "场景标识不能超过32个字符", nil))
    }
    qc.scene = trueScene
}

func (qc *qrCode) SetPage(page string) {
    truePage := strings.TrimSpace(page)
    if len(truePage) > 0 {
        qc.page = truePage
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "页面地址不能为空", nil))
    }
}

func (qc *qrCode) SetWidth(width int) {
    if width > 0 {
        qc.width = width
    } else {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "二维码宽度必须大于0", nil))
    }
}

func (qc *qrCode) SetAutoColor(autoColor bool) {
    qc.autoColor = autoColor
}

func (qc *qrCode) SetLineColor(red, green, blue int) {
    if (red < 0) || (red > 255) {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "线条颜色red不合法", nil))
    }
    if (green < 0) || (green > 255) {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "线条颜色green不合法", nil))
    }
    if (blue < 0) || (blue > 255) {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "线条颜色blue不合法", nil))
    }
    qc.lineColor["r"] = strconv.Itoa(red)
    qc.lineColor["g"] = strconv.Itoa(green)
    qc.lineColor["b"] = strconv.Itoa(blue)
}

func (qc *qrCode) SetIsHyaLine(isHyaLine bool) {
    qc.isHyaLine = isHyaLine
}

func (qc *qrCode) checkData() (*fasthttp.Client, *fasthttp.Request) {
    if len(qc.scene) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "场景标识必须填写", nil))
    }
    if len(qc.page) == 0 {
        panic(mperr.NewWxMini(errorcode.WxMiniParam, "页面地址必须填写", nil))
    }

    reqData := make(map[string]interface{})
    reqData["scene"] = qc.scene
    reqData["page"] = qc.page
    reqData["width"] = qc.width
    reqData["auto_color"] = qc.autoColor
    reqData["line_color"] = qc.lineColor
    reqData["is_hyaline"] = qc.isHyaLine
    reqBody := mpf.JsonMarshal(reqData)

    qc.ReqUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + wx.NewUtilWx().GetSingleCache(qc.appId, wx.SingleCacheTypeAccessToken)
    client, req := qc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func (qc *qrCode) SendRequest() api.ApiResult {
    client, req := qc.checkData()
    resp, result := qc.SendInner(client, req, errorcode.WxMiniRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData, err := mpf.JsonUnmarshalMap(resp.Content)
    if err != nil {
        imageData := make(map[string]string)
        imageData["image"] = base64.StdEncoding.EncodeToString(resp.Body)
        result.Data = imageData
    } else {
        result.Code = errorcode.WxMiniRequestPost
        result.Msg = respData["errmsg"].(string)
    }

    return result
}

func NewQrCode(appId string) *qrCode {
    qc := &qrCode{wx.NewBaseWxMini(), "", "", "", 0, false, make(map[string]string), false}
    qc.appId = appId
    qc.width = 430
    qc.autoColor = false
    qc.isHyaLine = false
    qc.lineColor["r"] = "0"
    qc.lineColor["g"] = "0"
    qc.lineColor["b"] = "0"
    qc.ReqContentType = project.HTTPContentTypeJSON
    qc.ReqMethod = fasthttp.MethodPost
    return qc
}
