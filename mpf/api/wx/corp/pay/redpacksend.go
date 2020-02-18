/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 13:06
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

// 发放企业红包
type redPackSend struct {
    wx.BaseWxCorp
    corpId              string
    agentTag            string
    nonceStr            string // 随机字符串
    mchId               string // 商户号
    mchBillNo           string // 商户订单号
    senderName          string // 发送者名称
    senderHeaderMediaId string // 发送者头像
    reOpenid            string // 用户openid
    totalAmount         int    // 付款金额
    wishing             string // 红包祝福语
    actName             string // 活动名称
    remark              string // 备注
    sceneId             int    // 场景id
    acceptKeys          []string
}

func (rps *redPackSend) SetMchBillNo(mchBillNo string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, mchBillNo)
    if match {
        rps.mchBillNo = mchBillNo
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不合法", nil))
    }
}

func (rps *redPackSend) SetSenderName(senderName string) {
    if len(senderName) > 0 {
        trueName := []rune(senderName)
        rps.ReqData["sender_name"] = string(trueName[:64])
        delete(rps.ReqData, "agentid")
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发送者名称不合法", nil))
    }
}

func (rps *redPackSend) SetSenderHeaderMediaId(senderHeaderMediaId string) {
    if len(senderHeaderMediaId) == 0 {
        delete(rps.ReqData, "sender_header_media_id")
    } else if len(senderHeaderMediaId) > 128 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "发送者头像不合法", nil))
    } else {
        rps.ReqData["sender_header_media_id"] = senderHeaderMediaId
    }
}

func (rps *redPackSend) SetReOpenid(reOpenid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, reOpenid)
    if match {
        rps.reOpenid = reOpenid
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不合法", nil))
    }
}

func (rps *redPackSend) SetTotalAmount(totalAmount int) {
    if totalAmount >= 100 {
        rps.totalAmount = totalAmount
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款金额不合法", nil))
    }
}

func (rps *redPackSend) SetWishing(wishing string) {
    if len(wishing) > 0 {
        trueWishing := []rune(wishing)
        rps.wishing = string(trueWishing[:64])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "红包祝福语不合法", nil))
    }
}

func (rps *redPackSend) SetActName(actName string) {
    if len(actName) > 0 {
        trueName := []rune(actName)
        rps.actName = string(trueName[:16])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "活动名称不合法", nil))
    }
}

func (rps *redPackSend) SetRemark(remark string) {
    if len(remark) > 0 {
        trueRemark := []rune(remark)
        rps.remark = string(trueRemark[:128])
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "备注不合法", nil))
    }
}

func (rps *redPackSend) SetSceneId(sceneId int) {
    if (sceneId >= 1) && (sceneId <= 8) {
        rps.sceneId = sceneId
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "场景id不合法", nil))
    }
}

func (rps *redPackSend) checkData() {
    if len(rps.mchBillNo) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不能为空", nil))
    }
    if len(rps.reOpenid) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "用户openid不能为空", nil))
    }
    if rps.totalAmount <= 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "付款金额不能为空", nil))
    }
    if len(rps.wishing) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "红包祝福语不能为空", nil))
    }
    if len(rps.actName) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "活动名称不能为空", nil))
    }
    if len(rps.remark) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "备注不能为空", nil))
    }
    if (rps.totalAmount >= 20000) && (rps.sceneId <= 0) {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "场景id不能为空", nil))
    }
    rps.ReqData["mch_billno"] = rps.mchBillNo
    rps.ReqData["re_openid"] = rps.reOpenid
    rps.ReqData["total_amount"] = strconv.Itoa(rps.totalAmount)
    rps.ReqData["wishing"] = rps.wishing
    rps.ReqData["act_name"] = rps.actName
    rps.ReqData["remark"] = rps.remark
    if rps.totalAmount >= 20000 {
        rps.ReqData["scene_id"] = "PRODUCT_" + strconv.Itoa(rps.sceneId)
    }
}

func (rps *redPackSend) SendRequest() api.ApiResult {
    rps.checkData()

    conf := wx.NewConfig().GetCorp(rps.corpId)
    agentInfo := conf.GetAgentInfo(rps.agentTag)
    workSign := wx.NewUtilWx().CreateCropSign(rps.ReqData, rps.acceptKeys, agentInfo["secret"])
    rps.ReqData["workwx_sign"] = workSign
    sign := wx.NewUtilWx().CreateCropPaySign(rps.ReqData, conf.GetPayKey())
    rps.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(rps.ReqData))

    rps.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendworkwxredpack"
    client, req := rps.GetRequest()
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

    resp, result := rps.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XmlMap)(&respData))
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

func NewRedPackSend(corpId, agentTag string) *redPackSend {
    conf := wx.NewConfig().GetCorp(corpId)
    agentInfo := conf.GetAgentInfo(agentTag)
    rps := &redPackSend{wx.NewBaseWxCorp(), "", "", "", "", "", "", "", "", 0, "", "", "", 0, make([]string, 0)}
    rps.corpId = corpId
    rps.agentTag = agentTag
    rps.ReqData["wxappid"] = conf.GetCorpId()
    rps.ReqData["mch_id"] = conf.GetPayMchId()
    rps.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    rps.ReqData["agentid"] = agentInfo["id"]
    rps.acceptKeys = append(rps.acceptKeys, "act_name", "mch_billno", "mch_id", "nonce_str", "re_openid", "total_amount", "wxappid")
    rps.ReqContentType = project.HTTPContentTypeXML
    rps.ReqMethod = fasthttp.MethodPost
    return rps
}
