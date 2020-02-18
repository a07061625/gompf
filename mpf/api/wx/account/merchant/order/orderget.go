/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 18:17
 */
package order

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

type orderGet struct {
    wx.BaseWxAccount
    appId   string
    orderId string // 订单ID
}

func (og *orderGet) SetOrderId(orderId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, orderId)
    if match {
        og.orderId = orderId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单ID不合法", nil))
    }
}

func (og *orderGet) checkData() {
    if len(og.orderId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单ID不能为空", nil))
    }
    og.ReqData["order_id"] = og.orderId
}

func (og *orderGet) SendRequest() api.ApiResult {
    og.checkData()

    reqBody := mpf.JSONMarshal(og.ReqData)
    og.ReqUrl = "https://api.weixin.qq.com/merchant/order/getbyid?access_token=" + wx.NewUtilWx().GetSingleAccessToken(og.appId)
    client, req := og.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := og.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewOrderGet(appId string) *orderGet {
    og := &orderGet{wx.NewBaseWxAccount(), "", ""}
    og.appId = appId
    og.ReqContentType = project.HTTPContentTypeJSON
    og.ReqMethod = fasthttp.MethodPost
    return og
}
