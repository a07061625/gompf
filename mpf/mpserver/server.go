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
    "github.com/a07061625/gompf/mpf/mpframe/mpserv"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
    "github.com/spf13/viper"
    "github.com/valyala/tcplisten"
)

type IServerSimple interface {
    bootServer()
    SetOuter(outer mpserv.IServerOuter)
    StartServer()
    AddRunConfig(configs ...iris.Configurator)
}

type serverSimple struct {
    outer      mpserv.IServerOuter
    runConfigs []iris.Configurator
    App        *iris.Application
}

func (s *serverSimple) SetOuter(outer mpserv.IServerOuter) {
    s.outer = outer
}

func (s *serverSimple) AddRunConfig(configs ...iris.Configurator) {
    if len(configs) > 0 {
        s.runConfigs = append(s.runConfigs, configs...)
    }
}

func (s *serverSimple) bootSimple(conf *viper.Viper) {
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_host", conf.GetString(confPrefix+"host")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_port", conf.GetInt(confPrefix+"port")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_type", conf.GetString(confPrefix+"type")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_request", conf.GetFloat64(confPrefix+"timeout.request")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_controller", conf.GetFloat64(confPrefix+"timeout.controller")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_action", conf.GetFloat64(confPrefix+"timeout.action")))

    s.runConfigs = append(s.runConfigs, iris.WithCharset("UTF-8"))
    s.runConfigs = append(s.runConfigs, iris.WithoutStartupLog)
    s.runConfigs = append(s.runConfigs, iris.WithOptimizations)
    s.runConfigs = append(s.runConfigs, iris.WithoutInterruptHandler)
    s.runConfigs = append(s.runConfigs, iris.WithoutAutoFireStatusCode)
    s.runConfigs = append(s.runConfigs, iris.WithoutServerError(iris.ErrServerClosed))

    mwGlobalList := s.outer.GetGlobalMiddleware()
    for eventType, handles := range mwGlobalList {
        handleNum := len(handles)
        if handleNum == 0 {
            continue
        }
        if eventType == frame.MWEventGlobalPrefix {
            for i := 0; i < handleNum; i++ {
                s.App.UseGlobal(handles[i])
            }
        } else {
            for i := 0; i < handleNum; i++ {
                s.App.DoneGlobal(handles[i])
            }
        }
    }

    s.App.ConfigureHost(func(host *iris.Supervisor) {
        host.RegisterOnShutdown(func() {
            mplog.LogInfo("server shut down")
        })
    })
}

func (s *serverSimple) startSimple() {
    go s.outer.GetNotify(s.App)()

    listenCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    listen, err := listenCfg.NewListener("tcp4", mpf.EnvServerDomain())
    if err != nil {
        log.Fatalln("listen error: " + err.Error())
    }
    s.App.Run(iris.Listener(listen), s.runConfigs...)
}

func newServerSimple() serverSimple {
    s := serverSimple{}
    s.App = iris.New()
    s.runConfigs = make([]iris.Configurator, 0)
    return s
}
