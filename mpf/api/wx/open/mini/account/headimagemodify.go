/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/10 0010
 * Time: 15:45
 */
package account

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改头像
type headImageModify struct {
    wx.BaseWxOpen
    appId          string  // 应用ID
    headImgMediaId string  // 头像素材
    x1             float32 // 起始点横坐标
    y1             float32 // 起始点纵坐标
    x2             float32 // 截止点横坐标
    y2             float32 // 截止点纵坐标
}

func (him *headImageModify) SetHeadImage(headImage string) {
    if len(headImage) > 0 {
        him.headImgMediaId = headImage
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "头像素材不合法", nil))
    }
}

func (him *headImageModify) SetX1(x1 float32) {
    if (x1 >= 0) && (x1 <= 1) {
        him.x1 = x1
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "起始点横坐标不合法", nil))
    }
}

func (him *headImageModify) SetY1(y1 float32) {
    if (y1 >= 0) && (y1 <= 1) {
        him.y1 = y1
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "起始点纵坐标不合法", nil))
    }
}

func (him *headImageModify) SetX2(x2 float32) {
    if (x2 >= 0) && (x2 <= 1) {
        him.x2 = x2
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "截止点横坐标不合法", nil))
    }
}

func (him *headImageModify) SetY2(y2 float32) {
    if (y2 >= 0) && (y2 <= 1) {
        him.y2 = y2
    } else {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "截止点纵坐标不合法", nil))
    }
}

func (him *headImageModify) checkData() {
    if len(him.headImgMediaId) == 0 {
        panic(mperr.NewWxOpenMini(errorcode.WxOpenParam, "头像素材不能为空", nil))
    }
}

func (him *headImageModify) SendRequest() api.APIResult {
    him.checkData()

    reqData := make(map[string]interface{})
    reqData["head_img_media_id"] = him.headImgMediaId
    reqData["x1"] = him.x1
    reqData["y1"] = him.y1
    reqData["x2"] = him.x2
    reqData["y2"] = him.y2
    reqBody := mpf.JSONMarshal(reqData)
    him.ReqURI = "https://api.weixin.qq.com/cgi-bin/account/modifyheadimage?access_token=" + wx.NewUtilWx().GetOpenAuthorizeAccessToken(him.appId)
    client, req := him.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := him.SendInner(client, req, errorcode.WxOpenRequestPost)
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

func NewHeadImageModify(appId string) *headImageModify {
    him := &headImageModify{wx.NewBaseWxOpen(), "", "", 0.00, 0.00, 0.00, 0.00}
    him.appId = appId
    him.x1 = 0.00
    him.y1 = 0.00
    him.x2 = 1.00
    him.y2 = 1.00
    him.ReqContentType = project.HTTPContentTypeJSON
    him.ReqMethod = fasthttp.MethodPost
    return him
}
