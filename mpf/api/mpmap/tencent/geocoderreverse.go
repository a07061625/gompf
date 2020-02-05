/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 13:14
 */
package tencent

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type geoCoderReverse struct {
    mpmap.BaseTencent
    location   string            // 坐标
    poiStatus  int               // poi状态 0:不返回 1:返回
    poiOptions map[string]string // poi选项列表
}

func (gcr *geoCoderReverse) SetLocation(lat, lng string) {
    gcr.ReqData["location"] = lat + "," + lng
}

func (gcr *geoCoderReverse) SetPoiStatus(poiStatus int) {
    if (poiStatus == 0) || (poiStatus == 1) {
        gcr.ReqData["get_poi"] = strconv.Itoa(poiStatus)
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "poi状态不合法", nil))
    }
}

func (gcr *geoCoderReverse) SetPoiOptions(poiOptions map[string]string) {
    if len(poiOptions) > 0 {
        optionStr := ""
        for k, v := range poiOptions {
            optionStr += ";" + k + "=" + v
        }
        gcr.ReqData["poi_options"] = optionStr[1:]
    }
}

func (gcr *geoCoderReverse) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := gcr.ReqData["location"]
    if !ok {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "坐标不能为空", nil))
    }

    return gcr.GetRequest()
}

func NewGeoCoderReverse() *geoCoderReverse {
    gcr := &geoCoderReverse{mpmap.NewBaseTencent(), "", 0, make(map[string]string)}
    gcr.SetServiceUrl("https://apis.map.qq.com/ws/geocoder/v1/")
    gcr.SetRespTag("result")
    gcr.ReqData["get_poi"] = "0"
    return gcr
}
