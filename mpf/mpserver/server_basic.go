package mpserver

import (
    "net"
    "net/http"
    "os"
    "sync"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mplog"
)

type IServerBasic interface {
    ForkProcess(serverMode string) int
    Start()
    Shutdown(timeout time.Duration)
}

type serverBasic struct {
    http.Server
    Listener    *serverListener
    connManager *connManager
}

// 创建服务进程 daemon:守护进程 child:子进程
func (s *serverBasic) ForkProcess(serverMode string) int {
    os.Setenv(mpf.GoEnvServerMode, serverMode)

    l := NewListener(s.Server.Addr)
    files := make([]uintptr, 0)
    if serverMode == mpf.EnvServerModeChild {
        lFd, err := l.File()
        if err != nil {
            mplog.LogFatal(serverMode + " server fork error: " + err.Error())
        }
        files = []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), lFd}
    } else {
        files = []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
    }

    attr := &syscall.ProcAttr{
        Env:   os.Environ(),
        Files: files,
    }
    pid, err := syscall.ForkExec(os.Args[0], os.Args, attr)
    if err != nil {
        mplog.LogFatal(serverMode + " server fork error: " + err.Error())
    }
    s.Listener = l

    return pid
}

func (s *serverBasic) Start() {
    err := s.Server.Serve(s.Listener)
    if err != nil {
        mplog.LogFatal(os.Getenv(mpf.GoEnvServerMode) + " server start error: " + err.Error())
    }
}

// 关闭服务
func (s *serverBasic) Shutdown(timeout time.Duration) {
    s.SetKeepAlivesEnabled(false)
    s.connManager.Close(timeout)
    logMsg := os.Getenv(mpf.GoEnvServerMode) + " server shutdown "
    err := s.Listener.Close()
    if err != nil {
        mplog.LogError(logMsg + "error: " + err.Error())
    } else {
        mplog.LogInfo(logMsg + "success")
    }
}

var (
    once sync.Once
    ins  *serverBasic
)

func New(addr string, app http.Handler) *serverBasic {
    once.Do(func() {
        ins = &serverBasic{
            Server: http.Server{
                Addr:    addr,
                Handler: app,
            },
            connManager: newConnManager(),
        }
        ins.ConnState = func(conn net.Conn, state http.ConnState) {
            switch state {
            case http.StateNew:
                ins.connManager.Add(1)
            case http.StateActive:
                ins.connManager.RemoveIdleConn(conn.LocalAddr().String())
            case http.StateIdle:
                ins.connManager.AddIdleConn(conn.LocalAddr().String(), conn)
            case http.StateHijacked, http.StateClosed:
                ins.connManager.Done()
            }
        }
    })

    return ins
}
