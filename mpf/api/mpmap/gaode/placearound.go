/**
 * 周边搜索
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 23:28
 */
package gaode

import (
    "regexp"
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type placeAround struct {
    mpmap.BaseGaoDe
    location   string // 中心点坐标
    keywords   string // 关键字
    types      string // POI类型
    city       string // 城市
    radius     uint   // 搜索半径,单位:米
    sortRule   string // 排序规则 distance:距离排序 weight:综合排序
    offset     uint   // 每页记录数
    page       uint   // 当前页数
    extensions string // 返回结果标识 base:返回基本地址信息 all:返回地址信息、附近POI、道路以及道路交叉口信息
}

func (pa *placeAround) SetLocation(lat, lng string) {
    pa.ReqData["location"] = lng + "," + lat
}

func (pa *placeAround) SetKeywords(keywords []string) {
    keywordList := make([]string, 0)
    for _, v := range keywords {
        if len(v) > 0 {
            keywordList = append(keywordList, v)
        }
    }
    if len(keywordList) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "关键字不能为空", nil))
    }
    pa.ReqData["keywords"] = strings.Join(keywordList, "|")
}

func (pa *placeAround) SetTypes(types []string) {
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
    pa.ReqData["types"] = strings.Join(typeList, "|")
}

func (pa *placeAround) SetCity(city string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, city)
    if match {
        pa.ReqData["city"] = city
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "城市不合法", nil))
    }
}

func (pa *placeAround) SetRadius(radius uint) {
    if (radius > 0) && (radius <= 50000) {
        pa.ReqData["radius"] = strconv.Itoa(int(radius))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "搜索半径不合法", nil))
    }
}

func (pa *placeAround) SetSortRule(sortRule string) {
    if (sortRule == "distance") || (sortRule == "weight") {
        pa.ReqData["sortrule"] = sortRule
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "排序规则不合法", nil))
    }
}

func (pa *placeAround) SetOffset(offset uint) {
    if (offset > 0) && (offset <= 25) {
        pa.ReqData["offset"] = strconv.Itoa(int(offset))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "每页记录数不合法", nil))
    }
}

func (pa *placeAround) SetPage(page uint) {
    if (page > 0) && (page <= 100) {
        pa.ReqData["page"] = strconv.Itoa(int(page))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "当前页数不合法", nil))
    }
}

func (pa *placeAround) SetExtensions(extensions string) {
    if (extensions == "base") || (extensions == "all") {
        pa.ReqData["extensions"] = extensions
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "返回结果标识不合法", nil))
    }
}

func (pa *placeAround) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := pa.ReqData["location"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "中心点坐标不能为空", nil))
    }

    return pa.GetRequest()
}

func NewPlaceAround() *placeAround {
    pa := &placeAround{mpmap.NewBaseGaoDe(), "", "", "", "", 0, "", 0, 0, ""}
    pa.SetServiceUri("/place/around")
    pa.ReqData["radius"] = "3000"
    pa.ReqData["sortrule"] = "distance"
    pa.ReqData["offset"] = "10"
    pa.ReqData["page"] = "1"
    pa.ReqData["extensions"] = "base"
    return pa
}
