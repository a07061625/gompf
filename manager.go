package main

import (
    "io/ioutil"
    "os"
    "strconv"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf/mplog"
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

//// 捕获系统信号
//func listenNotify() {
//    signals := make(chan os.Signal)
//    signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
//
//    for {
//        sig := <-signals
//        switch sig {
//        case syscall.SIGUSR2: // 重启
//            pid = server.ForkProcess(mpf.EnvServerModeChild)
//            savePid(pid)
//            server.Start()
//        case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
//            server.Shutdown(timeoutShutdown)
//            os.Exit(1)
//        }
//    }
//}

//func manageServer() {
//    if os.Getenv(mpf.GoEnvServerMode) != mpf.EnvServerModeChild { // master
//        switch *optionType {
//       case "start":
//           if checkRunning() {
//               fmt.Println("server " + serverTag + " is running")
//               return
//           }
//
//           pid = server.ForkProcess(mpf.EnvServerModeDaemon)
//           savePid(pid)
//           server.Start()
//       case "restart":
//           stopStatus := make(chan bool, 1)
//           go func() {
//               for {
//                   if !checkRunning() {
//                       stopStatus <- true
//                       break
//                   }
//                   time.Sleep(500 * time.Millisecond)
//               }
//           }()
//
//           select {
//           case status := <-stopStatus:
//               if status {
//                   syscall.Kill(pid, syscall.SIGHUP) //kill -HUP, daemon only时会直接退出
//
//                   pid = server.ForkProcess(mpf.EnvServerModeDaemon)
//                   savePid(pid)
//                   server.Start()
//               }
//           case <-time.After(10 * time.Second):
//               fmt.Println("server " + serverTag + " restart timeout")
//           }
//        case "stop":
//           if checkRunning() {
//               err := syscall.Kill(pid, syscall.SIGTERM)
//               if err != nil {
//                   fmt.Println("server " + serverTag + " stop error: " + err.Error())
//               }
//           } else {
//               fmt.Println("server " + serverTag + " already stop")
//           }
//        default:
//            fmt.Println("操作类型必须是以下其一: start|stop|restart")
//        }
//        // 主进程退出
//        os.Exit(0)
//    }
//}

var (
    pid             = 0           // 进程ID
    pidFile         = ""          // 进程ID文件
    timeoutShutdown time.Duration // 进程关闭超时时间
    //server          mpserver.IServerBasic // 服务实例
    serverTag = "" // 服务标识
)
