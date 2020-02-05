/**
 * 距离测量
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 22:37
 */
package gaode

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type distance struct {
    mpmap.BaseGaoDe
    origins      string // 出发点
    destination  string // 目的地
    distanceType int    // 距离计算类型 0:直线距离 1:驾车导航距离 2:公交规划距离 3:步行规划距离
}

func (d *distance) SetOrigins(origins []string) {
    if len(origins) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "出发点不能为空", nil))
    } else if len(origins) > 100 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "出发点不能超过100个", nil))
    }

    d.ReqData["origins"] = strings.Join(origins, "|")
}

func (d *distance) SetDestination(lat, lng string) {
    d.ReqData["destination"] = lng + "," + lat
}

func (d *distance) SetDistanceType(distanceType int) {
    if (distanceType >= 0) && (distanceType <= 3) {
        d.ReqData["type"] = strconv.Itoa(distanceType)
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "距离计算类型不合法", nil))
    }
}

func (d *distance) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := d.ReqData["origins"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "出发点不能为空", nil))
    }
    _, ok = d.ReqData["destination"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "目的地不能为空", nil))
    }

    return d.GetRequest()
}

func NewDistance() *distance {
    d := &distance{mpmap.NewBaseGaoDe(), "", "", 0}
    d.SetServiceUri("/distance")
    d.ReqData["type"] = "1"
    return d
}
