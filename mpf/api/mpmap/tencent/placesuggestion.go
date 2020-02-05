/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/27 0027
 * Time: 19:25
 */
package tencent

import (
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type placeSuggestion struct {
    mpmap.BaseTencent
    keyword     string            // 关键词
    region      string            // 地区
    regionLimit int               // 地区限制标识,0：当前城市无结果时,自动扩大范围到全国匹配 1：固定在当前城市
    location    string            // 定位坐标
    subLimit    int               // 子地点限制标识 0:不返回 1:返回
    policy      int               // 检索策略
    filters     map[string]string // 筛选条件
    page        uint              // 页码
    limit       uint              // 每页条数
}

func (ps *placeSuggestion) SetKeyword(keyword string) {
    if len(keyword) > 0 {
        ps.ReqData["keyword"] = keyword
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "关键词不能为空", nil))
    }
}

func (ps *placeSuggestion) SetRegion(region string) {
    if len(region) > 0 {
        ps.ReqData["region"] = region
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "地区不能为空", nil))
    }
}

func (ps *placeSuggestion) SetRegionLimit(regionLimit int) {
    if (regionLimit == 0) || (regionLimit == 1) {
        ps.ReqData["region_fix"] = strconv.Itoa(regionLimit)
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "地区限制不合法", nil))
    }
}

func (ps *placeSuggestion) SetLocation(lat, lng string) {
    ps.ReqData["location"] = lat + "," + lng
}

func (ps *placeSuggestion) SetSubLimit(subLimit int) {
    if (subLimit == 0) || (subLimit == 1) {
        ps.ReqData["get_subpois"] = strconv.Itoa(subLimit)
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "子地点限制不合法", nil))
    }
}

func (ps *placeSuggestion) SetPolicy(policy int) {
    ps.ReqData["policy"] = strconv.Itoa(policy)
}

func (ps *placeSuggestion) SetFilters(filters map[string]string) {
    if len(filters) > 0 {
        filterStr := ""
        for k, v := range filters {
            filterStr += "," + k + "=" + v
        }
        ps.ReqData["filter"] = filterStr[1:]
    }
}

func (ps *placeSuggestion) SetPage(page uint) {
    if page > 0 {
        ps.ReqData["page_index"] = strconv.Itoa(int(page))
    } else {
        ps.ReqData["page_index"] = "1"
    }
}

func (ps *placeSuggestion) SetLimit(limit uint) {
    if (limit > 0) || (limit <= 20) {
        ps.ReqData["page_size"] = strconv.Itoa(int(limit))
    } else {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "每页条数不合法", nil))
    }
}

func (ps *placeSuggestion) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := ps.ReqData["keyword"]
    if !ok {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "搜索关键字不能为空", nil))
    }
    _, ok = ps.ReqData["region"]
    if !ok {
        panic(mperr.NewMapTencent(errorcode.MapTencentParam, "地区不能为空", nil))
    }

    return ps.GetRequest()
}

func NewPlaceSuggestion() *placeSuggestion {
    ps := &placeSuggestion{mpmap.NewBaseTencent(), "", "", 0, "", 0, 0, make(map[string]string), 0, 0}
    ps.SetServiceUrl("https://apis.map.qq.com/ws/place/v1/suggestion")
    ps.ReqData["region_fix"] = "0"
    ps.ReqData["get_subpois"] = "0"
    ps.ReqData["policy"] = "0"
    ps.ReqData["page_index"] = "1"
    ps.ReqData["page_size"] = "10"
    return ps
}
