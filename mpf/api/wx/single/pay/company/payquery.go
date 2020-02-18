/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 19:27
 */
package company

import (
    "crypto/tls"
    "encoding/xml"
    "io/ioutil"
    "os"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type payQuery struct {
    wx.BaseWxAccount
    appId          string
    partnerTradeNo string // 商户订单号
}

func (pq *payQuery) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, partnerTradeNo)
    if match {
        pq.partnerTradeNo = partnerTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (pq *payQuery) checkData() {
    if len(pq.partnerTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不能为空", nil))
    }
    pq.ReqData["partner_trade_no"] = pq.partnerTradeNo
}

func (pq *payQuery) SendRequest() api.ApiResult {
    pq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(pq.ReqData, pq.appId, "md5")
    pq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(pq.ReqData))
    pq.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"
    client, req := pq.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(pq.appId)
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

    resp, result := pq.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPayQuery(appId string) *payQuery {
    conf := wx.NewConfig().GetAccount(appId)
    pq := &payQuery{wx.NewBaseWxAccount(), "", ""}
    pq.appId = appId
    pq.ReqData["appid"] = conf.GetAppId()
    pq.ReqData["mch_id"] = conf.GetPayMchId()
    pq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    pq.ReqContentType = project.HTTPContentTypeXML
    pq.ReqMethod = fasthttp.MethodPost
    return pq
}
