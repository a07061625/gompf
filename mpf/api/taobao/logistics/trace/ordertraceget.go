/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 8:57
 */
package trace

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 流转信息查询
type orderTraceGet struct {
    taobao.BaseTaoBao
    companyCode string // 物流公司编码
    mailNo      string // 运单号
    cacheFlag   int    // 缓存状态 1:缓存 0:不缓存
}

func (otg *orderTraceGet) SetCompanyCode(companyCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, companyCode)
    if match {
        otg.companyCode = companyCode
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "物流公司编码不合法", nil))
    }
}

func (otg *orderTraceGet) SetMailNo(mailNo string) {
    if len(mailNo) > 0 {
        otg.mailNo = mailNo
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "运单号不合法", nil))
    }
}

func (otg *orderTraceGet) SetCacheFlag(cacheFlag int) {
    if (cacheFlag == 0) || (cacheFlag == 1) {
        otg.ReqData["cache"] = strconv.Itoa(cacheFlag)
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "缓存状态不合法", nil))
    }
}

func (otg *orderTraceGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(otg.companyCode) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "物流公司编码不能为空", nil))
    }
    if len(otg.mailNo) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "运单号不能为空", nil))
    }
    otg.ReqData["company_code"] = otg.companyCode
    otg.ReqData["mail_no"] = otg.mailNo

    return otg.GetRequest()
}

func NewOrderTraceGet() *orderTraceGet {
    otg := &orderTraceGet{taobao.NewBaseTaoBao(), "", "", 0}
    conf := logistics.NewConfigTaoBao()
    otg.AppKey = conf.GetAppKey()
    otg.AppSecret = conf.GetAppSecret()
    otg.ReqData["cache"] = "0"
    otg.SetMethod("taobao.logistics.ordertrace.get")
    return otg
}
