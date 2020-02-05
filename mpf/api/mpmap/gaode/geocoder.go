/**
 * 地理编码
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 22:47
 */
package gaode

import (
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type geoCoder struct {
    mpmap.BaseGaoDe
    address  string // 地址
    city     string // 城市
    batchTag string // 批量查询标识
}

func (gc *geoCoder) SetAddress(addressList []string) {
    if len(addressList) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "地址不能为空", nil))
    } else if len(addressList) > 10 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "地址不能超过10个", nil))
    }

    gc.ReqData["address"] = strings.Join(addressList, "|")
}

func (gc *geoCoder) SetCity(city string) {
    if len(city) > 0 {
        gc.ReqData["city"] = city
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "城市不合法", nil))
    }
}

func (gc *geoCoder) SetBatchTag(batchTag string) {
    if (batchTag == "true") || (batchTag == "false") {
        gc.ReqData["batch"] = batchTag
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "批量查询标识不合法", nil))
    }
}

func (gc *geoCoder) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := gc.ReqData["address"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "地址不能为空", nil))
    }

    return gc.GetRequest()
}

func NewGeoCoder() *geoCoder {
    gc := &geoCoder{mpmap.NewBaseGaoDe(), "", "", ""}
    gc.SetServiceUri("/geocode/geo")
    gc.ReqData["batch"] = "false"
    return gc
}
