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

type payBankQuery struct {
    wx.BaseWxAccount
    appId          string
    partnerTradeNo string // 商户订单号
}

func (pbq *payBankQuery) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, partnerTradeNo)
    if match {
        pbq.partnerTradeNo = partnerTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (pbq *payBankQuery) checkData() {
    if len(pbq.partnerTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不能为空", nil))
    }
    pbq.ReqData["partner_trade_no"] = pbq.partnerTradeNo
}

func (pbq *payBankQuery) SendRequest() api.ApiResult {
    pbq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(pbq.ReqData, pbq.appId, "md5")
    pbq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(pbq.ReqData))
    pbq.ReqUrl = "https://api.mch.weixin.qq.com/mmpaysptrans/query_bank"
    client, req := pbq.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(pbq.appId)
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

    resp, result := pbq.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPayBankQuery(appId string) *payBankQuery {
    conf := wx.NewConfig().GetAccount(appId)
    pbq := &payBankQuery{wx.NewBaseWxAccount(), "", ""}
    pbq.appId = appId
    pbq.ReqData["mch_id"] = conf.GetPayMchId()
    pbq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    pbq.ReqContentType = project.HTTPContentTypeXML
    pbq.ReqMethod = fasthttp.MethodPost
    return pbq
}
