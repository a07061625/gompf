/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 13:21
 */
package mpserver

import (
    stdContext "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "reflect"
    "regexp"
    "runtime"
    "strconv"
    "strings"
    "sync"
    "syscall"
    "time"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/errorcode"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
    "github.com/spf13/viper"
    "github.com/valyala/tcplisten"
)

type IServerBasic interface {
    bootServer()
    StartServer()
    AddRunConfig(configs ...iris.Configurator)
    SetMwGlobal(isPrefix bool, mwList ...context.Handler)
    SetRoute(controllers ...controllers.IControllerBasic)
}

type basic struct {
    routeFlag  bool // 路由标识 true:已设置 false:未设置
    serverConf *viper.Viper
    app        *iris.Application
    runConfigs []iris.Configurator
}

func (s *basic) AddRunConfig(configs ...iris.Configurator) {
    if len(configs) > 0 {
        s.runConfigs = append(s.runConfigs, configs...)
    }
}

// 设置全局中间件
//   isPrefix: 中间件类型,true:前置 false:后置
func (s *basic) SetMwGlobal(isPrefix bool, mwList ...context.Handler) {
    if len(mwList) == 0 {
        return
    }

    if isPrefix {
        for _, v := range mwList {
            s.app.UseGlobal(v)
        }
    } else {
        for _, v := range mwList {
            s.app.DoneGlobal(v)
        }
    }
}

func (s *basic) formatUri(name string) string {
    match, _ := regexp.MatchString(`^[a-zA-Z]+$`, name)
    if !match {
        return ""
    }

    reg := regexp.MustCompile(`([A-Z])`)
    needStr := reg.ReplaceAllString(name, `-${1}`)
    return strings.ToLower(strings.TrimPrefix(needStr, "-"))
}

func (s *basic) registerActionRoute(groupUri string, controller controllers.IControllerBasic) {
    refControllerVal := reflect.ValueOf(controller)
    methodNum := refControllerVal.NumMethod()
    if methodNum <= 0 {
        return
    }

    groupRoute := s.app.Party(groupUri)
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

        actionTag := s.formatUri(strings.TrimPrefix(methodName, "Action"))
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
        groupRoute.HandleMany(iris.MethodGet, actionUri, actionMwList...)
        groupRoute.HandleMany(iris.MethodPost, actionUri, actionMwList...)
    }
}

func (s *basic) SetRoute(controllers ...controllers.IControllerBasic) {
    if s.routeFlag {
        return
    }
    s.routeFlag = true

    controllerNum := len(controllers)
    if controllerNum <= 0 {
        return
    }

    uriPrefix := ""
    blocks := s.serverConf.GetStringMapString(mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + ".mvc.block.accept")
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
        controllerUri := s.formatUri(strings.TrimSuffix(typeNameList[1], "Controller"))
        if len(controllerUri) == 0 {
            continue
        }
        uriPrefix += "/" + controllerUri
        s.registerActionRoute(uriPrefix, controllers[i])
    }
}

func (s *basic) bootBasic() {
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_host", s.serverConf.GetString(confPrefix+"host")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_port", s.serverConf.GetInt(confPrefix+"port")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_type", s.serverConf.GetString(confPrefix+"type")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("version_min", s.serverConf.GetString(confPrefix+"version.min")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("version_deprecated", s.serverConf.GetString(confPrefix+"version.deprecated")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("version_current", s.serverConf.GetString(confPrefix+"version.current")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("version_max", s.serverConf.GetString(confPrefix+"version.max")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_request", s.serverConf.GetFloat64(confPrefix+"timeout.request")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_controller", s.serverConf.GetFloat64(confPrefix+"timeout.controller")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_action", s.serverConf.GetFloat64(confPrefix+"timeout.action")))

    s.runConfigs = append(s.runConfigs, iris.WithCharset("UTF-8"))
    s.runConfigs = append(s.runConfigs, iris.WithoutStartupLog)
    s.runConfigs = append(s.runConfigs, iris.WithOptimizations)
    s.runConfigs = append(s.runConfigs, iris.WithoutInterruptHandler)
    s.runConfigs = append(s.runConfigs, iris.WithoutBodyConsumptionOnUnmarshal)
    s.runConfigs = append(s.runConfigs, iris.WithoutServerError(iris.ErrServerClosed))

    s.app.ConfigureHost(func(host *iris.Supervisor) {
        host.RegisterOnShutdown(func() {
            mplog.LogInfo("server shut down")
        })
    })
    s.app.OnAnyErrorCode(func(ctx context.Context) {
        logMsg := "HTTP ERROR CODE: " + strconv.Itoa(ctx.GetStatusCode()) + " URI: " + ctx.Path()
        result := mpresponse.NewResultBasic()
        switch ctx.GetStatusCode() {
        case iris.StatusNotFound:
            mplog.LogInfo(logMsg)
            result.Code = errorcode.CommonRequestResourceEmpty
            result.Msg = "接口不存在"
        case iris.StatusMethodNotAllowed:
            mplog.LogInfo(logMsg)
            result.Code = errorcode.CommonRequestMethod
            result.Msg = "请求方式不支持"
        default:
            mplog.LogError(logMsg)
            result.Code = errorcode.CommonBaseServer
            result.Msg = "服务出错"
        }
        ctx.WriteString(mpf.JsonMarshal(result))
        ctx.ContentType(project.HttpContentTypeJson)
        ctx.StopExecution()
    })
}

func (s *basic) listenNotify() {
    go func(app *iris.Application) {
        ch := make(chan os.Signal, 1)
        signal.Notify(ch,
            os.Interrupt,
            syscall.SIGINT,
            os.Kill,
            syscall.SIGKILL,
            syscall.SIGTERM,
        )

        select {
        case s := <-ch:
            mplog.LogInfo("shutdown on signal " + fmt.Sprintf("%#v", s))

            timeout := 5 * time.Second
            ctx, _ := stdContext.WithTimeout(stdContext.Background(), timeout)
            app.Shutdown(ctx)
        }
    }(s.app)
}

func (s *basic) startBasic() {
    s.listenNotify()

    listenCfg := tcplisten.Config{
        ReusePort:   true,
        DeferAccept: true,
        FastOpen:    true,
    }

    listen, err := listenCfg.NewListener("tcp4", mpf.EnvServerDomain())
    if err != nil {
        log.Fatalln("listen error: " + err.Error())
    }
    s.app.Run(iris.Listener(listen), s.runConfigs...)
}

func newBasic(conf *viper.Viper) basic {
    s := basic{}
    s.routeFlag = false
    s.serverConf = conf
    s.app = iris.New()
    s.runConfigs = make([]iris.Configurator, 0)
    return s
}

type basicHttp struct {
    basic
}

func (s *basicHttp) bootServer() {
    s.bootBasic()
}

func (s *basicHttp) StartServer() {
    s.bootServer()
    s.startBasic()
}

var (
    onceBasic sync.Once
    insBasic  IServerBasic
)

func NewBasic(conf *viper.Viper) IServerBasic {
    onceBasic.Do(func() {
        if mpf.EnvServerType() == mpf.EnvServerTypeApi {
            insBasic = &basicHttp{newBasic(conf)}
        }
    })

    return insBasic
}
