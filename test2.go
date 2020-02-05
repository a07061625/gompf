package main

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/kataras/iris/v12"
    "github.com/valyala/tcplisten"
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
    app := iris.New()
    app.Get("/", func(ctx iris.Context) {
        ctx.HTML("<h1>Hello World!</h1>")
    })

    listenerCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    l, err := listenerCfg.NewListener("tcp4", ":8080")
    if err != nil {
        mpf.NewLogger().Error(err.Error())
    }

    app.Run(iris.Listener(l))
}
