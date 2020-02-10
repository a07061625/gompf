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
    mwGlobalPrefix := make([]context.Handler, 0)
    mwGlobalPrefix = append(mwGlobalPrefix, mpreq.NewBasicInit())
    mwGlobalPrefix = append(mwGlobalPrefix, mpreq.NewBasicRecover())
    mwGlobalPrefix = append(mwGlobalPrefix, mpreq.NewBasicLog())
    mwGlobalPrefix = append(mwGlobalPrefix, mpversion.NewBasicError())
    server.SetGlobalMiddleware(true, mwGlobalPrefix...)

    // 全局后置中间件
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    versionDeprecated := conf.GetString(confPrefix + "version.deprecated")
    versionMax := conf.GetString(confPrefix + "version.max")
    versionKey1 := "< " + versionDeprecated
    versionKey2 := ">= " + versionDeprecated + ", < " + versionMax
    mwVersionList := make(map[string]context.Handler)
    mwVersionList[versionKey1] = mpversion.NewBasicDeprecated(mpresp.NewBasicSend(), "WARNING! You are using deprecated version of API", "Please use right version of API as soon as possible")
    mwVersionList[versionKey2] = mpresp.NewBasicSend()
    mwGlobalSuffix := make([]context.Handler, 0)
    mwGlobalSuffix = append(mwGlobalSuffix, mpversion.NewBasicMatcher(mwVersionList))
    mwGlobalSuffix = append(mwGlobalSuffix, mpresp.NewBasicEnd())
    server.SetGlobalMiddleware(false, mwGlobalSuffix...)

    // 注册路由
    mpRouter := controllers.NewRouter()
    mpRouter.RegisterGroup(index.NewRouter())
    mpRouter.RegisterGroup(frontend.NewRouter())
    mpRouter.RegisterGroup(backend.NewRouter())
    server.SetRoutes(mpRouter.GetControllers()...)
    server.StartServer()
}
