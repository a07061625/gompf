// Package backend index
// User: 姜伟
// Time: 2020-02-25 11:07:58
package backend

import (
    "github.com/kataras/iris/v12/context"
)

type indexController struct {
    common
}

func (c *indexController) ActionGetName(ctx context.Context) interface{} {
    result := make(map[string]string)
    result["myname2"] = ctx.URLParamDefault("name", "jiangwei2")
    result["directory2"] = ctx.Params().GetStringDefault("directory", "888888")
    return result
}
