/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 9:17
 */
package profitsharing

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

// 查询分账结果
type sharingQuery struct {
    wx.BaseWxAccount
    appId         string
    transactionId string // 微信单号
    outOrderNo    string // 商户分账单号
}

func (sq *sharingQuery) SetTransactionId(transactionId string) {
    match, _ := regexp.MatchString(project.RegexDigit, transactionId)
    if match {
        sq.transactionId = transactionId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信单号不合法", nil))
    }
}

func (sq *sharingQuery) SetOutOrderNo(outOrderNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outOrderNo)
    if match {
        sq.outOrderNo = outOrderNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户分账单号不合法", nil))
    }
}

func (sq *sharingQuery) checkData() {
    if len(sq.transactionId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "微信单号不能为空", nil))
    }
    if len(sq.outOrderNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户分账单号不能为空", nil))
    }
    sq.ReqData["transaction_id"] = sq.transactionId
    sq.ReqData["out_order_no"] = sq.outOrderNo
}

func (sq *sharingQuery) SendRequest() api.ApiResult {
    sq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(sq.ReqData, sq.appId, "sha256")
    sq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(sq.ReqData))
    sq.ReqUrl = "https://api.mch.weixin.qq.com/pay/profitsharingquery"
    client, req := sq.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := sq.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewSharingQuery(appId, merchantType string) *sharingQuery {
    conf := wx.NewConfig().GetAccount(appId)
    sq := &sharingQuery{wx.NewBaseWxAccount(), "", "", ""}
    sq.appId = appId
    sq.SetPayAccount(conf, merchantType)
    delete(sq.ReqData, "sub_appid")
    delete(sq.ReqData, "appid")
    sq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    sq.ReqData["sign_type"] = "HMAC-SHA256"
    sq.ReqContentType = project.HTTPContentTypeXML
    sq.ReqMethod = fasthttp.MethodPost
    return sq
}
