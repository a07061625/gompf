/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 14:17
 */
package pay

import (
    "crypto/tls"
    "encoding/xml"
    "io/ioutil"
    "os"
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type orderRefund struct {
    wx.BaseWxAccount
    appId         string
    deviceInfo    string // 设备号
    transactionId string // 微信订单号
    outTradeNo    string // 商户订单号
    outRefundNo   string // 商户退款单号
    totalFee      int    // 订单总金额,单位为分
    refundFee     int    // 退款总金额,单位为分
}

func (or *orderRefund) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        or.ReqData["device_info"] = deviceInfo
    }
}

func (or *orderRefund) SetTransactionId(transactionId string) {
    match, _ := regexp.MatchString(`^[0-9]{27}$`, transactionId)
    if match {
        or.transactionId = transactionId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信订单号不合法", nil))
    }
}

func (or *orderRefund) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outTradeNo)
    if match {
        or.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (or *orderRefund) SetOutRefundNo(outRefundNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outRefundNo)
    if match {
        or.outRefundNo = outRefundNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户退款单号不合法", nil))
    }
}

func (or *orderRefund) SetTotalFee(totalFee int) {
    if totalFee > 0 {
        or.totalFee = totalFee
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单金额不合法", nil))
    }
}

func (or *orderRefund) SetRefundFee(refundFee int) {
    if refundFee > 0 {
        or.refundFee = refundFee
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "退款金额不合法", nil))
    }
}

func (or *orderRefund) checkData() {
    if len(or.transactionId) > 0 {
        or.ReqData["transaction_id"] = or.transactionId
        delete(or.ReqData, "out_trade_no")
    } else if len(or.outTradeNo) > 0 {
        or.ReqData["out_trade_no"] = or.outTradeNo
        delete(or.ReqData, "transaction_id")
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信订单号与商户订单号不能同时为空", nil))
    }
    if len(or.outRefundNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户退款单号不能为空", nil))
    }
    if or.totalFee <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单金额必须大于0", nil))
    } else if or.refundFee <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "退款金额必须大于0", nil))
    } else if or.refundFee > or.totalFee {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "订单金额必须大于等于退款金额", nil))
    }
    or.ReqData["total_fee"] = strconv.Itoa(or.totalFee)
    or.ReqData["refund_fee"] = strconv.Itoa(or.refundFee)
}

func (or *orderRefund) SendRequest() api.ApiResult {
    or.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(or.ReqData, or.appId, "md5")
    or.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(or.ReqData))
    or.ReqUrl = "https://api.mch.weixin.qq.com/secapi/pay/refund"
    client, req := or.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(or.appId)
    keyFile, _ := ioutil.TempFile("", "tmpfile")
    defer os.Remove(keyFile.Name())
    if _, err := keyFile.Write([]byte(conf.GetSslKey())); err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "写入证书key文件失败", nil))
    }

    certFile, _ := ioutil.TempFile("", "tmpfile")
    defer os.Remove(certFile.Name())
    if _, err := certFile.Write([]byte(conf.GetSslCert())); err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "写入证书cert文件失败", nil))
    }

    certs, err := tls.LoadX509KeyPair(certFile.Name(), keyFile.Name())
    if err != nil {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "加载证书文件失败", nil))
    }
    client.TLSConfig.Certificates = []tls.Certificate{certs}

    resp, result := or.SendInner(client, req, errorcode.WxAccountRequestPost)
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
    }
    result.Data = respData
    return result
}

func NewOrderRefund(appId, merchantType string) *orderRefund {
    conf := wx.NewConfig().GetAccount(appId)
    or := &orderRefund{wx.NewBaseWxAccount(), "", "", "", "", "", 0, 0}
    or.appId = appId
    or.totalFee = 0
    or.refundFee = 0
    or.SetPayAccount(conf, merchantType)
    or.ReqData["op_user_id"] = conf.GetPayMchId()
    or.ReqData["sign_type"] = "MD5"
    or.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    or.ReqData["refund_fee_type"] = "CNY"
    or.ReqContentType = project.HTTPContentTypeXML
    or.ReqMethod = fasthttp.MethodPost
    return or
}
