/**
 * Created by GoLand.
 * User: 姜伟
 * Date: 2020/1/17 0017
 * Time: 16:20
 */
package mpcurrency

import (
    "github.com/a07061625/gompf/mpf/api"
)

type utilCurrency struct {
    api.UtilApi
}

var (
    insUtil *utilCurrency
)

func init() {
    insUtil = &utilCurrency{api.NewUtilApi()}
}

func NewUtilCurrency() *utilCurrency {
    return insUtil
}
