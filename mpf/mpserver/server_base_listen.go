/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 2:10
 */
package mpserver

import (
    stdContext "context"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

func (s *serverBase) listenErrorCode() {
    s.app.OnAnyErrorCode(func(ctx context.Context) {
        statusCode := ctx.GetStatusCode()
        logMsg := "HTTP ERROR CODE: " + strconv.Itoa(statusCode) + " URI: " + ctx.Path()
        result := mpresponse.NewResultProblem()
        result.Title = "服务错误"
        result.Status = statusCode

        switch statusCode {
        case iris.StatusNotFound:
            mplog.LogInfo(logMsg)
            result.Type = "internal-address-not-exist"
            result.Code = errorcode.CommonRequestResourceEmpty
            result.Msg = "接口地址不存在"
        case iris.StatusMethodNotAllowed:
            mplog.LogInfo(logMsg)
            result.Type = "internal-method-not-allow"
            result.Code = errorcode.CommonRequestMethod
            result.Msg = "请求方式不支持"
        default:
            mplog.LogError(logMsg)
            result.Type = "internal-other"
            result.Code = errorcode.CommonBaseServer
            result.Msg = "其他服务错误"
        }
        ctx.Problem(mpresp.GetProblemHandleBasic(result, 30*time.Second))
        mpresp.NewBasicEnd()(ctx)
    })
}

func (s *serverBase) listenNotify() {
    go func(app *iris.Application, timeout time.Duration) {
        ch := make(chan os.Signal, 1)
        signal.Notify(ch,
            os.Interrupt,
            syscall.SIGINT,
            os.Kill,
            syscall.SIGKILL,
            syscall.SIGTERM,
        )

        select {
        case <-ch:
            ctx, _ := stdContext.WithTimeout(stdContext.Background(), timeout)
            app.Shutdown(ctx)
        }
    }(s.app, s.timeoutShutdown)
}
