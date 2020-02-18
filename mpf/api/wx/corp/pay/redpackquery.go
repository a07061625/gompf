/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 13:00
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

// 查询红包记录
type redPackQuery struct {
    wx.BaseWxCorp
    corpId    string
    agentTag  string
    nonceStr  string // 随机字符串
    mchId     string // 商户号
    mchBillNo string // 商户订单号
}

func (rpq *redPackQuery) SetMchBillNo(mchBillNo string) {
    match, _ := regexp.MatchString(`^[0-9a-zA-Z]{1,32}$`, mchBillNo)
    if match {
        rpq.mchBillNo = mchBillNo
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不合法", nil))
    }
}

func (rpq *redPackQuery) checkData() {
    if len(rpq.mchBillNo) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不能为空", nil))
    }
    rpq.ReqData["mch_billno"] = rpq.mchBillNo
}

func (rpq *redPackQuery) SendRequest() api.APIResult {
    rpq.checkData()

    conf := wx.NewConfig().GetCorp(rpq.corpId)
    sign := wx.NewUtilWx().CreateCropPaySign(rpq.ReqData, conf.GetPayKey())
    rpq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(rpq.ReqData))

    rpq.ReqURI = "https://api.mch.weixin.qq.com/mmpaymkttransfers/queryworkwxredpack"
    client, req := rpq.GetRequest()
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

    resp, result := rpq.SendInner(client, req, errorcode.WxCorpRequestPost)
    if resp.RespCode > 0 {
        return result
    }

    respData := make(map[string]string)
    xml.Unmarshal(resp.Body, (*mpf.XMLMap)(&respData))
    if respData["return_code"] == "FAIL" {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["return_msg"]
    } else if respData["result_code"] == "FAIL" {
        result.Code = errorcode.WxCorpRequestPost
        result.Msg = respData["err_code_des"]
    } else {
        result.Data = respData
    }
    return result
}

func NewRedPackQuery(corpId, agentTag string) *redPackQuery {
    conf := wx.NewConfig().GetCorp(corpId)
    rpq := &redPackQuery{wx.NewBaseWxCorp(), "", "", "", "", ""}
    rpq.corpId = corpId
    rpq.agentTag = agentTag
    rpq.ReqData["appid"] = conf.GetCorpId()
    rpq.ReqData["mch_id"] = conf.GetPayMchId()
    rpq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    rpq.ReqContentType = project.HTTPContentTypeXML
    rpq.ReqMethod = fasthttp.MethodPost
    return rpq
}
