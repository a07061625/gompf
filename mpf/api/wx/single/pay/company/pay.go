/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 16:42
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
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type pay struct {
    wx.BaseWxAccount
    appId          string
    partnerTradeNo string // 商户订单号
    openid         string // 用户openid
    checkName      string // 校验用户姓名选项
    reUserName     string // 收款用户姓名
    amount         int    // 金额
    desc           string // 企业付款描述信息
}

func (p *pay) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, partnerTradeNo)
    if match {
        p.partnerTradeNo = partnerTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (p *pay) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        p.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (p *pay) SetUserName(checkName, userName string) {
    if checkName == "NO_CHECK" {
        p.checkName = checkName
        p.reUserName = ""
    } else if checkName == "FORCE_CHECK" {
        if len(userName) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "收款用户姓名不合法", nil))
        }

        trueName := []rune(userName)
        p.checkName = checkName
        p.reUserName = string(trueName[:32])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "校验用户姓名选项不合法", nil))
    }
}

func (p *pay) SetAmount(amount int) {
    if amount > 0 {
        p.amount = amount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额必须大于0", nil))
    }
}

func (p *pay) SetDesc(desc string) {
    if len(desc) > 0 {
        p.desc = desc
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款描述信息不合法", nil))
    }
}

func (p *pay) checkData() {
    if len(p.partnerTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不能为空", nil))
    }
    if len(p.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if p.amount <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额必须大于0", nil))
    }
    if len(p.desc) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款描述信息不能为空", nil))
    }
    p.ReqData["partner_trade_no"] = p.partnerTradeNo
    p.ReqData["openid"] = p.openid
    p.ReqData["check_name"] = p.checkName
    if len(p.reUserName) > 0 {
        p.ReqData["re_user_name"] = p.reUserName
    }
    p.ReqData["amount"] = strconv.Itoa(p.amount)
    p.ReqData["desc"] = p.desc
}

func (p *pay) SendRequest() api.ApiResult {
    p.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(p.ReqData, p.appId, "md5")
    p.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(p.ReqData))
    p.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
    client, req := p.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(p.appId)
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

    resp, result := p.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPay(appId string) *pay {
    conf := wx.NewConfig().GetAccount(appId)
    p := &pay{wx.NewBaseWxAccount(), "", "", "", "", "", 0, ""}
    p.appId = appId
    p.amount = 0
    p.checkName = "NO_CHECK"
    p.ReqData["mch_appid"] = conf.GetAppId()
    p.ReqData["mch_id"] = conf.GetPayMchId()
    p.ReqData["spbill_create_ip"] = conf.GetClientIp()
    p.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    p.ReqContentType = project.HttpContentTypeXml
    p.ReqMethod = fasthttp.MethodPost
    return p
}
