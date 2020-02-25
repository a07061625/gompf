// Package mpcontroller basic
// User: 姜伟
// Time: 2020-02-25 10:54:23
package mpcontroller

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12/context"
)

// NewBasicLog 控制器日志
func NewBasicLog() context.Handler {
    return func(ctx context.Context) {
        reqURL := ctx.Values().GetString(project.DataParamKeyReqURL)
        mplog.LogInfo(reqURL + " controller-enter")

        // 业务结束日志
        controllerStart := time.Now()
        defer func() {
            costTime := time.Since(controllerStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(reqURL + " controller-exit,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_controller"].(float64) {
                mplog.LogWarn("handle " + reqURL + " controller-timeout,cost_time: " + costTimeStr + "s")
            }
        }()
    }
}
