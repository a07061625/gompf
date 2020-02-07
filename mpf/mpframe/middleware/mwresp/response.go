/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 10:26
 */
package mwresp

import (
    "github.com/a07061625/gompf/mpf/mplog"
    "github.com/kataras/iris/v12"
)

func NewIrisBefore() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        ctx.Next()
    }
}

func NewIrisAfter() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        mplog.LogInfo("request: " + ctx.FullRequestURI() + " exit")
        ctx.Next()
    }
}
