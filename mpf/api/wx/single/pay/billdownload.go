/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 11:09
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

type billDownload struct {
    wx.BaseWxAccount
    appId      string
    deviceInfo string // 设备号
    billDate   string // 对账单日期
    billType   string // 账单类型
    outputFile string // 输出文件全名
}

func (bd *billDownload) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        bd.ReqData["device_info"] = deviceInfo
    }
}

func (bd *billDownload) SetBillDate(billDate string) {
    match, _ := regexp.MatchString(`^2[0-9]{7}$`, billDate)
    if match {
        bd.billDate = billDate
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "对账单日期不合法", nil))
    }
}

func (bd *billDownload) SetBillType(billType string) {
    if (billType == "ALL") || (billType == "SUCCESS") || (billType == "REFUND") || (billType == "RECHARGE_REFUND") {
        bd.ReqData["bill_type"] = billType
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "账单类型不合法", nil))
    }
}

func (bd *billDownload) SetOutputFile(outputFile string) {
    f, err := os.Stat(outputFile)
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出文件不合法", nil))
    }
    if f.IsDir() {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "输出文件不合法", nil))
    }
    bd.outputFile = outputFile
}

func (bd *billDownload) checkData() {
    if len(bd.billDate) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "对账单日期不能为空", nil))
    }
    if len(bd.outputFile) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "输出文件不能为空", nil))
    }
    bd.ReqData["bill_date"] = bd.billDate
}

func (bd *billDownload) SendRequest() api.ApiResult {
    bd.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(bd.ReqData, bd.appId, "md5")
    bd.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(bd.ReqData))
    bd.ReqUrl = "https://api.mch.weixin.qq.com/pay/downloadbill"
    client, req := bd.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(bd.appId)
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

    resp, result := bd.SendInner(client, req, errorcode.WxAccountRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    err = xml.Unmarshal(resp.Body, (*mpf.XmlMap)(&respData))
    if err != nil {
        f, err := os.Create(bd.outputFile)
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
    } else {
        result.Code = errorcode.WxAccountRequestPost
        result.Msg = respData["return_msg"]
    }
    return result
}

func NewBillDownload(appId string) *billDownload {
    bd := &billDownload{wx.NewBaseWxAccount(), "", "", "", "", ""}
    bd.appId = appId
    bd.ReqData["sign_type"] = "MD5"
    bd.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    bd.ReqData["tar_type"] = "GZIP"
    bd.ReqData["bill_type"] = "ALL"
    bd.ReqContentType = project.HTTPContentTypeXML
    bd.ReqMethod = fasthttp.MethodPost
    return bd
}
