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

type groupSend struct {
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

func (gs *groupSend) SetMchBillNo(mchBillNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, mchBillNo)
    if match {
        gs.mchBillNo = mchBillNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户订单号不合法", nil))
    }
}

func (gs *groupSend) SetSendName(sendName string) {
    if len(sendName) > 0 {
        trueName := []rune(sendName)
        gs.sendName = string(trueName[:16])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户名称不合法", nil))
    }
}

func (gs *groupSend) SetReOpenid(reOpenid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, reOpenid)
    if match {
        gs.reOpenid = reOpenid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (gs *groupSend) SetTotalAmount(totalAmount int) {
    if totalAmount > 0 {
        gs.totalAmount = totalAmount
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额不合法", nil))
    }
}

func (gs *groupSend) SetTotalNum(totalNum int) {
    if totalNum > 0 {
        gs.totalNum = totalNum
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "红包发放总人数不合法", nil))
    }
}

func (gs *groupSend) SetWishing(wishing string) {
    if len(wishing) > 0 {
        trueWishing := []rune(wishing)
        gs.wishing = string(trueWishing[:64])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "红包祝福语不合法", nil))
    }
}

func (gs *groupSend) SetActName(actName string) {
    if len(actName) > 0 {
        trueName := []rune(actName)
        gs.actName = string(trueName[:16])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "活动名称不合法", nil))
    }
}

func (gs *groupSend) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        gs.remark = string(trueRemark[:128])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "备注不合法", nil))
    }
}

func (gs *groupSend) SetSceneId(sceneId int) {
    if (sceneId > 0) && (sceneId <= 8) {
        gs.sceneId = sceneId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "场景id不合法", nil))
    }
}

func (gs *groupSend) SetRiskInfo(riskInfo map[string]string) {
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
        gs.ReqData["risk_info"] = infoStr[1:]
    }
}

func (gs *groupSend) SetConsumeMchId(consumeMchId string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, consumeMchId)
    if match {
        gs.ReqData["consume_mch_id"] = consumeMchId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "资金授权商户号不合法", nil))
    }
}

func (gs *groupSend) checkData() {
    if len(gs.mchBillNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户订单号不能为空", nil))
    }
    if len(gs.sendName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户名称不能为空", nil))
    }
    if len(gs.reOpenid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if gs.totalAmount <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "付款金额不能为空", nil))
    }
    if len(gs.wishing) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "红包祝福语不能为空", nil))
    }
    if len(gs.actName) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "活动名称不能为空", nil))
    }
    if len(gs.remark) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "备注不能为空", nil))
    }
    if (gs.totalAmount < 100) || (gs.totalAmount > 20000) {
        if gs.sceneId <= 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "场景id不能为空", nil))
        }
    }
    gs.ReqData["mch_billno"] = gs.mchBillNo
    gs.ReqData["send_name"] = gs.sendName
    gs.ReqData["re_openid"] = gs.reOpenid
    gs.ReqData["total_amount"] = strconv.Itoa(gs.totalAmount)
    gs.ReqData["wishing"] = gs.wishing
    gs.ReqData["act_name"] = gs.actName
    gs.ReqData["remark"] = gs.remark
    gs.ReqData["total_num"] = strconv.Itoa(gs.totalNum)
    if gs.sceneId > 0 {
        gs.ReqData["scene_id"] = "PRODUCT_" + strconv.Itoa(gs.sceneId)
    }
}

func (gs *groupSend) SendRequest() api.ApiResult {
    gs.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(gs.ReqData, gs.appId, "md5")
    gs.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(gs.ReqData))
    gs.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendgroupredpack"
    client, req := gs.GetRequest()
    req.SetBody([]byte(reqBody))

    conf := wx.NewConfig().GetAccount(gs.appId)
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

    resp, result := gs.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewGroupSend(appId string) *groupSend {
    conf := wx.NewConfig().GetAccount(appId)
    gs := &groupSend{wx.NewBaseWxAccount(), "", "", "", "", 0, 0, "", "", "", 0, "", ""}
    gs.appId = appId
    gs.totalNum = 1
    gs.ReqData["wxappid"] = conf.GetAppId()
    gs.ReqData["mch_id"] = conf.GetPayMchId()
    gs.ReqData["client_ip"] = conf.GetClientIp()
    gs.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    gs.ReqData["amt_type"] = "ALL_RAND"
    gs.ReqContentType = project.HTTPContentTypeXML
    gs.ReqMethod = fasthttp.MethodPost
    return gs
}
