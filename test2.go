package main

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpframe"
    "github.com/a07061625/gompf/mpf/mpresponse"
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
    server.BindErrHandles()
    server.App.Get("/", func(ctx iris.Context) {
        ctx.HTML("<h1>Hello World!</h1>")
        ctx.Next()
    })
    server.App.Get("/error/404", func(ctx iris.Context) {
        result := mpresponse.NewResult()
        result.Code = errorcode.CommonRequestResourceEmpty
        result.Msg = "接口不存在"
        ctx.JSON(result)
        ctx.Next()
    })
    server.App.Get("/error/500", func(ctx iris.Context) {
        result := mpresponse.NewResult()
        result.Code = errorcode.CommonBaseServer
        result.Msg = "服务出错"
        ctx.JSON(result)
        ctx.Next()
    })

    server.StartServer()
}
