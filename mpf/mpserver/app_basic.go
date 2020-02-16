package mpserver

import (
    "reflect"
    "regexp"
    "runtime"
    "strconv"
    "strings"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpresp"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/kataras/iris/v12/i18n"
)

type appBasic struct {
    instance *iris.Application // 应用实例
}

func (app *appBasic) GetInstance() *iris.Application {
    return app.instance
}

// 设置全局前置中间件
func (app *appBasic) SetGlobalMiddlewarePrefix(middlewareList ...context.Handler) {
    app.instance.UseGlobal(middlewareList...)
}

// 设置全局后置中间件
func (app *appBasic) SetGlobalMiddlewareSuffix(middlewareList ...context.Handler) {
    app.instance.DoneGlobal(middlewareList...)
}

// 设置应用配置
func (app *appBasic) SetConfApp() {
    conf := mpf.NewConfig().GetConfig("server")
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    appConf := make([]iris.Configurator, 20)
    appConf[0] = iris.WithoutStartupLog
    appConf[1] = iris.WithoutInterruptHandler
    appConf[2] = iris.WithoutPathCorrectionRedirection
    appConf[3] = iris.WithOptimizations
    appConf[4] = iris.WithoutBodyConsumptionOnUnmarshal
    appConf[5] = iris.WithCharset("UTF-8")
    appConf[6] = iris.WithoutServerError(iris.ErrServerClosed)
    appConf[7] = iris.WithRemoteAddrHeader("X-Real-Ip")
    appConf[8] = iris.WithRemoteAddrHeader("X-Forwarded-For")
    appConf[9] = iris.WithRemoteAddrHeader("CF-Connecting-IP")
    appConf[10] = iris.WithOtherValue("server_host", conf.GetString(confPrefix+"host"))
    appConf[11] = iris.WithOtherValue("server_port", conf.GetInt(confPrefix+"port"))
    appConf[12] = iris.WithOtherValue("server_type", conf.GetString(confPrefix+"type"))
    appConf[13] = iris.WithOtherValue("version_min", conf.GetString(confPrefix+"version.min"))
    appConf[14] = iris.WithOtherValue("version_deprecated", conf.GetString(confPrefix+"version.deprecated"))
    appConf[15] = iris.WithOtherValue("version_current", conf.GetString(confPrefix+"version.current"))
    appConf[16] = iris.WithOtherValue("version_max", conf.GetString(confPrefix+"version.max"))
    appConf[17] = iris.WithOtherValue("timeout_request", conf.GetFloat64(confPrefix+"timeout.request"))
    appConf[18] = iris.WithOtherValue("timeout_controller", conf.GetFloat64(confPrefix+"timeout.controller"))
    appConf[19] = iris.WithOtherValue("timeout_action", conf.GetFloat64(confPrefix+"timeout.action"))
    app.instance.Configure(appConf...)
}

// 设置国际化配置
func (app *appBasic) SetConfI18n(conf *i18n.I18n) {
    app.instance.I18n = conf
}

func (app *appBasic) formatUri(name string) string {
    match, _ := regexp.MatchString(`^[a-zA-Z]+$`, name)
    if !match {
        return ""
    }

    reg := regexp.MustCompile(`([A-Z])`)
    needStr := reg.ReplaceAllString(name, `-${1}`)
    return strings.ToLower(strings.TrimPrefix(needStr, "-"))
}

func (app *appBasic) registerRouteAction(groupUri string, controller controllers.IControllerBasic) {
    refControllerVal := reflect.ValueOf(controller)
    methodNum := refControllerVal.NumMethod()
    if methodNum <= 0 {
        return
    }

    groupRoute := app.instance.Party(groupUri)
    // 无需显式调用ctx.Next(),自动触发下一个handle
    groupRoute.SetExecutionRules(iris.ExecutionRules{
        Begin: iris.ExecutionOptions{true},
        Main:  iris.ExecutionOptions{true},
        Done:  iris.ExecutionOptions{true},
    })
    groupRoute.Use(controller.GetMwController(true)...)
    groupRoute.Done(controller.GetMwController(false)...)

    refControllerType := reflect.TypeOf(controller)
    for i := 0; i < methodNum; i++ {
        funcName := runtime.FuncForPC(refControllerType.Method(i).Func.Pointer()).Name()
        funcNameList := strings.Split(funcName, ".")
        methodName := funcNameList[len(funcNameList)-1]
        if len(methodName) <= 6 {
            continue
        }
        if !strings.HasPrefix(methodName, "Action") {
            continue
        }

        refAction := refControllerVal.Method(i)
        _, ok := refAction.Interface().(func(ctx context.Context) interface{})
        if !ok {
            continue
        }

        actionTag := app.formatUri(strings.TrimPrefix(methodName, "Action"))
        if len(actionTag) == 0 {
            continue
        }

        actionMwList := controller.GetMwAction(true, actionTag)
        actionMwList = append(actionMwList, func(ctx context.Context) {
            args := []reflect.Value{reflect.ValueOf(ctx)}
            callRes := refAction.Call(args)
            actionRes := callRes[0].Interface()
            if actionRes != nil {
                ctx.Values().Set(project.DataParamKeyRespData, actionRes)
            }
        })
        actionMwList = append(actionMwList, controller.GetMwAction(false, actionTag)...)
        actionUri := "/" + actionTag + " /" + actionTag + "/{directory:path}"
        groupRoute.HandleMany(iris.MethodGet+" "+iris.MethodPost, actionUri, actionMwList...)
    }
}

// 设置路由
func (app *appBasic) SetRouters(blocks map[string]string, controllers []controllers.IControllerBasic) {
    if len(blocks) == 0 {
        return
    }

    controllerNum := len(controllers)
    if controllerNum <= 0 {
        return
    }

    uriPrefix := ""
    for i := 0; i < controllerNum; i++ {
        objType := reflect.TypeOf(controllers[i])
        typeNameList := strings.Split(objType.String(), ".")
        if len(typeNameList) < 2 {
            continue
        }

        // 校验版块
        packageName := strings.TrimPrefix(typeNameList[0], "*")
        match, _ := regexp.MatchString(`^[a-z]+$`, packageName)
        if !match {
            continue
        }
        _, ok := blocks[packageName]
        if !ok {
            continue
        }
        uriPrefix = "/" + packageName

        // 校验控制器
        if !strings.HasSuffix(typeNameList[1], "Controller") {
            continue
        }
        controllerUri := app.formatUri(strings.TrimSuffix(typeNameList[1], "Controller"))
        if len(controllerUri) == 0 {
            continue
        }
        uriPrefix += "/" + controllerUri
        app.registerRouteAction(uriPrefix, controllers[i])
    }
}

// 设置错误处理
func (app *appBasic) SetErrorHandler() {
    app.instance.OnAnyErrorCode(func(ctx context.Context) {
        statusCode := ctx.GetStatusCode()
        logMsg := "HTTP ERROR CODE: " + strconv.Itoa(statusCode) + " URI: " + ctx.Path()
        result := mpresponse.NewResultProblem()
        result.Title = "服务错误"
        result.Status = statusCode

        switch statusCode {
        case iris.StatusNotFound:
            mplog.LogInfo(logMsg)
            result.Tag = "resource-not-exist"
            result.Code = errorcode.CommonRequestResourceEmpty
            result.Msg = "接口地址不存在"
        case iris.StatusMethodNotAllowed:
            mplog.LogInfo(logMsg)
            result.Tag = "method-not-allow"
            result.Code = errorcode.CommonRequestMethod
            result.Msg = "请求方式不支持"
        case iris.StatusInternalServerError:
            mplog.LogError(logMsg)
            result.Tag = "internal-error"
            result.Code = errorcode.CommonBaseServer
            result.Msg = "内部服务出错"
        default:
            mplog.LogError(logMsg)
            result.Tag = "others"
            result.Code = errorcode.CommonBaseServer
            result.Msg = "其他错误"
        }
        ctx.Problem(mpresp.GetProblemHandleBasic(result, 30*time.Second))
        mpresp.NewBasicEnd()(ctx)
    })
}

var (
    insApp *appBasic
)

func init() {
    insApp = &appBasic{}
    insApp.instance = iris.New()
}

func NewApp() *appBasic {
    return insApp
}
