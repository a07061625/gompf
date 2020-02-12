package mpserver

import (
    "context"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "sync"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mplog"
)

type serverBasic struct {
    server    http.Server
    serverTag string                   // 服务标识
    appFlag   bool                     // app初始化标识 true:已初始化 false:未初始化
    timeout   map[string]time.Duration // 超时时间列表
    pidFile   string                   // 进程ID文件
    pid       int                      // 进程ID
}

func (s *serverBasic) getPid() int {
    pid := 0
    if f, err := os.Open(s.pidFile); err == nil {
        pidStr, _ := ioutil.ReadAll(f)
        pid, _ = strconv.Atoi(string(pidStr))
        defer f.Close()
    }

    return pid
}

func (s *serverBasic) savePid(pid int) {
    f, err := os.OpenFile(s.pidFile, os.O_CREATE|os.O_WRONLY, 0664)
    if err != nil {
        mplog.LogInfo("write pid file error: " + err.Error())
        return
    }
    defer f.Close()
    f.WriteString(strconv.Itoa(pid))
}

// 发一个信号为0到指定进程ID,如果没有错误发生,表示进程存活
func (s *serverBasic) checkRunning() bool {
    if s.pid <= 0 {
        return false
    }

    err := syscall.Kill(s.pid, 0)
    return err == nil
}

func (s *serverBasic) GetPidFile() string {
    return s.pidFile
}

func (s *serverBasic) SetApp(app http.Handler) {
    if s.appFlag {
        return
    }

    s.server.Handler = app
    s.server.Addr = mpf.EnvServerDomain()
}

// 捕获系统信号
func (s *serverBasic) listenNotify() {
    signals := make(chan os.Signal)
    signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

    select {
    case sig := <-signals:
        mplog.LogInfo("server " + s.serverTag + " shutdown by signal " + sig.String())
        ctx, _ := context.WithTimeout(context.Background(), insBasic.timeout["shutdown"])
        s.server.Shutdown(ctx)
    }
}

// 启动服务
//   startType: int 1:启动 2:重启
func (s *serverBasic) Start(startType int) {
    logMsg := "server " + s.serverTag
    if s.checkRunning() {
        fmt.Println(logMsg + " is running")
        return
    }
    if s.server.Handler == nil {
        fmt.Println(logMsg + " not bind app")
        return
    }

    if startType == 1 {
        logMsg += " start "
    } else {
        logMsg += " restart "
    }
    err := s.server.ListenAndServe()
    if err != nil {
        fmt.Println(logMsg + "error: " + err.Error())
        return
    }

    fmt.Println(logMsg + "success")
    s.savePid(os.Getpid())
    s.listenNotify()
}

func (s *serverBasic) Stop() {
    logMsg := "server " + s.serverTag
    if !s.checkRunning() {
        fmt.Println(logMsg + " already stoped")
        return
    }

    err := syscall.Kill(s.pid, syscall.SIGTERM)
    if err == nil {
        fmt.Println(logMsg + " stop success")
    } else {
        fmt.Println(logMsg + " stop error: " + err.Error())
    }
}

func (s *serverBasic) Restart() {
    stopStatus := make(chan bool, 1)
    go func() {
        for {
            if !s.checkRunning() {
                stopStatus <- true
                break
            }
            time.Sleep(500 * time.Millisecond)
        }
    }()

    select {
    case status := <-stopStatus:
        if status {
            s.Start(2)
        }
    case <-time.After(10 * time.Second):
        fmt.Println("server " + s.serverTag + " restart timeout")
    }
}

var (
    onceBasic sync.Once
    insBasic  *serverBasic
)

func init() {
    insBasic = &serverBasic{}
    insBasic.appFlag = false
}

func NewBasic() *serverBasic {
    onceBasic.Do(func() {
        insBasic.serverTag = mpf.EnvProjectKey() + strconv.Itoa(mpf.EnvServerPort())
        insBasic.pidFile = mpf.EnvDirRoot() + "/pid/" + insBasic.serverTag + ".pid"
        insBasic.pid = insBasic.getPid()

        conf := mpf.NewConfig().GetConfig("server")
        confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + ".timeout."
        insBasic.timeout = make(map[string]time.Duration)
        insBasic.timeout["read"] = time.Duration(conf.GetInt(confPrefix+"read")) * time.Second
        insBasic.timeout["readheader"] = time.Duration(conf.GetInt(confPrefix+"readheader")) * time.Second
        insBasic.timeout["write"] = time.Duration(conf.GetInt(confPrefix+"write")) * time.Second
        insBasic.timeout["idle"] = time.Duration(conf.GetInt(confPrefix+"idle")) * time.Second
        insBasic.timeout["shutdown"] = time.Duration(conf.GetInt(confPrefix+"shutdown")) * time.Second

        insBasic.server = http.Server{
            ReadTimeout:       insBasic.timeout["read"],
            ReadHeaderTimeout: insBasic.timeout["readheader"],
            WriteTimeout:      insBasic.timeout["write"],
            IdleTimeout:       insBasic.timeout["idle"],
        }
        insBasic.server.SetKeepAlivesEnabled(false)
    })

    return insBasic
}
