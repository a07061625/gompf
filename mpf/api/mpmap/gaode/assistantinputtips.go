/**
 * 输入提示
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 18:35
 */
package gaode

import (
    "regexp"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type assistantInputTips struct {
    mpmap.BaseGaoDe
    keywords  string   // 关键字
    poiType   []string // POI类型
    location  string   // 坐标
    city      string   // 城市
    cityLimit string   // 指定城市标识
    dataType  string   // 返回数据类型
}

func (ait *assistantInputTips) SetKeyword(keyword string) {
    if len(keyword) > 0 {
        ait.ReqData["keywords"] = keyword
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "关键字不合法", nil))
    }
}

func (ait *assistantInputTips) SetPoiType(poiType []string) {
    for _, v := range poiType {
        match, _ := regexp.MatchString(`^[0-9]{6}$`, v)
        if match {
            ait.poiType = append(ait.poiType, v)
        }
    }
    if len(ait.poiType) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "POI类型不能为空", nil))
    }
    ait.ReqData["type"] = strings.Join(ait.poiType, "|")
}

func (ait *assistantInputTips) SetLocation(lat, lng string) {
    ait.ReqData["location"] = lng + "," + lat
}

func (ait *assistantInputTips) SetCity(city string) {
    match, _ := regexp.MatchString(project.RegexDigitLower, city)
    if match {
        ait.ReqData["city"] = city
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "城市不合法", nil))
    }
}

// cityLimit string false:不限制城市 true:限制城市
func (ait *assistantInputTips) SetCityLimit(cityLimit string) {
    if (cityLimit == "false") || (cityLimit == "true") {
        ait.ReqData["citylimit"] = cityLimit
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "指定城市标识不合法", nil))
    }
}

func (ait *assistantInputTips) SetDataType(dataTypes []string) {
    trueTypes := make([]string, 0)
    for _, v := range dataTypes {
        _, ok := inputTipDataTypes[v]
        if ok {
            trueTypes = append(trueTypes, v)
        }
    }
    if len(trueTypes) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "返回数据类型不能为空", nil))
    }
    ait.ReqData["datatype"] = strings.Join(trueTypes, "|")
}

func (ait *assistantInputTips) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := ait.ReqData["keywords"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "关键字不能为空", nil))
    }

    return ait.GetRequest()
}

func NewAssistantInputTips() *assistantInputTips {
    ait := &assistantInputTips{mpmap.NewBaseGaoDe(), "", make([]string, 0), "", "", "", ""}
    ait.SetServiceUri("/assistant/inputtips")
    ait.ReqData["citylimit"] = "false"
    ait.ReqData["datatype"] = "all"
    return ait
}
