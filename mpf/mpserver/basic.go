/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 13:21
 */
package mpserver

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "reflect"
    "regexp"
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
    "github.com/spf13/viper"
    "github.com/valyala/tcplisten"
)

type IServerBasic interface {
    bootServer()
    StartServer()
    AddRunConfig(configs ...iris.Configurator)
    SetMwGlobal(mwType bool, mwList ...func(ctx iris.Context))
    SetRoute(controllers ...controllers.IControllerBasic)
}

type basic struct {
    once       sync.Once
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
//   mwType: bool 中间件类型,true:前置 false:后置
func (s *basic) SetMwGlobal(mwType bool, mwList ...func(ctx iris.Context)) {
    if len(mwList) == 0 {
        return
    }

    if mwType {
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

func (s *basic) getRouteParams(tag, val string) (string, string, string, bool) {
    trueTag := strings.TrimSpace(tag)
    if (len(trueTag) == 0) || (len(val) == 0) {
        return "", "", "", false
    }
    valList := strings.Split(val, "|")
    if len(valList) != 2 {
        return "", "", "", false
    }
    reqMethod := strings.TrimSpace(valList[0])
    if len(reqMethod) == 0 {
        return "", "", "", false
    }
    uri := strings.TrimSpace(valList[1])
    if len(uri) == 0 {
        return "", "", "", false
    }
    actionName := "Action" + strings.Title(trueTag)

    return actionName, strings.ToUpper(reqMethod), uri, true
}

func (s *basic) registerActionRoute(groupUri string, controller controllers.IControllerBasic) {
    registers := controller.GetRegisters()
    if len(registers) == 0 {
        return
    }

    refController := reflect.ValueOf(controller)
    routeGroup := s.app.Party(groupUri)
    for tag, val := range registers {
        actionName, reqMethod, uri, ok := s.getRouteParams(tag, val)
        if !ok {
            continue
        }
        refAction := refController.MethodByName(actionName)
        if !refAction.IsValid() {
            continue
        }
        res := refAction.Call(make([]reflect.Value, 0))
        mwList := controller.GetMwActionBefore(tag)
        mwList = append(mwList, res[0].Interface().(func(ctx iris.Context)))
        mwAfter := controller.GetMwActionAfter(tag)
        mwList = append(mwList, mwAfter...)
        mwList = append(mwList, controller.SendResp())
        mwNum := len(mwList)
        for i := 0; i < mwNum; i++ {
            routeGroup.HandleMany(reqMethod, uri, mwList[i])
        }
    }
}

func (s *basic) SetRoute(controllers ...controllers.IControllerBasic) {
    s.once.Do(func() {
        controllerNum := len(controllers)
        if controllerNum == 0 {
            return
        }

        controllerUri := ""
        blocks := s.serverConf.GetStringMapString(mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + ".mvc.block.accept")
        for i := 0; i < controllerNum; i++ {
            objType := reflect.TypeOf(controllers[i])
            typeNameList := strings.Split(objType.String(), ".")
            if len(typeNameList) < 2 {
                continue
            }

            // 校验版块
            packageName := strings.TrimPrefix(typeNameList[0], "*")
            _, ok := blocks[packageName]
            if !ok {
                continue
            }
            controllerUri = "/" + packageName

            // 校验控制器
            if !strings.HasSuffix(typeNameList[1], "Controller") {
                continue
            }
            controllerUri += "/" + s.formatUri(strings.TrimSuffix(typeNameList[1], "Controller"))
            s.registerActionRoute(controllerUri, controllers[i])
        }

        s.app.Any("/{directory:path}", func(ctx iris.Context) {
            result := mpresponse.NewResultBasic()
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
            ctx.ContentType(project.HttpContentTypeJson)
            ctx.Next()
        })
    })
}

func (s *basic) bootBasic() {
    confPrefix := mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + "."
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_host", s.serverConf.GetString(confPrefix+"host")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_port", s.serverConf.GetInt(confPrefix+"port")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("server_type", s.serverConf.GetString(confPrefix+"type")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_request", s.serverConf.GetFloat64(confPrefix+"timeout.request")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_controller", s.serverConf.GetFloat64(confPrefix+"timeout.controller")))
    s.runConfigs = append(s.runConfigs, iris.WithOtherValue("timeout_action", s.serverConf.GetFloat64(confPrefix+"timeout.action")))

    s.runConfigs = append(s.runConfigs, iris.WithCharset("UTF-8"))
    s.runConfigs = append(s.runConfigs, iris.WithoutStartupLog)
    s.runConfigs = append(s.runConfigs, iris.WithOptimizations)
    s.runConfigs = append(s.runConfigs, iris.WithoutInterruptHandler)
    s.runConfigs = append(s.runConfigs, iris.WithoutAutoFireStatusCode)
    s.runConfigs = append(s.runConfigs, iris.WithoutServerError(iris.ErrServerClosed))

    s.app.ConfigureHost(func(host *iris.Supervisor) {
        host.RegisterOnShutdown(func() {
            mplog.LogInfo("server shut down")
        })
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
            ctx, _ := context.WithTimeout(context.Background(), timeout)
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
    s.serverConf = conf
    s.app = iris.New()
    s.runConfigs = make([]iris.Configurator, 0)
    return s
}
