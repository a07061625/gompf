/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 8:36
 */
package frontend

import (
    "github.com/kataras/iris/v12/context"
)

type indexController struct {
    common
}

func (c *indexController) ActionGetName(ctx context.Context) interface{} {
    result := make(map[string]string)
    result["myname"] = ctx.URLParamDefault("name", "jiangwei")
    result["directory"] = ctx.Params().GetStringDefault("directory", "999999")
    return result
}
