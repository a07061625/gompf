package main

import (
    "os"

    "log"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpframe"
    "github.com/a07061625/gompf/mpf/mpserver"
    "github.com/kataras/iris/v12"
)

func init() {
    dirRoot, _ := os.Getwd()
    bs := mpf.NewBootstrap()
    bs.SetDirRoot(dirRoot)
    bs.SetDirConfigs(dirRoot + "/configs")
    bs.SetDirLogs(dirRoot + "/logs")
    mpf.LoadBoot(bs)
}

func main() {
    outer := mpframe.NewOuterHttp()
    server := mpserver.NewServerHttp()
    server.SetOuter(outer)
    server.App.Get("/", func(ctx iris.Context) {
        ctx.HTML("<h1>Hello World!</h1>")
    })

    go outer.GetNotify(server.App)()

    if server.Runner == nil {
        log.Fatalln("Runner")
    }
    if len(server.RunConfigs) == 0 {
        log.Fatalln("RunConfigs")
    }
    server.App.Run(server.Runner, server.RunConfigs...)
}
