/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 12:48
 */
package mpcontroller

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12/context"
)

// 控制器日志
func NewBasicLog() context.Handler {
    return func(ctx context.Context) {
        reqUrl := ctx.Values().GetString(project.DataParamKeyReqUrl)
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
    }
}
