// Package frontend index
// User: 姜伟
// Time: 2020-02-25 11:10:59
package frontend

import (
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/context"
)

type indexController struct {
    common
}

func (c *indexController) ActionGetName(ctx context.Context) interface{} {
    result := make(map[string]string)
    result["myname"] = ctx.URLParamDefault("name", "jiangwei")
    result["local_name"] = ctx.Tr("mp.name", iris.Map{"Name": "jjj"})
    result["directory"] = ctx.Params().GetStringDefault("directory", "999999")
    return result
}
