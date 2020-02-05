/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/14 0014
 * Time: 11:36
 */
package pay

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/wx"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
)

type payNativeReturn struct {
    wx.BaseWxAccount
    appId    string
    nonceStr string // 微信返回的随机字符串
    prepayId string // 预支付ID
}

func (pnr *payNativeReturn) SetNonceStr(nonceStr string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, nonceStr)
    if match {
        pnr.nonceStr = nonceStr
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "随机字符串不合法", nil))
    }
}

func (pnr *payNativeReturn) SetPrepayId(prepayId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, prepayId)
    if match {
        pnr.nonceStr = prepayId
    } else {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "预支付ID不合法", nil))
    }
}

func (pnr *payNativeReturn) SetErrorMsg(returnMsg, resultMsg string) {
    if len(returnMsg) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "返回信息不能为空", nil))
    }
    if len(resultMsg) == 0 {
        panic(mperr.NewWxAccount(errorcode.WxAccountParam, "错误描述不能为空", nil))
    }
    trueReturn := []rune(returnMsg)
    trueResult := []rune(resultMsg)
    pnr.ReqData["return_code"] = "FAIL"
    pnr.ReqData["return_msg"] = string(trueReturn[:40])
    pnr.ReqData["result_code"] = "FAIL"
    pnr.ReqData["err_code_des"] = string(trueResult[:40])
}

func (pnr *payNativeReturn) checkData() {
    if pnr.ReqData["return_code"] == "SUCCESS" {
        if len(pnr.nonceStr) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "随机字符串不能为空", nil))
        }
        if len(pnr.prepayId) == 0 {
            panic(mperr.NewWxAccount(errorcode.WxAccountParam, "预支付ID不能为空", nil))
        }
        pnr.ReqData["nonce_str"] = pnr.nonceStr
        pnr.ReqData["prepay_id"] = pnr.prepayId
    }
}

func (pnr *payNativeReturn) GetResult() map[string]string {
    pnr.checkData()

    sign := wx.NewUtilWx().CreateSinglePaySign(pnr.ReqData, pnr.appId, "md5")
    pnr.ReqData["sign"] = sign
    return pnr.ReqData
}

func NewPayNativeReturn(appId string) *payNativeReturn {
    conf := wx.NewConfig().GetAccount(appId)
    pnr := &payNativeReturn{wx.NewBaseWxAccount(), "", "", ""}
    pnr.appId = appId
    pnr.ReqData["appid"] = conf.GetAppId()
    pnr.ReqData["mch_id"] = conf.GetPayMchId()
    pnr.ReqData["result_code"] = "SUCCESS"
    pnr.ReqData["return_code"] = "SUCCESS"
    return pnr
}
