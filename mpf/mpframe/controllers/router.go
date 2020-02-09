/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/9 0009
 * Time: 15:46
 */
package controllers

import (
    "reflect"
    "runtime"
    "sort"

    "github.com/a07061625/gompf/mpf"
)

type controllerItem struct {
    Key int
    Val IControllerBasic
}

type controllerSorter struct {
    Items []controllerItem
}

func (cs controllerSorter) Len() int {
    return len(cs.Items)
}

func (cs controllerSorter) Swap(i, j int) {
    cs.Items[i], cs.Items[j] = cs.Items[j], cs.Items[i]
}

func (cs controllerSorter) Less(i, j int) bool {
    return cs.Items[i].Key < cs.Items[j].Key
}

func newControllerSorter(data map[int]IControllerBasic) *controllerSorter {
    cs := &controllerSorter{}
    cs.Items = make([]controllerItem, 0, len(data))
    for k, v := range data {
        cs.Items = append(cs.Items, controllerItem{k, v})
    }
    return cs
}

type IRouterBasic interface {
    GetControllers() map[int]IControllerBasic
}

type routerBasic struct {
    controllers []IControllerBasic
    routers     map[string]IRouterBasic
}

func (r *routerBasic) GetControllers() []IControllerBasic {
    return r.controllers
}

func (r *routerBasic) RegisterGroup(router IRouterBasic) {
    refRouter := reflect.TypeOf(router)
    method, _ := refRouter.MethodByName("GetControllers")
    methodName := runtime.FuncForPC(method.Func.Pointer()).Name()
    routerKey := mpf.HashMd5(methodName, "")
    _, ok := r.routers[routerKey]
    if ok {
        return
    }
    r.routers[routerKey] = router

    controllers := router.GetControllers()
    controllerNum := len(controllers)
    if controllerNum == 0 {
        return
    }

    cs := newControllerSorter(controllers)
    sort.Sort(cs)
    for i := 0; i < controllerNum; i++ {
        r.controllers = append(r.controllers, cs.Items[i].Val)
    }
}

var (
    insRouter *routerBasic
)

func init() {
    insRouter = &routerBasic{}
    insRouter.controllers = make([]IControllerBasic, 0)
    insRouter.routers = make(map[string]IRouterBasic)
}

func NewRouter() *routerBasic {
    return insRouter
}
