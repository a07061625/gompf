/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/8 0008
 * Time: 12:31
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

// 查询付款记录
type pocketQuery struct {
    wx.BaseWxCorp
    corpId         string
    agentTag       string
    mchId          string // 商户号
    nonceStr       string // 随机字符串
    partnerTradeNo string // 商户订单号
}

func (pq *pocketQuery) SetPartnerTradeNo(partnerTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, partnerTradeNo)
    if match {
        pq.partnerTradeNo = partnerTradeNo
    } else {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不合法", nil))
    }
}

func (pq *pocketQuery) checkData() {
    if len(pq.partnerTradeNo) == 0 {
        panic(mperr.NewWxCorp(errorcode.WxCorpParam, "商户订单号不能为空", nil))
    }
    pq.ReqData["partner_trade_no"] = pq.partnerTradeNo
}

func (pq *pocketQuery) SendRequest() api.ApiResult {
    pq.checkData()

    conf := wx.NewConfig().GetCorp(pq.corpId)
    sign := wx.NewUtilWx().CreateCropPaySign(pq.ReqData, conf.GetPayKey())
    pq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(pq.ReqData))

    pq.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/querywwsptrans2pocket"
    client, req := pq.GetRequest()
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

    resp, result := pq.SendInner(client, req, errorcode.WxCorpRequestPost)
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

func NewPocketQuery(corpId, agentTag string) *pocketQuery {
    conf := wx.NewConfig().GetCorp(corpId)
    pq := &pocketQuery{wx.NewBaseWxCorp(), "", "", "", "", ""}
    pq.corpId = corpId
    pq.agentTag = agentTag
    pq.ReqData["appid"] = conf.GetCorpId()
    pq.ReqData["mch_id"] = conf.GetPayMchId()
    pq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    pq.ReqContentType = project.HTTPContentTypeXML
    pq.ReqMethod = fasthttp.MethodPost
    return pq
}
