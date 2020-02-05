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

// 获取定时任务对应的所有消息ID
type msgIdList struct {
    mppush.BaseJPush
    scheduleId string // 任务ID
}

func (mil *msgIdList) SetScheduleId(scheduleId string) {
    if len(scheduleId) > 0 {
        mil.scheduleId = scheduleId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不合法", nil))
    }
}

func (mil *msgIdList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(mil.scheduleId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不能为空", nil))
    }
    mil.ServiceUri = "/v3/schedules/" + mil.scheduleId + "/msg_ids"

    mil.ReqUrl = mil.GetServiceUrl()

    return mil.GetRequest()
}

func NewMsgIdList(key string) *msgIdList {
    mil := &msgIdList{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), ""}
    return mil
}
