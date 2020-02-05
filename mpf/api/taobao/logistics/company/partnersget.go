/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 9:21
 */
package company

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/taobao"
    "github.com/a07061625/gompf/mpf/api/taobao/logistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询支持起始地到目的地范围的物流公司
type partnersGet struct {
    taobao.BaseTaoBao
    sourceId       string // 揽货地地区码
    targetId       string // 派送地地区码
    serviceType    string // 服务类型 cod:货到付款 online:在线下单 offline:自己联系 limit:限时物流
    goodsValue     int    // 货物价格
    isNeedCarriage int    // 揽收资费标识 1:需要揽收资费 0:不需要揽收资费
}

func (pg *partnersGet) SetSourceId(sourceId string) {
    match, _ := regexp.MatchString(project.RegexDigit, sourceId)
    if match {
        pg.ReqData["source_id"] = sourceId
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "揽货地地区码不合法", nil))
    }
}

func (pg *partnersGet) SetTargetId(targetId string) {
    match, _ := regexp.MatchString(project.RegexDigit, targetId)
    if match {
        pg.ReqData["target_id"] = targetId
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "派送地地区码不合法", nil))
    }
}

func (pg *partnersGet) SetServiceType(serviceType string) {
    if (serviceType == "cod") || (serviceType == "online") || (serviceType == "offline") || (serviceType == "limit") {
        pg.ReqData["service_type"] = serviceType
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "服务类型不合法", nil))
    }
}

func (pg *partnersGet) SetGoodsValue(goodsValue float32) {
    if goodsValue >= 0 {
        pg.ReqData["goods_value"] = fmt.Sprintf("%.2f", goodsValue)
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "货物价格不合法", nil))
    }
}

func (pg *partnersGet) SetIsNeedCarriage(isNeedCarriage int) {
    if (isNeedCarriage == 0) || (isNeedCarriage == 1) {
        pg.ReqData["is_need_carriage"] = strconv.Itoa(isNeedCarriage)
    } else {
        panic(mperr.NewLogisticsTaoBao(errorcode.LogisticsTaoBaoParam, "揽收资费标识不合法", nil))
    }
}

func (pg *partnersGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return pg.GetRequest()
}

func NewPartnersGet() *partnersGet {
    pg := &partnersGet{taobao.NewBaseTaoBao(), "", "", "", 0, 0}
    conf := logistics.NewConfigTaoBao()
    pg.AppKey = conf.GetAppKey()
    pg.AppSecret = conf.GetAppSecret()
    pg.ReqData["source_id"] = "0"
    pg.ReqData["target_id"] = "0"
    pg.ReqData["goods_value"] = "0.00"
    pg.ReqData["is_need_carriage"] = "0"
    pg.SetMethod("taobao.logistics.partners.get")
    return pg
}
