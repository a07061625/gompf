/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 12:18
 */
package baidu

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type geoCoderReverse struct {
    mpmap.BaseBaiDu
    location        string // 坐标地址
    coordType       string // 坐标类型
    returnCoordType string // 返回的坐标类型
    poiStatus       int    // poi召回状态,0为不召回,1为召回
    poiRadius       uint   // poi召回半径,单位为米
}

func (gc *geoCoderReverse) SetLocation(lat string, lng string) {
    gc.ReqData["location"] = lat + "," + lng
}

func (gc *geoCoderReverse) SetCoordType(coordType string) {
    gc.ReqData["coordtype"] = coordType
    gc.ReqData["ret_coordtype"] = coordType
}

func (gc *geoCoderReverse) SetReturnCoordType(returnCoordType string) {
    gc.ReqData["ret_coordtype"] = returnCoordType
}

func (gc *geoCoderReverse) SetPoiStatusAndRadius(poiStatus int, poiRadius uint) {
    if (poiStatus < 0) || (poiStatus > 1) {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "poi召回状态不合法", nil))
    } else if (poiRadius < 0) || (poiRadius > 1000) {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "poi召回半径不合法", nil))
    }

    gc.ReqData["pois"] = strconv.Itoa(poiStatus)
    gc.ReqData["radius"] = strconv.Itoa(int(poiRadius))
}

func (gc *geoCoderReverse) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := gc.ReqData["location"]
    if !ok {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "坐标地址不能为空", nil))
    }

    return gc.GetRequest()
}

func NewGeoCoderReverse() *geoCoderReverse {
    gc := &geoCoderReverse{mpmap.NewBaseBaiDu(), "", "", "", 0, 0}
    gc.SetServiceUri("/geocoder/v2/")
    return gc
}
