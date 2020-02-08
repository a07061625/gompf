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

func NewIndex() *indexController {
    return &indexController{newCommon()}
}

func (c *indexController) ActionGetName(ctx context.Context) interface{} {
    name := ctx.Params().GetStringDefault("name", "jiangwei")
    result := make(map[string]string)
    result["myname"] = name
    return result
}
