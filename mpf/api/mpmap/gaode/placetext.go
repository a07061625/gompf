/**
 * 关键字搜索
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 8:46
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

type placeText struct {
    mpmap.BaseGaoDe
    keywords   string // 关键字
    types      string // POI类型
    city       string // 城市
    cityLimit  string // 指定城市标识
    children   int    // 层级展示标识 0:子POI都会显示 1:子POI会归类到父POI之中
    offset     uint   // 每页记录数
    page       uint   // 当前页数
    building   string // POI编号
    floor      int    // 楼层
    extensions string // 返回结果标识 base:返回基本地址信息 all:返回地址信息、附近POI、道路以及道路交叉口信息
}

func (pt *placeText) SetKeywords(keywords []string) {
    keywordList := make([]string, 0)
    for _, v := range keywords {
        if len(v) > 0 {
            keywordList = append(keywordList, v)
        }
    }
    if len(keywordList) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "关键字不能为空", nil))
    }
    pt.keywords = strings.Join(keywordList, "|")
}

func (pt *placeText) SetTypes(types string) {
    match, _ := regexp.MatchString(`^[0-9]{6}$`, types)
    if match {
        pt.types = types
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "POI类型不合法", nil))
    }
}

func (pt *placeText) SetCity(city string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, city)
    if match {
        pt.ReqData["city"] = city
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "城市不合法", nil))
    }
}

// cityLimit string false:不限制城市 true:限制城市
func (pt *placeText) SetCityLimit(cityLimit string) {
    if (cityLimit == "false") || (cityLimit == "true") {
        pt.ReqData["citylimit"] = cityLimit
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "指定城市标识不合法", nil))
    }
}

func (pt *placeText) SetChildren(children int) {
    if (children == 0) || (children == 1) {
        pt.ReqData["children"] = strconv.Itoa(children)
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "层级展示标识不合法", nil))
    }
}

func (pt *placeText) SetOffset(offset uint) {
    if (offset > 0) && (offset <= 25) {
        pt.ReqData["offset"] = strconv.Itoa(int(offset))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "每页记录数不合法", nil))
    }
}

func (pt *placeText) SetPage(page uint) {
    if (page > 0) && (page <= 100) {
        pt.ReqData["page"] = strconv.Itoa(int(page))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "当前页数不合法", nil))
    }
}

func (pt *placeText) SetBuilding(building string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, building)
    if match {
        pt.ReqData["building"] = building
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "POI编号不合法", nil))
    }
}

func (pt *placeText) SetFloor(floor int) {
    pt.ReqData["floor"] = strconv.Itoa(floor)
}

func (pt *placeText) SetExtensions(extensions string) {
    if (extensions == "base") || (extensions == "all") {
        pt.ReqData["extensions"] = extensions
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "返回结果标识不合法", nil))
    }
}

func (pt *placeText) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if (len(pt.keywords) == 0) && (len(pt.types) == 0) {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "关键字和POI类型不能同时为空", nil))
    }
    if len(pt.keywords) > 0 {
        pt.ReqData["keywords"] = pt.keywords
    }
    if len(pt.types) > 0 {
        pt.ReqData["types"] = pt.types
    }

    return pt.GetRequest()
}

func NewPlaceText() *placeText {
    pt := &placeText{mpmap.NewBaseGaoDe(), "", "", "", "", 0, 0, 0, "", 0, ""}
    pt.SetServiceUri("/place/text")
    pt.ReqData["citylimit"] = "false"
    pt.ReqData["children"] = "0"
    pt.ReqData["offset"] = "10"
    pt.ReqData["page"] = "1"
    pt.ReqData["extensions"] = "base"
    return pt
}
