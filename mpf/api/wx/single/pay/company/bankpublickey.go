/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 1:17
 */
package company

import (
    "crypto/tls"
    "encoding/xml"
    "io/ioutil"
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type bankPublicKey struct {
    wx.BaseWxAccount
    appId string
}

func (bpk *bankPublicKey) SendRequest() api.APIResult {
    sign := wx.NewUtilWx().CreateSinglePaySign(bpk.ReqData, bpk.appId, "md5")
    bpk.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(bpk.ReqData))
    bpk.ReqURI = "https://fraud.mch.weixin.qq.com/risk/getpublickey"
    client, req := bpk.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(bpk.appId)
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

    resp, result := bpk.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewBankPublicKey(appId string) *bankPublicKey {
    conf := wx.NewConfig().GetAccount(appId)
    bpk := &bankPublicKey{wx.NewBaseWxAccount(), ""}
    bpk.appId = appId
    bpk.ReqData["mch_id"] = conf.GetPayMchId()
    bpk.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    bpk.ReqData["sign_type"] = "MD5"
    bpk.ReqContentType = project.HTTPContentTypeXML
    bpk.ReqMethod = fasthttp.MethodPost
    return bpk
}
