/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package timer

import (
    "regexp"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 取消定时任务
type timerCancel struct {
    mppush.BaseBaiDu
    timerId string // 定时任务ID
}

func (tc *timerCancel) SetTimerId(timerId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, timerId)
    if match {
        tc.timerId = timerId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "定时任务ID不合法", nil))
    }
}

func (tc *timerCancel) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(tc.timerId) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "定时任务ID不能为空", nil))
    }
    tc.ReqData["timer_id"] = tc.timerId

    return tc.GetRequest()
}

func NewTimerCancel() *timerCancel {
    tc := &timerCancel{mppush.NewBaseBaiDu(), ""}
    tc.ServiceUri = "/timer/cancel"
    return tc
}
