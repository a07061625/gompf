/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/13 0013
 * Time: 22:15
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

type deliverySet struct {
    wx.BaseWxAccount
    appId           string
    orderId         string // 订单ID
    deliveryCompany string // 物流公司
    deliveryTrackNo string // 运单ID
    needDelivery    int    // 物流状态(0-不需要 1-需要)
    isOthers        int    // 其它物流公司状态(0-否 1-是)
}

func (ds *deliverySet) SetOrderId(orderId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, orderId)
    if match {
        ds.orderId = orderId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单ID不合法", nil))
    }
}

func (ds *deliverySet) SetDeliveryCompany(deliveryCompany string) {
    if len(deliveryCompany) > 0 {
        ds.deliveryCompany = deliveryCompany
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "物流公司不合法", nil))
    }
}

func (ds *deliverySet) SetDeliveryTrackNo(deliveryTrackNo string) {
    if len(deliveryTrackNo) > 0 {
        ds.deliveryTrackNo = deliveryTrackNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "运单ID不合法", nil))
    }
}

func (ds *deliverySet) SetNeedDelivery(needDelivery int) {
    if (needDelivery == 0) || (needDelivery == 1) {
        ds.needDelivery = needDelivery
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "物流状态不合法", nil))
    }
}

func (ds *deliverySet) SetIsOthers(isOthers int) {
    if (isOthers == 0) || (isOthers == 1) {
        ds.isOthers = isOthers
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "其它物流公司状态不合法", nil))
    }
}

func (ds *deliverySet) checkData() {
    if len(ds.orderId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单ID不能为空", nil))
    }
    if ds.needDelivery == 1 {
        if len(ds.deliveryCompany) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "物流公司不能为空", nil))
        }
        if len(ds.deliveryTrackNo) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "运单ID不能为空", nil))
        }
    }
}

func (ds *deliverySet) SendRequest() api.APIResult {
    ds.checkData()

    reqData := make(map[string]interface{})
    reqData["order_id"] = ds.orderId
    reqData["need_delivery"] = ds.needDelivery
    reqData["is_others"] = ds.isOthers
    if ds.needDelivery == 1 {
        reqData["delivery_company"] = ds.deliveryCompany
        reqData["delivery_track_no"] = ds.deliveryTrackNo
    }
    reqBody := mpf.JSONMarshal(reqData)
    ds.ReqURI = "https://api.weixin.qq.com/merchant/order/setdelivery?access_token=" + wx.NewUtilWx().GetSingleAccessToken(ds.appId)
    client, req := ds.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ds.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewDeliverySet(appId string) *deliverySet {
    ds := &deliverySet{wx.NewBaseWxAccount(), "", "", "", "", 0, 0}
    ds.appId = appId
    ds.needDelivery = 1
    ds.isOthers = 0
    ds.ReqContentType = project.HTTPContentTypeJSON
    ds.ReqMethod = fasthttp.MethodPost
    return ds
}
