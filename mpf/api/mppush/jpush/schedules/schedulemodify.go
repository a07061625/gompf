/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 12:19
 */
package schedules

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 修改定时任务
type scheduleModify struct {
    mppush.BaseJPush
    scheduleId string // 任务ID
}

func (sm *scheduleModify) SetScheduleId(scheduleId string) {
    if len(scheduleId) > 0 {
        sm.scheduleId = scheduleId
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不合法", nil))
    }
}

func (sm *scheduleModify) SetName(name string) {
    if len(name) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务名不能为空", nil))
    } else if len(name) > 255 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务名不能超过255字节", nil))
    }

    sm.ExtendData["name"] = name
}

func (sm *scheduleModify) SetTrigger(trigger map[string]interface{}) {
    if len(trigger) > 0 {
        sm.ExtendData["trigger"] = trigger
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "触发条件不合法", nil))
    }
}

func (sm *scheduleModify) SetPushContent(pushContent map[string]interface{}) {
    if len(pushContent) > 0 {
        sm.ExtendData["push"] = pushContent
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送内容不合法", nil))
    }
}

func (sm *scheduleModify) SetEnabled(enabled bool) {
    sm.ExtendData["enabled"] = enabled
}

func (sm *scheduleModify) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sm.scheduleId) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务ID不能为空", nil))
    }
    if len(sm.ExtendData) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "修改内容不能为空", nil))
    }
    sm.ServiceUri = "/v3/schedules/" + sm.scheduleId

    sm.ReqUrl = sm.GetServiceUrl()

    reqBody := mpf.JSONMarshal(sm.ExtendData)
    client, req := sm.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewScheduleModify(key string) *scheduleModify {
    sm := &scheduleModify{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), ""}
    sm.ReqContentType = project.HTTPContentTypeJSON
    sm.ReqMethod = fasthttp.MethodPut
    return sm
}
