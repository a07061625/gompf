/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/2/8 0008
 * Time: 9:32
 */
package index

import (
    "github.com/a07061625/gompf/mpf/mpframe/controllers"
)

type common struct {
    controllers.ControllerBasic
}

func newCommon() common {
    return common{controllers.NewControllerBasic()}
}
