/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 14:30
 */
package baidu

import (
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type placeSearch struct {
    mpmap.BaseBaiDu
    keywords            []string // 检索关键字数组,最多支持10个关键字检索
    tags                []string // 标签数组
    scope               string   // 结果详细程度
    filter              string   // 过滤条件
    coordType           string   // 坐标类型
    pageSize            uint     // 每页条目数,最大限制为20条
    pageIndex           uint     // 页数,默认第1页
    searchType          string   // 搜索类型
    areaRegionName      string   // 区域搜索地区名称,市级以上行政区域
    areaRegionCityLimit string   // 区域搜索是否只返回指定region（城市）内的POI
    areaNearbyLng       string   // 圆形范围搜索中心点经度
    areaNearbyLat       string   // 圆形范围搜索中心点纬度
    areaNearbyRadius    uint     // 圆形范围搜索半径,单位为米
    areaRectangleLng1   string   // 矩形范围搜索西南角经度
    areaRectangleLat1   string   // 矩形范围搜索西南角纬度
    areaRectangleLng2   string   // 矩形范围搜索东北角经度
    areaRectangleLat2   string   // 矩形范围搜索东北角纬度
}

func (ps *placeSearch) AddKeyword(keyword string) {
    if len(ps.keywords) >= 10 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "检索关键字数量超过最大限制", nil))
    }

    trueKeyword := strings.TrimSpace(keyword)
    if len(trueKeyword) > 0 {
        ps.keywords = append(ps.keywords, trueKeyword)
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "检索关键字不合法", nil))
    }
}

func (ps *placeSearch) AddTag(tag string) {
    trueTag := strings.TrimSpace(tag)
    if len(trueTag) > 0 {
        ps.tags = append(ps.tags, trueTag)
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "标签不合法", nil))
    }
}

// scope string 结果详细程度 1:基本信息 2:POI详细信息
func (ps *placeSearch) SetScope(scope string) {
    if (scope == "1") || (scope == "2") {
        ps.ReqData["scope"] = scope
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "结果详细程度不合法", nil))
    }
}

func (ps *placeSearch) SetFilter(filter string) {
    trueFilter := strings.TrimSpace(filter)
    if len(trueFilter) > 0 {
        ps.ReqData["filter"] = trueFilter
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "过滤条件不能为空", nil))
    }
}

func (ps *placeSearch) SetCoordType(coordType string) {
    _, ok := placeSearchCoordTypes[coordType]
    if ok {
        ps.ReqData["coord_type"] = coordType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "坐标类型不合法", nil))
    }
}

func (ps *placeSearch) SetPageSize(pageSize uint) {
    if (pageSize > 0) && (pageSize <= 20) {
        ps.ReqData["page_size"] = strconv.Itoa(int(pageSize))
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "每页条目数只能在1-20之间", nil))
    }
}

func (ps *placeSearch) SetPageIndex(pageIndex uint) {
    if pageIndex > 0 {
        ps.ReqData["page_num"] = strconv.Itoa(int(pageIndex))
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "页数必须大于0", nil))
    }
}

// searchType string 区域搜索类型 region:地区 nearby:圆形区域 rectangle:矩形区域
func (ps *placeSearch) SetSearchType(searchType string) {
    _, ok := placeSearchTypes[searchType]
    if ok {
        ps.searchType = searchType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "区域搜索类型不支持", nil))
    }
}

func (ps *placeSearch) SetAreaRegionName(areaRegionName string) {
    trueName := strings.TrimSpace(areaRegionName)
    if len(trueName) > 0 {
        ps.areaRegionName = trueName
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "区域名称不能为空", nil))
    }
}

// areaRegionCityLimit string 城市限制标识 false:不限定城市 true:限定城市
func (ps *placeSearch) SetAreaRegionCityLimit(areaRegionCityLimit string) {
    if (areaRegionCityLimit == "false") || (areaRegionCityLimit == "true") {
        ps.areaRegionCityLimit = areaRegionCityLimit
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "城市限制标识不合法", nil))
    }
}

func (ps *placeSearch) SetAreaNearbyLngAndLat(lng string, lat string) {
    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "圆形范围搜索中心点经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "圆形范围搜索中心点纬度不合法", nil))
    }
    ps.areaNearbyLng = lng
    ps.areaNearbyLat = lat
}

func (ps *placeSearch) SetAreaNearbyRadius(areaNearbyRadius uint) {
    if areaNearbyRadius > 0 {
        ps.areaNearbyRadius = areaNearbyRadius
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "圆形范围搜索半径必须大于0", nil))
    }
}

// lng1 string 西南角经度
// lat1 string 西南角纬度
// lng2 string 东北角经度
// lat2 string 东北角纬度
func (ps *placeSearch) SetAreaRectangleLngAndLat(lng1, lat1, lng2, lat2 string) {
    match, _ := regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng1)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "矩形范围搜索西南角经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat1)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "矩形范围搜索西南角纬度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[-]?(\d(\.\d+)?|[1-9]\d(\.\d+)?|1[0-7]\d(\.\d+)?|180)$`, lng2)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "矩形范围搜索东北角经度不合法", nil))
    }
    match, _ = regexp.MatchString(`^[\-]?(\d(\.\d+)?|[1-8]\d(\.\d+)?|90)$`, lat2)
    if !match {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "矩形范围搜索东北角纬度不合法", nil))
    }

    num1, _ := strconv.ParseFloat(lat1, 64)
    num2, _ := strconv.ParseFloat(lat2, 64)
    if num1 >= num2 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "矩形范围搜索东北角纬度必须大于西南角纬度", nil))
    }

    ps.areaRectangleLng1 = lng1
    ps.areaRectangleLat1 = lat1
    ps.areaRectangleLng2 = lng2
    ps.areaRectangleLat2 = lat2
}

func (ps *placeSearch) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(ps.keywords) == 0 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "检索关键字不能为空", nil))
    }

    switch ps.searchType {
    case "region":
        if len(ps.areaRegionName) == 0 {
            panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "区域名称不能为空", nil))
        }
        ps.ReqData["region"] = ps.areaRegionName
        ps.ReqData["city_limit"] = ps.areaRegionCityLimit
    case "nearby":
        if len(ps.areaNearbyLng) == 0 {
            panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "中心点经度和纬度都不能为空", nil))
        }
        if ps.areaNearbyRadius <= 0 {
            panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "搜索半径必须大于0", nil))
        }
        ps.ReqData["location"] = ps.areaNearbyLat + "," + ps.areaNearbyLng
        ps.ReqData["radius"] = strconv.Itoa(int(ps.areaNearbyRadius))
    case "rectangle":
        if len(ps.areaRectangleLng1) == 0 {
            panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "矩形范围搜索经度和纬度都不能为空", nil))
        }
        ps.ReqData["bounds"] = ps.areaRectangleLat1 + "," + ps.areaRectangleLng1 + "," + ps.areaRectangleLat2 + "," + ps.areaRectangleLng2
    default:
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "区域搜索类型不能为空", nil))
    }
    ps.ReqData["query"] = strings.Join(ps.keywords, " ")
    if len(ps.tags) > 0 {
        ps.ReqData["tag"] = strings.Join(ps.tags, ",")
    }

    return ps.GetRequest()
}

func NewPlaceSearch() *placeSearch {
    ps := &placeSearch{mpmap.NewBaseBaiDu(), make([]string, 0), make([]string, 0), "", "", "", 0, 0, "", "", "", "", "", 0, "", "", "", ""}
    ps.SetServiceUri("/place/v2/search")
    ps.ReqData["scope"] = "1"
    ps.ReqData["coord_type"] = "3"
    ps.ReqData["page_size"] = "10"
    ps.ReqData["page_num"] = "0"
    ps.ReqData["timestamp"] = strconv.Itoa(time.Now().Second())
    ps.areaRegionCityLimit = "false"
    return ps
}
