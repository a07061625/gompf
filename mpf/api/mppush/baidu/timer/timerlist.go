/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package timer

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询定时任务列表
type timerList struct {
    mppush.BaseBaiDu
}

func (tl *timerList) SetTimerId(timerId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, timerId)
    if match {
        tl.ReqData["timer_id"] = timerId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "定时任务ID不合法", nil))
    }
}

func (tl *timerList) SetStart(start int) {
    if start >= 0 {
        tl.ReqData["start"] = strconv.Itoa(start)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "起始索引位置不合法", nil))
    }
}

func (tl *timerList) SetLimit(limit int) {
    if (limit > 0) && (limit <= 10) {
        tl.ReqData["limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "记录条数不合法", nil))
    }
}

func (tl *timerList) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    return tl.GetRequest()
}

func NewTimerList() *timerList {
    tl := &timerList{mppush.NewBaseBaiDu()}
    tl.ServiceUri = "/timer/query_list"
    tl.ReqData["start"] = "0"
    tl.ReqData["limit"] = "10"
    return tl
}
