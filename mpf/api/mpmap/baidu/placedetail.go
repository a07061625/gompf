/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 13:15
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

type placeDetail struct {
    mpmap.BaseBaiDu
    uids  []string // poi的uid数组
    scope string   // 结果详细程度
}

func (pd *placeDetail) AddUid(uid string) {
    if len(pd.uids) >= 10 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "uid数量超过限制", nil))
    }

    match, _ := regexp.MatchString(`^[0-9a-z]{24}$`, uid)
    if match {
        pd.uids = append(pd.uids, uid)
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "uid不合法", nil))
    }
}

// scope string 结果详细程度 1:基本信息 2:POI详细信息
func (pd *placeDetail) SetScope(scope string) {
    if (scope == "1") || (scope == "2") {
        pd.ReqData["scope"] = scope
    } else {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "结果详细程度不合法", nil))
    }
}

func (pd *placeDetail) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(pd.uids) == 0 {
        panic(mperr.NewMapBaiDu(errorcode.MapBaiDuParam, "uid不能为空", nil))
    }
    pd.ReqData["uids"] = strings.Join(pd.uids, ",")

    return pd.GetRequest()
}

func NewPlaceDetail() *placeDetail {
    pd := &placeDetail{mpmap.NewBaseBaiDu(), make([]string, 0), ""}
    pd.SetServiceUri("/place/v2/detail")
    pd.SetRespTag("result")
    pd.ReqData["scope"] = "1"
    pd.ReqData["timestamp"] = strconv.Itoa(time.Now().Second())
    return pd
}
