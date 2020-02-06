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
    initServer()
}

type serverWeb struct {
    App        *iris.Application
    outer      mpframe.IOuterWeb
    RunConfigs []iris.Configurator
    Runner     iris.Runner
}

func (s serverWeb) SetOuter(outer mpframe.IOuterWeb) {
    s.outer = outer
}

func (s *serverWeb) initBase() {
    s.RunConfigs = append(s.RunConfigs, iris.WithoutInterruptHandler)

    listenCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    listen, err := listenCfg.NewListener("tcp4", mpf.EnvServerDomain())
    if err != nil {
        log.Fatalln("listen error:" + err.Error())
    }
    s.Runner = iris.Listener(listen)
}

func newServerWeb() serverWeb {
    s := serverWeb{}
    s.App = iris.New()
    s.RunConfigs = make([]iris.Configurator, 0)
    return s
}
