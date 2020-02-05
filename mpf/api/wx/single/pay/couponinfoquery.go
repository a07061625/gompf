/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 9:06
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

type couponInfoQuery struct {
    wx.BaseWxAccount
    appId      string
    couponId   string // 代金券id
    openid     string // 用户openid
    stockId    string // 批次号
    opUserId   string // 操作员
    deviceInfo string // 设备号
}

func (ciq *couponInfoQuery) SetCouponId(couponId string) {
    match, _ := regexp.MatchString(project.RegexDigit, couponId)
    if match {
        ciq.couponId = couponId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "代金券id不合法", nil))
    }
}

func (ciq *couponInfoQuery) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        ciq.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (ciq *couponInfoQuery) SetStockId(stockId string) {
    match, _ := regexp.MatchString(`^[0-9]{1,64}$`, stockId)
    if match {
        ciq.stockId = stockId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "批次号不合法", nil))
    }
}

func (ciq *couponInfoQuery) SetOpUserId(opUserId string) {
    match, _ := regexp.MatchString(project.RegexDigit, opUserId)
    if match {
        ciq.ReqData["op_user_id"] = opUserId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "操作员不合法", nil))
    }
}

func (ciq *couponInfoQuery) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        ciq.ReqData["device_info"] = deviceInfo
    }
}

func (ciq *couponInfoQuery) checkData() {
    if len(ciq.couponId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "代金券id不能为空", nil))
    }
    if len(ciq.openid) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
    }
    if len(ciq.stockId) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "批次号不能为空", nil))
    }
    ciq.ReqData["coupon_id"] = ciq.couponId
    ciq.ReqData["openid"] = ciq.openid
    ciq.ReqData["stock_id"] = ciq.stockId
}

func (ciq *couponInfoQuery) SendRequest() api.ApiResult {
    ciq.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(ciq.ReqData, ciq.appId, "md5")
    ciq.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XmlMap(ciq.ReqData))
    ciq.ReqUrl = "https://api.mch.weixin.qq.com/mmpaymkttransfers/querycouponsinfo"
    client, req := ciq.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := ciq.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewCouponInfoQuery(appId string) *couponInfoQuery {
    conf := wx.NewConfig().GetAccount(appId)
    ciq := &couponInfoQuery{wx.NewBaseWxAccount(), "", "", "", "", "", ""}
    ciq.appId = appId
    ciq.ReqData["appid"] = conf.GetAppId()
    ciq.ReqData["mch_id"] = conf.GetPayMchId()
    ciq.ReqData["op_user_id"] = conf.GetPayMchId()
    ciq.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    ciq.ReqData["version"] = "1.0"
    ciq.ReqData["type"] = "XML"
    ciq.ReqContentType = project.HttpContentTypeXml
    ciq.ReqMethod = fasthttp.MethodPost
    return ciq
}
