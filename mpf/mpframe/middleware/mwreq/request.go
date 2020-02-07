/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 10:26
 */
package mwreq

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
)

func NewIrisBefore() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        mplog.LogInfo(ctx.FullRequestURI() + " request-enter,params: " + string(ctx.Values().Serialize()))
        reqId := ctx.PostValueDefault("_req_id", "")
        mpf.ToolCreateReqId(reqId)

        // 兜底错误捕获
        defer func() {
            if r := recover(); r != nil {
                if ctx.IsStopped() {
                    return
                }

                errMsg := ""
                result := mpresponse.NewResult()
                if err1, ok := r.(*mperr.ErrorCommon); ok {
                    if err1.Type != mperr.TypeInnerValidator {
                        errMsg = err1.Msg
                    }
                    result.Code = err1.Code
                    result.Msg = err1.Msg
                } else if err2, ok := r.(error); ok {
                    errMsg = err2.Error()
                    result.Code = errorcode.CommonBaseServer
                    result.Msg = errMsg
                } else {
                    errMsg = "请求出错"
                    result.Code = errorcode.CommonBaseServer
                    result.Msg = errMsg
                }

                if len(errMsg) > 0 {
                    mplog.LogError(errMsg)
                }
                ctx.JSON(result)
                ctx.StatusCode(iris.StatusOK)
                ctx.StopExecution()
            }
        }()

        // 请求结束日志
        reqStart := time.Now()
        defer func() {
            costTime := time.Since(reqStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(ctx.FullRequestURI() + " request-exist,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_request"].(float64) {
                mplog.LogWarn("handle " + ctx.FullRequestURI() + " request-timeout,cost_time: " + costTimeStr + "s,params: " + string(ctx.Values().Serialize()))
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
