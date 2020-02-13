package mpserver

import (
    "net"
    "os"
    "syscall"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/valyala/tcplisten"
)

type listenerTcp struct {
    net.Listener
}

func (l *listenerTcp) Accept() (conn net.Conn, err error) {
    tl, err := l.Listener.(*net.TCPListener).AcceptTCP()
    if err != nil {
        return nil, err
    }

    conf := newTcpConf()
    tl.SetKeepAlive(conf.keepAlive)
    tl.SetKeepAlivePeriod(conf.keepAlivePeriod)
    tl.SetNoDelay(conf.noDelay)
    tl.SetLinger(conf.linger)
    tl.SetDeadline(conf.deadline)
    tl.SetReadDeadline(conf.readDeadline)
    tl.SetWriteDeadline(conf.writeDeadline)
    return tl, nil
}

//获取sock文件句柄
func (l *listenerTcp) File() (uintptr, error) {
    f, err := l.Listener.(*net.TCPListener).File()
    if err != nil {
        return 0, err
    }
    return f.Fd(), nil
}

func NewListenerTcp(addr string) *listenerTcp {
    listener := &listenerTcp{}
    if os.Getenv(mpf.GoEnvServerMode) == mpf.EnvServerModeChild { // 子进程
        f := os.NewFile(3, "")
        l, err := net.FileListener(f)
        syscall.Kill(syscall.Getppid(), syscall.SIGTERM) //发信号给父进程,让父进程停止服务
        if err != nil {
            mplog.LogFatal("child server create listener error: " + err.Error())
        }
        listener.Listener = l
    } else { // 守护进程
        conf := newTcpConf()
        cfg := tcplisten.Config{
            ReusePort:   conf.reusePort,
            DeferAccept: conf.deferAccept,
            FastOpen:    conf.fastOpen,
            Backlog:     conf.backlog,
        }
        l, err := cfg.NewListener("tcp4", addr)
        if err != nil {
            mplog.LogFatal("daemon server create listener error: " + err.Error())
        }
        listener.Listener = l
    }

    return listener
}
