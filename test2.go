package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpapp"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/backend"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/frontend"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/index"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpreq"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpversion"
    "github.com/kataras/iris/v12/context"
    "github.com/kataras/iris/v12/i18n"
)

var (
    envType       = flag.String("mpet", mpf.EnvTypeProduct, "环境类型,只能是dev或product")
    projectTag    = flag.String("mppt", "", "项目标识,由小写字母和数字组成的3位长度字符串")
    projectModule = flag.String("mppm", "", "项目模块,由字母和数字组成的字符串")
    optionType    = flag.String("mpot", "", "操作类型,start:启动服务 stop:停止服务 restart:重启服务")
)

func init() {
    flag.Parse()

    dirRoot, _ := os.Getwd()
    bs := mpf.NewBootstrap()
    bs.SetDirRoot(dirRoot)
    bs.SetDirConfigs(dirRoot + "/configs")
    bs.SetDirLogs(dirRoot + "/logs")
    bs.SetEnvType(*envType)
    bs.SetProjectTag(*projectTag)
    bs.SetProjectModule(*projectModule)
    mpf.LoadBoot(bs)
}

func main() {
    app := mpapp.New()

    // 全局前置中间件
    middlewarePrefix := make([]context.Handler, 0)
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicBegin())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicInit())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicRecover())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicLog())
    middlewarePrefix = append(middlewarePrefix, mpversion.NewBasicError())
    app.SetMiddleware(true, middlewarePrefix...)

    // 全局后置中间件
    conf := mpf.NewConfig().GetConfig("server")
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    versionKey1 := "< " + conf.GetString(confPrefix+"version.deprecated")
    versionKey2 := ">= " + conf.GetString(confPrefix+"version.deprecated") + ", < " + conf.GetString(confPrefix+"version.max")
    middlewareVersion := make(map[string]context.Handler)
    middlewareVersion[versionKey1] = mpversion.NewBasicDeprecated(mpresp.NewBasicSend(), "WARNING! You are using deprecated version of API", "Please use right version of API as soon as possible")
    middlewareVersion[versionKey2] = mpresp.NewBasicSend()
    middlewareSuffix := make([]context.Handler, 0)
    middlewareSuffix = append(middlewareSuffix, mpversion.NewBasicMatcher(middlewareVersion))
    middlewareSuffix = append(middlewareSuffix, mpresp.NewBasicEnd())
    app.SetMiddleware(false, middlewarePrefix...)

    configOther := make(map[string]interface{})
    configOther["server_host"] = conf.GetString(confPrefix + "host")
    configOther["server_port"] = conf.GetInt(confPrefix + "port")
    configOther["server_type"] = conf.GetString(confPrefix + "type")
    configOther["version_min"] = conf.GetString(confPrefix + "version.min")
    configOther["version_deprecated"] = conf.GetString(confPrefix + "version.deprecated")
    configOther["version_current"] = conf.GetString(confPrefix + "version.current")
    configOther["version_max"] = conf.GetString(confPrefix + "version.max")
    configOther["timeout_request"] = conf.GetFloat64(confPrefix + "timeout.request")
    configOther["timeout_controller"] = conf.GetFloat64(confPrefix + "timeout.controller")
    configOther["timeout_action"] = conf.GetFloat64(confPrefix + "timeout.action")
    app.SetConfOther(configOther)

    confI18n := &i18n.I18n{}
    confI18n.Load("./configs/i18n/*/*.ini", "zh-CN", "en-US")
    confI18n.PathRedirect = false
    confI18n.URLParameter = conf.GetString(confPrefix + "reqparam.i18n")
    app.SetConfI18n(confI18n)

    app.SetRouterBlocks(conf.GetStringMapString(confPrefix + "mvc.block.accept"))
    routers := controllers.NewRouter()
    routers.RegisterGroup(index.NewRouter())
    routers.RegisterGroup(frontend.NewRouter())
    routers.RegisterGroup(backend.NewRouter())
    app.SetRouters(routers.GetControllers()...)

    app.Build()

    if os.Getenv(mpf.GoEnvServerMode) != mpf.EnvServerModeChild { // 主进程
        switch *optionType {
        case "start":
            server.Start()
        case "stop":
            server.Stop()
        case "restart":
            server.Restart()
        default:
            fmt.Println("操作类型必须是以下其一: start|stop|restart")
        }
        // 主进程退出
        os.Exit(0)
    }

    go server.ListenNotify()
}
