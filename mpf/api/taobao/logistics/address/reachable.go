/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 0:35
 */
package address

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 判定服务是否可达
type reachable struct {
    taobao.BaseTaoBao
    areaCode       string   // 区域编码
    address        string   // 详细地址
    partnerIdList  []string // 物流公司编码ID
    serviceType    int      // 服务编码
    sourceAreaCode string   // 发货地编码
}

func (r *reachable) SetAreaCode(areaCode string) {
    match, _ := regexp.MatchString(project.RegexDigit, areaCode)
    if match {
        r.ReqData["area_code"] = areaCode
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "区域编码不合法", nil))
    }
}

func (r *reachable) SetAddress(address string) {
    if len(address) > 0 {
        r.ReqData["address"] = address
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "详细地址不合法", nil))
    }
}

func (r *reachable) SetPartnerIdList(partnerIdList []string) {
    r.partnerIdList = make([]string, 0)
    for _, v := range partnerIdList {
        match, _ := regexp.MatchString(project.RegexDigit, v)
        if match {
            r.partnerIdList = append(r.partnerIdList, v)
        }
    }
}

func (r *reachable) SetServiceType(serviceType int) {
    if serviceType > 0 {
        r.serviceType = serviceType
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "服务编码不合法", nil))
    }
}

func (r *reachable) SetSourceAreaCode(sourceAreaCode string) {
    match, _ := regexp.MatchString(project.RegexDigit, sourceAreaCode)
    if match {
        r.ReqData["source_area_code"] = sourceAreaCode
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "发货地编码不合法", nil))
    }
}

func (r *reachable) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if r.serviceType <= 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "服务编码不能为空", nil))
    }
    if len(r.partnerIdList) == 0 {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "物流公司编码ID列表不能为空", nil))
    }
    r.ReqData["service_type"] = strconv.Itoa(r.serviceType)
    r.ReqData["partner_ids"] = strings.Join(r.partnerIdList, ",")

    return r.GetRequest()
}

func NewReachable() *reachable {
    r := &reachable{taobao.NewBaseTaoBao(), "", "", make([]string, 0), 0, ""}
    conf := logistics.NewConfigTaoBao()
    r.AppKey = conf.GetAppKey()
    r.AppSecret = conf.GetAppSecret()
    r.SetMethod("taobao.logistics.address.reachable")
    return r
}
