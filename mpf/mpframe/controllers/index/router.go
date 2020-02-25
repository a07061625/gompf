// Package index router
// User: 姜伟
// Time: 2020-02-25 11:09:33
package index

import (
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
)

type router struct {
    controllers map[int]controllers.IControllerBasic
}

func (r *router) GetControllers() map[int]controllers.IControllerBasic {
    return r.controllers
}

var (
    insRouter *router
)

func init() {
    insRouter = &router{}
    insRouter.controllers = make(map[int]controllers.IControllerBasic)
    insRouter.controllers[1] = &indexController{newCommon()}
}

// NewRouter 实例化路由
func NewRouter() *router {
    return insRouter
}
