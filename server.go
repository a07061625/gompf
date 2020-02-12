package main

import (
    "flag"
    "fmt"
    "io/ioutil"
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
    f, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY, 0666)
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

    for {
        sig := <-signals
        switch sig {
        case syscall.SIGUSR2: // 重启
            pid = server.ForkProcess(mpf.EnvServerModeChild)
            savePid(pid)
            server.Start()
        case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
            server.Shutdown(timeoutShutdown)
            os.Exit(1)
        }
    }
}

var (
    pid             = 0                   // 进程ID
    pidFile         = ""                  // 进程ID文件
    timeoutShutdown time.Duration         // 进程关闭超时时间
    server          mpserver.IServerBasic // 服务实例
    serverTag       = ""                  // 服务标识
    envType         = flag.String("mpet", mpf.EnvTypeProduct, "环境类型,只能是dev或product")
    projectTag      = flag.String("mppt", "", "项目标识,由小写字母和数字组成的3位长度字符串")
    projectModule   = flag.String("mppm", "", "项目模块,由字母和数字组成的字符串")
    optionType      = flag.String("mpot", "", "操作类型,start:启动服务 stop:停止服务 restart:重启服务")
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

    timeoutShutdown = 10 * time.Second
    serverTag = mpf.EnvProjectKey() + strconv.Itoa(mpf.EnvServerPort())
    pidFile = mpf.EnvDirRoot() + "/pid/" + serverTag + ".pid"
    pid = getPid()
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

    timeoutShutdown = time.Duration(conf.GetInt(mpf.EnvType()+"."+mpf.EnvProjectKeyModule()+"."+"timeout.shutdown")) * time.Second
    server = mpserver.New(mpf.EnvServerDomain(), app.GetInstance())
    if os.Getenv(mpf.GoEnvServerMode) != mpf.EnvServerModeChild { // master
        switch *optionType {
        case "start":
            if checkRunning() {
                fmt.Println("server " + serverTag + " is running")
                return
            }

            pid = server.ForkProcess(mpf.EnvServerModeDaemon)
            savePid(pid)
            server.Start()
        case "restart":
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
                    syscall.Kill(pid, syscall.SIGHUP) //kill -HUP, daemon only时会直接退出

                    pid = server.ForkProcess(mpf.EnvServerModeDaemon)
                    savePid(pid)
                    server.Start()
                }
            case <-time.After(10 * time.Second):
                fmt.Println("server " + serverTag + " restart timeout")
            }
        case "stop":
            if checkRunning() {
                err := syscall.Kill(pid, syscall.SIGTERM)
                if err != nil {
                    fmt.Println("server " + serverTag + " stop error: " + err.Error())
                }
            } else {
                fmt.Println("server " + serverTag + " already stop")
            }
        default:
            fmt.Println("操作类型必须是以下其一: start|stop|restart")
        }
        // 主进程退出
        os.Exit(0)
    }

    go listenNotify()
}
