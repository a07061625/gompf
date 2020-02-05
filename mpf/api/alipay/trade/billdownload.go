/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/20 0020
 * Time: 9:30
 */
package trade

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询对账单下载地址
type billDownload struct {
    alipay.BaseAliPay
    billType string // 账单类型
    billDate string // 账单时间：日账单格式为yyyy-MM-dd，月账单格式为yyyy-MM
}

func (bd *billDownload) SetBillType(billType string) {
    if (billType == "trade") || (billType == "signcustomer") {
        bd.billType = billType
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "账单类型不合法", nil))
    }
}

func (bd *billDownload) SetBillDate(billDate string) {
    match, _ := regexp.MatchString(`^\d{4}(-\d{2}){1,2}$`, billDate)
    if match {
        bd.billDate = billDate
    } else {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "账单时间不合法", nil))
    }
}

func (bd *billDownload) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(bd.billType) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "账单类型不能为空", nil))
    }
    if len(bd.billDate) == 0 {
        panic(mperr.NewAliPayTrade(errorcode.AliPayTradeParam, "账单时间不能为空", nil))
    }
    bd.BizContent["bill_type"] = bd.billType
    bd.BizContent["bill_date"] = bd.billDate

    return bd.GetRequest()
}

func NewBillDownload(appId string) *billDownload {
    bd := &billDownload{alipay.NewBase(appId), "", ""}
    bd.SetMethod("alipay.data.dataservice.bill.downloadurl.query")
    return bd
}
