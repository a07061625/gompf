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
    "github.com/a07061625/gompf/mpf/mpframe"
    "github.com/kataras/iris/v12"
    "github.com/valyala/tcplisten"
)

type IServerWeb interface {
    SetOuter(outer mpframe.IOuterWeb)
    StartServer()
}

type serverWeb struct {
    App        *iris.Application
    outer      mpframe.IOuterWeb
    runConfigs []iris.Configurator
}

func (s *serverWeb) SetOuter(outer mpframe.IOuterWeb) {
    s.outer = outer
}

func (s *serverWeb) AddRunConfig(conf iris.Configurator) {
    s.runConfigs = append(s.runConfigs, conf)
}

func (s *serverWeb) baseStart() {
    s.runConfigs = append(s.runConfigs, iris.WithoutInterruptHandler)

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
    return s
}
