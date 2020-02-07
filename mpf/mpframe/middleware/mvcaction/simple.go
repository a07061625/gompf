/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 23:43
 */
package mvcaction

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

// 动作日志
func NewSimpleLog() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        reqUrl := ctx.Values().GetStringDefault(project.DataParamKeyUrl, "")
        mplog.LogInfo(reqUrl + " action-enter")

        // 业务结束日志
        actionStart := time.Now()
        defer func() {
            costTime := time.Since(actionStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(reqUrl + " action-exist,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_action"].(float64) {
                mplog.LogWarn("handle " + reqUrl + " action-timeout,cost_time: " + costTimeStr + "s")
            }
        }()

        ctx.Next()
    }
}
