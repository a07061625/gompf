/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 10:14
 */
package mpserver

import (
    "log"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/frame"
    "github.com/a07061625/gompf/mpf/mpframe"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
    "github.com/valyala/tcplisten"
)

type IServerWeb interface {
    initConfig()
    initMiddleware()
    SetOuter(outer mpframe.IOuterWeb)
    StartServer()
}

type serverWeb struct {
    outer      mpframe.IOuterWeb
    runConfigs []iris.Configurator
    errHandles map[int]func(ctx iris.Context)
    mwHandles  map[int][]func(ctx iris.Context)
    App        *iris.Application
}

func (s *serverWeb) SetOuter(outer mpframe.IOuterWeb) {
    s.outer = outer
}

func (s *serverWeb) AddRunConfig(conf iris.Configurator) {
    s.runConfigs = append(s.runConfigs, conf)
}

func (s *serverWeb) AddMiddlewareHandle(event int, handler func(ctx iris.Context)) {
    _, ok := frame.TotalMiddlewareEvent[event]
    if !ok {
        return
    }

    _, ok = s.mwHandles[event]
    if !ok {
        s.mwHandles[event] = make([]func(ctx iris.Context), 0)
    }
    s.mwHandles[event] = append(s.mwHandles[event], handler)
}

func (s *serverWeb) AddErrHandle(statusCode int, handler func(ctx iris.Context)) {
    if (statusCode >= 100) && (statusCode < 600) {
        s.errHandles[statusCode] = handler
    }
}

func (s *serverWeb) BindErrHandles() {
    _, ok := s.errHandles[iris.StatusNotFound]
    if !ok {
        s.errHandles[iris.StatusNotFound] = func(ctx iris.Context) {
            mplog.LogInfo("uri:" + ctx.RequestPath(false) + " not exist")
            ctx.Redirect("/error/404")
        }
    }
    _, ok = s.errHandles[iris.StatusInternalServerError]
    if !ok {
        s.errHandles[iris.StatusInternalServerError] = func(ctx iris.Context) {
            mplog.LogError("uri:" + ctx.RequestPath(false) + " error")
            ctx.Redirect("/error/500")
        }
    }

    for k, v := range s.errHandles {
        s.App.OnErrorCode(k, v)
    }
}

func (s *serverWeb) baseStart() {
    go s.outer.GetNotify(s.App)()

    s.App.ConfigureHost(func(host *iris.Supervisor) {
        host.RegisterOnShutdown(func() {
            mplog.LogInfo("server shut down")
        })
    })

    listenCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    listen, err := listenCfg.NewListener("tcp4", mpf.EnvServerDomain())
    if err != nil {
        log.Fatalln("listen error:" + err.Error())
    }
    s.App.Run(iris.Listener(listen), s.runConfigs...)
}

func newServerWeb() serverWeb {
    s := serverWeb{}
    s.App = iris.New()
    s.runConfigs = make([]iris.Configurator, 0)
    s.errHandles = make(map[int]func(ctx iris.Context))
    s.mwHandles = make(map[int][]func(ctx iris.Context))
    return s
}
