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
    "github.com/kataras/iris/v12/context"
)

type IControllerBasic interface {
    GetMwController(isPrefix bool) []context.Handler         // 获取动作的中间件
    GetMwAction(isPrefix bool, tag string) []context.Handler // 获取动作的中间件
}

type ControllerBasic struct {
    MwControllerPrefix []context.Handler            // 控制器前置中间件
    MwControllerSuffix []context.Handler            // 控制器后置中间件
    MwActionPrefix     map[string][]context.Handler // 动作前置中间件
    MwActionSuffix     map[string][]context.Handler // 动作后置中间件
}

func (c *ControllerBasic) GetMwController(isPrefix bool) []context.Handler {
    if isPrefix {
        return c.MwControllerPrefix
    } else {
        return c.MwControllerSuffix
    }
}

func (c *ControllerBasic) GetMwAction(isPrefix bool, tag string) []context.Handler {
    if isPrefix {
        handles, ok := c.MwActionPrefix[tag]
        if ok {
            return handles
        } else {
            return make([]context.Handler, 0)
        }
    } else {
        handles, ok := c.MwActionSuffix[tag]
        if ok {
            return handles
        } else {
            return make([]context.Handler, 0)
        }
    }
}

func NewControllerBasic() ControllerBasic {
    c := ControllerBasic{}
    c.MwControllerPrefix = make([]context.Handler, 0)
    c.MwControllerPrefix = append(c.MwControllerPrefix, mpcontroller.NewBasicLog(), mpaction.NewBasicLog())
    c.MwControllerSuffix = make([]context.Handler, 0)
    c.MwActionPrefix = make(map[string][]context.Handler)
    c.MwActionSuffix = make(map[string][]context.Handler)
    return c
}
