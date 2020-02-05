/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/16 0016
 * Time: 15:04
 */
package amali

import (
    "regexp"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mplogistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 快递物流节点跟踪
type expressInfo struct {
    mplogistics.BaseAMAli
    com           string // 公司字母简称
    nu            string // 快递单号
    receiverPhone string // 收件人手机号后四位(顺丰需要)
    senderPhone   string // 寄件人手机号后四位(顺丰需要)
}

func (ei *expressInfo) SetCom(com string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, com)
    if match {
        ei.com = com
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "公司字母简称不合法", nil))
    }
}

func (ei *expressInfo) SetNu(nu string) {
    if len(nu) > 0 {
        ei.nu = nu
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "快递单号不合法", nil))
    }
}

func (ei *expressInfo) SetReceiverPhone(receiverPhone string) {
    match, _ := regexp.MatchString(`^[0-9]{4}$`, receiverPhone)
    if match {
        ei.ReqData["receiverPhone"] = receiverPhone
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "收件人手机号后四位不合法", nil))
    }
}

func (ei *expressInfo) SetSenderPhone(senderPhone string) {
    match, _ := regexp.MatchString(`^[0-9]{4}$`, senderPhone)
    if match {
        ei.ReqData["senderPhone"] = senderPhone
    } else {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "寄件人手机号后四位不合法", nil))
    }
}

func (ei *expressInfo) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ei.com) == 0 {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "公司字母简称不能为空", nil))
    }
    if len(ei.nu) == 0 {
        panic(mperr.NewLogisticsAMAli(errorcode.LogisticsAMAliParam, "快递单号不能为空", nil))
    }
    ei.ReqData["com"] = ei.com
    ei.ReqData["nu"] = ei.nu
    ei.ServiceUri = "/showapi_expInfo?" + mpf.HttpCreateParams(ei.ReqData, "none", 1)

    return ei.GetRequest()
}

func NewExpressInfo() *expressInfo {
    ei := &expressInfo{mplogistics.NewBaseAMAli(), "", "", "", ""}
    return ei
}
