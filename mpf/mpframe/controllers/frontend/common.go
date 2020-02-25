// Package frontend common
// User: 姜伟
// Time: 2020-02-25 11:10:11
package frontend

import (
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
)

type common struct {
    controllers.ControllerBasic
}

func newCommon() common {
    c := common{controllers.NewControllerBasic()}
    return c
}
