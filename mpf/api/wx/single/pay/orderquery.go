/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 14:10
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

type orderQuery struct {
    wx.BaseWxAccount
    appId         string
    transactionId string // 微信订单号
    outTradeNo    string // 商户订单号
}

func (oq *orderQuery) SetTransactionId(transactionId string) {
    match, _ := regexp.MatchString(`^[0-9]{27}$`, transactionId)
    if match {
        oq.transactionId = transactionId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信订单号不合法", nil))
    }
}

func (oq *orderQuery) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outTradeNo)
    if match {
        oq.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (oq *orderQuery) checkData() {
    if len(oq.transactionId) > 0 {
        oq.ReqData["transaction_id"] = oq.transactionId
        delete(oq.ReqData, "out_trade_no")
    } else if len(oq.outTradeNo) > 0 {
        oq.ReqData["out_trade_no"] = oq.outTradeNo
        delete(oq.ReqData, "transaction_id")
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信订单号与商户订单号不能同时为空", nil))
    }
}

func (oq *orderQuery) SendRequest() api.ApiResult {
    oq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(oq.ReqData, oq.appId, "md5")
    oq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(oq.ReqData))
    oq.ReqUrl = "https://api.mch.weixin.qq.com/pay/orderquery"
    client, req := oq.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := oq.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XmlMap)(&respData))
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

func NewOrderQuery(appId, merchantType string) *orderQuery {
    conf := wx.NewConfig().GetAccount(appId)
    oq := &orderQuery{wx.NewBaseWxAccount(), "", "", ""}
    oq.appId = appId
    oq.SetPayAccount(conf, merchantType)
    oq.ReqData["sign_type"] = "MD5"
    oq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    oq.ReqContentType = project.HTTPContentTypeXML
    oq.ReqMethod = fasthttp.MethodPost
    return oq
}
