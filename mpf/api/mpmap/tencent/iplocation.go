/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 13:26
 */
package tencent

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type ipLocation struct {
    mpmap.BaseTencent
    ip  string // IP
}

func (il *ipLocation) SetIp(ip string) {
    match, _ := regexp.MatchString(project.RegexIP, "."+ip)
    if match {
        il.ReqData["ip"] = ip
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "ip不合法", nil))
    }
}

func (il *ipLocation) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := il.ReqData["ip"]
    if !ok {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "ip不能为空", nil))
    }

    return il.GetRequest()
}

func NewIpLocation() *ipLocation {
    il := &ipLocation{mpmap.NewBaseTencent(), ""}
    il.SetServiceUrl("https://apis.map.qq.com/ws/location/v1/ip")
    il.SetRespTag("result")
    return il
}
