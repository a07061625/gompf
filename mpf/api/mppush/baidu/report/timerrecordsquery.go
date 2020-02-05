/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/3 0003
 * Time: 22:30
 */
package report

import (
    "regexp"
    "strconv"

    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 查询定时消息的发送记录
type timerRecordsQuery struct {
    mppush.BaseBaiDu
    timerId string // 定时任务ID
}

func (trq *timerRecordsQuery) SetTimerId(timerId string) {
    match, _ := regexp.MatchString(project.RegexDigitAlpha, timerId)
    if match {
        trq.timerId = timerId
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "定时任务ID不合法", nil))
    }
}

func (trq *timerRecordsQuery) SetStart(start int) {
    if start >= 0 {
        trq.ReqData["start"] = strconv.Itoa(start)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "起始索引位置不合法", nil))
    }
}

func (trq *timerRecordsQuery) SetLimit(limit int) {
    if (limit > 0) && (limit <= 100) {
        trq.ReqData["limit"] = strconv.Itoa(limit)
    } else {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "记录条数不合法", nil))
    }
}

func (trq *timerRecordsQuery) SetRangeTime(startTime, endTime int) {
    if startTime < 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "起始时间不合法", nil))
    } else if endTime < 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "截止时间不合法", nil))
    }
    if (startTime > 0) && (endTime > 0) {
        if endTime < startTime {
            panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "起始时间不能大于截止时间", nil))
        }
        trq.ReqData["range_start"] = strconv.Itoa(startTime)
        trq.ReqData["range_end"] = strconv.Itoa(endTime)
    } else if startTime > 0 {
        trq.ReqData["range_start"] = strconv.Itoa(startTime)
    } else if endTime > 0 {
        trq.ReqData["range_end"] = strconv.Itoa(endTime)
    }
}

func (trq *timerRecordsQuery) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if len(trq.timerId) == 0 {
        panic(mperr.NewPushBaiDu(errorcode.PushBaiDuParam, "定时任务ID不能为空", nil))
    }
    trq.ReqData["timer_id"] = trq.timerId

    return trq.GetRequest()
}

func NewTimerRecordsQuery() *timerRecordsQuery {
    trq := &timerRecordsQuery{mppush.NewBaseBaiDu(), ""}
    trq.ServiceUri = "/report/query_timer_records"
    trq.ReqData["start"] = "0"
    trq.ReqData["limit"] = "100"
    return trq
}
