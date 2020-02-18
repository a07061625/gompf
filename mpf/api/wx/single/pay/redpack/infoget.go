/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 9:12
 */
package redpack

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

type infoGet struct {
    wx.BaseWxAccount
    appId     string
    mchBillNo string // 商户订单号
}

func (ig *infoGet) SetMchBillNo(mchBillNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, mchBillNo)
    if match {
        ig.mchBillNo = mchBillNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户订单号不合法", nil))
    }
}

func (ig *infoGet) checkData() {
    if len(ig.mchBillNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户订单号不能为空", nil))
    }
    ig.ReqData["mch_billno"] = ig.mchBillNo
}

func (ig *infoGet) SendRequest() api.ApiResult {
    ig.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(ig.ReqData, ig.appId, "md5")
    ig.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(ig.ReqData))
    ig.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gethbinfo"
    client, req := ig.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(ig.appId)
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

    resp, result := ig.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewInfoGet(appId string) *infoGet {
    conf := wx.NewConfig().GetAccount(appId)
    ig := &infoGet{wx.NewBaseWxAccount(), "", ""}
    ig.appId = appId
    ig.ReqData["appid"] = conf.GetAppId()
    ig.ReqData["mch_id"] = conf.GetPayMchId()
    ig.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    ig.ReqData["bill_type"] = "MCHT"
    ig.ReqContentType = project.HTTPContentTypeXML
    ig.ReqMethod = fasthttp.MethodPost
    return ig
}
