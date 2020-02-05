/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 18:23
 */
package order

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type orderList struct {
    wx.BaseWxAccount
    appId       string
    orderStatus int // 订单状态
    beginTime   int // 创建起始时间
    endTime     int // 创建终止时间
}

func (ol *orderList) SetOrderStatus(orderStatus int) {
    switch orderStatus {
    case 2:
        ol.orderStatus = orderStatus
    case 3:
        ol.orderStatus = orderStatus
    case 5:
        ol.orderStatus = orderStatus
    case 8:
        ol.orderStatus = orderStatus
    default:
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单状态不合法", nil))
    }
}

func (ol *orderList) SetTime(beginTime, endTime int) {
    if beginTime < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "创建起始时间不合法", nil))
    } else if endTime < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "创建终止时间不合法", nil))
    } else if (beginTime > 0) && (endTime > 0) && (beginTime > endTime) {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "创建起始时间不能大于终止时间", nil))
    }
    ol.beginTime = beginTime
    ol.endTime = endTime
}

func (ol *orderList) SendRequest() api.ApiResult {
    reqData := make(map[string]interface{})
    if ol.beginTime > 0 {
        reqData["begintime"] = ol.beginTime
    }
    if ol.endTime > 0 {
        reqData["endtime"] = ol.endTime
    }
    if ol.orderStatus > 0 {
        reqData["status"] = ol.orderStatus
    }
    reqBody := mpf.JsonMarshal(ol.ReqData)
    ol.ReqUrl = "https://api.weixin.qq.com/merchant/order/getbyfilter?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ol.appId)
    client, req := ol.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ol.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewOrderList(appId string) *orderList {
    ol := &orderList{wx.NewBaseWxAccount(), "", 0, 0, 0}
    ol.appId = appId
    ol.ReqContentType = project.HttpContentTypeJson
    ol.ReqMethod = fasthttp.MethodPost
    return ol
}
