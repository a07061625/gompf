/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 16:12
 */
package baidu

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type placeSuggestion struct {
    mpmap.BaseBaiDu
    keyword         string // 关键词
    region          string // 区域
    cityLimit       string // 区域限制标识
    location        string // 地址
    coordType       string // 坐标类型
    returnCoordType string // 返回的坐标类型
}

func (ps *placeSuggestion) SetKeyword(keyword string) {
    if len(keyword) > 0 {
        ps.ReqData["query"] = keyword
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "关键词不能为空", nil))
    }
}

func (ps *placeSuggestion) SetRegion(region string) {
    if len(region) > 0 {
        ps.ReqData["region"] = region
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "区域不能为空", nil))
    }
}

// cityLimit string 区域限制 false:不限制城市 true:限制城市
func (ps *placeSuggestion) SetCityLimit(cityLimit string) {
    if (cityLimit == "false") || (cityLimit == "true") {
        ps.ReqData["city_limit"] = cityLimit
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "区域限制不合法", nil))
    }
}

func (ps *placeSuggestion) SetLocation(lat, lng string) {
    ps.ReqData["location"] = lat + "," + lng
}

func (ps *placeSuggestion) SetCoordType(coordType int) {
    if (coordType > 0) || (coordType <= 4) {
        ps.ReqData["coord_type"] = strconv.Itoa(coordType)
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "坐标类型不合法", nil))
    }
}

func (ps *placeSuggestion) SetReturnCoordType(returnCoordType string) {
    if len(returnCoordType) > 0 {
        ps.ReqData["ret_coordtype"] = returnCoordType
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "返回坐标类型不能为空", nil))
    }
}

func (ps *placeSuggestion) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := ps.ReqData["query"]
    if !ok {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "关键词不能为空", nil))
    }
    _, ok = ps.ReqData["region"]
    if !ok {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "地区不能为空", nil))
    }

    return ps.GetRequest()
}

func NewPlaceSuggestion() *placeSuggestion {
    ps := &placeSuggestion{mpmap.NewBaseBaiDu(), "", "", "", "", "", ""}
    ps.SetServiceUri("/place/v2/suggestion")
    ps.ReqData["city_limit"] = "true"
    ps.ReqData["coord_type"] = "3"
    ps.ReqData["ret_coordtype"] = "bd09ll"
    ps.ReqData["timestamp"] = strconv.Itoa(time.Now().Second())
    return ps
}
