/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/16 0016
 * Time: 18:39
 */
package kd100

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mplogistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type expressQuery struct {
    mplogistics.BaseKd100
    com      string // 快递公司编码
    num      string // 快递单号
    phone    string // 手机号码
    from     string // 出发地
    to       string // 目的地
    resultV2 int    // 行政区域解析开通状态
}

func (eq *expressQuery) SetCom(com string) {
    match, _ := regexp.MatchString(project.RegexAlpha, com)
    if match {
        eq.com = com
    } else {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "快递公司编码不合法", nil))
    }
}

func (eq *expressQuery) SetNum(num string) {
    if (len(num) > 0) && (len(num) <= 32) {
        eq.num = num
    } else {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "快递单号不合法", nil))
    }
}

func (eq *expressQuery) SetPhone(phone string) {
    match, _ := regexp.MatchString(project.RegexPhone, phone)
    if match {
        eq.ExtendData["phone"] = phone
    } else {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "手机号码不合法", nil))
    }
}

func (eq *expressQuery) SetFrom(from string) {
    eq.ExtendData["from"] = strings.TrimSpace(from)
}

func (eq *expressQuery) SetTo(to string) {
    eq.ExtendData["to"] = strings.TrimSpace(to)
}

func (eq *expressQuery) SetResultV2(resultV2 int) {
    if (resultV2 == 0) || (resultV2 == 1) {
        eq.ExtendData["resultv2"] = resultV2
    } else {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "行政区域解析开通状态不合法", nil))
    }
}

func (eq *expressQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(eq.com) == 0 {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "快递公司编码不能为空", nil))
    }
    if len(eq.num) == 0 {
        panic(mperr.NewLogisticsKd100(errorcode.LogisticsKd100Param, "快递单号不能为空", nil))
    }
    eq.ExtendData["com"] = eq.com
    eq.ExtendData["num"] = eq.num

    return eq.GetRequest()
}

func NewExpressQuery() *expressQuery {
    eq := &expressQuery{mplogistics.NewBaseKd100(), "", "", "", "", "", 0}
    eq.ExtendData["resultv2"] = 0
    return eq
}
