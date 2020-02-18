/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 9:17
 */
package profitsharing

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

// 完结分账
type sharingFinish struct {
    wx.BaseWxAccount
    appId         string
    transactionId string // 微信单号
    outOrderNo    string // 商户分账单号
    description   string // 分账完结描述
}

func (sf *sharingFinish) SetTransactionId(transactionId string) {
    match, _ := regexp.MatchString(project.RegexDigit, transactionId)
    if match {
        sf.transactionId = transactionId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信单号不合法", nil))
    }
}

func (sf *sharingFinish) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outOrderNo)
    if match {
        sf.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户分账单号不合法", nil))
    }
}

func (sf *sharingFinish) SetDescription(description string) {
    if len(description) > 0 {
        trueDesc := []rune(description)
        sf.description = string(trueDesc[:40])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账完结描述不能为空", nil))
    }
}

func (sf *sharingFinish) checkData() {
    if len(sf.transactionId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信单号不能为空", nil))
    }
    if len(sf.outOrderNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户分账单号不能为空", nil))
    }
    if len(sf.description) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账完结描述不能为空", nil))
    }
    sf.ReqData["transaction_id"] = sf.transactionId
    sf.ReqData["out_order_no"] = sf.outOrderNo
    sf.ReqData["description"] = sf.description
}

func (sf *sharingFinish) SendRequest() api.ApiResult {
    sf.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(sf.ReqData, sf.appId, "sha256")
    sf.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(sf.ReqData))
    sf.ReqUrl = "https://api.mch.weixin.qq.com/secapi/pay/profitsharingfinish"
    client, req := sf.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(sf.appId)
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

    resp, result := sf.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewSharingFinish(appId, merchantType string) *sharingFinish {
    conf := wx.NewConfig().GetAccount(appId)
    sf := &sharingFinish{wx.NewBaseWxAccount(), "", "", "", ""}
    sf.appId = appId
    sf.SetPayAccount(conf, merchantType)
    delete(sf.ReqData, "sub_appid")
    sf.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    sf.ReqData["sign_type"] = "HMAC-SHA256"
    sf.ReqContentType = project.HTTPContentTypeXML
    sf.ReqMethod = fasthttp.MethodPost
    return sf
}
