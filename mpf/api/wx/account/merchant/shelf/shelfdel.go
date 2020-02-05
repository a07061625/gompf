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

type shelfDel struct {
    wx.BaseWxAccount
    appId   string
    shelfId int // 货架ID
}

func (sd *shelfDel) SetShelfId(shelfId int) {
    if shelfId > 0 {
        sd.shelfId = shelfId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架ID不合法", nil))
    }
}

func (sd *shelfDel) checkData() {
    if sd.shelfId <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "货架ID不能为空", nil))
    }
}

func (sd *shelfDel) SendRequest() api.ApiResult {
    sd.checkData()

    reqData := make(map[string]interface{})
    reqData["shelf_id"] = sd.shelfId
    reqBody := mpf.JsonMarshal(reqData)
    sd.ReqUrl = "https://api.weixin.qq.com/merchant/shelf/del?access_token=" + wx.NewUtilWx().GetSingleAccessToken(sd.appId)
    client, req := sd.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sd.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewShelfDel(appId string) *shelfDel {
    sd := &shelfDel{wx.NewBaseWxAccount(), "", 0}
    sd.appId = appId
    sd.ReqContentType = project.HttpContentTypeJson
    sd.ReqMethod = fasthttp.MethodPost
    return sd
}
