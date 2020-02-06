/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 11:11
 */
package mpframe

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

type outerHttp struct {
    outerWeb
}

func (oh *outerHttp) GetNotify(app *iris.Application) func() {
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

func NewOuterHttp() *outerHttp {
    return &outerHttp{newOuterWeb()}
}
