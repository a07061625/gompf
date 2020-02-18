/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 9:23
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

type couponStockQuery struct {
    wx.BaseWxAccount
    appId      string
    stockId    string // 代金券批次id
    opUserId   string // 操作员
    deviceInfo string // 设备号
}

func (csq *couponStockQuery) SetStockId(stockId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,64}$`, stockId)
    if match {
        csq.stockId = stockId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "代金券批次id不合法", nil))
    }
}

func (csq *couponStockQuery) SetOpUserId(opUserId string) {
    match, _ := regexp.MatchString(project.RegexDigit, opUserId)
    if match {
        csq.ReqData["op_user_id"] = opUserId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "操作员不合法", nil))
    }
}

func (csq *couponStockQuery) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        csq.ReqData["device_info"] = deviceInfo
    }
}

func (csq *couponStockQuery) checkData() {
    if len(csq.stockId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "代金券批次id不能为空", nil))
    }
    csq.ReqData["coupon_stock_id"] = csq.stockId
}

func (csq *couponStockQuery) SendRequest() api.APIResult {
    csq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(csq.ReqData, csq.appId, "md5")
    csq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(csq.ReqData))
    csq.ReqURI = "https://api.mch.weixin.qq.com/mmpaymkttransfers/query_coupon_stock"
    client, req := csq.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := csq.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewCouponStockQuery(appId string) *couponStockQuery {
    conf := wx.NewConfig().GetAccount(appId)
    csq := &couponStockQuery{wx.NewBaseWxAccount(), "", "", "", ""}
    csq.appId = appId
    csq.ReqData["appid"] = conf.GetAppId()
    csq.ReqData["mch_id"] = conf.GetPayMchId()
    csq.ReqData["op_user_id"] = conf.GetPayMchId()
    csq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    csq.ReqData["version"] = "1.0"
    csq.ReqData["type"] = "XML"
    csq.ReqContentType = project.HTTPContentTypeXML
    csq.ReqMethod = fasthttp.MethodPost
    return csq
}
