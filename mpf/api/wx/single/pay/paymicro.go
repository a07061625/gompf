/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 23:49
 */
package pay

import (
    "encoding/xml"
    "regexp"
    "strconv"

    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type payMicro struct {
    wx.BaseWxAccount
    appId      string
    deviceInfo string                 // 设备号
    body       string                 // 商品描述
    detail     string                 // 商品详情
    attach     string                 // 附加数据
    outTradeNo string                 // 商户订单号
    totalFee   int                    // 订单金额
    goodsTag   string                 // 商品标记
    limitPay   string                 // 指定支付方式
    startTime  int64                  // 交易起始时间
    expireTime int64                  // 交易结束时间
    receipt    string                 // 电子发票入口开放标识
    authCode   string                 // 授权码
    sceneInfo  map[string]interface{} // 场景信息
}

func (pm *payMicro) SetDeviceInfo(deviceInfo string) {
    if len(deviceInfo) > 0 {
        pm.ReqData["device_info"] = deviceInfo
    }
}

func (pm *payMicro) SetBody(body string) {
    if len(body) > 0 {
        trueBody := []rune(body)
        pm.body = string(trueBody[:40])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不能为空", nil))
    }
}

func (pm *payMicro) SetDetail(detail string) {
    if len(detail) > 0 {
        pm.ReqData["detail"] = detail
    }
}

func (pm *payMicro) SetAttach(attach string) {
    if len(attach) <= 127 {
        pm.ReqData["attach"] = attach
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "附加数据不合法", nil))
    }
}

func (pm *payMicro) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outTradeNo)
    if match {
        pm.outTradeNo = outTradeNo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (pm *payMicro) SetTotalFee(totalFee int) {
    if totalFee > 0 {
        pm.totalFee = totalFee
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付金额不能小于0", nil))
    }
}

func (pm *payMicro) SetGoodsTag(goodsTag string) {
    if len(goodsTag) > 0 {
        pm.ReqData["goods_tag"] = goodsTag
    }
}

func (pm *payMicro) SetLimitPay(limitPay string) {
    if (limitPay == "") || (limitPay == "no_credit") {
        pm.limitPay = limitPay
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "指定支付方式不合法", nil))
    }
}

func (pm *payMicro) SetTime(startTime, expireTime int64) {
    nowTime := time.Now().Unix()
    if startTime < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "交易起始时间不合法", nil))
    } else if expireTime < 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "交易结束时间不合法", nil))
    } else if (expireTime > 0) && (expireTime <= nowTime) {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "交易结束时间不能小于当前时间", nil))
    } else if (startTime > 0) && (expireTime > 0) && (startTime >= expireTime) {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "交易起始时间必须小于交易结束时间", nil))
    }
    pm.startTime = startTime
    pm.expireTime = expireTime
}

func (pm *payMicro) SetReceipt(receipt string) {
    if (receipt == "") || (receipt == "Y") {
        pm.receipt = receipt
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "电子发票入口开放标识不合法", nil))
    }
}

func (pm *payMicro) SetAuthCode(authCode string) {
    match, _ := regexp.MatchString(`^1[0-9]{17}$`, authCode)
    if match {
        pm.authCode = authCode
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不合法", nil))
    }
}

func (pm *payMicro) SetSceneInfo(sceneInfo map[string]interface{}) {
    pm.sceneInfo = sceneInfo
}

func (pm *payMicro) checkData() {
    if len(pm.body) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不能为空", nil))
    }
    if len(pm.outTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不能为空", nil))
    }
    if pm.totalFee <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付金额必须大于0", nil))
    }
    if len(pm.authCode) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "授权码不能为空", nil))
    }
    pm.ReqData["body"] = pm.body
    pm.ReqData["out_trade_no"] = pm.outTradeNo
    pm.ReqData["total_fee"] = strconv.Itoa(pm.totalFee)
    pm.ReqData["auth_code"] = pm.authCode
    if len(pm.limitPay) > 0 {
        pm.ReqData["limit_pay"] = pm.limitPay
    }
    if pm.startTime > 0 {
        st := time.Unix(int64(pm.startTime), 0)
        pm.ReqData["time_start"] = st.Format("20060102030405")
    }
    if pm.expireTime > 0 {
        et := time.Unix(int64(pm.expireTime), 0)
        pm.ReqData["time_expire"] = et.Format("20060102030405")
    }
    if len(pm.receipt) > 0 {
        pm.ReqData["receipt"] = pm.receipt
    }
    if len(pm.sceneInfo) > 0 {
        pm.ReqData["scene_info"] = mpf.JSONMarshal(pm.sceneInfo)
    }
}

func (pm *payMicro) SendRequest() api.ApiResult {
    pm.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(pm.ReqData, pm.appId, "md5")
    pm.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(pm.ReqData))
    pm.ReqUrl = "https://api.mch.weixin.qq.com/pay/micropay"
    client, req := pm.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := pm.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewPayMicro(appId string) *payMicro {
    conf := wx.NewConfig().GetAccount(appId)
    pm := &payMicro{wx.NewBaseWxAccount(), "", "", "", "", "", "", 0, "", "", 0, 0, "", "", make(map[string]interface{})}
    pm.appId = appId
    pm.totalFee = 0
    pm.ReqData["appid"] = conf.GetAppId()
    pm.ReqData["mch_id"] = conf.GetPayMchId()
    pm.ReqData["spbill_create_ip"] = conf.GetClientIp()
    pm.ReqData["sign_type"] = "MD5"
    pm.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    pm.ReqData["fee_type"] = "CNY"
    pm.ReqContentType = project.HTTPContentTypeXML
    pm.ReqMethod = fasthttp.MethodPost
    return pm
}
