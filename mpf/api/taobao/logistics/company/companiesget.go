/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 1:18
 */
package company

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询物流公司信息
type companiesGet struct {
    taobao.BaseTaoBao
    fields        []string // 返回字段列表
    isRecommended int      // 推荐标识,默认为1 1:所有支持电话联系的物流公司 0:所有
    orderMode     string   // 推荐物流公司的下单方式 offline:电话联系/自己联系 online:在线下单 all: 既电话联系又在线下单
}

func (cg *companiesGet) SetFields(fields []string) {
    if len(fields) > 0 {
        cg.fields = fields
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "返回字段列表不能为空", nil))
    }
}

func (cg *companiesGet) SetIsRecommended(isRecommended int) {
    if (isRecommended == 0) || (isRecommended == 1) {
        cg.ReqData["is_recommended"] = strconv.Itoa(isRecommended)
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "推荐标识不合法", nil))
    }
}

func (cg *companiesGet) SetOrderMode(orderMode string) {
    if (orderMode == "offline") || (orderMode == "online") || (orderMode == "all") {
        cg.ReqData["order_mode"] = orderMode
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "下单方式不合法", nil))
    }
}

func (cg *companiesGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cg.fields) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "返回字段列表不能为空", nil))
    }
    cg.ReqData["fields"] = strings.Join(cg.fields, ",")

    return cg.GetRequest()
}

func NewCompaniesGet() *companiesGet {
    cg := &companiesGet{taobao.NewBaseTaoBao(), make([]string, 0), 0, ""}
    conf := logistics.NewConfigTaoBao()
    cg.AppKey = conf.GetAppKey()
    cg.AppSecret = conf.GetAppSecret()
    cg.ReqData["is_recommended"] = "1"
    cg.ReqData["order_mode"] = "offline"
    cg.SetMethod("taobao.logistics.companies.get")
    return cg
}
