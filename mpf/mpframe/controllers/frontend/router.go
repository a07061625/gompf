/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/9 0009
 * Time: 16:11
 */
package frontend

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

func NewRouter() *router {
    return insRouter
}
