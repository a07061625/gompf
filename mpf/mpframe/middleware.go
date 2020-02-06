/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 22:56
 */
package mpframe

import (
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

type mwGlobal struct {
}

func (mw *mwGlobal) BeforeRequest() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        mplog.LogInfo("request: " + ctx.FullRequestURI() + " enter")
        ctx.Next()
    }
}

func (mw *mwGlobal) AfterRequest() func(ctx iris.Context) {
    return func(ctx iris.Context) {
    }
}

func (mw *mwGlobal) BeforeRouter() func(ctx iris.Context) {
    return func(ctx iris.Context) {
    }
}

func (mw *mwGlobal) AfterRouter() func(ctx iris.Context) {
    return func(ctx iris.Context) {
    }
}

func (mw *mwGlobal) BeforeAction() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        mplog.LogInfo("request: " + ctx.FullRequestURI() + "action: " + ctx.Path() + " enter")
        ctx.Next()
    }
}

func (mw *mwGlobal) AfterAction() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        mplog.LogInfo("request: " + ctx.FullRequestURI() + "action: " + ctx.Path() + " exit")
        ctx.Next()
    }
}

func (mw *mwGlobal) BeforeResponse() func(ctx iris.Context) {
    return func(ctx iris.Context) {
    }
}

func (mw *mwGlobal) SendResponse() func(ctx iris.Context) {
    return func(ctx iris.Context) {
    }
}

func (mw *mwGlobal) AfterResponse() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        mplog.LogInfo("request: " + ctx.FullRequestURI() + " exit")
        ctx.Next()
    }
}

var (
    insMw *mwGlobal
)

func init() {
    insMw = &mwGlobal{}
}

func NewMiddleWare() *mwGlobal {
    return insMw
}
