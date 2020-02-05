/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 16:27
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

type placeSearch struct {
    mpmap.BaseTencent
    keyword              string // 搜索关键字
    pageSize             uint   // 每页条目数,最大限制为20条
    pageIndex            uint   // 页数,默认第1页
    filter               string // 筛选条件
    orderBy              string // 排序方式
    searchType           string // 搜索类型
    areaRegionCityName   string // 区域搜索城市名称
    areaRegionAutoExtend int    // 区域搜索是否自动扩大范围 0:仅在当前城市搜索 1:若当前城市搜索无结果,则自动扩大范围
    areaRegionLng        string // 区域搜索经度
    areaRegionLat        string // 区域搜索纬度
    areaNearbyLng        string // 圆形范围搜索中心点经度
    areaNearbyLat        string // 圆形范围搜索中心点纬度
    areaNearbyRadius     uint   // 圆形范围搜索半径,单位为米
    areaRectangleLng1    string // 矩形范围搜索西南角经度
    areaRectangleLat1    string // 矩形范围搜索西南角纬度
    areaRectangleLng2    string // 矩形范围搜索东北角经度
    areaRectangleLat2    string // 矩形范围搜索东北角纬度
}

func (ps *placeSearch) SetKeyword(keyword string) {
    trueWord := strings.TrimSpace(keyword)
    if len(trueWord) > 0 {
        ps.ReqData["keyword"] = trueWord
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "搜索关键字不能为空", nil))
    }
}

func (ps *placeSearch) SetPageSize(pageSize uint) {
    if (pageSize > 0) && (pageSize <= 20) {
        ps.ReqData["page_size"] = strconv.Itoa(int(pageSize))
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "每页条目数只能在1-20之间", nil))
    }
}

func (ps *placeSearch) SetPageIndex(pageIndex uint) {
    if pageIndex > 0 {
        ps.ReqData["page_index"] = strconv.Itoa(int(pageIndex))
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "页数必须大于0", nil))
    }
}

func (ps *placeSearch) SetFilter(filter string) {
    trueFilter := strings.TrimSpace(filter)
    if len(trueFilter) > 0 {
        ps.ReqData["filter"] = trueFilter
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "筛选条件不能为空", nil))
    }
}

func (ps *placeSearch) SetOrderBy(field string, isAsc bool) {
    trueField := strings.TrimSpace(field)
    if len(trueField) > 0 {
        if isAsc {
            ps.ReqData["orderby"] = trueField + " asc"
        } else {
            ps.ReqData["orderby"] = trueField + " desc"
        }
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "排序字段不能为空", nil))
    }
}

func (ps *placeSearch) SetAreaRegionCityName(areaRegionCityName string) {
    cityName := strings.TrimSpace(areaRegionCityName)
    if len(cityName) > 0 {
        ps.areaRegionCityName = cityName
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域搜索城市名称不能为空", nil))
    }
}

func (ps *placeSearch) SetAreaRegionAutoExtend(areaRegionAutoExtend int) {
    if (areaRegionAutoExtend == 0) || (areaRegionAutoExtend == 1) {
        ps.areaRegionAutoExtend = areaRegionAutoExtend
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域搜索自动扩大标识不合法", nil))
    }
}

func (ps *placeSearch) SetAreaRegionLngAndLat(lat, lng string) {
    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域搜索经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域搜索纬度不合法", nil))
    }

    ps.areaRegionLat = lat
    ps.areaRegionLng = lng
}

func (ps *placeSearch) SetAreaNearbyLngAndLat(lat, lng string) {
    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "圆形范围搜索中心点经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "圆形范围搜索中心点纬度不合法", nil))
    }

    ps.areaNearbyLat = lat
    ps.areaNearbyLng = lng
}

func (ps *placeSearch) SetAreaNearbyRadius(areaNearbyRadius uint) {
    if areaNearbyRadius > 0 {
        ps.areaNearbyRadius = areaNearbyRadius
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "圆形范围搜索半径必须大于0", nil))
    }
}

func (ps *placeSearch) SetAreaRectangleLngAndLat(lat1, lng1, lat2, lng2 string) {
    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng1)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "矩形范围搜索西南角经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat1)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "矩形范围搜索西南角纬度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng2)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "矩形范围搜索东北角经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat2)
    if !match {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "矩形范围搜索东北角纬度不合法", nil))
    }

    num1, _ := strconv.ParseFloat(lat1, 64)
    num2, _ := strconv.ParseFloat(lat2, 64)
    if num1 >= num2 {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "矩形范围搜索东北角纬度必须大于西南角纬度", nil))
    }

    ps.areaRectangleLng1 = lng1
    ps.areaRectangleLat1 = lat1
    ps.areaRectangleLng2 = lng2
    ps.areaRectangleLat2 = lat2
}

func (ps *placeSearch) SetSearchType(searchType string) {
    _, ok := placeSearchTypes[searchType]
    if ok {
        ps.searchType = searchType
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域搜索类型不支持", nil))
    }
}

func (ps *placeSearch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := ps.ReqData["keyword"]
    if !ok {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "搜索关键字不能为空", nil))
    }

    switch ps.searchType {
    case PlaceSearchTypeRegion:
        if len(ps.areaRegionCityName) == 0 {
            panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域名称不能为空", nil))
        }
        if (len(ps.areaRegionLng) > 0) && (len(ps.areaRegionLat) > 0) {
            ps.ReqData["boundary"] = "region(" + ps.areaRegionCityName + "," + strconv.Itoa(ps.areaRegionAutoExtend) + "," + ps.areaRegionLat + "," + ps.areaRegionLng + ")"
        } else {
            ps.ReqData["boundary"] = "region(" + ps.areaRegionCityName + "," + strconv.Itoa(ps.areaRegionAutoExtend) + ")"
        }
    case PlaceSearchTypeNearby:
        if (len(ps.areaNearbyLat) == 0) || (len(ps.areaNearbyLng) == 0) {
            panic(mperr.NewMapTencent(errorcode.MapTencentParam, "中心点经度和纬度都不能为空", nil))
        }
        if ps.areaNearbyRadius <= 0 {
            panic(mperr.NewMapTencent(errorcode.MapTencentParam, "搜索半径必须大于0", nil))
        }
        ps.ReqData["boundary"] = "nearby(" + ps.areaNearbyLat + "," + ps.areaNearbyLng + "," + strconv.Itoa(int(ps.areaNearbyRadius)) + ")"
    case PlaceSearchTypeRectangle:
        if len(ps.areaRectangleLng1) == 0 {
            panic(mperr.NewMapTencent(errorcode.MapTencentParam, "矩形范围搜索经度和纬度都不能为空", nil))
        }
        ps.ReqData["boundary"] = "rectangle(" + ps.areaRectangleLat1 + "," + ps.areaRectangleLng1 + "," + ps.areaRectangleLat2 + "," + ps.areaRectangleLng2 + ")"
    default:
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "区域搜索类型不能为空", nil))
    }

    return ps.GetRequest()
}

func NewPlaceSearch() *placeSearch {
    ps := &placeSearch{mpmap.NewBaseTencent(), "", 0, 0, "", "", "", "", 0, "", "", "", "", 0, "", "", "", ""}
    ps.SetServiceUrl("https://apis.map.qq.com/ws/place/v1/search")
    ps.ReqData["page_size"] = "10"
    ps.ReqData["page_index"] = "1"
    ps.areaRegionAutoExtend = 1
    return ps
}
