/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/16 0016
 * Time: 19:08
 */
package kdbird

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mplogistics"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 实时查询物流轨迹
type expressInfo struct {
    mplogistics.BaseKdBird
    orderCode    string // 订单编号
    shipperCode  string // 快递公司编码
    logisticCode string // 物流单号
}

func (ei *expressInfo) SetOrderCode(orderCode string) {
    if len(orderCode) > 0 {
        ei.ExtendData["OrderCode"] = orderCode
    } else {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "订单编号不合法", nil))
    }
}

func (ei *expressInfo) SetShipperCode(shipperCode string) {
    if len(shipperCode) > 0 {
        ei.shipperCode = shipperCode
    } else {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "快递公司编码不合法", nil))
    }
}

func (ei *expressInfo) SetLogisticCode(logisticCode string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, logisticCode)
    if match {
        ei.logisticCode = logisticCode
    } else {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "物流单号不合法", nil))
    }
}

func (ei *expressInfo) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ei.shipperCode) == 0 {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "快递公司编码不能为空", nil))
    }
    if len(ei.logisticCode) == 0 {
        panic(mperr.NewLogisticsKdBird(errorcode.LogisticsKdBirdParam, "物流单号不能为空", nil))
    }
    ei.ExtendData["ShipperCode"] = ei.shipperCode
    ei.ExtendData["LogisticCode"] = ei.logisticCode

    return ei.GetRequest()
}

func NewExpressInfo() *expressInfo {
    ei := &expressInfo{mplogistics.NewBaseKdBird(), "", "", ""}
    ei.ReqData["RequestType"] = "1002"
    ei.ReqUrl = "http://api.kdniao.com/Ebusiness/EbusinessOrderHandle.aspx"
    return ei
}
