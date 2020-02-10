/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 1:15
 */
package mpserver

import (
    "log"
    "sync"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/spf13/viper"
    "github.com/valyala/tcplisten"
)

type IServerBase interface {
    AddIrisConf(configs ...iris.Configurator)
    SetGlobalMiddleware(isPrefix bool, middlewareList ...context.Handler)
    SetRoutes(controllers ...controllers.IControllerBasic)
    StartServer()
}

type serverBase struct {
    app        *iris.Application
    confIris   []iris.Configurator
    confServer *viper.Viper
    routeFlag  bool // 路由标识 true:已设置 false:未设置
}

// 设置全局中间件
//   isPrefix: 中间件类型,true:前置 false:后置
//   middlewareList: 中间件列表
func (s *serverBase) SetGlobalMiddleware(isPrefix bool, middlewareList ...context.Handler) {
    if len(middlewareList) == 0 {
        return
    }

    if isPrefix {
        s.app.UseGlobal(middlewareList...)
    } else {
        s.app.DoneGlobal(middlewareList...)
    }
}

func (s *serverBase) bootstrap() {
    s.initConf()
    s.listenErrorCode()
    s.listenNotify()

    listenCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    listen, err := listenCfg.NewListener("tcp4", mpf.EnvServerDomain())
    if err != nil {
        log.Fatalln("listen error: " + err.Error())
    }
    s.app.Run(iris.Listener(listen), s.confIris...)
}

func newServerBase(conf *viper.Viper) serverBase {
    s := serverBase{}
    s.app = iris.New()
    s.confServer = conf
    s.routeFlag = false
    s.initConf()
    return s
}

var (
    once sync.Once
    ins  IServerBase
)

func NewServer(conf *viper.Viper) IServerBase {
    once.Do(func() {
        if mpf.EnvServerType() == mpf.EnvServerTypeApi {
            ins = &serverHttp{newServerBase(conf)}
        }
    })

    return ins
}
