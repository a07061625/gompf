/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 20:24
 */
package mpserv

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/frame"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

type IServerOuter interface {
    GetNotify(app *iris.Application) func()
    AddGlobalMiddleware(event int, handlers ...func(ctx iris.Context))
    GetGlobalMiddleware() map[int][]func(ctx iris.Context)
}

type simpleOuter struct {
    mwGlobal map[int][]func(ctx iris.Context)
}

func (so *simpleOuter) GetNotify(app *iris.Application) func() {
    return func() {
        ch := make(chan os.Signal, 1)
        signal.Notify(ch,
            os.Interrupt,
            syscall.SIGINT,
            os.Kill,
            syscall.SIGKILL,
            syscall.SIGTERM,
        )

        select {
        case s := <-ch:
            mplog.LogInfo("shutdown on signal " + fmt.Sprintf("%#v", s))

            timeout := 5 * time.Second
            ctx, _ := context.WithTimeout(context.Background(), timeout)
            app.Shutdown(ctx)
        }
    }
}

func (so *simpleOuter) AddGlobalMiddleware(event int, handlers ...func(ctx iris.Context)) {
    if (event == frame.MWEventGlobalPrefix) || (event == frame.MWEventGlobalSuffix) {
        so.mwGlobal[event] = handlers
    }
}

func (so *simpleOuter) GetGlobalMiddleware() map[int][]func(ctx iris.Context) {
    return so.mwGlobal
}

func NewSimpleOuter() *simpleOuter {
    so := &simpleOuter{}
    so.mwGlobal = make(map[int][]func(ctx iris.Context))
    return so
}
