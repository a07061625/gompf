/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 1:15
 */
package mpserver

import (
    "strconv"
    "sync"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/spf13/viper"
)

type IServerBase interface {
    init() // 初始化
    AddAppConf(configs ...iris.Configurator)
    SetGlobalMiddleware(isPrefix bool, middlewareList ...context.Handler)
    SetRouters(controllers ...controllers.IControllerBasic)
    ListenNotify() // 监听信号
    Restart()      // 重启服务
    Start()        // 启动服务
    Stop()         // 停止服务
}

type serverBase struct {
    app             *iris.Application   // 应用实例
    appConf         []iris.Configurator // 应用配置
    serverConf      *viper.Viper        // 服务配置
    serverTag       string              // 服务标识
    pid             int                 // 服务进程ID
    pidFile         string              // 服务进程ID文件
    timeoutShutdown time.Duration       // 关闭服务超时时间
}

// 设置全局中间件
//   isPrefix: 中间件类型,true:前置 false:后置
//   middlewareList: 中间件列表
func (s *serverBase) SetGlobalMiddleware(isPrefix bool, middlewareList ...context.Handler) {
    if len(middlewareList) == 0 {
        return
    }

    if isPrefix {
        s.app.UseGlobal(middlewareList...)
    } else {
        s.app.DoneGlobal(middlewareList...)
    }
}

func (s *serverBase) AddAppConf(configs ...iris.Configurator) {
    if len(configs) > 0 {
        s.appConf = append(s.appConf, configs...)
    }
}

func (s *serverBase) refreshConf() {
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    s.appConf = make([]iris.Configurator, 20)
    s.appConf[0] = iris.WithoutStartupLog
    s.appConf[1] = iris.WithoutInterruptHandler
    s.appConf[2] = iris.WithoutPathCorrectionRedirection
    s.appConf[3] = iris.WithOptimizations
    s.appConf[4] = iris.WithoutBodyConsumptionOnUnmarshal
    s.appConf[5] = iris.WithCharset("UTF-8")
    s.appConf[6] = iris.WithoutServerError(iris.ErrServerClosed)
    s.appConf[7] = iris.WithRemoteAddrHeader("X-Real-Ip")
    s.appConf[8] = iris.WithRemoteAddrHeader("X-Forwarded-For")
    s.appConf[9] = iris.WithRemoteAddrHeader("CF-Connecting-IP")
    s.appConf[10] = iris.WithOtherValue("server_host", s.serverConf.GetString(confPrefix+"host"))
    s.appConf[11] = iris.WithOtherValue("server_port", s.serverConf.GetInt(confPrefix+"port"))
    s.appConf[12] = iris.WithOtherValue("server_type", s.serverConf.GetString(confPrefix+"type"))
    s.appConf[13] = iris.WithOtherValue("version_min", s.serverConf.GetString(confPrefix+"version.min"))
    s.appConf[14] = iris.WithOtherValue("version_deprecated", s.serverConf.GetString(confPrefix+"version.deprecated"))
    s.appConf[15] = iris.WithOtherValue("version_current", s.serverConf.GetString(confPrefix+"version.current"))
    s.appConf[16] = iris.WithOtherValue("version_max", s.serverConf.GetString(confPrefix+"version.max"))
    s.appConf[17] = iris.WithOtherValue("timeout_request", s.serverConf.GetFloat64(confPrefix+"timeout.request"))
    s.appConf[18] = iris.WithOtherValue("timeout_controller", s.serverConf.GetFloat64(confPrefix+"timeout.controller"))
    s.appConf[19] = iris.WithOtherValue("timeout_action", s.serverConf.GetFloat64(confPrefix+"timeout.action"))

    // 国际化配置文件只能是以./开始,否则会报错
    s.app.I18n.Load("./configs/i18n/*/*.ini", "zh-CN", "en-US")
    s.app.I18n.PathRedirect = false
    s.app.I18n.URLParameter = s.serverConf.GetString(confPrefix + "reqparam.i18n")

    // 错误码处理
    s.app.OnAnyErrorCode(func(ctx context.Context) {
        statusCode := ctx.GetStatusCode()
        logMsg := "HTTP ERROR CODE: " + strconv.Itoa(statusCode) + " URI: " + ctx.Path()
        result := mpresponse.NewResultProblem()
        result.Title = "服务错误"
        result.Status = statusCode

        switch statusCode {
        case iris.StatusNotFound:
            mplog.LogInfo(logMsg)
            result.Type = "internal-address-not-exist"
            result.Code = errorcode.CommonRequestResourceEmpty
            result.Msg = "接口地址不存在"
        case iris.StatusMethodNotAllowed:
            mplog.LogInfo(logMsg)
            result.Type = "internal-method-not-allow"
            result.Code = errorcode.CommonRequestMethod
            result.Msg = "请求方式不支持"
        default:
            mplog.LogError(logMsg)
            result.Type = "internal-other"
            result.Code = errorcode.CommonBaseServer
            result.Msg = "其他服务错误"
        }
        ctx.Problem(mpresp.GetProblemHandleBasic(result, 30*time.Second))
        mpresp.NewBasicEnd()(ctx)
    })
}

// 初始化
func (s *serverBase) initBase() {
    s.app = iris.New()
    s.timeoutShutdown = time.Duration(s.serverConf.GetInt(mpf.EnvType()+"."+mpf.EnvProjectKeyModule()+"."+"timeout.shutdown")) * time.Second
    s.serverTag = mpf.EnvProjectKey() + strconv.Itoa(mpf.EnvServerPort())
    s.pidFile = mpf.EnvDirRoot() + "/pid/" + s.serverTag + ".pid"
    s.pid = s.getPid()
    s.refreshConf()
}

func newServerBase(conf *viper.Viper) serverBase {
    s := serverBase{}
    s.serverConf = conf
    return s
}

var (
    once sync.Once
    ins  IServerBase
)

func NewServer(conf *viper.Viper) IServerBase {
    once.Do(func() {
        if mpf.EnvServerType() == mpf.EnvServerTypeApi {
            ins = &serverHttp{newServerBase(conf)}
            ins.init()
        }
    })

    return ins
}
