// Package backend common
// User: 姜伟
// Time: 2020-02-25 11:07:41
package backend

import (
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
    "github.com/a07061625/gompf/mpf/mpframe/middleware/mpaction"
)

type common struct {
    controllers.ControllerBasic
}

func newCommon() common {
    c := common{controllers.NewControllerBasic()}
    c.MwControllerPrefix = append(c.MwControllerPrefix, mpaction.NewBasicSignSimple())
    return c
}
