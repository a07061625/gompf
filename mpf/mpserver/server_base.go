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
    "time"

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
    SetRouters(controllers ...controllers.IControllerBasic) // 设置路由
    ReStart()                                               // 重启服务
    Start()                                                 // 启动服务
    Stop()                                                  // 停止服务
}

type serverBase struct {
    app             *iris.Application
    timeoutShutdown time.Duration
    confIris        []iris.Configurator
    confServer      *viper.Viper
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
    s.timeoutShutdown = 0
    s.confServer = conf
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
