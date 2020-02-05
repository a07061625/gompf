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

// 创建定时任务
type scheduleCreate struct {
    mppush.BaseJPush
    name        string                 // 任务名
    trigger     map[string]interface{} // 触发条件
    pushContent map[string]interface{} // 推送内容
}

func (sc *scheduleCreate) SetName(name string) {
    if len(name) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务名不能为空", nil))
    } else if len(name) > 255 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务名不能超过255字节", nil))
    }

    sc.name = name
}

func (sc *scheduleCreate) SetTrigger(trigger map[string]interface{}) {
    if len(trigger) > 0 {
        sc.trigger = trigger
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "触发条件不合法", nil))
    }
}

func (sc *scheduleCreate) SetPushContent(pushContent map[string]interface{}) {
    if len(pushContent) > 0 {
        sc.pushContent = pushContent
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送内容不合法", nil))
    }
}

func (sc *scheduleCreate) SetCid(cid string) {
    if len(cid) > 0 {
        sc.ExtendData["cid"] = cid
    } else {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送标识符不合法", nil))
    }
}

func (sc *scheduleCreate) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(sc.name) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "任务名不能为空", nil))
    }
    if len(sc.trigger) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "触发条件不能为空", nil))
    }
    if len(sc.pushContent) == 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "推送内容不能为空", nil))
    }
    sc.ExtendData["name"] = sc.name
    sc.ExtendData["trigger"] = sc.trigger
    sc.ExtendData["push"] = sc.pushContent

    sc.ReqUrl = sc.GetServiceUrl()

    reqBody := mpf.JsonMarshal(sc.ExtendData)
    client, req := sc.GetRequest()
    req.SetBody([]byte(reqBody))

    return client, req
}

func NewScheduleCreate(key string) *scheduleCreate {
    sc := &scheduleCreate{mppush.NewBaseJPush(mppush.JPushServiceDomainApi, key, "app"), "", make(map[string]interface{}), make(map[string]interface{})}
    sc.ServiceUri = "/v3/schedules"
    sc.ExtendData["enabled"] = true
    sc.ReqContentType = project.HttpContentTypeJson
    sc.ReqMethod = fasthttp.MethodPost
    return sc
}
