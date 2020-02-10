/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 2:09
 */
package mpserver

import (
    "reflect"
    "regexp"
    "runtime"
    "strings"

    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

func (s *serverBase) formatUri(name string) string {
    match, _ := regexp.MatchString(`^[a-zA-Z]+$`, name)
    if !match {
        return ""
    }

    reg := regexp.MustCompile(`([A-Z])`)
    needStr := reg.ReplaceAllString(name, `-${1}`)
    return strings.ToLower(strings.TrimPrefix(needStr, "-"))
}

func (s *serverBase) registerRouteAction(groupUri string, controller controllers.IControllerBasic) {
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

func (s *serverBase) SetRoutes(controllers ...controllers.IControllerBasic) {
    if s.routeFlag {
        return
    }
    s.routeFlag = true

    controllerNum := len(controllers)
    if controllerNum <= 0 {
        return
    }

    uriPrefix := ""
    blocks := s.confServer.GetStringMapString(mpf.EnvType() + "." + mpf.EnvProjectKeyModule() + ".mvc.block.accept")
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
        s.registerRouteAction(uriPrefix, controllers[i])
    }
}
