/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 16:50
 */
package schedules

import (
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 获取定时任务
type scheduleGet struct {
    mppush.BaseJPush
    scheduleId string // 任务ID
}

func (sg *scheduleGet) SetScheduleId(scheduleId string) {
    if len(scheduleId) > 0 {
        sg.scheduleId = scheduleId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不合法", nil))
    }
}

func (sg *scheduleGet) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sg.scheduleId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不能为空", nil))
    }
    sg.ServiceUri = "/v3/schedules/" + sg.scheduleId

    sg.ReqURI = sg.GetServiceUrl()

    return sg.GetRequest()
}

func NewScheduleGet(key string) *scheduleGet {
    sg := &scheduleGet{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), ""}
    return sg
}
