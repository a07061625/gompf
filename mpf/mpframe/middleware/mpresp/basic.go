/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 12:47
 */
package mpresp

import (
    "github.com/a07061625/gompf/mpf/mpconstant/project"
    "github.com/kataras/iris/v12"
)

func NewBasicClear() func(ctx iris.Context) {
    return func(ctx iris.Context) {
        ctx.Values().Remove(project.DataParamKeyUrl)
    }
}
