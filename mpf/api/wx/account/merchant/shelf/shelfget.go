/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 17:24
 */
package shelf

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type shelfGet struct {
    wx.BaseWxAccount
    appId   string
    shelfId int // 货架ID
}

func (sg *shelfGet) SetShelfId(shelfId int) {
    if shelfId > 0 {
        sg.shelfId = shelfId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架ID不合法", nil))
    }
}

func (sg *shelfGet) checkData() {
    if sg.shelfId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架ID不能为空", nil))
    }
}

func (sg *shelfGet) SendRequest() api.APIResult {
    sg.checkData()

    reqData := make(map[string]interface{})
    reqData["shelf_id"] = sg.shelfId
    reqBody := mpf.JSONMarshal(reqData)
    sg.ReqURI = "https://api.weixin.qq.com/merchant/shelf/getbyid?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sg.appId)
    client, req := sg.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sg.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewShelfGet(appId string) *shelfGet {
    sg := &shelfGet{wx.NewBaseWxAccount(), "", 0}
    sg.appId = appId
    sg.ReqContentType = project.HTTPContentTypeJSON
    sg.ReqMethod = fasthttp.MethodPost
    return sg
}
