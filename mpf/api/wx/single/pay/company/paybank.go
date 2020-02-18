/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 0:50
 */
package company

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
    "github.com/a07061625/gompf/mpf/api/wx/single"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpencrypt"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type payBank struct {
    wx.BaseWxAccount
    appId          string
    partnerTradeNo string // 付款单号
    encBankNo      string // 收款方银行卡号
    encTrueName    string // 收款方用户名
    bankCode       string // 收款方开户行
    amount         int    // 付款金额
    desc           string // 付款说明
}

func (pcb *payBank) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, partnerTradeNo)
    if match {
        pcb.partnerTradeNo = partnerTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款单号不合法", nil))
    }
}

func (pcb *payBank) SetBankInfo(bankCode, accountNo, accountName string) {
    _, ok := single.CompanyBankCodes[bankCode]
    if !ok {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "收款方开户行不支持", nil))
    }
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, accountNo)
    if !match {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "收款方银行卡号不合法", nil))
    }
    if len(accountName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "收款方用户名不合法", nil))
    }
    pcb.encBankNo = accountNo
    pcb.encTrueName = accountName
    pcb.bankCode = bankCode
}

func (pcb *payBank) SetAmount(amount int) {
    if (amount > 0) && (amount <= 2000000) {
        pcb.amount = amount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额不合法", nil))
    }
}

func (pcb *payBank) SetDesc(desc string) {
    if len(desc) > 0 {
        trueDesc := []rune(desc)
        pcb.ReqData["desc"] = string(trueDesc[:50])
    }
}

func (pcb *payBank) checkData() {
    if len(pcb.partnerTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款单号不能为空", nil))
    }
    if len(pcb.encBankNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "银行信息不能为空", nil))
    }
    if pcb.amount <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额不能为空", nil))
    }
    pcb.ReqData["partner_trade_no"] = pcb.partnerTradeNo
    pcb.ReqData["amount"] = strconv.Itoa(pcb.amount)
    pcb.ReqData["bank_code"] = pcb.bankCode
}

func (pcb *payBank) SendRequest() api.APIResult {
    pcb.checkData()

    conf := wx.NewConfig().GetAccount(pcb.appId)
    keyContent := conf.GetSslCompanyBank()
    if len(keyContent) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "银行卡公钥不能为空", nil))
    }
    encNo, _ := mpencrypt.RsaEncrypt([]byte(pcb.encBankNo), []byte(keyContent))
    encName, _ := mpencrypt.RsaEncrypt([]byte(pcb.encTrueName), []byte(keyContent))
    pcb.ReqData["enc_bank_no"] = string(encNo)
    pcb.ReqData["enc_true_name"] = string(encName)
    sign := wx.NewUtilWx().CreateSinglePaySign(pcb.ReqData, pcb.appId, "md5")
    pcb.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(pcb.ReqData))
    pcb.ReqURI = "https://api.mch.weixin.qq.com/mmpaysptrans/pay_bank"
    client, req := pcb.GetRequest()
    req.SetBody([]byte(reqBody))

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

    resp, result := pcb.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPayBank(appId string) *payBank {
    conf := wx.NewConfig().GetAccount(appId)
    pcb := &payBank{wx.NewBaseWxAccount(), "", "", "", "", "", 0, ""}
    pcb.appId = appId
    pcb.amount = 0
    pcb.ReqData["mch_id"] = conf.GetPayMchId()
    pcb.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    pcb.ReqContentType = project.HTTPContentTypeXML
    pcb.ReqMethod = fasthttp.MethodPost
    return pcb
}
