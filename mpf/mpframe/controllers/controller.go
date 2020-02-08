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
    GetMwActionBefore(tag string) []func(ctx iris.Context) // 获取动作的前置中间件
    GetMwActionAfter(tag string) []func(ctx iris.Context)  // 获取动作的后置中间件
}

type ControllerBasic struct {
    MwControllerBefore []func(ctx iris.Context)            // 控制器前置中间件
    MwControllerAfter  []func(ctx iris.Context)            // 控制器后置中间件
    MwActionBefore     map[string][]func(ctx iris.Context) // 动作前置中间件
    MwActionAfter      map[string][]func(ctx iris.Context) // 动作后置中间件
}

func (c *ControllerBasic) GetMwActionBefore(tag string) []func(ctx iris.Context) {
    mwList := make([]func(ctx iris.Context), 0)
    mwList = append(mwList, c.MwControllerBefore...)

    handles, ok := c.MwActionBefore[tag]
    if ok && (len(handles) > 0) {
        mwList = append(mwList, handles...)
    }

    return mwList
}

func (c *ControllerBasic) GetMwActionAfter(tag string) []func(ctx iris.Context) {
    mwList := make([]func(ctx iris.Context), 0)

    handles, ok := c.MwActionAfter[tag]
    if ok && (len(handles) > 0) {
        mwList = append(mwList, handles...)
    }
    mwList = append(mwList, c.MwControllerAfter...)

    return mwList
}

func NewControllerBasic() ControllerBasic {
    c := ControllerBasic{}
    c.MwControllerBefore = make([]func(ctx iris.Context), 0)
    c.MwControllerBefore = append(c.MwControllerBefore, mpcontroller.NewBasicLog(), mpaction.NewBasicLog())
    c.MwControllerAfter = make([]func(ctx iris.Context), 0)
    c.MwActionBefore = make(map[string][]func(ctx iris.Context))
    c.MwActionAfter = make(map[string][]func(ctx iris.Context))
    return c
}
