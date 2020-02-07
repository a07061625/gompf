/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/7 0007
 * Time: 10:29
 */
package mwroute

import (
    "github.com/kataras/iris/v12"
)

func NewIrisBefore() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        ctx.Next()
    }
}

func NewIrisAfter() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        ctx.Next()
    }
}
