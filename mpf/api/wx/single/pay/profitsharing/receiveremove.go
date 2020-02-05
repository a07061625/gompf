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

// 删除分账接收方
type receiverRemove struct {
    wx.BaseWxAccount
    appId    string
    receiver map[string]interface{} // 分账接收方
}

func (rr *receiverRemove) SetReceiver(receiver map[string]interface{}) {
    if len(receiver) > 0 {
        rr.receiver = receiver
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账接收方不合法", nil))
    }
}

func (rr *receiverRemove) checkData() {
    if len(rr.receiver) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "分账接收方不能为空", nil))
    }
    rr.ReqData["receiver"] = mpf.JsonMarshal(rr.receiver)
}

func (rr *receiverRemove) SendRequest() api.ApiResult {
    rr.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(rr.ReqData, rr.appId, "sha256")
    rr.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(rr.ReqData))
    rr.ReqUrl = "https://api.mch.weixin.qq.com/pay/profitsharingremovereceiver"
    client, req := rr.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := rr.SendInner(client, req, errorcode.WxAccountRequestPost)
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
    }
    result.Data = respData
    return result
}

func NewReceiverRemove(appId, merchantType string) *receiverRemove {
    conf := wx.NewConfig().GetAccount(appId)
    rr := &receiverRemove{wx.NewBaseWxAccount(), "", make(map[string]interface{})}
    rr.appId = appId
    rr.SetPayAccount(conf, merchantType)
    rr.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    rr.ReqData["sign_type"] = "HMAC-SHA256"
    rr.ReqContentType = project.HttpContentTypeXml
    rr.ReqMethod = fasthttp.MethodPost
    return rr
}
