/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 11:13
 */
package baidu

import (
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type geoCoder struct {
    mpmap.BaseBaiDu
    address         string // 地址
    cityName        string // 城市名
    returnCoordType string // 返回的坐标类型
}

func (gc *geoCoder) SetAddress(address string) {
    if len(address) > 0 {
        addressRune := []rune(address)
        gc.ReqData["address"] = string(addressRune[0:42])
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "地址不能为空", nil))
    }
}

func (gc *geoCoder) SetCityName(cityName string) {
    gc.ReqData["city"] = strings.TrimSpace(cityName)
}

func (gc *geoCoder) SetReturnCoordType(returnCoordType string) {
    gc.ReqData["ret_coordtype"] = strings.TrimSpace(returnCoordType)
}

func (gc *geoCoder) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := gc.ReqData["address"]
    if !ok {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "地址不能为空", nil))
    }

    return gc.GetRequest()
}

func NewGeoCoder() *geoCoder {
    gc := &geoCoder{mpmap.NewBaseBaiDu(), "", "", ""}
    gc.SetServiceUri("/geocoder/v2/")
    return gc
}
