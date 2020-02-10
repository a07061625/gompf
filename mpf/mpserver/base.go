/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/11 0011
 * Time: 1:15
 */
package mpserver

import (
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

type IBase interface {
    bootServer()
    StartServer()
    AddRunConfig(configs ...iris.Configurator)
    SetMwGlobal(isPrefix bool, mwList ...context.Handler)
    SetRoute(controllers ...controllers.IControllerBasic)
}
