/**
 * 天气查询
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 23:21
 */
package gaode

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type weatherInfo struct {
    mpmap.BaseGaoDe
    city       string // 城市编码
    extensions string // 气象类型 base:实况天气 all:预报天气
}

func (wi *weatherInfo) SetCity(city string) {
    match, _ := regexp.MatchString(`^[0-9]{6}$`, city)
    if match {
        wi.ReqData["city"] = city
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "城市编码不合法", nil))
    }
}

func (wi *weatherInfo) SetExtensions(extensions string) {
    if (extensions == "base") || (extensions == "all") {
        wi.ReqData["extensions"] = extensions
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "气象类型不合法", nil))
    }
}

func (wi *weatherInfo) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := wi.ReqData["city"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "城市编码不能为空", nil))
    }

    return wi.GetRequest()
}

func NewWeatherInfo() *weatherInfo {
    wi := &weatherInfo{mpmap.NewBaseGaoDe(), "", ""}
    wi.SetServiceUri("/weather/weatherInfo")
    wi.ReqData["extensions"] = "base"
    return wi
}
