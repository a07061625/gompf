/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 0:09
 */
package baidu

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type coordTranslate struct {
    mpmap.BaseBaiDu
    coords   []string // 源坐标数组
    fromType string   // 源坐标类型
    toType   string   // 目的坐标类型
}

func (ct *coordTranslate) AddCoord(lng string, lat string) {
    if len(ct.coords) >= 100 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "源坐标数量超过限制", nil))
    }

    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "源坐标经度不合法", nil))
    }

    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "源坐标纬度不合法", nil))
    }

    ct.coords = append(ct.coords, lng+","+lat)
}

func (ct *coordTranslate) SetFromType(fromType string) {
    _, ok := coordTranslateTypes[fromType]
    if ok {
        ct.ReqData["from"] = fromType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "源坐标类型不合法", nil))
    }
}

func (ct *coordTranslate) SetToType(toType string) {
    _, ok := coordTranslateTypes[toType]
    if ok {
        ct.ReqData["to"] = toType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "目的坐标类型不合法", nil))
    }
}

func (ct *coordTranslate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ct.coords) == 0 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "源坐标不能为空", nil))
    }
    ct.ReqData["coords"] = strings.Join(ct.coords, ";")

    return ct.GetRequest()
}

func NewCoordTranslate() *coordTranslate {
    ct := &coordTranslate{mpmap.NewBaseBaiDu(), make([]string, 0), "", ""}
    ct.SetServiceUri("/geoconv/v1/")
    ct.SetRespTag("result")
    ct.ReqData["from"] = "1"
    ct.ReqData["to"] = "5"
    return ct
}
