/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 16:50
 */
package schedules

import (
    "strconv"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取有效的定时任务列表
type scheduleList struct {
    mppush.BaseJPush
}

func (sl *scheduleList) SetPage(page int) {
    if page > 0 {
        sl.ReqData["page"] = strconv.Itoa(page)
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "页数不合法", nil))
    }
}

func (sl *scheduleList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    sl.ReqURI = sl.GetServiceUrl() + "?" + mpf.HTTPCreateParams(sl.ReqData, "none", 1)

    return sl.GetRequest()
}

func NewScheduleList(key string) *scheduleList {
    sl := &scheduleList{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app")}
    sl.ServiceUri = "/v3/schedules"
    sl.ReqData["page"] = "1"
    return sl
}
