/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2019/12/25 0025
 * Time: 10:40
 */
package mpmap

import (
    "github.com/a07061625/gompf/mpf/api"
)

type IMapBase interface {
    api.IApiOuter
    GetRespTag() string
}

type utilMap struct {
    api.UtilApi
}

var (
    insUtil *utilMap
)

func init() {
    insUtil = &utilMap{api.NewUtilApi()}
}

func NewUtilMap() *utilMap {
    return insUtil
}
