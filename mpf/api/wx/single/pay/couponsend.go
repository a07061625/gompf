/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 10:06
 */
package pay

import (
    "crypto/tls"
    "encoding/xml"
    "io/ioutil"
    "os"
    "regexp"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type couponSend struct {
    wx.BaseWxAccount
    appId          string
    stockId        string // 代金券批次id
    openidCount    int    // openid记录数
    partnerTradeNo string // 商户单据号
    openid         string // 用户openid
    opUserId       string // 操作员
    deviceInfo     string // 设备号
}

func (cs *couponSend) SetStockId(stockId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,64}$`, stockId)
    if match {
        cs.stockId = stockId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "代金券批次id", nil))
    }
}

func (cs *couponSend) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, partnerTradeNo)
    if match {
        nowTime := time.Now()
        cs.partnerTradeNo = cs.ReqData["mch_id"] + nowTime.Format("20060102") + partnerTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单据号不合法", nil))
    }
}

func (cs *couponSend) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        cs.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (cs *couponSend) SetOpUserId(opUserId string) {
    match, _ := regexp.MatchString(project.RegexDigit, opUserId)
    if match {
        cs.ReqData["op_user_id"] = opUserId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "操作员不合法", nil))
    }
}

func (cs *couponSend) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        cs.ReqData["device_info"] = deviceInfo
    }
}

func (cs *couponSend) checkData() {
    if len(cs.stockId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "代金券批次id不能为空", nil))
    }
    if len(cs.partnerTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单据号不能为空", nil))
    }
    if len(cs.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    cs.ReqData["coupon_stock_id"] = cs.stockId
    cs.ReqData["partner_trade_no"] = cs.partnerTradeNo
    cs.ReqData["openid"] = cs.openid
}

func (cs *couponSend) SendRequest() api.ApiResult {
    cs.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(cs.ReqData, cs.appId, "md5")
    cs.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(cs.ReqData))
    cs.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/send_coupon"
    client, req := cs.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(cs.appId)
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

    resp, result := cs.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewCouponSend(appId string) *couponSend {
    conf := wx.NewConfig().GetAccount(appId)
    cs := &couponSend{wx.NewBaseWxAccount(), "", "", 0, "", "", "", ""}
    cs.appId = appId
    cs.ReqData["appid"] = conf.GetAppId()
    cs.ReqData["mch_id"] = conf.GetPayMchId()
    cs.ReqData["op_user_id"] = conf.GetPayMchId()
    cs.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    cs.ReqData["openid_count"] = "1"
    cs.ReqData["version"] = "1.0"
    cs.ReqData["type"] = "XML"
    cs.ReqContentType = project.HTTPContentTypeXML
    cs.ReqMethod = fasthttp.MethodPost
    return cs
}
