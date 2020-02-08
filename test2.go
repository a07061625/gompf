package main

import (
    "os"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpframe/controllers/frontend"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpreq"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mpserver"
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
    conf := mpf.NewConfig().GetConfig("server")
    server := mpserver.NewBasicHttp(conf)
    server.SetMwGlobal(true, mpreq.NewBasicLog(), mpreq.NewBasicRecover())
    server.SetMwGlobal(false, mpresp.NewBasicClear())
    server.SetRoute(frontend.NewIndex())
    server.StartServer()
}
