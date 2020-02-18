/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 10:45
 */
package pay

import (
    "encoding/xml"
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api"
    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/api/wx/single"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type unifiedOrder struct {
    wx.BaseWxAccount
    appId         string
    merchantType  string                 // 商户类型
    body          string                 // 商品描述
    detail        string                 // 商品详情
    attach        string                 // 附加数据
    outTradeNo    string                 // 商户订单号
    totalFee      int                    // 标价金额,单位为分
    startTime     string                 // 交易起始时间,格式为yyyyMMddHHmmss
    expireTime    string                 // 交易结束时间,格式为yyyyMMddHHmmss
    goodsTag      string                 // 商品标记,使用代金券或立减优惠功能时需要的参数
    tradeType     string                 // 交易类型
    openid        string                 // 用户标识 trade_type=JSAPI时(即公众号支付),此参数必传
    productId     string                 // 商品ID trade_type=NATIVE时(即扫码支付),此参数必传
    sceneInfo     map[string]interface{} // 场景信息,json格式
    profitSharing string                 // 服务商分账状态,默认不分账 Y:需要分账 N:不分账
    clientIp      string                 // 客户端IP
}

func (uo *unifiedOrder) SetBody(body string) {
    if len(body) > 0 {
        trueBody := []rune(body)
        uo.body = string(trueBody[:40])
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不能为空", nil))
    }
}

func (uo *unifiedOrder) SetAttach(attach string) {
    if len(attach) <= 127 {
        uo.ReqData["attach"] = attach
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "附加数据不合法", nil))
    }
}

func (uo *unifiedOrder) SetOutTradeNo(outTradeNo string) {
    match, _ := regexp.MatchString(`^[0-9]{1,32}$`, outTradeNo)
    if match {
        uo.outTradeNo = outTradeNo
        if uo.tradeType == single.TradeTypeNative {
            uo.productId = outTradeNo
        }
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不合法", nil))
    }
}

func (uo *unifiedOrder) SetTotalFee(totalFee int) {
    if totalFee > 0 {
        uo.totalFee = totalFee
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付金额不能小于0", nil))
    }
}

func (uo *unifiedOrder) SetOpenid(openid string) {
    match, _ := regexp.MatchString(project.RegexWxOpenid, openid)
    if match {
        uo.openid = openid
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不合法", nil))
    }
}

func (uo *unifiedOrder) SetDetail(detail string) {
    if len(detail) > 0 {
        uo.ReqData["detail"] = detail
    }
}

func (uo *unifiedOrder) SetStartTime(startTime string) {
    match, _ := regexp.MatchString(`^2[0-9]{13}$`, startTime)
    if match {
        uo.ReqData["time_start"] = startTime
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "交易起始时间不合法", nil))
    }
}

func (uo *unifiedOrder) SetExpireTime(expireTime string) {
    match, _ := regexp.MatchString(`^2[0-9]{13}$`, expireTime)
    if match {
        uo.ReqData["time_expire"] = expireTime
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "交易结束时间不合法", nil))
    }
}

func (uo *unifiedOrder) SetGoodsTag(goodsTag string) {
    if len(goodsTag) > 0 {
        uo.ReqData["goods_tag"] = goodsTag
    }
}

func (uo *unifiedOrder) SetSceneInfo(sceneInfo map[string]interface{}) {
    if len(sceneInfo) > 0 {
        uo.sceneInfo = sceneInfo
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "场景信息不合法", nil))
    }
}

func (uo *unifiedOrder) SetProfitSharing(profitSharing string) {
    if uo.merchantType != wx.AccountMerchantTypeSub {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "非服务商支付", nil))
    }

    if (profitSharing == "Y") || (profitSharing == "N") {
        uo.ReqData["profit_sharing"] = profitSharing
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "服务商分账状态不合法", nil))
    }
}

func (uo *unifiedOrder) SetClientIp(clientIp string) {
    match, _ := regexp.MatchString(project.RegexIP, "."+clientIp)
    if match {
        uo.clientIp = clientIp
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "客户端IP不合法", nil))
    }
}

func (uo *unifiedOrder) checkData() {
    if len(uo.body) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品名称不能为空", nil))
    }
    if len(uo.outTradeNo) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商户单号不能为空", nil))
    }
    if uo.totalFee <= 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "支付金额不能小于0", nil))
    }
    if uo.tradeType == single.TradeTypeNative {
        if len(uo.productId) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "商品ID不能为空", nil))
        }
        uo.ReqData["product_id"] = uo.productId
    } else if uo.tradeType == single.TradeTypeMobileWeb {
        if len(uo.clientIp) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "终端IP不能为空", nil))
        }
        if len(uo.sceneInfo) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "场景信息不能为空", nil))
        }
        uo.ReqData["spbill_create_ip"] = uo.clientIp
        uo.ReqData["scene_info"] = mpf.JSONMarshal(uo.sceneInfo)
    } else if uo.tradeType == single.TradeTypeJsApi {
        if len(uo.openid) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "用户openid不能为空", nil))
        }
        if uo.merchantType == wx.AccountMerchantTypeSelf {
            uo.ReqData["openid"] = uo.openid
        } else {
            uo.ReqData["sub_openid"] = uo.openid
        }
    }
    uo.ReqData["trade_type"] = uo.tradeType
    uo.ReqData["body"] = uo.body
    uo.ReqData["out_trade_no"] = uo.outTradeNo
    uo.ReqData["total_fee"] = strconv.Itoa(uo.totalFee)
}

func (uo *unifiedOrder) SendRequest() api.APIResult {
    uo.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(uo.ReqData, uo.appId, "md5")
    uo.ReqData["sign"] = sign
    reqBody, _ := xml.Marshal(mpf.XMLMap(uo.ReqData))
    uo.ReqURI = "https://api.mch.weixin.qq.com/pay/unifiedorder"
    client, req := uo.GetRequest()
    req.SetBody([]byte(reqBody))

    resp, result := uo.SendInner(client, req, errorcode.WxAccountRequestPost)
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

func NewUnifiedOrder(appId, tradeType, merchantType string) *unifiedOrder {
    conf := wx.NewConfig().GetAccount(appId)
    uo := &unifiedOrder{wx.NewBaseWxAccount(), "", "", "", "", "", "", 0, "", "", "", "", "", "", make(map[string]interface{}), "", ""}
    uo.appId = appId
    uo.tradeType = tradeType
    uo.merchantType = merchantType
    uo.totalFee = 0
    uo.SetPayAccount(conf, merchantType)
    uo.ReqData["notify_url"] = conf.GetPayNotifyUrl()
    uo.ReqData["fee_type"] = "CNY"
    uo.ReqData["nonce_str"] = mpf.ToolCreateNonceStr(32, "numlower")
    uo.ReqData["device_info"] = "WEB"
    uo.ReqData["sign_type"] = "MD5"
    if tradeType != single.TradeTypeMobileWeb {
        uo.ReqData["spbill_create_ip"] = conf.GetClientIp()
    }
    uo.ReqContentType = project.HTTPContentTypeXML
    uo.ReqMethod = fasthttp.MethodPost
    return uo
}
