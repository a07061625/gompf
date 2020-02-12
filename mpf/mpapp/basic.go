package mpapp

import (
    "strconv"
    "time"

    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

type appBasic struct {
    initFlag    bool                   // 初始化标识 true:已初始化 false:未初始化
    instance    *iris.Application      // 应用实例
    confApp     []iris.Configurator    // 应用配置
    confOther   map[string]interface{} // 应用附加数据
    mwPrefix    []context.Handler      // 前置全局中间件
    mwSuffix    []context.Handler      // 后置全局中间件
    routeBlocks map[string]string      // 路由支持的版块列表
}

func (app *appBasic) GetInstance() *iris.Application {
    return app.instance
}

// 设置全局中间件
//   isPrefix: 中间件类型,true:前置 false:后置
//   middlewareList: 中间件列表
func (app *appBasic) SetMiddleware(isPrefix bool, middlewareList ...context.Handler) {
    if app.initFlag {
        return
    }

    if isPrefix {
        app.mwPrefix = middlewareList
    } else {
        app.mwSuffix = middlewareList
    }
}

func (app *appBasic) Build() {
    app.initFlag = true
    app.initConf()
    app.instance.UseGlobal(app.mwPrefix...)
    app.instance.DoneGlobal(app.mwSuffix...)

    // 错误码处理
    app.instance.OnAnyErrorCode(func(ctx context.Context) {
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

    err := app.instance.Build()
    if err != nil {
        mplog.LogFatal("app build error: " + err.Error())
    }
}

var (
    app *appBasic
)

func init() {
    app = &appBasic{}
    app.initFlag = false
    app.instance = iris.New()
    app.confApp = make([]iris.Configurator, 10)
    app.confApp[0] = iris.WithoutStartupLog
    app.confApp[1] = iris.WithoutInterruptHandler
    app.confApp[2] = iris.WithoutPathCorrectionRedirection
    app.confApp[3] = iris.WithOptimizations
    app.confApp[4] = iris.WithoutBodyConsumptionOnUnmarshal
    app.confApp[5] = iris.WithCharset("UTF-8")
    app.confApp[6] = iris.WithoutServerError(iris.ErrServerClosed)
    app.confApp[7] = iris.WithRemoteAddrHeader("X-Real-Ip")
    app.confApp[8] = iris.WithRemoteAddrHeader("X-Forwarded-For")
    app.confApp[9] = iris.WithRemoteAddrHeader("CF-Connecting-IP")
    app.confOther = make(map[string]interface{})
    app.mwPrefix = make([]context.Handler, 0)
    app.mwSuffix = make([]context.Handler, 0)
    app.routeBlocks = make(map[string]string)
}

func New() *appBasic {
    return app
}
