/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 10:09
 */
package controllers

import (
    "github.com/a07061625/gompf/mpf"
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpaction"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpcontroller"
    "github.com/a07061625/gompf/mpf/mpresponse"
    "github.com/kataras/iris/v12"
)

type IControllerBasic interface {
    GetRegisters() map[string]string
    SendResp() func(ctx iris.Context)                      // 发送响应数据
    GetMwActionBefore(tag string) []func(ctx iris.Context) // 获取动作的前置中间件
    GetMwActionAfter(tag string) []func(ctx iris.Context)  // 获取动作的后置中间件
}

type ControllerBasic struct {
    Registers          map[string]string
    RespResult         *mpresponse.ResultBasic
    MwControllerBefore []func(ctx iris.Context)            // 控制器前置中间件
    MwControllerAfter  []func(ctx iris.Context)            // 控制器后置中间件
    MwActionBefore     map[string][]func(ctx iris.Context) // 动作前置中间件
    MwActionAfter      map[string][]func(ctx iris.Context) // 动作后置中间件
}

// 获取注册列表
//   返回数据格式(tag => method|uri)
//   参考样例 get-list => GET POST|/aaa /bbb
//   表示 ActionGetList方法支持GET,POST两种请求方式,并且绑定了/aaa和/bbb两个路由
func (c *ControllerBasic) GetRegisters() map[string]string {
    return c.Registers
}

func (c *ControllerBasic) SendResp() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        ctx.ResponseWriter().WriteHeader(iris.StatusOK)
        body, _ := ctx.GetBody()
        if len(body) == 0 {
            ctx.ResponseWriter().WriteString(mpf.JsonMarshal(c.RespResult))
            ctx.ContentType(project.HttpContentTypeJson)
        } else {
            ctx.ContentType(project.HttpContentTypeText)
        }

        ctx.Next()
    }
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
    c.Registers = make(map[string]string)
    c.RespResult = mpresponse.NewResultBasic()
    c.MwControllerBefore = make([]func(ctx iris.Context), 0)
    c.MwControllerBefore = append(c.MwControllerBefore, mpcontroller.NewBasicLog(), mpaction.NewBasicLog())
    c.MwControllerAfter = make([]func(ctx iris.Context), 0)
    c.MwActionBefore = make(map[string][]func(ctx iris.Context))
    c.MwActionAfter = make(map[string][]func(ctx iris.Context))
    return c
}
