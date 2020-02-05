/**
 * ID查询
 * User: 姜伟
 * Date: 2019/12/26 0026
 * Time: 23:52
 */
package gaode

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mpmap"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

type placeDetail struct {
    mpmap.BaseGaoDe
    id  string // 兴趣点ID
}

func (pd *placeDetail) SetId(id string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, id)
    if match {
        pd.ReqData["id"] = id
    } else {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "兴趣点ID不合法", nil))
    }
}

func (pd *placeDetail) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    _, ok := pd.ReqData["id"]
    if !ok {
        panic(mperr.NewMapGaoDe(errorcode.MapGaoDeParam, "兴趣点ID不能为空", nil))
    }

    return pd.GetRequest()
}

func NewPlaceDetail() *placeDetail {
    pd := &placeDetail{mpmap.NewBaseGaoDe(), ""}
    pd.SetServiceUri("/place/detail")
    return pd
}
