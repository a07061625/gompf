package main

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/backend"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/frontend"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/index"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpreq"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpversion"
    "github.com/a07061625/gompf/mpf/mpserver"
    "github.com/kataras/iris/v12/context"
)

func init() {
    dirRoot, _ := os.Getwd()
    bs := mpf.NewBootstrap()
    bs.SetDirRoot(dirRoot)
    bs.SetDirConfigs(dirRoot + "/configs")
    bs.SetDirLogs(dirRoot + "/logs")
    mpf.LoadBoot(bs)
}

func main() {
    conf := mpf.NewConfig().GetConfig("server")
    server := mpserver.NewServer(conf)

    // 全局前置中间件
    middlewarePrefix := make([]context.Handler, 0)
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicInit())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicRecover())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicLog())
    middlewarePrefix = append(middlewarePrefix, mpversion.NewBasicError())
    server.SetGlobalMiddleware(true, middlewarePrefix...)

    // 全局后置中间件
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    versionKey1 := "< " + conf.GetString(confPrefix+"version.deprecated")
    versionKey2 := ">= " + conf.GetString(confPrefix+"version.deprecated") + ", < " + conf.GetString(confPrefix+"version.max")
    middlewareVersion := make(map[string]context.Handler)
    middlewareVersion[versionKey1] = mpversion.NewBasicDeprecated(mpresp.NewBasicSend(), "WARNING! You are using deprecated version of API", "Please use right version of API as soon as possible")
    middlewareVersion[versionKey2] = mpresp.NewBasicSend()
    middlewareSuffix := make([]context.Handler, 0)
    middlewareSuffix = append(middlewareSuffix, mpversion.NewBasicMatcher(middlewareVersion))
    middlewareSuffix = append(middlewareSuffix, mpresp.NewBasicEnd())
    server.SetGlobalMiddleware(false, middlewareSuffix...)

    // 注册路由
    routers := controllers.NewRouter()
    routers.RegisterGroup(index.NewRouter())
    routers.RegisterGroup(frontend.NewRouter())
    routers.RegisterGroup(backend.NewRouter())
    server.SetRouters(routers.GetControllers()...)
    server.StartServer()
}
