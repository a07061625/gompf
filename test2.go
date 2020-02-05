package main

import (
    "log"
    "os"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/kataras/iris/v12"
    "github.com/valyala/tcplisten"
)

func init() {
    dir, err := os.Getwd()
    if err != nil {
        log.Fatalln("get root dir error")
    }

    dirRoot := strings.Replace(dir, "\\", "/", -1)
    configDir := dirRoot + "/configs"
    mpf.LoadConfig(configDir)
    serviceConfig := mpf.NewConfig().GetConfig("service")
    mpf.LoadEnv(serviceConfig)
    projectConfig := mpf.NewConfig().GetConfig("project")
    project.LoadProject(projectConfig)
    logDir := dirRoot + "/logs"
    logConfig := mpf.NewConfig().GetConfig("log")
    mpf.LoadLog(logDir, logConfig)
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

    l, err := listenerCfg.NewListener("tcp", ":8080")
    if err != nil {
        app.Logger().Fatal(err)
    }

    app.Run(iris.Listener(l))
}
