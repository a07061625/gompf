// Package index common
// User: 姜伟
// Time: 2020-02-25 11:09:02
package index

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
