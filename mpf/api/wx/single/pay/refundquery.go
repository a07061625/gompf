/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 20:07
 */
package pay

import (
    "encoding/xml"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type refundQuery struct {
    wx.BaseWxAccount
    appId         string
    deviceInfo    string // 设备号
    transactionId string // 微信订单号
    outTradeNo    string // 商户订单号
    refundId      string // 微信退款单号
    outRefundNo   string // 商户退款单号
}

func (rq *refundQuery) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        rq.ReqData["device_info"] = deviceInfo
    }
}

func (rq *refundQuery) SetTransactionId(transactionId string) {
    match, _ := regexp.MatchString(`^[0-9]{27}$`, transactionId)
    if match {
        rq.transactionId = transactionId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信订单号不合法", nil))
    }
}

func (rq *refundQuery) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outTradeNo)
    if match {
        rq.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (rq *refundQuery) SetRefundId(refundId string) {
    match, _ := regexp.MatchString(`^[0-9]{28}$`, refundId)
    if match {
        rq.refundId = refundId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信退款单号不合法", nil))
    }
}

func (rq *refundQuery) SetOutRefundNo(outRefundNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outRefundNo)
    if match {
        rq.outRefundNo = outRefundNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户退款单号不合法", nil))
    }
}

func (rq *refundQuery) checkData() {
    if len(rq.refundId) > 0 {
        rq.ReqData["refund_id"] = rq.refundId
    } else if len(rq.outRefundNo) > 0 {
        rq.ReqData["out_refund_no"] = rq.outRefundNo
    } else if len(rq.transactionId) > 0 {
        rq.ReqData["transaction_id"] = rq.transactionId
    } else if len(rq.outTradeNo) > 0 {
        rq.ReqData["out_trade_no"] = rq.outTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信订单号,商户订单号,微信退款单号,商户退款单号必须设置其中一个", nil))
    }
}

func (rq *refundQuery) SendRequest() api.APIResult {
    rq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(rq.ReqData, rq.appId, "md5")
    rq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(rq.ReqData))
    rq.ReqURI = "https://api.mch.weixin.qq.com/pay/refundquery"
    client, req := rq.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := rq.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XMLMap)(&respData))
    if respData["return_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["return_msg"]
    } else if respData["result_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["err_code_des"]
    } else {
        result.Data = respData
    }
    return result
}

func NewRefundQuery(appId, merchantType string) *refundQuery {
    conf := wx.NewConfig().GetAccount(appId)
    rq := &refundQuery{wx.NewBaseWxAccount(), "", "", "", "", "", ""}
    rq.appId = appId
    rq.SetPayAccount(conf, merchantType)
    rq.ReqData["sign_type"] = "MD5"
    rq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    rq.ReqContentType = project.HTTPContentTypeXML
    rq.ReqMethod = fasthttp.MethodPost
    return rq
}
