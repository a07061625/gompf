/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 2:10
 */
package mpserver

import (
    stdContext "context"
    "io/ioutil"
    "net"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
    "github.com/valyala/tcplisten"
)

func (s *serverBase) getPid() int {
    pid := 0
    if f, err := os.Open(s.pidFile); err == nil {
        pidStr, _ := ioutil.ReadAll(f)
        pid, _ = strconv.Atoi(string(pidStr))
        defer f.Close()
    }

    return pid
}

func (s *serverBase) savePid(pid int) {
    f, err := os.OpenFile(s.pidFile, os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        mplog.LogInfo("write pid file error: " + err.Error())
        return
    }
    defer f.Close()
    f.WriteString(strconv.Itoa(pid))
}

// 发一个信号为0到指定进程ID,如果没有错误发生,表示进程存活
func (s *serverBase) checkRunning() bool {
    if s.pid <= 0 {
        return false
    }

    err := syscall.Kill(s.pid, 0)
    return err == nil
}

// 创建子进程
func (s *serverBase) forkChild() (net.Listener, error) {
    os.Setenv(mpf.GoEnvServerMode, mpf.EnvServerModeChild)
    f := os.NewFile(3, "")
    listener, err := net.FileListener(f)
    if err != nil {
        return nil, err
    }

    syscall.Kill(syscall.Getppid(), syscall.SIGTERM) // 发信号给父进程，让父进程停止服务
    processAttr := &syscall.ProcAttr{
        Env:   os.Environ(),
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), f.Fd()},
    }
    pid, err := syscall.ForkExec(os.Args[0], os.Args, processAttr)
    if err != nil {
        return nil, err
    }

    s.savePid(pid)
    s.pid = pid
    return listener, nil
}

// 创建守护进程
func (s *serverBase) forkDaemon() (net.Listener, error) {
    os.Setenv(mpf.GoEnvServerMode, mpf.EnvServerModeDaemon)
    listenCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    listener, err := listenCfg.NewListener("tcp4", mpf.EnvServerDomain())
    if err != nil {
        return nil, err
    }
    processAttr := &syscall.ProcAttr{
        Env:   os.Environ(),
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
    }
    pid, err := syscall.ForkExec(os.Args[0], os.Args, processAttr)
    if err != nil {
        return nil, err
    }

    s.savePid(pid)
    s.pid = pid
    return listener, nil
}

func (s *serverBase) startBase() {
    if s.checkRunning() {
        mplog.LogInfo("server " + s.serverTag + " is running")
        return
    }

    listener, err := s.forkDaemon()
    if err != nil {
        mplog.LogFatal("server " + s.serverTag + " start error: " + err.Error())
    }

    err = s.app.Run(iris.Listener(listener), s.appConf...)
    if err != nil {
        mplog.LogFatal("server " + s.serverTag + " start run error: " + err.Error())
    }
}

func (s *serverBase) stopBase() {
    if s.checkRunning() {
        err := syscall.Kill(s.getPid(), syscall.SIGTERM)
        if err != nil {
            mplog.LogFatal("server " + s.serverTag + " stop error: " + err.Error())
        }
    } else {
        mplog.LogInfo("server " + s.serverTag + " already stop")
    }
}

func (s *serverBase) restartBase() {
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
            syscall.Kill(s.pid, syscall.SIGHUP) //kill -HUP, daemon only时会直接退出

            listener, err := s.forkDaemon()
            if err != nil {
                mplog.LogFatal("server " + s.serverTag + " restart error: " + err.Error())
            }

            err = s.app.Run(iris.Listener(listener), s.appConf...)
            if err != nil {
                mplog.LogFatal("server " + s.serverTag + " restart run error: " + err.Error())
            }
        }
    case <-time.After(10 * time.Second):
        mplog.LogFatal("server " + s.serverTag + " restart timeout")
    }
}

func (s *serverBase) ListenNotify() {
    signals := make(chan os.Signal)
    signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
    for {
        sig := <-signals
        switch sig {
        case syscall.SIGUSR2: //重启
            listener, err := s.forkChild()
            if err != nil {
                mplog.LogFatal("server " + s.serverTag + " reload error: " + err.Error())
            }

            err = s.app.Run(iris.Listener(listener), s.appConf...)
            if err != nil {
                mplog.LogFatal("server " + s.serverTag + " reload run error: " + err.Error())
            }
            os.Exit(2)
        case syscall.SIGHUP:
            fallthrough
        case syscall.SIGINT:
            fallthrough
        case syscall.SIGTERM:
            ctx, _ := stdContext.WithTimeout(stdContext.Background(), s.timeoutShutdown)
            s.app.Shutdown(ctx)
            mplog.LogInfo("server " + s.serverTag + " shutdown")
            os.Exit(1)
        }
    }
}
