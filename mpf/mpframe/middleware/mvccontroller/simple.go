/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 20:26
 */
package mwcontroller

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

// 控制器日志
func NewSimpleLog() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        reqUrl := ctx.Values().GetStringDefault(project.DataParamKeyUrl, "")
        mplog.LogInfo(reqUrl + " controller-enter")

        // 业务结束日志
        controllerStart := time.Now()
        defer func() {
            costTime := time.Since(controllerStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(reqUrl + " controller-exist,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_controller"].(float64) {
                mplog.LogWarn("handle " + reqUrl + " controller-timeout,cost_time: " + costTimeStr + "s")
            }
        }()

        ctx.Next()
    }
}
