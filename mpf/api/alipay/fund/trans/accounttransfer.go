/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/19 0019
 * Time: 15:13
 */
package fund

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/alipay"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 单笔转账到支付宝账户接口
type accountTransfer struct {
    alipay.BaseAliPay
    outBizNo     string  // 商户转账单号
    payeeType    string  // 收款方账户类型
    payeeAccount string  // 收款方账户
    amount       float32 // 转账金额,单位为元
}

func (at *accountTransfer) SetOutBizNo(outBizNo string) {
    match, _ := regexp.MatchString(project.RegexDigit, outBizNo)
    if match {
        at.outBizNo = outBizNo
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户转账单号不合法", nil))
    }
}

func (at *accountTransfer) SetPayeeType(payeeType string) {
    if (payeeType == "ALIPAY_USERID") || (payeeType == "ALIPAY_LOGONID") {
        at.payeeType = payeeType
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "账户类型不合法", nil))
    }
}

func (at *accountTransfer) SetPayeeAccount(payeeAccount string) {
    if len(payeeAccount) > 0 {
        at.payeeAccount = payeeAccount
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "账户不合法", nil))
    }
}

func (at *accountTransfer) SetAmount(amount float32) {
    if amount >= 0.1 {
        nowAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount), 64)
        at.amount = float32(nowAmount)
    } else {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "转账金额不合法", nil))
    }
}

func (at *accountTransfer) SetPayerShowName(payerShowName string) {
    at.BizContent["payer_show_name"] = payerShowName
}

func (at *accountTransfer) SetPayeeRealName(payeeRealName string) {
    at.BizContent["payee_real_name"] = payeeRealName
}

func (at *accountTransfer) SetRemark(remark string) {
    at.BizContent["remark"] = remark
}

func (at *accountTransfer) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(at.outBizNo) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "商户转账单号不能为空", nil))
    }
    if at.amount <= 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "转账金额不能为空", nil))
    }
    if len(at.payeeType) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "账户类型不能为空", nil))
    }
    if len(at.payeeAccount) == 0 {
        panic(mperr.NewAliPayFund(errorcode.AliPayFundParam, "账户不能为空", nil))
    }
    at.BizContent["out_biz_no"] = at.outBizNo
    at.BizContent["amount"] = at.amount
    at.BizContent["payee_type"] = at.payeeType
    at.BizContent["payee_account"] = at.payeeAccount

    return at.GetRequest()
}

func NewAccountTransfer(appId string) *accountTransfer {
    at := &accountTransfer{alipay.NewBase(appId), "", "", "", 0.00}
    at.amount = 0.00
    at.SetMethod("alipay.fund.trans.toaccount.transfer")
    return at
}
