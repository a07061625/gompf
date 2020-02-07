package main

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/frame"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/globalprefix"
    "github.com/a07061625/gompf/mpf/mpframe/mpserv"
    "github.com/a07061625/gompf/mpf/mplog"
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
    outer := mpserv.NewSimpleOuter()
    outer.AddGlobalMiddleware(frame.MWEventGlobalPrefix, globalprefix.NewSimpleRequestLog(), globalprefix.NewSimpleRecover())
    server := mpserver.NewServerHttp()
    server.SetOuter(outer)
    server.App.Get("/", func(ctx iris.Context) {
        mplog.LogError("wwwwww")
        ctx.HTML("<h1>Hello World!</h1>")
        ctx.Next()
    })
    server.App.Any("/{directory:path}", func(ctx iris.Context) {
        result := mpresponse.NewResult()
        directory := ctx.Params().Get("directory")
        if directory == "error/500" {
            result.Code = errorcode.CommonBaseServer
            result.Msg = "服务出错"
        } else {
            mplog.LogInfo("uri: /" + directory + " not exist")
            result.Code = errorcode.CommonRequestResourceEmpty
            result.Msg = "接口不存在"
        }
        ctx.JSON(result)
        ctx.Next()
    })

    server.StartServer()
}
