/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 0:52
 */
package address

import (
    "strings"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询地址区域
type areasGet struct {
    taobao.BaseTaoBao
    fields []string // 返回字段列表
}

func (ag *areasGet) SetFields(fields []string) {
    if len(fields) > 0 {
        ag.fields = fields
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "返回字段列表不能为空", nil))
    }
}

func (ag *areasGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ag.fields) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "返回字段列表不能为空", nil))
    }
    ag.ReqData["fields"] = strings.Join(ag.fields, ",")

    return ag.GetRequest()
}

func NewAreasGet() *areasGet {
    ag := &areasGet{taobao.NewBaseTaoBao(), make([]string, 0)}
    conf := logistics.NewConfigTaoBao()
    ag.AppKey = conf.GetAppKey()
    ag.AppSecret = conf.GetAppSecret()
    ag.SetMethod("taobao.areas.get")
    return ag
}
