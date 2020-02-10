/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 9:32
 */
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
