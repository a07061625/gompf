/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 10:26
 */
package mwaction

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

func NewIrisBefore() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        logPrefix := ctx.FullRequestURI()
        if len(ctx.Request().URL.RawQuery) > 0 {
            logPrefix += "?" + ctx.Request().URL.RawQuery
        }

        mplog.LogInfo(logPrefix + " action-enter")

        // 业务结束日志
        actionStart := time.Now()
        defer func() {
            costTime := time.Since(actionStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(ctx.FullRequestURI() + " action-exist,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_action"].(float64) {
                mplog.LogWarn("handle " + logPrefix + " action-timeout,cost_time: " + costTimeStr + "s")
            }
        }()
        ctx.Next()
    }
}

func NewIrisAfter() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        ctx.Next()
    }
}
