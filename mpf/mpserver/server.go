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
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpframe"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/valyala/tcplisten"
)

type IServerWeb interface {
    SetOuter(outer mpframe.IOuterWeb)
    StartServer()
}

type serverWeb struct {
    outer      mpframe.IOuterWeb
    runConfigs []iris.Configurator
    errHandles map[int]func(ctx iris.Context)
    App        *iris.Application
}

func (s *serverWeb) SetOuter(outer mpframe.IOuterWeb) {
    s.outer = outer
}

func (s *serverWeb) AddRunConfig(conf iris.Configurator) {
    s.runConfigs = append(s.runConfigs, conf)
}

func (s *serverWeb) AddErrHandle(statusCode int, handler func(ctx iris.Context)) {
    if (statusCode >= 100) && (statusCode < 600) {
        s.errHandles[statusCode] = handler
    }
}

func (s *serverWeb) bindErrHandles() {
    _, ok := s.errHandles[iris.StatusNotFound]
    if !ok {
        s.errHandles[iris.StatusNotFound] = func(ctx iris.Context) {
            mpf.NewLogger().Info("uri:" + ctx.RequestPath(false) + " not exist")
            result := mpresponse.NewResult("")
            result.Code = errorcode.CommonRequestResourceEmpty
            result.Msg = "接口不存在"
            ctx.WriteString(mpf.JsonMarshal(result))
        }
    }
    _, ok = s.errHandles[iris.StatusInternalServerError]
    if !ok {
        s.errHandles[iris.StatusInternalServerError] = func(ctx iris.Context) {
            mpf.NewLogger().Error("uri:" + ctx.RequestPath(false) + " error")
            result := mpresponse.NewResult("")
            result.Code = errorcode.CommonBaseServer
            result.Msg = "服务出错"
            ctx.WriteString(mpf.JsonMarshal(result))
        }
    }

    for k, v := range s.errHandles {
        s.App.OnErrorCode(k, v)
    }
}

func (s *serverWeb) baseStart() {
    s.runConfigs = append(s.runConfigs, iris.WithCharset("UTF-8"))
    s.runConfigs = append(s.runConfigs, iris.WithoutInterruptHandler)
    s.runConfigs = append(s.runConfigs, iris.WithoutStartupLog)
    s.runConfigs = append(s.runConfigs, iris.WithoutServerError(iris.ErrServerClosed))

    s.App.ConfigureHost(func(host *iris.Supervisor) {
        host.RegisterOnShutdown(func() {
            mpf.NewLogger().Info("server shut down")
        })
    })

    go s.outer.GetNotify(s.App)()

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
    return s
}
