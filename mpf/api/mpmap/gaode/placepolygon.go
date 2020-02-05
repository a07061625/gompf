/**
 * 多边形搜索
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 23:57
 */
package gaode

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type placePolygon struct {
    mpmap.BaseGaoDe
    polygon    string // 经纬度坐标对
    keywords   string // 关键字
    types      string // POI类型
    offset     uint   // 每页记录数
    page       uint   // 当前页数
    extensions string // 返回结果标识 base:返回基本地址信息 all:返回地址信息、附近POI、道路以及道路交叉口信息
}

func (pp *placePolygon) SetPolygon(polygonList []string) {
    if len(polygonList) <= 2 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标对不合法", nil))
    }
    pp.ReqData["polygon"] = strings.Join(polygonList, "|")
}

func (pp *placePolygon) SetKeywords(keywords []string) {
    keywordList := make([]string, 0)
    for _, v := range keywords {
        if len(v) > 0 {
            keywordList = append(keywordList, v)
        }
    }
    if len(keywordList) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "关键字不能为空", nil))
    }
    pp.ReqData["keywords"] = strings.Join(keywordList, "|")
}

func (pp *placePolygon) SetTypes(types []string) {
    typeList := make([]string, 0)
    for _, v := range types {
        match, _ := regexp.MatchString(`^[0-9]{6}$`, v)
        if match {
            typeList = append(typeList, v)
        }
    }
    if len(typeList) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "POI类型不能为空", nil))
    }
    pp.ReqData["types"] = strings.Join(typeList, "|")
}

func (pp *placePolygon) SetOffset(offset uint) {
    if (offset > 0) && (offset <= 25) {
        pp.ReqData["offset"] = strconv.Itoa(int(offset))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "每页记录数不合法", nil))
    }
}

func (pp *placePolygon) SetPage(page uint) {
    if (page > 0) && (page <= 100) {
        pp.ReqData["page"] = strconv.Itoa(int(page))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "当前页数不合法", nil))
    }
}

func (pp *placePolygon) SetExtensions(extensions string) {
    if (extensions == "base") || (extensions == "all") {
        pp.ReqData["extensions"] = extensions
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "返回结果标识不合法", nil))
    }
}

func (pp *placePolygon) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := pp.ReqData["polygon"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标对不能为空", nil))
    }

    return pp.GetRequest()
}

func NewPlacePolygon() *placePolygon {
    pp := &placePolygon{mpmap.NewBaseGaoDe(), "", "", "", 0, 0, ""}
    pp.SetServiceUri("/place/polygon")
    pp.ReqData["offset"] = "10"
    pp.ReqData["page"] = "1"
    pp.ReqData["extensions"] = "base"
    return pp
}
