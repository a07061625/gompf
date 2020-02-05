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

// 删除定时任务
type scheduleDelete struct {
    mppush.BaseJPush
    scheduleId string // 任务ID
}

func (sd *scheduleDelete) SetScheduleId(scheduleId string) {
    if len(scheduleId) > 0 {
        sd.scheduleId = scheduleId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不合法", nil))
    }
}

func (sd *scheduleDelete) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sd.scheduleId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不能为空", nil))
    }
    sd.ServiceUri = "/v3/schedules/" + sd.scheduleId

    sd.ReqUrl = sd.GetServiceUrl()

    return sd.GetRequest()
}

func NewScheduleDelete(key string) *scheduleDelete {
    sd := &scheduleDelete{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), ""}
    sd.ReqMethod = fasthttp.MethodDelete
    return sd
}
