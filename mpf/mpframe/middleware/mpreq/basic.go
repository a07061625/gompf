/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 12:48
 */
package mpreq

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mperr"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12/context"
)

// 请求初始化
func NewBasicInit() context.Handler {
    return func(ctx context.Context) {
        ctx.Record()

        reqId := ""
        if mpf.EnvServerTypeRpc == ctx.Application().ConfigurationReadOnly().GetOther()["server_type"].(string) {
            reqId = ctx.PostValueDefault(project.DataParamKeyReqId, "")
        }
        mpf.ToolCreateReqId(reqId)

        reqUrl := ctx.FullRequestURI()
        if len(ctx.Request().URL.RawQuery) > 0 {
            reqUrl += "?" + ctx.Request().URL.RawQuery
        }
        ctx.Values().Set(project.DataParamKeyReqUrl, reqUrl)
    }
}

// 请求日志
func NewBasicLog() context.Handler {
    return func(ctx context.Context) {
        reqUrl := ctx.Values().GetString(project.DataParamKeyReqUrl)
        mplog.LogInfo(reqUrl + " request-enter")

        reqStart := time.Now()
        defer func() {
            costTime := time.Since(reqStart).Seconds()
            costTimeStr := strconv.FormatFloat(costTime, 'f', 6, 64)
            mplog.LogInfo(reqUrl + " request-exist,cost_time: " + costTimeStr + "s")
            if costTime >= ctx.Application().ConfigurationReadOnly().GetOther()["timeout_request"].(float64) {
                mplog.LogWarn("handle " + reqUrl + " request-timeout,cost_time: " + costTimeStr + "s")
            }
        }()
    }
}

// 错误捕获
func NewBasicRecover() context.Handler {
    return func(ctx context.Context) {
        defer func() {
            if r := recover(); r != nil {
                if ctx.IsStopped() {
                    return
                }

                errMsg := ""
                result := mpresponse.NewResultBasic()
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

                ctx.Recorder().Header().Set(project.HttpHeadKeyContentType, project.HttpContentTypeJson)
                ctx.Recorder().SetBodyString(mpf.JsonMarshal(result))
                ctx.StopExecution()
            }
        }()
    }
}
