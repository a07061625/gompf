/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/15 0015
 * Time: 13:17
 */
package mplogistics

import (
    "github.com/a07061625/gompf/mpf/api"
)

type utilLogistics struct {
    api.UtilApi
}

var (
    insUtil *utilLogistics
)

func init() {
    insUtil = &utilLogistics{api.NewUtilApi()}
}

func NewUtilLogistics() *utilLogistics {
    return insUtil
}
