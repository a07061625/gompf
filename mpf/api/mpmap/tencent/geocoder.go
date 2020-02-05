/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 13:08
 */
package tencent

import (
    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type geoCoder struct {
    mpmap.BaseTencent
    address string // 地址
    region  string // 地区
}

func (gc *geoCoder) SetAddress(address string) {
    if len(address) > 0 {
        gc.ReqData["address"] = address
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "地址不能为空", nil))
    }
}

func (gc *geoCoder) SetRegion(region string) {
    gc.ReqData["region"] = region
}

func (gc *geoCoder) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := gc.ReqData["address"]
    if !ok {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "地址不能为空", nil))
    }

    return gc.GetRequest()
}

func NewGeoCoder() *geoCoder {
    gc := &geoCoder{mpmap.NewBaseTencent(), "", ""}
    gc.SetServiceUrl("https://apis.map.qq.com/ws/geocoder/v1/")
    gc.SetRespTag("result")
    return gc
}
