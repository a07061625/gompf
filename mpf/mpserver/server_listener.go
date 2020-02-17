package mpserver

import (
    "net"
    "os"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/valyala/tcplisten"
)

type serverListener struct {
    net.Listener
}

func (sl *serverListener) Accept() (conn net.Conn, err error) {
    tc, err := sl.Listener.(*net.TCPListener).AcceptTCP()
    if err != nil {
        return nil, err
    }
    tc.SetKeepAlive(true)
    tc.SetKeepAlivePeriod(3 * time.Minute)
    tc.SetNoDelay(true)
    return tc, nil
}

// 获取sock文件句柄
func (sl *serverListener) File() (uintptr, error) {
    f, err := sl.Listener.(*net.TCPListener).File()
    if err != nil {
        return 0, err
    }
    return f.Fd(), nil
}

func NewListener(addr string) *serverListener {
    listener := &serverListener{}
    if os.Getenv(mpf.GoEnvServerMode) == mpf.EnvServerModeChild { // 子进程
        f := os.NewFile(3, "")
        l, err := net.FileListener(f)
        syscall.Kill(syscall.Getppid(), syscall.SIGTERM) // 发信号给父进程,让父进程停止服务
        if err != nil {
            mplog.LogFatal("child server create listener error: " + err.Error())
        }
        listener.Listener = l
    } else { // 守护进程
        cfg := tcplisten.Config{
            ReusePort:   true,
            DeferAccept: true,
            FastOpen:    true,
        }
        l, err := cfg.NewListener("tcp4", addr)
        if err != nil {
            mplog.LogFatal("daemon server create listener error: " + err.Error())
        }
        listener.Listener = l
    }

    return listener
}
