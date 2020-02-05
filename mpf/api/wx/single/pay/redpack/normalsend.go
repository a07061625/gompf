/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 8:37
 */
package redpack

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

type normalSend struct {
    wx.BaseWxAccount
    appId        string
    mchBillNo    string // 商户订单号
    sendName     string // 商户名称
    reOpenid     string // 用户openid
    totalAmount  int    // 付款金额
    totalNum     int    // 红包发放总人数
    wishing      string // 红包祝福语
    actName      string // 活动名称
    remark       string // 备注
    sceneId      int    // 场景id
    riskInfo     string // 活动信息
    consumeMchId string // 资金授权商户号
}

func (ns *normalSend) SetMchBillNo(mchBillNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, mchBillNo)
    if match {
        ns.mchBillNo = mchBillNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户订单号不合法", nil))
    }
}

func (ns *normalSend) SetSendName(sendName string) {
    if len(sendName) > 0 {
        trueName := []rune(sendName)
        ns.sendName = string(trueName[:16])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户名称不合法", nil))
    }
}

func (ns *normalSend) SetReOpenid(reOpenid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, reOpenid)
    if match {
        ns.reOpenid = reOpenid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ns *normalSend) SetTotalAmount(totalAmount int) {
    if totalAmount > 0 {
        ns.totalAmount = totalAmount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额不合法", nil))
    }
}

func (ns *normalSend) SetTotalNum(totalNum int) {
    if totalNum > 0 {
        ns.totalNum = totalNum
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "红包发放总人数不合法", nil))
    }
}

func (ns *normalSend) SetWishing(wishing string) {
    if len(wishing) > 0 {
        trueWishing := []rune(wishing)
        ns.wishing = string(trueWishing[:64])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "红包祝福语不合法", nil))
    }
}

func (ns *normalSend) SetActName(actName string) {
    if len(actName) > 0 {
        trueName := []rune(actName)
        ns.actName = string(trueName[:16])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "活动名称不合法", nil))
    }
}

func (ns *normalSend) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        ns.remark = string(trueRemark[:128])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "备注不合法", nil))
    }
}

func (ns *normalSend) SetSceneId(sceneId int) {
    if (sceneId > 0) && (sceneId <= 8) {
        ns.sceneId = sceneId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "场景id不合法", nil))
    }
}

func (ns *normalSend) SetRiskInfo(riskInfo map[string]string) {
    infoStr := ""
    for k, v := range riskInfo {
        if len(v) == 0 {
            continue
        }
        if (k == "posttime") || (k == "mobile") || (k == "deviceid") || (k == "clientversion") {
            infoStr += "&" + k + "=" + v
        }
    }
    if len(infoStr) > 0 {
        ns.ReqData["risk_info"] = infoStr[1:]
    }
}

func (ns *normalSend) SetConsumeMchId(consumeMchId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, consumeMchId)
    if match {
        ns.ReqData["consume_mch_id"] = consumeMchId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "资金授权商户号不合法", nil))
    }
}

func (ns *normalSend) checkData() {
    if len(ns.mchBillNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户订单号不能为空", nil))
    }
    if len(ns.sendName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户名称不能为空", nil))
    }
    if len(ns.reOpenid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if ns.totalAmount <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额不能为空", nil))
    }
    if len(ns.wishing) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "红包祝福语不能为空", nil))
    }
    if len(ns.actName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "活动名称不能为空", nil))
    }
    if len(ns.remark) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "备注不能为空", nil))
    }
    if (ns.totalAmount < 100) || (ns.totalAmount > 20000) {
        if ns.sceneId <= 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "场景id不能为空", nil))
        }
    }
    ns.ReqData["mch_billno"] = ns.mchBillNo
    ns.ReqData["send_name"] = ns.sendName
    ns.ReqData["re_openid"] = ns.reOpenid
    ns.ReqData["total_amount"] = strconv.Itoa(ns.totalAmount)
    ns.ReqData["wishing"] = ns.wishing
    ns.ReqData["act_name"] = ns.actName
    ns.ReqData["remark"] = ns.remark
    ns.ReqData["total_num"] = strconv.Itoa(ns.totalNum)
    if ns.sceneId > 0 {
        ns.ReqData["scene_id"] = "PRODUCT_" + strconv.Itoa(ns.sceneId)
    }
}

func (ns *normalSend) SendRequest() api.ApiResult {
    ns.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(ns.ReqData, ns.appId, "md5")
    ns.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(ns.ReqData))
    ns.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
    client, req := ns.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(ns.appId)
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

    resp, result := ns.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewNormalSend(appId string) *normalSend {
    conf := wx.NewConfig().GetAccount(appId)
    ns := &normalSend{wx.NewBaseWxAccount(), "", "", "", "", 0, 0, "", "", "", 0, "", ""}
    ns.appId = appId
    ns.totalNum = 1
    ns.ReqData["wxappid"] = conf.GetAppId()
    ns.ReqData["mch_id"] = conf.GetPayMchId()
    ns.ReqData["client_ip"] = conf.GetClientIp()
    ns.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    ns.ReqContentType = project.HttpContentTypeXml
    ns.ReqMethod = fasthttp.MethodPost
    return ns
}
