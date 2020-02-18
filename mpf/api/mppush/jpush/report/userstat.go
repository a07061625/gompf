/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/4 0004
 * Time: 15:54
 */
package report

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/api/mppush"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/valyala/fasthttp"
)

// 用户统计
type userStat struct {
    mppush.BaseJPush
    duration int // 持续时长
}

func (us *userStat) SetStatTime(timeUnit string, timestamp int, duration int) {
    if timestamp <= 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "起始时间不合法", nil))
    }
    if duration <= 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "持续时长不合法", nil))
    }

    st := time.Unix(int64(timestamp), 0)
    switch timeUnit {
    case "HOUR":
        if duration > 24 {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "持续时长不合法", nil))
        }
        us.ReqData["start"] = st.Format("2006-01-02 03")
    case "DAY":
        if duration > 60 {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "持续时长不合法", nil))
        }
        us.ReqData["start"] = st.Format("2006-01-02")
    case "MONTH":
        if duration > 2 {
            panic(mperr.NewPushJPush(errorcode.PushJPushParam, "持续时长不合法", nil))
        }
        us.ReqData["start"] = st.Format("2006-01")
    default:
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "时间单位不合法", nil))
    }

    us.ReqData["time_unit"] = timeUnit
    us.duration = duration
}

func (us *userStat) CheckData() (*fasthttp.Client, *fasthttp.Request) {
    if us.duration <= 0 {
        panic(mperr.NewPushJPush(errorcode.PushJPushParam, "持续时长不能为空", nil))
    }
    us.ReqData["duration"] = strconv.Itoa(us.duration)

    us.ReqURI = us.GetServiceUrl() + "?" + mpf.HTTPCreateParams(us.ReqData, "none", 1)

    return us.GetRequest()
}

func NewUserStat(key string) *userStat {
    us := &userStat{mppush.NewBaseJPush(mppush.JPushServiceDomainReport, key, "app"), 0}
    us.ServiceUri = "/v3/users"
    return us
}
