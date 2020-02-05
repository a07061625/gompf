/**
 * 坐标转换
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 19:05
 */
package gaode

import (
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type coordConvert struct {
    mpmap.BaseGaoDe
    locations []string // 经纬度坐标
    coordsys  string   // 原坐标系
}

func (cc *coordConvert) SetLocations(locations []string) {
    if len(locations) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标不能为空", nil))
    } else if len(locations) > 40 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标不能超过40个", nil))
    }

    cc.ReqData["locations"] = strings.Join(locations, "|")
}

func (cc *coordConvert) SetCoordSys(coordSys string) {
    _, ok := coordConvertSysList[coordSys]
    if ok {
        cc.ReqData["coordsys"] = coordSys
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "原坐标系不合法", nil))
    }
}

func (cc *coordConvert) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(cc.locations) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标不能为空", nil))
    }

    return cc.GetRequest()
}

func NewCoordConvert() *coordConvert {
    cc := &coordConvert{mpmap.NewBaseGaoDe(), make([]string, 0), ""}
    cc.SetServiceUri("/assistant/coordinate/convert")
    cc.ReqData["coordsys"] = "autonavi"
    return cc
}
