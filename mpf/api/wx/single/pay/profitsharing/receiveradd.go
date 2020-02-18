/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 10:32
 */
package profitsharing

import (
    "encoding/xml"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 添加分账接收方
type receiverAdd struct {
    wx.BaseWxAccount
    appId    string
    receiver map[string]interface{} // 分账接收方
}

func (ra *receiverAdd) SetReceiver(receiver map[string]interface{}) {
    if len(receiver) > 0 {
        ra.receiver = receiver
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账接收方不合法", nil))
    }
}

func (ra *receiverAdd) checkData() {
    if len(ra.receiver) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账接收方不能为空", nil))
    }
    ra.ReqData["receiver"] = mpf.JSONMarshal(ra.receiver)
}

func (ra *receiverAdd) SendRequest() api.APIResult {
    ra.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(ra.ReqData, ra.appId, "sha256")
    ra.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(ra.ReqData))
    ra.ReqURI = "https://api.mch.weixin.qq.com/pay/profitsharingaddreceiver"
    client, req := ra.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ra.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewReceiverAdd(appId, merchantType string) *receiverAdd {
    conf := wx.NewConfig().GetAccount(appId)
    ra := &receiverAdd{wx.NewBaseWxAccount(), "", make(map[string]interface{})}
    ra.appId = appId
    ra.SetPayAccount(conf, merchantType)
    ra.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    ra.ReqData["sign_type"] = "HMAC-SHA256"
    ra.ReqContentType = project.HTTPContentTypeXML
    ra.ReqMethod = fasthttp.MethodPost
    return ra
}
