/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 12:34
 */
package baidu

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type ipLocation struct {
    mpmap.BaseBaiDu
    ip              string // IP
    returnCoordType string // 返回坐标类型
}

func (il *ipLocation) SetIp(ip string) {
    match, _ := regexp.MatchString(project.RegexIp, "."+ip)
    if match {
        il.ReqData["ip"] = ip
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "ip不合法", nil))
    }
}

// returnCoordType string 空字符串:百度墨卡托 bd09ll:百度 gcj02:国测局
func (il *ipLocation) SetReturnCoordType(returnCoordType string) {
    if (returnCoordType == "") || (returnCoordType == "bd09ll") || (returnCoordType == "gcj02") {
        il.ReqData["coor"] = returnCoordType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "返回坐标类型不支持", nil))
    }
}

func (il *ipLocation) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := il.ReqData["ip"]
    if !ok {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "ip不能为空", nil))
    }

    return il.GetRequest()
}

func NewIpLocation() *ipLocation {
    il := &ipLocation{mpmap.NewBaseBaiDu(), "", ""}
    il.SetServiceUri("/location/ip")
    return il
}
