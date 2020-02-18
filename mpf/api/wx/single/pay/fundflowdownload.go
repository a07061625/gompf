/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 10:29
 */
package pay

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

type fundFlowDownload struct {
    wx.BaseWxAccount
    appId       string
    billDate    string // 资金账单日期
    accountType string // 资金账户类型
    outputFile  string // 输出文件全名
}

func (ffd *fundFlowDownload) SetBillDate(billDate string) {
    match, _ := regexp.MatchString(`^2[0-9]{7}$`, billDate)
    if match {
        ffd.billDate = billDate
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "资金账单日期不合法", nil))
    }
}

func (ffd *fundFlowDownload) SetAccountType(accountType string) {
    if (accountType == "Basic") || (accountType == "Operation") || (accountType == "Fees") {
        ffd.accountType = accountType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "资金账户类型不合法", nil))
    }
}

func (ffd *fundFlowDownload) SetOutputFile(outputFile string) {
    f, err := os.Stat(outputFile)
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出文件不合法", nil))
    }
    ffd.outputFile = outputFile
}

func (ffd *fundFlowDownload) checkData() {
    if len(ffd.billDate) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "资金账单日期不能为空", nil))
    }
    if len(ffd.accountType) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "资金账户类型不能为空", nil))
    }
    if len(ffd.outputFile) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "输出文件不能为空", nil))
    }
    ffd.ReqData["bill_date"] = ffd.billDate
    ffd.ReqData["account_type"] = ffd.accountType
}

func (ffd *fundFlowDownload) SendRequest() api.ApiResult {
    ffd.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(ffd.ReqData, ffd.appId, "sha256")
    ffd.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(ffd.ReqData))
    ffd.ReqUrl = "https://api.mch.weixin.qq.com/pay/downloadfundflow"
    client, req := ffd.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(ffd.appId)
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

    resp, result := ffd.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    err = xml.Unmarshal(resp.Body, (*mpf.XmlMap)(&respData))
    if err != nil {
        f, err := os.Create(ffd.outputFile)
        defer f.Close()
        if err != nil {
            result.Code = errorcode.WxAccountRequestPost
            result.Msg = err.Error()
        } else {
            f.Write(resp.Body)
            resultData := make(map[string]string)
            resultData["return_code"] = "SUCCESS"
            result.Data = resultData
        }
    } else if respData["return_code"] == "FAIL" {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["return_msg"]
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["err_code_des"]
    }
    return result
}

func NewFundFlowDownload(appId string) *fundFlowDownload {
    conf := wx.NewConfig().GetAccount(appId)
    ffd := &fundFlowDownload{wx.NewBaseWxAccount(), "", "", "", ""}
    ffd.appId = appId
    ffd.ReqData["appid"] = conf.GetAppId()
    ffd.ReqData["mch_id"] = conf.GetPayMchId()
    ffd.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    ffd.ReqData["sign_type"] = "HMAC-SHA256"
    ffd.ReqData["tar_type"] = "GZIP"
    ffd.ReqContentType = project.HTTPContentTypeXML
    ffd.ReqMethod = fasthttp.MethodPost
    return ffd
}
