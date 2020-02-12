package main

import (
    stdContext "context"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpapp"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/backend"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/frontend"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/index"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpreq"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpversion"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/kataras/iris/v12/i18n"
)

func getPid() int {
    pid := 0
    if f, err := os.Open(pidFile); err == nil {
        pidStr, _ := ioutil.ReadAll(f)
        pid, _ = strconv.Atoi(string(pidStr))
        defer f.Close()
    }

    return pid
}

func savePid(pid int) {
    f, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY, 0664)
    if err != nil {
        mplog.LogInfo("write pid file error: " + err.Error())
        return
    }
    defer f.Close()
    f.WriteString(strconv.Itoa(pid))
}

// 发一个信号为0到指定进程ID,如果没有错误发生,表示进程存活
func checkRunning() bool {
    if pid <= 0 {
        return false
    }

    err := syscall.Kill(pid, 0)
    return err == nil
}

// 捕获系统信号
func listenNotify() {
    signals := make(chan os.Signal)
    signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

    select {
    case sig := <-signals:
        mplog.LogInfo("server " + serverTag + " shutdown by signal " + sig.String())
        ctx, _ := stdContext.WithTimeout(stdContext.Background(), timeoutList["shutdown"])
        app.Shutdown(ctx)
    }
}

func stop() {
    logMsg := "server " + serverTag
    if !checkRunning() {
        fmt.Println(logMsg + " already stoped")
        return
    }

    err := syscall.Kill(pid, syscall.SIGTERM)
    if err == nil {
        fmt.Println(logMsg + " stop success")
    } else {
        fmt.Println(logMsg + " stop error: " + err.Error())
    }
}

var (
    app         *iris.Application
    timeoutList map[string]time.Duration // 超时时间列表
    serverTag   = ""                     // 服务标识
    pidFile     = ""                     // 进程ID文件
    pid         = 0                      // 进程ID

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

    timeoutList = make(map[string]time.Duration)
}

func main() {
    serverTag = mpf.EnvProjectKey() + strconv.Itoa(mpf.EnvServerPort())
    pidFile = mpf.EnvDirRoot() + "/pid/" + serverTag + ".pid"
    pid = getPid()

    switch *optionType {
    case "stop":
        stop()
        os.Exit(0)
    case "start":
    case "restart":
    default:
        fmt.Println("操作类型必须是以下其一: start|stop|restart")
        os.Exit(1)
    }

    conf := mpf.NewConfig().GetConfig("server")
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    timeoutList["read"] = time.Duration(conf.GetInt(confPrefix+"timeout.read")) * time.Second
    timeoutList["readheader"] = time.Duration(conf.GetInt(confPrefix+"timeout.readheader")) * time.Second
    timeoutList["write"] = time.Duration(conf.GetInt(confPrefix+"timeout.write")) * time.Second
    timeoutList["idle"] = time.Duration(conf.GetInt(confPrefix+"timeout.idle")) * time.Second
    timeoutList["shutdown"] = time.Duration(conf.GetInt(confPrefix+"timeout.shutdown")) * time.Second

    appBasic := mpapp.New()
    // 全局前置中间件
    middlewarePrefix := make([]context.Handler, 0)
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicBegin())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicInit())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicRecover())
    middlewarePrefix = append(middlewarePrefix, mpreq.NewBasicLog())
    middlewarePrefix = append(middlewarePrefix, mpversion.NewBasicError())
    appBasic.SetMiddleware(true, middlewarePrefix...)

    // 全局后置中间件
    versionKey1 := "< " + conf.GetString(confPrefix+"version.deprecated")
    versionKey2 := ">= " + conf.GetString(confPrefix+"version.deprecated") + ", < " + conf.GetString(confPrefix+"version.max")
    middlewareVersion := make(map[string]context.Handler)
    middlewareVersion[versionKey1] = mpversion.NewBasicDeprecated(mpresp.NewBasicSend(), "WARNING! You are using deprecated version of API", "Please use right version of API as soon as possible")
    middlewareVersion[versionKey2] = mpresp.NewBasicSend()
    middlewareSuffix := make([]context.Handler, 0)
    middlewareSuffix = append(middlewareSuffix, mpversion.NewBasicMatcher(middlewareVersion))
    middlewareSuffix = append(middlewareSuffix, mpresp.NewBasicEnd())
    appBasic.SetMiddleware(false, middlewarePrefix...)

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
    appBasic.SetConfOther(configOther)

    confI18n := &i18n.I18n{}
    confI18n.Load("./configs/i18n/*/*.ini", "zh-CN", "en-US")
    confI18n.PathRedirect = false
    confI18n.URLParameter = conf.GetString(confPrefix + "reqparam.i18n")
    appBasic.SetConfI18n(confI18n)

    appBasic.SetRouterBlocks(conf.GetStringMapString(confPrefix + "mvc.block.accept"))
    routers := controllers.NewRouter()
    routers.RegisterGroup(index.NewRouter())
    routers.RegisterGroup(frontend.NewRouter())
    routers.RegisterGroup(backend.NewRouter())
    appBasic.SetRouters(routers.GetControllers()...)

    appBasic.Init()
    app = appBasic.GetInstance()

    go listenNotify()

    server := &http.Server{
        Addr:              mpf.EnvServerDomain(),
        ReadTimeout:       timeoutList["read"],
        ReadHeaderTimeout: timeoutList["readheader"],
        WriteTimeout:      timeoutList["write"],
        IdleTimeout:       timeoutList["idle"],
    }

    if *optionType == "start" {
        logMsg := "server " + serverTag
        if checkRunning() {
            fmt.Println(logMsg + " is running")
            return
        }

        pid = os.Getpid()
        savePid(pid)
        err := app.Run(iris.Server(server))
        if err != nil {
            log.Fatalln(logMsg + " start error: " + err.Error())
        }
    } else {
        syscall.Kill(pid, syscall.SIGTERM)
        stopStatus := make(chan bool, 1)
        go func() {
            for {
                if !checkRunning() {
                    stopStatus <- true
                    break
                }
                time.Sleep(500 * time.Millisecond)
            }
        }()

        select {
        case status := <-stopStatus:
            if status {
                pid = os.Getpid()
                savePid(pid)
                err := app.Run(iris.Server(server))
                if err != nil {
                    log.Fatalln("server " + serverTag + " restart error: " + err.Error())
                }
            }
        case <-time.After(10 * time.Second):
            fmt.Println("server " + serverTag + " restart timeout")
        }
    }
}
