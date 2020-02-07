/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 10:17
 */
package mpserver

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mwreq"
    "github.com/kataras/iris/v12"
)

type serverHttp struct {
    serverWeb
}

func (s *serverHttp) initConfig() {
    conf := mpf.NewConfig().GetConfig("server")
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_host", conf.GetString(confPrefix+"host")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_port", conf.GetInt(confPrefix+"port")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_type", conf.GetString(confPrefix+"type")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_request", conf.GetFloat64(confPrefix+"timeout.request")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_action", conf.GetFloat64(confPrefix+"timeout.action")))

    s.runConfigs = append(s.runConfigs, iris.WithCharset("UTF-8"))
    s.runConfigs = append(s.runConfigs, iris.WithoutStartupLog)
    s.runConfigs = append(s.runConfigs, iris.WithOptimizations)
    s.runConfigs = append(s.runConfigs, iris.WithoutInterruptHandler)
    s.runConfigs = append(s.runConfigs, iris.WithoutServerError(iris.ErrServerClosed))
}

func (s *serverHttp) initMiddleware() {
    s.App.UseGlobal(mwreq.NewIrisBefore())
}

func (s *serverHttp) StartServer() {
    s.initConfig()
    s.initMiddleware()
    s.baseStart()
}

func NewServerHttp() *serverHttp {
    return &serverHttp{newServerWeb()}
}
