/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 19:31
 */
package mpmvc

import (
    "github.com/kataras/iris/v12"
)

type IControllerSimple interface {
    BeginRequest(ctx iris.Context)
    EndRequest(ctx iris.Context)
    GetMwController() []func(ctx iris.Context)              // 获取控制器中间件
    GetMwActionPrefix() map[string][]func(ctx iris.Context) // 获取动作前置中间件
    GetMwActionSuffix() map[string][]func(ctx iris.Context) // 获取动作后置中间件
}

type controllerSimple struct {
    mwController   []func(ctx iris.Context)
    mwActionPrefix map[string][]func(ctx iris.Context)
    mwActionSuffix map[string][]func(ctx iris.Context)
}

func (c *controllerSimple) GetMwController() []func(ctx iris.Context) {
    return c.mwController
}

func (c *controllerSimple) GetMwActionPrefix() map[string][]func(ctx iris.Context) {
    return c.mwActionPrefix
}

func (c *controllerSimple) GetMwActionSuffix() map[string][]func(ctx iris.Context) {
    return c.mwActionSuffix
}

func (c *controllerSimple) BeginRequest(ctx iris.Context) {
}

func (c *controllerSimple) EndRequest(ctx iris.Context) {
}

func NewControllerSimple() *controllerSimple {
    c := &controllerSimple{}
    c.mwController = make([]func(ctx iris.Context), 0)
    c.mwActionPrefix = make(map[string][]func(ctx iris.Context))
    c.mwActionSuffix = make(map[string][]func(ctx iris.Context))
    return c
}
