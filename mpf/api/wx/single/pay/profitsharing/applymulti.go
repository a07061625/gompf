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

// 请求多次分账
type applyMulti struct {
    wx.BaseWxAccount
    appId         string
    transactionId string                   // 微信单号
    outOrderNo    string                   // 商户分账单号
    receivers     []map[string]interface{} // 分账接收方列表
}

func (am *applyMulti) SetTransactionId(transactionId string) {
    match, _ := regexp.MatchString(project.RegexDigit, transactionId)
    if match {
        am.transactionId = transactionId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信单号不合法", nil))
    }
}

func (am *applyMulti) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outOrderNo)
    if match {
        am.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户分账单号不合法", nil))
    }
}

func (am *applyMulti) SetReceivers(receivers []map[string]interface{}) {
    am.receivers = receivers
}

func (am *applyMulti) checkData() {
    if len(am.transactionId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信单号不能为空", nil))
    }
    if len(am.outOrderNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户分账单号不能为空", nil))
    }
    if len(am.receivers) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账接收方不能为空", nil))
    }
    am.ReqData["transaction_id"] = am.transactionId
    am.ReqData["out_order_no"] = am.outOrderNo
    am.ReqData["receivers"] = mpf.JsonMarshal(am.receivers)
}

func (am *applyMulti) SendRequest() api.ApiResult {
    am.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(am.ReqData, am.appId, "sha256")
    am.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(am.ReqData))
    am.ReqUrl = "https://api.mch.weixin.qq.com/secapi/pay/multiprofitsharing"
    client, req := am.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(am.appId)
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

    resp, result := am.SendInner(client, req, errorcode.WxAccountRequestPost)
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
    }
    result.Data = respData
    return result
}

func NewApplyMulti(appId, merchantType string) *applyMulti {
    conf := wx.NewConfig().GetAccount(appId)
    am := &applyMulti{wx.NewBaseWxAccount(), "", "", "", make([]map[string]interface{}, 0)}
    am.appId = appId
    am.SetPayAccount(conf, merchantType)
    am.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    am.ReqData["sign_type"] = "HMAC-SHA256"
    am.ReqContentType = project.HttpContentTypeXml
    am.ReqMethod = fasthttp.MethodPost
    return am
}
