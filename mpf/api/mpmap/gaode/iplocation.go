/**
 * IP定位
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 23:15
 */
package gaode

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type ipLocation struct {
    mpmap.BaseGaoDe
    ip  string // IP
}

func (il *ipLocation) SetIp(ip string) {
    match, _ := regexp.MatchString(project.RegexIP, "."+ip)
    if match {
        il.ReqData["ip"] = ip
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "ip不合法", nil))
    }
}

func (il *ipLocation) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := il.ReqData["ip"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "ip不能为空", nil))
    }

    return il.GetRequest()
}

func NewIpLocation() *ipLocation {
    il := &ipLocation{mpmap.NewBaseGaoDe(), ""}
    il.SetServiceUri("/ip")
    return il
}
