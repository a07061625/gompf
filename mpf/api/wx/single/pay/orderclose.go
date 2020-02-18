/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 12:28
 */
package pay

import (
    "encoding/xml"
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type orderClose struct {
    wx.BaseWxAccount
    appId      string
    outTradeNo string // 商户订单号
}

func (oc *orderClose) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outTradeNo)
    if match {
        oc.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (oc *orderClose) checkData() {
    if len(oc.outTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不能为空", nil))
    }
    oc.ReqData["out_trade_no"] = oc.outTradeNo
}

func (oc *orderClose) SendRequest() api.ApiResult {
    oc.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(oc.ReqData, oc.appId, "md5")
    oc.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(oc.ReqData))
    oc.ReqUrl = "https://api.mch.weixin.qq.com/pay/closeorder"
    client, req := oc.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := oc.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewOrderClose(appId, merchantType string) *orderClose {
    conf := wx.NewConfig().GetAccount(appId)
    oc := &orderClose{wx.NewBaseWxAccount(), "", ""}
    oc.appId = appId
    oc.SetPayAccount(conf, merchantType)
    oc.ReqData["sign_type"] = "MD5"
    oc.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    oc.ReqContentType = project.HTTPContentTypeXML
    oc.ReqMethod = fasthttp.MethodPost
    return oc
}
