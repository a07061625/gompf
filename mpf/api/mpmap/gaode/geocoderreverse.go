/**
 * 逆地理编码
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 22:56
 */
package gaode

import (
    "strconv"
    "strings"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type geoCoderReverse struct {
    mpmap.BaseGaoDe
    location   string // 经纬度坐标
    poiType    string // POI类型
    radius     uint   // 搜索半径,单位:米
    extensions string // 返回结果标识
    batchTag   string // 批量查询标识
    roadLevel  int    // 道路等级 0:显示所有道路 1:仅输出主干道路数据
    homeOrCorp int    // POI返回排序 0:不干扰 1:居家相关内容优先返回 2:公司相关内容优先返回
}

func (gcr *geoCoderReverse) SetLocation(locations []string) {
    if len(locations) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标不能为空", nil))
    } else if len(locations) > 20 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标不能超过20个", nil))
    }

    gcr.ReqData["location"] = strings.Join(locations, "|")
}

func (gcr *geoCoderReverse) SetPoiType(poiTypeList []string) {
    if len(poiTypeList) == 0 {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "POI类型不能为空", nil))
    }

    gcr.ReqData["poitype"] = strings.Join(poiTypeList, "|")
}

func (gcr *geoCoderReverse) SetRadius(radius uint) {
    if (radius > 0) && (radius <= 3000) {
        gcr.ReqData["radius"] = strconv.Itoa(int(radius))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "搜索半径不合法", nil))
    }
}

func (gcr *geoCoderReverse) SetExtensions(extensions string) {
    if (extensions == "base") || (extensions == "all") {
        gcr.ReqData["extensions"] = extensions
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "返回结果标识不合法", nil))
    }
}

func (gcr *geoCoderReverse) SetBatchTag(batchTag string) {
    if (batchTag == "true") || (batchTag == "false") {
        gcr.ReqData["batch"] = batchTag
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "批量查询标识不合法", nil))
    }
}

func (gcr *geoCoderReverse) SetRoadLevel(roadLevel int) {
    if (roadLevel == 0) || (roadLevel == 1) {
        gcr.ReqData["roadlevel"] = strconv.Itoa(int(roadLevel))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "道路等级不合法", nil))
    }
}

func (gcr *geoCoderReverse) SetHomeOrCorp(homeOrCorp uint) {
    if (homeOrCorp > 0) && (homeOrCorp <= 2) {
        gcr.ReqData["homeorcorp"] = strconv.Itoa(int(homeOrCorp))
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "POI返回排序不合法", nil))
    }
}

func (gcr *geoCoderReverse) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := gcr.ReqData["location"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "经纬度坐标不能为空", nil))
    }

    return gcr.GetRequest()
}

func NewGeoCoderReverse() *geoCoderReverse {
    gcr := &geoCoderReverse{mpmap.NewBaseGaoDe(), "", "", 0, "", "", 0, 0}
    gcr.SetServiceUri("/geocode/regeo")
    gcr.ReqData["radius"] = "1000"
    gcr.ReqData["extensions"] = "base"
    gcr.ReqData["batch"] = "false"
    gcr.ReqData["roadlevel"] = "0"
    gcr.ReqData["homeorcorp"] = "0"
    return gcr
}
