/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 9:42
 */
package tencent

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type coordTranslate struct {
    mpmap.BaseTencent
    coords   []string // 源坐标数组
    fromType uint     // 源坐标类型
}

func (ct *coordTranslate) AddCoord(lat, lng string) {
    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "源坐标经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "源坐标纬度不合法", nil))
    }
    ct.coords = append(ct.coords, lat+","+lng)
}

// fromType uint 坐标类型 1:GPS 2:搜狗 3:百度 4:mapbar 5:google 6:搜狗墨卡托
func (ct *coordTranslate) SetFromType(fromType uint) {
    if (fromType > 0) && (fromType <= 6) {
        ct.ReqData["type"] = strconv.Itoa(int(fromType))
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "源坐标类型不合法", nil))
    }
}

func (ct *coordTranslate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ct.coords) == 0 {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "源坐标不能为空", nil))
    }
    ct.ReqData["locations"] = strings.Join(ct.coords, ";")

    return ct.GetRequest()
}

func NewCoordTranslate() *coordTranslate {
    ct := &coordTranslate{mpmap.NewBaseTencent(), make([]string, 0), 0}
    ct.SetServiceUrl("https://apis.map.qq.com/ws/coord/v1/translate")
    ct.SetRespTag("locations")
    ct.ReqData["type"] = "5"
    return ct
}
