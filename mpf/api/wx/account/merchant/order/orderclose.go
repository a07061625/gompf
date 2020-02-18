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

type orderClose struct {
    wx.BaseWxAccount
    appId   string
    orderId string // 订单ID
}

func (oc *orderClose) SetOrderId(orderId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, orderId)
    if match {
        oc.orderId = orderId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单ID不合法", nil))
    }
}

func (oc *orderClose) checkData() {
    if len(oc.orderId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单ID不能为空", nil))
    }
    oc.ReqData["order_id"] = oc.orderId
}

func (oc *orderClose) SendRequest() api.ApiResult {
    oc.checkData()

    reqBody := mpf.JSONMarshal(oc.ReqData)
    oc.ReqUrl = "https://api.weixin.qq.com/merchant/order/close?access_token=" + wx.NewUtilWx().GetSingleAccessToken(oc.appId)
    client, req := oc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := oc.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewOrderClose(appId string) *orderClose {
    oc := &orderClose{wx.NewBaseWxAccount(), "", ""}
    oc.appId = appId
    oc.ReqContentType = project.HTTPContentTypeJSON
    oc.ReqMethod = fasthttp.MethodPost
    return oc
}
