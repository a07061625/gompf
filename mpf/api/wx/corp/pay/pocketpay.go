/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 9:00
 */
package pay

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
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/valyala/fasthttp"
)

// 向员工付款
type pocketPay struct {
    wx.BaseWxCorp
    corpId         string
    agentTag       string
    mchId          string // 商户号
    deviceInfo     string // 设备号
    nonceStr       string // 随机字符串
    partnerTradeNo string // 商户订单号
    openid         string // 用户openid
    checkName      string // 校验用户姓名选项
    reUserName     string // 收款用户姓名
    amount         int    // 金额
    payDesc        string // 付款说明
    spBillCreateIp string // Ip地址
    msgType        string // 付款消息类型
    approvalNumber string // 审批单号
    approvalType   int    // 审批类型
    actName        string // 项目名称
    acceptKeys     []string
}

func (pp *pocketPay) SetDeviceInfo(deviceInfo string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, deviceInfo)
    if match {
        pp.ReqData["device_info"] = deviceInfo
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "设备号不合法", nil))
    }
}

func (pp *pocketPay) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, partnerTradeNo)
    if match {
        pp.partnerTradeNo = partnerTradeNo
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不合法", nil))
    }
}

func (pp *pocketPay) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        pp.openid = openid
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不合法", nil))
    }
}

func (pp *pocketPay) SetCheckName(checkName string) {
    if checkName == "NO_CHECK" {
        pp.checkName = checkName
        pp.reUserName = ""
    } else if checkName == "FORCE_CHECK" {
        pp.checkName = checkName
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "校验用户姓名选项不合法", nil))
    }
}

func (pp *pocketPay) SetReUserName(reUserName string) {
    if (pp.checkName == "FORCE_CHECK") && (len(reUserName) > 0) {
        trueName := []rune(reUserName)
        pp.reUserName = string(trueName[:32])
    }
}

func (pp *pocketPay) SetAmount(amount int) {
    if amount > 0 {
        pp.amount = amount
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款金额必须大于0", nil))
    }
}

func (pp *pocketPay) SetPayDesc(payDesc string) {
    if len(payDesc) > 0 {
        pp.payDesc = payDesc
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款说明不合法", nil))
    }
}

func (pp *pocketPay) SetMsgType(msgType string) {
    if msgType == "NORMAL_MSG" {
        pp.msgType = msgType
        pp.approvalType = 0
        pp.approvalNumber = ""
    } else if msgType == "APPROVAL_MSG" {
        pp.msgType = msgType
        pp.approvalType = 1
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款消息类型不合法", nil))
    }
}

func (pp *pocketPay) SetApprovalNumber(approvalNumber string) {
    if pp.msgType == "APPROVAL_MSG" {
        match, _ := regexp.MatchString(project.RegexDigitAlpha, approvalNumber)
        if match {
            pp.approvalNumber = approvalNumber
        } else {
            panic(mperr.NewWxCorp(errorcode.WxCorpParam, "审批单号不合法", nil))
        }
    }
}

func (pp *pocketPay) SetActName(actName string) {
    if len(actName) > 0 {
        trueName := []rune(actName)
        pp.actName = string(trueName[:25])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "项目名称不合法", nil))
    }
}

func (pp *pocketPay) checkData() {
    if len(pp.partnerTradeNo) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不能为空", nil))
    }
    if len(pp.openid) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不能为空", nil))
    }
    if (pp.checkName == "FORCE_CHECK") && (len(pp.reUserName) == 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "收款用户姓名不能为空", nil))
    }
    if pp.amount <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款金额必须大于0", nil))
    }
    if len(pp.payDesc) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款描述信息不能为空", nil))
    }
    if (pp.msgType == "APPROVAL_MSG") && (len(pp.approvalNumber) == 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "审批单号不能为空", nil))
    }
    if len(pp.actName) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "项目名称不能为空", nil))
    }
    pp.ReqData["partner_trade_no"] = pp.partnerTradeNo
    pp.ReqData["openid"] = pp.openid
    pp.ReqData["check_name"] = pp.checkName
    if pp.checkName == "FORCE_CHECK" {
        pp.ReqData["re_user_name"] = pp.reUserName
    }
    pp.ReqData["amount"] = strconv.Itoa(pp.amount)
    pp.ReqData["desc"] = pp.payDesc
    pp.ReqData["ww_msg_type"] = pp.msgType
    if pp.msgType == "NORMAL_MSG" {
        delete(pp.ReqData, "approval_number")
        delete(pp.ReqData, "approval_type")
    } else if pp.msgType == "APPROVAL_MSG" {
        pp.ReqData["approval_type"] = "1"
        pp.ReqData["approval_number"] = pp.approvalNumber
    }
    pp.ReqData["act_name"] = pp.actName
}

func (pp *pocketPay) SendRequest() api.APIResult {
    pp.checkData()

    conf := wx.NewConfig().GetCorp(pp.corpId)
    agentInfo := conf.GetAgentInfo(pp.agentTag)
    pp.ReqData["agentid"] = agentInfo["id"]
    workSign := wx.NewUtilWx().CreateCropSign(pp.ReqData, pp.acceptKeys, agentInfo["secret"])
    pp.ReqData["workwx_sign"] = workSign
    sign := wx.NewUtilWx().CreateCropPaySign(pp.ReqData, conf.GetPayKey())
    pp.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(pp.ReqData))

    pp.ReqURI = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/paywwsptrans2pocket"
    client, req := pp.GetRequest()
    req.SetBody(reqBody)

    keyFile, _ := ioutil.TempFile("", "tmpfile")
    defer os.Remove(keyFile.Name())
    if _, err := keyFile.Write([]byte(conf.GetSslKey())); err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "写入证书key文件失败", nil))
    }

    certFile, _ := ioutil.TempFile("", "tmpfile")
    defer os.Remove(certFile.Name())
    if _, err := certFile.Write([]byte(conf.GetSslCert())); err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "写入证书cert文件失败", nil))
    }

    certs, err := tls.LoadX509KeyPair(certFile.Name(), keyFile.Name())
    if err != nil {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "加载证书文件失败", nil))
    }
    client.TLSConfig.Certificates = []tls.Certificate{certs}

    resp, result := pp.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XMLMap)(&respData))
    if respData["return_code"] == "FAIL" {
        mplog.LogError(respData["return_msg"])
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["return_msg"]
    } else if respData["result_code"] == "FAIL" {
        mplog.LogError(respData["err_code"])
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["err_code_des"]
    } else {
        result.Data = respData
    }
    return result
}

func NewPocketPay(corpId, agentTag string) *pocketPay {
    conf := wx.NewConfig().GetCorp(corpId)
    pp := &pocketPay{wx.NewBaseWxCorp(), "", "", "", "", "", "", "", "", "", 0, "", "", "", "", 0, "", make([]string, 0)}
    pp.corpId = corpId
    pp.agentTag = agentTag
    pp.ReqData["appid"] = conf.GetCorpId()
    pp.ReqData["mch_id"] = conf.GetPayMchId()
    pp.ReqData["spbill_create_ip"] = conf.GetClientIp()
    pp.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    pp.checkName = "NO_CHECK"
    pp.msgType = "NORMAL_MSG"
    pp.amount = 0
    pp.acceptKeys = append(pp.acceptKeys, "amount", "appid", "desc", "mch_id", "nonce_str", "openid", "partner_trade_no", "ww_msg_type")
    pp.ReqContentType = project.HTTPContentTypeXML
    pp.ReqMethod = fasthttp.MethodPost
    return pp
}
