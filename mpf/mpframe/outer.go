/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/6 0006
 * Time: 11:10
 */
package mpframe

import (
    "github.com/kataras/iris/v12"
)

type IOuterWeb interface {
    GetNotify(app *iris.Application) func()
}

type outerWeb struct {
}

func newOuterWeb() outerWeb {
    return outerWeb{}
}
