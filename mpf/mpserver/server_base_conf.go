/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 2:09
 */
package mpserver

import (
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/kataras/iris/v12"
)

func (s *serverBase) AddIrisConf(configs ...iris.Configurator) {
    if len(configs) > 0 {
        s.confIris = append(s.confIris, configs...)
    }
}

func (s *serverBase) initConf() {
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    s.confIris = make([]iris.Configurator, 20)
    s.confIris[0] = iris.WithoutStartupLog
    s.confIris[1] = iris.WithoutInterruptHandler
    s.confIris[2] = iris.WithoutPathCorrectionRedirection
    s.confIris[3] = iris.WithOptimizations
    s.confIris[4] = iris.WithoutBodyConsumptionOnUnmarshal
    s.confIris[5] = iris.WithCharset("UTF-8")
    s.confIris[6] = iris.WithoutServerError(iris.ErrServerClosed)
    s.confIris[7] = iris.WithRemoteAddrHeader("X-Real-Ip")
    s.confIris[8] = iris.WithRemoteAddrHeader("X-Forwarded-For")
    s.confIris[9] = iris.WithRemoteAddrHeader("CF-Connecting-IP")
    s.confIris[10] = iris.WithOtherValue("server_host", s.confServer.GetString(confPrefix+"host"))
    s.confIris[11] = iris.WithOtherValue("server_port", s.confServer.GetInt(confPrefix+"port"))
    s.confIris[12] = iris.WithOtherValue("server_type", s.confServer.GetString(confPrefix+"type"))
    s.confIris[13] = iris.WithOtherValue("version_min", s.confServer.GetString(confPrefix+"version.min"))
    s.confIris[14] = iris.WithOtherValue("version_deprecated", s.confServer.GetString(confPrefix+"version.deprecated"))
    s.confIris[15] = iris.WithOtherValue("version_current", s.confServer.GetString(confPrefix+"version.current"))
    s.confIris[16] = iris.WithOtherValue("version_max", s.confServer.GetString(confPrefix+"version.max"))
    s.confIris[17] = iris.WithOtherValue("timeout_request", s.confServer.GetFloat64(confPrefix+"timeout.request"))
    s.confIris[18] = iris.WithOtherValue("timeout_controller", s.confServer.GetFloat64(confPrefix+"timeout.controller"))
    s.confIris[19] = iris.WithOtherValue("timeout_action", s.confServer.GetFloat64(confPrefix+"timeout.action"))

    s.timeoutShutdown = time.Duration(s.confServer.GetInt(confPrefix+"timeout.shutdown")) * time.Second

    // 国际化配置文件只能是以./开始,否则会报错
    s.app.I18n.Load("./configs/i18n/*/*.ini", "zh-CN", "en-US")
    s.app.I18n.PathRedirect = false
    s.app.I18n.URLParameter = s.confServer.GetString(confPrefix + "reqparam.i18n")
}
