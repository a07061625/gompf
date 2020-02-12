package mpserver

import (
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "strconv"
    "syscall"
    "time"
)

var (
    TimeDeadLine = 10 * time.Second
    srv          *server
    appName      string
    pidFile      string
    pidVal       int
)

//improvement http.Server
type server struct {
    http.Server
    listener *listener
    cm       *ConnManager
}

//用来重载net.Listener的方法
type listener struct {
    net.Listener
    server *server
}

func init() {
    file, _ := filepath.Abs(os.Args[0])
    appPath := filepath.Dir(file)
    appName = filepath.Base(file)
    pidFile = appPath + "/" + appName + ".pid"
    if os.Getenv("__Daemon") != "true" { //master
        cmd := "start" //缺省为start
        if l := len(os.Args); l > 1 {
            cmd = os.Args[l-1]
        }
        switch cmd {
        case "start":
            if isRunning() {
                log.Printf("[%d] %s is running\n", pidVal, appName)
            } else { //fork daemon进程
                if err := forkDaemon(); err != nil {
                    log.Fatal(err)
                }
            }
        case "restart": //重启:
            if !isRunning() {
                log.Printf("%s not running\n", appName)
            } else {
                log.Printf("[%d] %s restart now\n", pidVal, appName)
                restart(pidVal)
            }
        case "stop": //停止
            if !isRunning() {
                log.Printf("%s not running\n", appName)
            } else {

                syscall.Kill(pidVal, syscall.SIGTERM) //kill
            }
        case "-h":
            fmt.Printf("Usage: %s start|restart|stop\n", appName)
        default: //其它不识别的参数
            return //返回至调用方
        }
        //主进程退出
        os.Exit(0)
    }
    go handleSignals()
}

//检查pidFile是否存在以及文件里的pid是否存活
func isRunning() bool {
    if mf, err := os.Open(pidFile); err == nil {
        pid, _ := ioutil.ReadAll(mf)
        pidVal, _ = strconv.Atoi(string(pid))
    }
    running := false
    if pidVal > 0 {
        if err := syscall.Kill(pidVal, 0); err == nil { //发一个信号为0到指定进程ID，如果没有错误发生，表示进程存活
            running = true
        }
    }
    return running
}

//保存pid
func savePid(pid int) error {
    file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
    if err != nil {
        return err
    }
    defer file.Close()
    file.WriteString(strconv.Itoa(pid))
    return nil
}

//捕获系统信号
func handleSignals() {
    signals := make(chan os.Signal)
    signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
    var err error
    for {
        sig := <-signals
        switch sig {
        case syscall.SIGHUP: //重启
            if srv != nil {
                err = srv.fork()
            } else { //only deamon时不支持kill -HUP,因为可能监听地址会占用
                log.Printf("[%d] %s stopped.", os.Getpid(), appName)
                os.Remove(pidFile)
                os.Exit(2)
            }
            if err != nil {
                log.Fatalln(err)
            }
        case syscall.SIGINT:
            fallthrough
        case syscall.SIGTERM:
            log.Printf("[%d] %s stop graceful", os.Getpid(), appName)
            if srv != nil {
                srv.shutdown()
            } else {
                log.Printf("[%d] %s stopped.", os.Getpid(), appName)
            }
            os.Exit(1)
        }
    }
}

//forkDaemon,当checkPid为true时，检查是否有存活的，有则不执行
func forkDaemon() error {
    args := os.Args
    os.Setenv("__Daemon", "true")
    procAttr := &syscall.ProcAttr{
        Env:   os.Environ(),
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
    }
    pid, err := syscall.ForkExec(os.Args[0], args, procAttr)
    if err != nil {
        return err
    }
    log.Printf("[%d] %s start daemon\n", pid, appName)
    savePid(pid)
    return nil
}

//重启(先发送kill -HUP到运行进程，手工重启daemon ...当有运行的进程时，daemon不启动)
func restart(pid int) {
    syscall.Kill(pid, syscall.SIGHUP) //kill -HUP, daemon only时，会直接退出
    fork := make(chan bool, 1)
    go func() { //循环，查看pidFile是否存在，不存在或值已改变，发送消息
        for {
            f, err := os.Open(pidFile)
            if err != nil || os.IsNotExist(err) { //文件已不存在
                fork <- true
                break
            } else {
                pidVal, _ := ioutil.ReadAll(f)
                if strconv.Itoa(pid) != string(pidVal) {
                    fork <- false
                    break
                }
            }
            time.Sleep(500 * time.Millisecond)
        }
    }()
    //处理结果
    select {
    case r := <-fork:
        if r {
            forkDaemon()
        }
    case <-time.After(time.Second * 5):
        log.Fatalln("restart timeout")
    }

}

//处理http.Server，使支持graceful stop/restart
func Graceful(s http.Server) error {
    os.Setenv("__GRACEFUL", "true")
    srv = &server{
        cm:     newConnManager(),
        Server: s,
    }
    srv.ConnState = func(conn net.Conn, state http.ConnState) {
        switch state {
        case http.StateNew:
            srv.cm.add(1)
        case http.StateActive:
            srv.cm.rmIdleConn(conn.LocalAddr().String())
        case http.StateIdle:
            srv.cm.addIdleConn(conn.LocalAddr().String(), conn)
        case http.StateHijacked, http.StateClosed:
            srv.cm.done()
        }
    }
    l, err := srv.getListener()
    if err == nil {
        err = srv.Server.Serve(l)
    }
    return err
}

//使用addr和handler来启动一个支持graceful的服务
func GracefulServe(addr string, handler http.Handler) error {
    s := http.Server{
        Addr:    addr,
        Handler: handler,
    }
    return Graceful(s)
}

//获取listener
func (this *server) getListener() (*listener, error) {
    var l net.Listener
    var err error
    if os.Getenv("_GRACEFUL_RESTART") == "true" { //grace restart出来的进程，从FD FILE获取
        f := os.NewFile(3, "")
        l, err = net.FileListener(f)
        syscall.Kill(syscall.Getppid(), syscall.SIGTERM) //发信号给父进程，让父进程停止服务
    } else { //初始启动，监听addr
        l, err = net.Listen("tcp", this.Addr)
    }
    if err == nil {
        this.listener = &listener{
            Listener: l,
            server:   this,
        }
    }
    return this.listener, err
}

//fork一个新的进程
func (this *server) fork() error {
    os.Setenv("_GRACEFUL_RESTART", "true")
    lFd, err := this.listener.File()
    if err != nil {
        return err
    }
    execSpec := &syscall.ProcAttr{
        Env:   os.Environ(),
        Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), lFd},
    }
    pid, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
    if err != nil {
        return err
    }
    savePid(pid)
    log.Printf("[%d] %s fork ok\n", pid, appName)
    return nil
}

//关闭服务
func (this *server) shutdown() {
    this.SetKeepAlivesEnabled(false)
    this.cm.close(TimeDeadLine)
    this.listener.Close()
    log.Printf("[%d] %s stopped.", os.Getpid(), appName)
}
