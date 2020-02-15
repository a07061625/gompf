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
    "github.com/a07061625/gompf/mpf/mpserver"
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
        timeout := 10 * time.Second
        ctx, _ := stdContext.WithTimeout(stdContext.Background(), timeout)
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
    app       *iris.Application
    ws        mpserver.IServerWebSocket
    serverTag = "" // 服务标识
    pidFile   = "" // 进程ID文件
    pid       = 0  // 进程ID

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
    appBasic := mpapp.New()
    // 全局前置中间件
    appBasic.SetGlobalMiddlewarePrefix(mpreq.NewBasicBegin(), mpreq.NewBasicInit(), mpreq.NewBasicRecover(), mpreq.NewBasicLog(), mpversion.NewBasicError())
    // 全局后置中间件
    versionKey1 := "< " + conf.GetString(confPrefix+"version.deprecated")
    versionKey2 := ">= " + conf.GetString(confPrefix+"version.deprecated") + ", < " + conf.GetString(confPrefix+"version.max")
    middlewareVersion := make(map[string]context.Handler)
    middlewareVersion[versionKey1] = mpversion.NewBasicDeprecated(mpresp.NewBasicSend(), "WARNING! You are using deprecated version of API", "Please use right version of API as soon as possible")
    middlewareVersion[versionKey2] = mpresp.NewBasicSend()
    appBasic.SetGlobalMiddlewareSuffix(mpversion.NewBasicMatcher(middlewareVersion), mpresp.NewBasicEnd())

    // 应用配置
    appBasic.SetConfApp()

    // 国际化配置,配置文件路径只能是以./开始,代表从main.go文件所在目录开始,否则会报错
    i18nConf := &i18n.I18n{}
    i18nConf.Load("./configs/i18n/*/*.ini", "zh-CN", "en-US")
    i18nConf.PathRedirect = false
    i18nConf.URLParameter = conf.GetString(confPrefix + "reqparam.i18n")
    appBasic.SetConfI18n(i18nConf)

    // 路由
    blocks := conf.GetStringMapString(confPrefix + "mvc.block.accept")
    routers := controllers.NewRouter()
    routers.RegisterGroup(index.NewRouter())
    routers.RegisterGroup(frontend.NewRouter())
    routers.RegisterGroup(backend.NewRouter())
    appBasic.SetRouters(blocks, routers.GetControllers())

    // 错误处理
    appBasic.SetErrorHandler()

    app = appBasic.GetInstance()

    // 对接第三方应用
    ws = mpserver.NewServerWebSocket() // web socket
    app.WrapRouter(func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
        path := r.URL.Path
        if path == "/websocket/demo" {
            ctx := app.ContextPool.Acquire(w, r)
            defer app.ContextPool.Release(ctx)
            ws.Handler()(ctx)
            return
        }

        router.ServeHTTP(w, r)
    })

    go listenNotify()

    listener := mpserver.NewListenerTcp(mpf.EnvServerDomain())
    if *optionType == "start" {
        logMsg := "server " + serverTag
        if checkRunning() {
            fmt.Println(logMsg + " is running")
            return
        }

        pid = os.Getpid()
        savePid(pid)
        err := app.Run(iris.Listener(listener))
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
                err := app.Run(iris.Listener(listener))
                if err != nil {
                    log.Fatalln("server " + serverTag + " restart error: " + err.Error())
                }
            }
        case <-time.After(10 * time.Second):
            fmt.Println("server " + serverTag + " restart timeout")
        }
    }
}
