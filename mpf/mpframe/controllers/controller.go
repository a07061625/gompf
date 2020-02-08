/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 10:09
 */
package controllers

import (
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpaction"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpcontroller"
    "github.com/kataras/iris/v12"
)

type IControllerBasic interface {
    GetMwController(isPrefix bool) []func(ctx iris.Context)         // 获取动作的中间件
    GetMwAction(isPrefix bool, tag string) []func(ctx iris.Context) // 获取动作的中间件
}

type ControllerBasic struct {
    MwControllerPrefix []func(ctx iris.Context)            // 控制器前置中间件
    MwControllerSuffix []func(ctx iris.Context)            // 控制器后置中间件
    MwActionPrefix     map[string][]func(ctx iris.Context) // 动作前置中间件
    MwActionSuffix     map[string][]func(ctx iris.Context) // 动作后置中间件
}

func (c *ControllerBasic) GetMwController(isPrefix bool) []func(ctx iris.Context) {
    if isPrefix {
        return c.MwControllerPrefix
    } else {
        return c.MwControllerSuffix
    }
}

func (c *ControllerBasic) GetMwAction(isPrefix bool, tag string) []func(ctx iris.Context) {
    if isPrefix {
        handles, ok := c.MwActionPrefix[tag]
        if ok {
            return handles
        } else {
            return make([]func(ctx iris.Context), 0)
        }
    } else {
        handles, ok := c.MwActionSuffix[tag]
        if ok {
            return handles
        } else {
            return make([]func(ctx iris.Context), 0)
        }
    }
}

func NewControllerBasic() ControllerBasic {
    c := ControllerBasic{}
    c.MwControllerPrefix = make([]func(ctx iris.Context), 0)
    c.MwControllerPrefix = append(c.MwControllerPrefix, mpcontroller.NewBasicLog(), mpaction.NewBasicLog())
    c.MwControllerSuffix = make([]func(ctx iris.Context), 0)
    c.MwActionPrefix = make(map[string][]func(ctx iris.Context))
    c.MwActionSuffix = make(map[string][]func(ctx iris.Context))
    return c
}
